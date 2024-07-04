# Installation

**From github releases**

Find latest suitable tarball link on the [release page](https://github.com/hotenet/airos-exporter/releases)

**Using go install**

 `go install https://github.com/hotenet/airos-exporter@latest`

# Configuration

You can find an documented example in the [config.yml.sample](./config.yml.sample) file.

# Run

```
usage: airos-exporter --config=CONFIG [<flags>]


Flags:
  -h, --[no-]help      Show context-sensitive help (also try --help-long and --help-man).
      --config=CONFIG  Configuration file path
      --[no-]version   Show application version.
```

# Usage

## Exporter http parameter

- **address**: Set device address (host:port). **Required**
- **scheme**: Set protocol used to connect the device (https or http). Default is https.
- **skip-ssl**: When given, skip ssl certificate validation when connecting to the device. Default is unset.

## Prometheus integration

```yaml
  - job_name: airos
    static_configs:
    - targets: 192.168.2.6:443
    metrics_path: /metrics
    scrape_timeout: 1m
    params:
      skip-ssl: [1]
      scheme: https
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_address
      - target_label: __address__
        replacement: localhost:7539
```

## Output

```
# HELP ubnt_cnx_info Connected peer information
# TYPE ubnt_cnx_info gauge
ubnt_cnx_info{device_id="c9419699cbe8dbe77c16c51b40b85357",index="0",lastip="192.168.2.7",mac="F4:E2:C6:86:9F:0C",name="B1E"} 1

# HELP ubnt_cnx_peer_base_signal_dbm Signal strength for peer on the current connection (in dBm)
# TYPE ubnt_cnx_peer_base_signal_dbm gauge
ubnt_cnx_peer_base_signal_dbm{index="0",peer="local"} -46
ubnt_cnx_peer_base_signal_dbm{index="0",peer="remote"} -45

# HELP ubnt_cnx_peer_chain_signal_dbm Signal strength for peer and chain on the current connection (in dBm)
# TYPE ubnt_cnx_peer_chain_signal_dbm gauge
ubnt_cnx_peer_chain_signal_dbm{chain="Chain 0",index="0",peer="local"} -47
ubnt_cnx_peer_chain_signal_dbm{chain="Chain 0",index="0",peer="remote"} -46
ubnt_cnx_peer_chain_signal_dbm{chain="Chain 1",index="0",peer="local"} -52
ubnt_cnx_peer_chain_signal_dbm{chain="Chain 1",index="0",peer="remote"} -51

# HELP ubnt_cnx_peer_cinr_dbm Carrier to Interference + Noise Ratio (CINR) for peer on the current connection (in dBm, requires AirMax)
# TYPE ubnt_cnx_peer_cinr_dbm gauge
ubnt_cnx_peer_cinr_dbm{index="0",peer="local"} 37
ubnt_cnx_peer_cinr_dbm{index="0",peer="remote"} 35

# HELP ubnt_cnx_peer_distance Distance to peer on the current connection
# TYPE ubnt_cnx_peer_distance gauge
ubnt_cnx_peer_distance{index="0",peer="local"} 0
ubnt_cnx_peer_distance{index="0",peer="remote"} 1

# HELP ubnt_cnx_peer_interference_dbm Interference + Noise level for peer on the current connection (in dBm, requires AirMax)
# TYPE ubnt_cnx_peer_interference_dbm gauge
ubnt_cnx_peer_interference_dbm{index="0",peer="local"} -83
ubnt_cnx_peer_interference_dbm{index="0",peer="remote"} -80

# HELP ubnt_cnx_peer_link_info Radio link information for peer on the current connection
# TYPE ubnt_cnx_peer_link_info gauge
ubnt_cnx_peer_link_info{coding="5/6",index="0",modulation="256QAM",multiplexer="MIMO",multiplier="x8",peer="local"} 1
ubnt_cnx_peer_link_info{coding="5/6",index="0",modulation="256QAM",multiplexer="MIMO",multiplier="x8",peer="remote"} 1

# HELP ubnt_cnx_peer_noise_floor_dbm Noise floor for peer on the current connection (in dBm)
# TYPE ubnt_cnx_peer_noise_floor_dbm gauge
ubnt_cnx_peer_noise_floor_dbm{index="0",peer="local"} -85
ubnt_cnx_peer_noise_floor_dbm{index="0",peer="remote"} -88

# HELP ubnt_cnx_peer_signal_dbm Signal strength for peer on the current connection (in dBm, requires AirMax)
# TYPE ubnt_cnx_peer_signal_dbm gauge
ubnt_cnx_peer_signal_dbm{index="0",peer="local"} -48
ubnt_cnx_peer_signal_dbm{index="0",peer="remote"} -49

# HELP ubnt_cnx_peer_tx_power_dbm Transmission power for peer on the current connection (in dBm)
# TYPE ubnt_cnx_peer_tx_power_dbm gauge
ubnt_cnx_peer_tx_power_dbm{index="0",peer="local"} -4
ubnt_cnx_peer_tx_power_dbm{index="0",peer="remote"} -4

# HELP ubnt_cnx_rx_bytes Received data on the current connection (in Bytes)
# TYPE ubnt_cnx_rx_bytes counter
ubnt_cnx_rx_bytes{index="0"} 3.1656288e+07

# HELP ubnt_cnx_rx_capacity_bps Reception capacity on the current connection (in bit/s, requires AirMax)
# TYPE ubnt_cnx_rx_capacity_bps gauge
ubnt_cnx_rx_capacity_bps{index="0"} 3.096576e+08

# HELP ubnt_cnx_rx_packets Received data on the current connection (in Bytes)
# TYPE ubnt_cnx_rx_packets counter
ubnt_cnx_rx_packets{index="0"} 38019

# HELP ubnt_cnx_tx_bytes Transmitted data on the current connection (in Bytes)
# TYPE ubnt_cnx_tx_bytes counter
ubnt_cnx_tx_bytes{index="0"} 4.7931777e+07

# HELP ubnt_cnx_tx_capacity_bps Transmission capacity on the current connection, (in bit/s, requires AirMax)
# TYPE ubnt_cnx_tx_capacity_bps gauge
ubnt_cnx_tx_capacity_bps{index="0"} 3.096576e+08

# HELP ubnt_cnx_tx_latency_ms Transmission latency on the current connection (in millisecond)
# TYPE ubnt_cnx_tx_latency_ms gauge
ubnt_cnx_tx_latency_ms{index="0"} 0

# HELP ubnt_cnx_tx_packets Transmitted packets on the current connection (in Bytes)
# TYPE ubnt_cnx_tx_packets counter
ubnt_cnx_tx_packets{index="0"} 119204

# HELP ubnt_config_info Various settings
# TYPE ubnt_config_info gauge
ubnt_config_info{setting="firewall.eb6tables"} 0
ubnt_config_info{setting="firewall.ebtables"} 0
ubnt_config_info{setting="firewall.ip6tables"} 0
ubnt_config_info{setting="firewall.iptables"} 0
ubnt_config_info{setting="portfw"} 0
ubnt_config_info{setting="services.airview"} 2
ubnt_config_info{setting="services.dhcp6d_stateful"} 0
ubnt_config_info{setting="services.dhcpc"} 1
ubnt_config_info{setting="services.dhcpd"} 0
ubnt_config_info{setting="services.pppoe"} 0
ubnt_config_info{setting="wireless.aprepeater"} 0
ubnt_config_info{setting="wireless.band"} 2
ubnt_config_info{setting="wireless.hide_essid"} 0

# HELP ubnt_local_antenna_gain_dbi Antenna gain (in dBi)
# TYPE ubnt_local_antenna_gain_dbi gauge
ubnt_local_antenna_gain_dbi 16

# HELP ubnt_local_atpc Automatic Transmission Power Control status. (0: Disabled, 1: Adjusting, 2: Automatic, 3: Automatic failure, 4: Automatic limit reached)
# TYPE ubnt_local_atpc gauge
ubnt_local_atpc 4

# HELP ubnt_local_channel_bw_mhz Channel band width (in MHz)
# TYPE ubnt_local_channel_bw_mhz gauge
ubnt_local_channel_bw_mhz 40

# HELP ubnt_local_channel_center_mhz Channel center frequency (in MHz)
# TYPE ubnt_local_channel_center_mhz gauge
ubnt_local_channel_center_mhz 5210

# HELP ubnt_local_channel_mhz Channel frequency (in MHz)
# TYPE ubnt_local_channel_mhz gauge
ubnt_local_channel_mhz 5200

# HELP ubnt_local_connections Number of connected peers
# TYPE ubnt_local_connections gauge
ubnt_local_connections 1

# HELP ubnt_local_dfs Dynamic Frequency Scanning status. (0: inactive, 1: active)
# TYPE ubnt_local_dfs gauge
ubnt_local_dfs 0

# HELP ubnt_local_info Wireless information
# TYPE ubnt_local_info gauge
ubnt_local_info{essid="b1r",ieeemode="11ACVHT40",mac="F4:E2:C6:86:9E:B0",mode="ap-ptp",security="WPA2"} 1

# HELP ubnt_local_noise_floor_dbm Noise floor (in dBi)
# TYPE ubnt_local_noise_floor_dbm gauge
ubnt_local_noise_floor_dbm -85

# HELP ubnt_scrape_error 1 when error occurred while fetching device data
# TYPE ubnt_scrape_error gauge
ubnt_scrape_error 0

# HELP ubnt_sys_cpu_percent CPU usage (in percentage)
# TYPE ubnt_sys_cpu_percent gauge
ubnt_sys_cpu_percent 3.960396

# HELP ubnt_sys_host_info Information about device
# TYPE ubnt_sys_host_info gauge
ubnt_sys_host_info{id="86721c212591f31cba1e70786d4c281e",model="NanoStation 5AC",name="B1R",role="bridge",version="v8.7.5"} 1

# HELP ubnt_sys_load Load average
# TYPE ubnt_sys_load gauge
ubnt_sys_load 0

# HELP ubnt_sys_localtime Internal time (seconds since Epoch)
# TYPE ubnt_sys_localtime gauge
ubnt_sys_localtime 1.720124407e+09

# HELP ubnt_sys_powertime Time elapsed since device is powered on (in seconds)
# TYPE ubnt_sys_powertime gauge
ubnt_sys_powertime 351110

# HELP ubnt_sys_ram_free_bytes Available memory (in Bytes)
# TYPE ubnt_sys_ram_free_bytes gauge
ubnt_sys_ram_free_bytes 2.0656128e+07

# HELP ubnt_sys_ram_total_bytes Total memory (in Bytes)
# TYPE ubnt_sys_ram_total_bytes gauge
ubnt_sys_ram_total_bytes 6.344704e+07

# HELP ubnt_sys_ram_used_bytes Used memory (in Bytes)
# TYPE ubnt_sys_ram_used_bytes gauge
ubnt_sys_ram_used_bytes 4.2790912e+07

# HELP ubnt_sys_uptime Time elapsed since main interface's last (re)initialization (in seconds)
# TYPE ubnt_sys_uptime gauge
ubnt_sys_uptime 342247
```
