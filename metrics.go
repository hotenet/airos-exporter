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
	CnxRxBytes   *prometheus.CounterVec
	CnxRxPackets *prometheus.CounterVec
	CnxRxCap     *prometheus.GaugeVec
	CnxTxBytes   *prometheus.CounterVec
	CnxTxPackets *prometheus.CounterVec
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
				Namespace: namespace, Subsystem: "sys", Name: "uptime", Help: "Time elapsed since main interface's last (re)initialization (in seconds)",
			},
		),
		PowerTime: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "powertime", Help: "Time elapsed since device is powered on (in seconds)",
			},
		),
		Localtime: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "localtime", Help: "Internal time (seconds since Epoch)",
			},
		),
		LoadAvg: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "load", Help: "Load average",
			},
		),
		RamTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "ram_total_bytes", Help: "Total memory (in Bytes)",
			},
		),
		RamFree: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "ram_free_bytes", Help: "Available memory (in Bytes)",
			},
		),
		RamUsed: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "ram_used_bytes", Help: "Used memory (in Bytes)",
			},
		),
		CPU: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "sys", Name: "cpu_percent", Help: "CPU usage (in percentage)",
			},
		),
		Config: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "config", Name: "info", Help: "Various settings",
			}, []string{"setting"},
		),

		
		LocalInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "info", Help: "Wireless information",
			}, []string{"essid", "mode", "ieeemode", "mac", "security"},
		),
		LocalCnx: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "connections", Help: "Number of connected peers",
			},
		),
		LocalChannel: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "channel_mhz", Help: "Channel frequency (in MHz)",
			},
		),
		LocalChanelCenter: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "channel_center_mhz", Help: "Channel center frequency (in MHz)",
			},
		),
		LocalChannelBw: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "channel_bw_mhz", Help: "Channel band width (in MHz)",
			},
		),
		LocalGain: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "antenna_gain_dbi", Help: "Antenna gain (in dBi)",
			},
		),
		LocalNoiseF: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "noise_floor_dbm", Help: "Noise floor (in dBi)",
			},
		),
		LocalDFS: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "dfs", Help: "Dynamic Frequency Scanning status. (0: inactive, 1: active)",
			},
		),
		LocalATPC: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "local", Name: "atpc", Help: "Automatic Transmission Power Control status. (0: Disabled, 1: Adjusting, 2: Automatic, 3: Automatic failure, 4: Automatic limit reached)",
			},
		),


		CnxInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "info", Help: "Connected peer information",
			}, []string{"index", "mac", "lastip", "name", "device_id"},
		),
		CnxRxBytes: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "rx_bytes", Help: "Received data on the current connection (in Bytes)",
			}, []string{"index"},
		),
		CnxRxPackets: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "rx_packets", Help: "Received data on the current connection (in Bytes)",
			}, []string{"index"},
		),
		CnxTxBytes: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "tx_bytes", Help: "Transmitted data on the current connection (in Bytes)",
			}, []string{"index"},
		),
		CnxTxPackets: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "tx_packets", Help: "Transmitted packets on the current connection (in Bytes)",
			}, []string{"index"},
		),
		CnxTxLatency: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "tx_latency_ms", Help: "Transmission latency on the current connection (in millisecond)",
			}, []string{"index"},
		),
		CnxTxCap: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "tx_capacity_bps", Help: "Transmission capacity on the current connection, (in bit/s, requires AirMax)",
			}, []string{"index"},
		),
		CnxRxCap: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "rx_capacity_bps", Help: "Reception capacity on the current connection (in bit/s, requires AirMax)",
			}, []string{"index"},
		),

		CnxPeerLink: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_link_info", Help: "Radio link information for peer on the current connection",
			}, []string{"index", "peer", "modulation", "coding", "multiplier", "multiplexer"},
		),
		CnxPeerTxPower: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_tx_power_dbm", Help: "Transmission power for peer on the current connection (in dBm)",
			}, []string{"index", "peer"},
		),
		CnxPeerChainSignal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_chain_signal_dbm", Help: "Signal strength for peer and chain on the current connection (in dBm)",
			}, []string{"index", "peer", "chain"},
		),
		CnxPeerSignal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_base_signal_dbm", Help: "Signal strength for peer on the current connection (in dBm)",
			}, []string{"index", "peer"},
		),
		CnxPeerNoiseF: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_noise_floor_dbm", Help: "Noise floor for peer on the current connection (in dBm)",
			}, []string{"index", "peer"},
		),
		CnxPeerDistance: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_distance", Help: "Distance to peer on the current connection",
			}, []string{"index", "peer"},
		),


		CnxPeerAmInterference: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_interference_dbm", Help: "Interference + Noise level for peer on the current connection (in dBm, requires AirMax)",
			}, []string{"index", "peer"},
		),
		CnxPeerAmSignal: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_signal_dbm", Help: "Signal strength for peer on the current connection (in dBm, requires AirMax)",
			}, []string{"index", "peer"},
		),
		CnxPeerAmCinr: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace, Subsystem: "cnx", Name: "peer_cinr_dbm", Help: "Carrier to Interference + Noise Ratio (CINR) for peer on the current connection (in dBm, requires AirMax)",
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
