package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"reflect"
)

type metricSet struct {
	reg         *prometheus.Registry
	ScrapeError prometheus.Gauge
	HostInfo    *prometheus.GaugeVec
	UpTime      prometheus.Gauge
	PowerTime   prometheus.Gauge
	Localtime   prometheus.Gauge
	LoadAvg     prometheus.Gauge
	RamTotal    prometheus.Gauge
	RamFree     prometheus.Gauge
	RamUsed     prometheus.Gauge
	CPU         prometheus.Gauge
	Config      *prometheus.GaugeVec

	LocalInfo         *prometheus.GaugeVec
	LocalCnx          prometheus.Gauge
	LocalChannel      prometheus.Gauge
	LocalChanelCenter prometheus.Gauge
	LocalChannelBw    prometheus.Gauge
	LocalGain         prometheus.Gauge
	LocalPower        prometheus.Gauge
	LocalNoiseF       prometheus.Gauge
	LocalDFS          prometheus.Gauge
	LocalATPC         prometheus.Gauge

	CnxInfo      *prometheus.GaugeVec
	CnxRxBytes   *prometheus.GaugeVec
	CnxRxPackets *prometheus.GaugeVec
	CnxRxCap     *prometheus.GaugeVec
	CnxTxBytes   *prometheus.GaugeVec
	CnxTxPackets *prometheus.GaugeVec
	CnxTxLatency *prometheus.GaugeVec
	CnxTxCap     *prometheus.GaugeVec

	CnxPeerTxPower     *prometheus.GaugeVec
	CnxPeerLink        *prometheus.GaugeVec
	CnxPeerChainSignal *prometheus.GaugeVec
	CnxPeerSignal      *prometheus.GaugeVec
	CnxPeerNoiseF      *prometheus.GaugeVec
	CnxPeerDistance    *prometheus.GaugeVec

	CnxPeerAmInterference *prometheus.GaugeVec
	CnxPeerAmSignal       *prometheus.GaugeVec
	CnxPeerAmCinr         *prometheus.GaugeVec
}

func NewMetricSet(namespace string) *metricSet {
	if namespace == "" {
		namespace = "ubnt"
	}
	s := &metricSet{
		reg: prometheus.NewRegistry(),

		ScrapeError: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "scrape", Name: "error", Help: "1 when error occurred while fetching device data",
			},
		),
		HostInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "host_info", Help: "Information about device",
			}, []string{"name", "id", "version", "model", "role"},
		),
		UpTime: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "uptime", Help: "Number of seconds since the device's main interface has been (re)initialized",
			},
		),
		PowerTime: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "powertime", Help: "Number of seconds since the device's is powered on",
			},
		),
		Localtime: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "localtime", Help: "Device's internal time given as seconds since epoch",
			},
		),
		LoadAvg: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "load", Help: "Device's load average",
			},
		),
		RamTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "ram_total_mb", Help: "Device's total memory in bytes",
			},
		),
		RamFree: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "ram_free_mb", Help: "Device's available memory in bytes",
			},
		),
		RamUsed: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "ram_used_mb", Help: "Device's available memory in bytes",
			},
		),
		CPU: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "cpu_percent", Help: "Device's CPU usage percentage",
			},
		),
		Config: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "config", Name: "info", Help: "Device's settings",
			}, []string{"setting"},
		),



		LocalInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "info", Help: "Information about local configuration",
			}, []string{"essid", "mode", "ieeemode", "mac", "security"},
		),
		LocalCnx: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "connections", Help: "Number of connected peers",
			},
		),
		LocalChannel: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "channel_mhz", Help: "Channel frequency in MHz",
			},
		),
		LocalChanelCenter: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "channel_center_mhz", Help: "Channel center frequency in MHz",
			},
		),
		LocalChannelBw: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "channel_bw_mhz", Help: "Channel band width in MHz",
			},
		),
		LocalGain: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "antenna_gain_dbi", Help: "Antenna gain in dBi",
			},
		),
		LocalNoiseF: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "noise_floor_dbm", Help: "Antenna gain in dBi",
			},
		),
		LocalDFS: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "dfs", Help: "1 when Dynamic Frequency Scanning is currently active",
			},
		),
		LocalATPC: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "atpc", Help: "Automatic Transmission Power Control status. 0: Disabled, 1: Adjusting, 2: Automatic, 3: Automatic failure, 4: Automatic limit reached",
			},
		),



		CnxInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "info", Help: "Number of connected peers",
			}, []string{"index", "mac", "lastip", "name", "device_id"},
		),
		CnxRxBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "rx_bytes", Help: "Numbers of received bytes on the current connection",
			}, []string{"index"},
		),
		CnxRxPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "rx_packets", Help: "Numbers of received packets on the current connection",
			}, []string{"index"},
		),
		CnxTxBytes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "tx_bytes", Help: "Numbers of transmitted bytes on the current connection",
			}, []string{"index"},
		),
		CnxTxPackets: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "tx_packets", Help: "Numbers of transmitted packets on the current connection",
			}, []string{"index"},
		),
		CnxTxLatency: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "tx_latency_ms", Help: "Transmission latency in ms for the current connection",
			}, []string{"index"},
		),
		CnxTxCap: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "tx_capacity_bps", Help: "Airmax transmission capacity on the current connection in bit/s",
			}, []string{"index"},
		),
		CnxRxCap: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "rx_capacity_bps", Help: "Airmax reception capacity in bytes on the current connection in bit/s",
			}, []string{"index"},
		),

		CnxPeerLink: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_link_info", Help: "Radio link status for the current connection",
			}, []string{"index", "peer", "modulation", "coding", "multiplier", "multiplexer"},
		),
		CnxPeerTxPower: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_tx_power", Help: "Transmission power for given peer in current connection in dBm",
			}, []string{"index", "peer"},
		),
		CnxPeerChainSignal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_chain_signal_dbm", Help: "Signal strength of the current chain in the current connection seen by given peer end in dBm",
			}, []string{"index", "peer", "chain"},
		),
		CnxPeerSignal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_base_signal_dbm", Help: "Signal strength in the current connection seen by given peer end in dBm",
			}, []string{"index", "peer"},
		),
		CnxPeerNoiseF: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_noise_floor_dbm", Help: "Noise floor in the current connection seen by given peer end in dBm",
			}, []string{"index", "peer"},
		),
		CnxPeerDistance: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_distance", Help: "Distance to peer the current connection",
			}, []string{"index", "peer"},
		),


		CnxPeerAmInterference: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_interference_dbm", Help: "Airmax interefence + Noise level in the current connection seen by given peer end in dBm",
			}, []string{"index", "peer"},
		),
		CnxPeerAmSignal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_signal_dbm", Help: "Airmax signal strength in the current connection seen by given peer end in dBm",
			}, []string{"index", "peer"},
		),
		CnxPeerAmCinr: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_cinr_dbm", Help: "Airmax CINR in the current connection seen by given peer end in dBm",
			}, []string{"index", "peer"},
		),
	}

	return s
}

func (s *metricSet) Registry() *prometheus.Registry {
	xref := reflect.ValueOf(*s)
	for idx := 1; idx < xref.NumField(); idx++ {
		fref := xref.Field(idx)
		if !fref.CanInterface() {
			continue
		}
		fval := fref.Interface()
		metric, ok := fval.(prometheus.Collector)
		if ok {
			s.reg.MustRegister(metric)
		}
	}
	return s.reg
}
