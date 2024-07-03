package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
	"github.com/tidwall/gjson"
)

type handler struct {
	config    UbntConfig
	cookies   []string
}

func NewHandler(config UbntConfig) *handler {
	return &handler{
		config: config,
		cookies: []string{},
	}
}

func (h *handler) buildClient(skipSSL bool) *http.Client {
	return &http.Client{
		Timeout: time.Minute,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: skipSSL,
			},
		},
	}
}

func (h *handler) buildLoginRequest(t target) (*http.Request, error) {
	fields := url.Values{
		"username": []string{h.config.Username},
		"password": []string{h.config.Password},
	}
	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s://%s:%s/api/auth", t.Scheme, t.Host, t.Port),
		strings.NewReader(fields.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (h *handler) buildStatusRequest(t target) (*http.Request, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("%s://%s:%s/status.cgi", t.Scheme, t.Host, t.Port),
		nil,
	)
	if err != nil {
		return nil, err
	}
	cookie_parts := []string{}
	for _, c := range h.cookies {
		parts := strings.Split(c, ";")
		cookie_parts = append(cookie_parts, parts[0])
	}
	req.Header.Set("Cookie", strings.Join(cookie_parts, "; "))
	return req, nil
}


func (h *handler) login(t target) error {
	log := slog.With("endpoint", "login", "target", t)
	log.Debug("starting query")
	req, err := h.buildLoginRequest(t)
	if err != nil {
		log.Error("could not create request", "cause", err)
		return err
	}
	resp, err := h.buildClient(t.InsecureSkipVerify).Do(req)
	if err != nil {
		log.Error("could not create client", "cause", err)
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected http status code '%d'", resp.StatusCode)
		log.Error("error on query", "cause", err)
		return err
	}
	h.cookies = resp.Header.Values("Set-Cookie")
	log.Debug("query done")
	return nil
}


func (h *handler) getStatus(t target, retriable bool) (string, error) {
	log := slog.With("endpoint", "status", "target", t, "retriable", retriable)
	log.Debug("starting query")
	req, err := h.buildStatusRequest(t)
	if err != nil {
		log.Error("could not create request", "cause", err)
		return "", err
	}
	resp, err := h.buildClient(t.InsecureSkipVerify).Do(req)
	if err != nil {
		log.Error("could not create client", "cause", err)
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("unexpected http status %d", resp.StatusCode)
		if resp.StatusCode == http.StatusForbidden && retriable {
			err = h.login(t)
			if err == nil {
				return h.getStatus(t, false)
			}
		}
		log.Error("error on query", "cause", err)
		return "", err
	}
	defer resp.Body.Close()
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("could not read response", "cause", err)
		return "", err
	}
	return string(content), nil
}

func boolVal(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

func getRadioSettings(idx int64, nss int64) []string {
	mux := "SISO"
	if nss > 1 {
		mux = "MIMO"
	}
	switch idx {
	case 0:
		return []string{"BPSK", "1/2", "x1", mux}
	case 1:
		return []string{"QPSK", "1/2", "x2", mux}
	case 2:
		return []string{"QPSK", "3/4", "x2", mux}
	case 3:
		return []string{"16QAM", "1/4", "x4", mux}
	case 4:
		return []string{"16QAM", "3/4", "x4", mux}
	case 5:
		return []string{"64QAM", "2/3", "x6", mux}
	case 6:
		return []string{"64QAM", "3/4", "x6", mux}
	case 7:
		return []string{"64QAM", "5/6", "x6", mux}
	case 8:
		return []string{"256QAM", "3/4", "x8", mux}
	case 9:
		return []string{"256QAM", "5/6", "x8", mux}
	default:
		return []string{"undef", "undef", "undef", "undef"}
	}
}

func (h *handler) Update(set *metricSet, t target) {
	data, err := h.getStatus(t, true)

	set.ScrapeError.Set(0)
	if err != nil {
		set.ScrapeError.Set(1)
		return
	}

	// sys_info
	set.HostInfo.WithLabelValues(
		gjson.Get(data, "host.hostname").String(),
		gjson.Get(data, "host.device_id").String(),
		gjson.Get(data, "host.fwversion").String(),
		gjson.Get(data, "host.devmodel").String(),
		gjson.Get(data, "host.netrole").String(),
	).Set(1)

	// sys
	set.UpTime   .Set(gjson.Get(data, "host.uptime").Float())
	set.PowerTime.Set(gjson.Get(data, "host.power_time").Float())
	set.LoadAvg  .Set(gjson.Get(data, "host.loadavg").Float())
	set.RamTotal .Set(gjson.Get(data, "host.totalram").Float())
	set.RamFree  .Set(gjson.Get(data, "host.freeram").Float())
	set.RamUsed  .Set(gjson.Get(data, "host.totalram").Float() - gjson.Get(data, "host.freeram").Float())
	set.CPU      .Set(gjson.Get(data, "host.cpuload").Float())
	val, _ := time.Parse(time.DateTime, gjson.Get(data, "host.time").String())
	set.Localtime.Set(float64(val.Unix()))

	// config
	set.Config.WithLabelValues("services.dhcpc").Set(boolVal(gjson.Get(data, "services.dhcpc").Bool()))
	set.Config.WithLabelValues("services.dhcpd").Set(boolVal(gjson.Get(data, "services.dhcpd").Bool()))
	set.Config.WithLabelValues("services.dhcp6d_stateful").Set(boolVal(gjson.Get(data, "services.dhcp6d_stateful").Bool()))
	set.Config.WithLabelValues("services.pppoe").Set(boolVal(gjson.Get(data, "services.pppoe").Bool()))
	set.Config.WithLabelValues("services.airview").Set(gjson.Get(data, "services.airview").Float())
	set.Config.WithLabelValues("firewall.iptables").Set(boolVal(gjson.Get(data, "firewall.iptables").Bool()))
	set.Config.WithLabelValues("firewall.ebtables").Set(boolVal(gjson.Get(data, "firewall.ebtables").Bool()))
	set.Config.WithLabelValues("firewall.ip6tables").Set(boolVal(gjson.Get(data, "firewall.ip6tables").Bool()))
	set.Config.WithLabelValues("firewall.eb6tables").Set(boolVal(gjson.Get(data, "firewall.eb6tables").Bool()))
	set.Config.WithLabelValues("portfw").Set(boolVal(gjson.Get(data, "portfw").Bool()))
	set.Config.WithLabelValues("wireless.aprepeater").Set(boolVal(gjson.Get(data, "wireless.aprepeater").Bool()))
	set.Config.WithLabelValues("wireless.hide_essid").Set(gjson.Get(data, "wireless.hide_essid").Float())
	set.Config.WithLabelValues("wireless.hide_essid").Set(gjson.Get(data, "wireless.hide_essid").Float())
	set.Config.WithLabelValues("wireless.band").Set(gjson.Get(data, "wireless.band").Float())

	// wireless info
	set.LocalInfo.WithLabelValues(
		gjson.Get(data, "wireless.essid").String(),
		gjson.Get(data, "wireless.mode").String(),
		gjson.Get(data, "wireless.ieeemode").String(),
		gjson.Get(data, "wireless.apmac").String(),
		gjson.Get(data, "wireless.security").String(),
	).Set(1)

	set.LocalChannel.Set(gjson.Get(data, "wireless.frequency").Float())
	set.LocalChanelCenter.Set(gjson.Get(data, "wireless.center1_freq").Float())
	set.LocalChannelBw.Set(gjson.Get(data, "wireless.chanbw").Float())
	set.LocalGain.Set(gjson.Get(data, "wireless.antenna_gain").Float())
	set.LocalNoiseF.Set(gjson.Get(data, "wireless.noisef").Float())
	set.LocalDFS.Set(gjson.Get(data, "wireless.dfs").Float())
	set.LocalATPC.Set(gjson.Get(data, "wireless.polling.atpc_status").Float())

	// chains
	chains := []string{}
	gjson.Get(data, "chain_names.#.name").ForEach(func(key, value gjson.Result) bool {
		chains = append(chains, value.String())
		return true;
	})

	// cnx
	set.LocalCnx.Set(gjson.Get(data, "wireless.sta.#").Float())

	gjson.Get(data, "wireless.sta").ForEach(func(key, value gjson.Result) bool {
		set.CnxInfo.WithLabelValues(
			key.String(),
			value.Get("mac").String(),
			value.Get("lastip").String(),
			value.Get("remote.hostname").String(),
			value.Get("remote.device_id").String(),
		).Set(1)

		// link
		labels := []string{key.String(), "local"}
		labels = append(labels, getRadioSettings(
			gjson.Get(data, "wireless.rx_idx").Int(),
			gjson.Get(data, "wireless.rx_nss").Int(),
		)...)
		set.CnxPeerLink.WithLabelValues(labels...).Set(1)
		labels = []string{key.String(), "remote"}
		labels = append(labels, getRadioSettings(
			value.Get("rx_idx").Int(),
			value.Get("rx_nss").Int(),
		)...)
		set.CnxPeerLink.WithLabelValues(labels...).Set(1)

		// distance
		set.CnxPeerDistance.WithLabelValues(key.String(), "local") .Set(gjson.Get(data, "wireless.distance").Float())
		set.CnxPeerDistance.WithLabelValues(key.String(), "remote").Set(value.Get("remote.distance").Float())

		// rx/tx stats
		set.CnxRxBytes.WithLabelValues(key.String()).Set(value.Get("stats.rx_bytes").Float())
		set.CnxRxPackets.WithLabelValues(key.String()).Set(value.Get("stats.rx_packets").Float())
		set.CnxTxBytes.WithLabelValues(key.String()).Set(value.Get("stats.tx_bytes").Float())
		set.CnxTxPackets.WithLabelValues(key.String()).Set(value.Get("stats.tx_packets").Float())
		set.CnxTxLatency.WithLabelValues(key.String()).Set(value.Get("tx_latency").Float())

		// airmax
		set.CnxRxCap.WithLabelValues(key.String()).Set(value.Get("airmax.downlink_capacity").Float() * 1024)
		set.CnxTxCap.WithLabelValues(key.String()).Set(value.Get("airmax.uplink_capacity").Float() * 1024)

		// peer signal/interference/cinr
		rxEvm := value.Get("airmax.rx.evm.0|@reverse|0").Float()
		rxRssi := value.Get("airmax.rx.evm.1|@reverse|0").Float()
		rxCINR := value.Get("airmax.rx.cinr").Float()
		txEvm := value.Get("airmax.tx.evm.0|@reverse|0").Float()
		txRssi := value.Get("airmax.tx.evm.1|@reverse|0").Float()
		txCINR := value.Get("airmax.tx.cinr").Float()
		set.CnxPeerAmInterference.WithLabelValues(key.String(), "local").Set(rxRssi - rxEvm - 96)
		set.CnxPeerAmSignal.      WithLabelValues(key.String(), "local").Set(- 96 + rxRssi)
		set.CnxPeerAmCinr.        WithLabelValues(key.String(), "local").Set(rxCINR)
		set.CnxPeerAmInterference.WithLabelValues(key.String(), "remote").Set(txRssi - txEvm - 96)
		set.CnxPeerAmSignal.      WithLabelValues(key.String(), "remote").Set(- 96 + txRssi)
		set.CnxPeerAmCinr.        WithLabelValues(key.String(), "remote").Set(txCINR)

		// peer signals
		set.CnxPeerSignal.WithLabelValues(key.String(),  "local").Set(value.Get("signal").Float())
		set.CnxPeerNoiseF.WithLabelValues(key.String(),  "local").Set(value.Get("noisefloor").Float())
		set.CnxPeerTxPower.WithLabelValues(key.String(), "local").Set(gjson.Get(data, "wireless.txpower").Float())
		set.CnxPeerSignal.WithLabelValues(key.String(), "remote").Set(value.Get("remote.signal").Float())
		set.CnxPeerNoiseF.WithLabelValues(key.String(), "remote").Set(value.Get("remote.noisefloor").Float())
		set.CnxPeerTxPower.WithLabelValues(key.String(), "remote").Set(value.Get("remote.tx_power").Float())

		// peer chain signals
		value.Get("chainrssi").ForEach(func(ckey, jChainRSSI gjson.Result) bool{
			chainRSSI := jChainRSSI.Float()
			if chainRSSI == 0 {
				return false
			}
			rssi := value.Get("rssi").Float()
			signal := value.Get("signal").Float()
			r := signal - rssi
			set.CnxPeerChainSignal.WithLabelValues(
				key.String(),
				"local",
				chains[ckey.Int()],
			).Set(chainRSSI + r)
			return true;
		})

		value.Get("remote.chainrssi").ForEach(func(ckey, jChainRSSI gjson.Result) bool{
			chainRSSI := jChainRSSI.Float()
			if chainRSSI == 0 {
				return false
			}
			rssi := value.Get("remote.rssi").Float()
			signal := value.Get("remote.signal").Float()
			r := signal - rssi
			set.CnxPeerChainSignal.WithLabelValues(
				key.String(),
				"remote",
				chains[ckey.Int()],
			).Set(chainRSSI + r)
			return true;
		})
		return true;
	})
}

// local chain signal to peer
// sta.rssi = 59
// sta.signal = -37
// sta.chainrssi = [58, 52, 0]
// r = signal - rssi = -37 - 59 = -96
// for chain
//   chain_signal = chainrssi[idx] + r = 58 + -96 = -38


// peer chain signal to local
// sta.remote.rssi = 61
// sta.remote.signal = -35
// sta.remote.chainrssi = [61, 51, 0]
// r = signal - rssi // -35 - 61 = -96
// for chain
//  chain_signal = chainrssi[idx] + r // 61 + -96 = -35

// format: signal (chain_rssi[0] / chain_rssi[1]) delta(|chain_rssi[0] - chain_rssi[1]|) dBm
