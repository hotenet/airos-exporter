# Installation

**From github releases**

Find latest suiatable tarball link on the [release page](https://github.com/hotenet/airos-exporter/releases)

**Using go install**

 `go install https://github.com/hotenet/airos-exporter@latest`

# Configuration

You can find an documented exemple in the [config.yml.sample](./config.yml.sample) file.

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

- **host**: Set device hostname. **Required**
- **scheme**: Set protocol used to connect the device (https or http). Default is https.
- **port**: Set port used to connect the device. Default is 443.
- **skip-ssl**: When given, skip ssl certificate validation when connecting to the device. Default is false.

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

```metrics
# HELP ubnt_cnx_info Number of connected peers
# TYPE ubnt_cnx_info gauge
ubnt_cnx_info{device_id="c9419699cbe8dbe77c16c51b40b85357",index="0",lastip="192.168.2.7",mac="F4:E2:C6:86:9F:0C",name="B1E"} 1

# HELP ubnt_cnx_peer_base_signal_dbm Signal strength in the current connection seen by given peer end in dBm
# TYPE ubnt_cnx_peer_base_signal_dbm gauge
ubnt_cnx_peer_base_signal_dbm{index="0",peer="local"} -70
ubnt_cnx_peer_base_signal_dbm{index="0",peer="remote"} -70

# HELP ubnt_cnx_peer_chain_signal_dbm Signal strength of the current chain in the current connection seen by given peer end in dBm
# TYPE ubnt_cnx_peer_chain_signal_dbm gauge
ubnt_cnx_peer_chain_signal_dbm{chain="Chain 0",index="0",peer="local"} -72
ubnt_cnx_peer_chain_signal_dbm{chain="Chain 0",index="0",peer="remote"} -75
ubnt_cnx_peer_chain_signal_dbm{chain="Chain 1",index="0",peer="local"} -74
ubnt_cnx_peer_chain_signal_dbm{chain="Chain 1",index="0",peer="remote"} -72

# HELP ubnt_cnx_peer_cinr_dbm Airmax CINR in the current connection seen by given peer end in dBm
# TYPE ubnt_cnx_peer_cinr_dbm gauge
ubnt_cnx_peer_cinr_dbm{index="0",peer="local"} 27
ubnt_cnx_peer_cinr_dbm{index="0",peer="remote"} 26

# HELP ubnt_cnx_peer_distance Distance to peer the current connection
# TYPE ubnt_cnx_peer_distance gauge
ubnt_cnx_peer_distance{index="0",peer="local"} 0
ubnt_cnx_peer_distance{index="0",peer="remote"} 1

# HELP ubnt_cnx_peer_interference_dbm Airmax interefence + Noise level in the current connection seen by given peer end in dBm
# TYPE ubnt_cnx_peer_interference_dbm gauge
ubnt_cnx_peer_interference_dbm{index="0",peer="local"} -101
ubnt_cnx_peer_interference_dbm{index="0",peer="remote"} -96

# HELP ubnt_cnx_peer_link_info Radio link status for the current connection
# TYPE ubnt_cnx_peer_link_info gauge
ubnt_cnx_peer_link_info{coding="5/6",index="0",modulation="64QAM",multiplexer="MIMO",multiplier="x6",peer="local"} 1
ubnt_cnx_peer_link_info{coding="5/6",index="0",modulation="64QAM",multiplexer="MIMO",multiplier="x6",peer="remote"} 1

# HELP ubnt_cnx_peer_noise_floor_dbm Noise floor in the current connection seen by given peer end in dBm
# TYPE ubnt_cnx_peer_noise_floor_dbm gauge
ubnt_cnx_peer_noise_floor_dbm{index="0",peer="local"} -86
ubnt_cnx_peer_noise_floor_dbm{index="0",peer="remote"} -88

# HELP ubnt_cnx_peer_signal_dbm Airmax signal strength in the current connection seen by given peer end in dBm
# TYPE ubnt_cnx_peer_signal_dbm gauge
ubnt_cnx_peer_signal_dbm{index="0",peer="local"} -73
ubnt_cnx_peer_signal_dbm{index="0",peer="remote"} -73

# HELP ubnt_cnx_peer_tx_power Transmission power for given peer in current connection in dBm
# TYPE ubnt_cnx_peer_tx_power gauge
ubnt_cnx_peer_tx_power{index="0",peer="local"} -4
ubnt_cnx_peer_tx_power{index="0",peer="remote"} -4

# HELP ubnt_cnx_rx_bytes Numbers of received bytes on the current connection
# TYPE ubnt_cnx_rx_bytes gauge
ubnt_cnx_rx_bytes{index="0"} 3.3324538e+07

# HELP ubnt_cnx_rx_capacity_bps Airmax reception capacity in bytes on the current connection in bit/s
# TYPE ubnt_cnx_rx_capacity_bps gauge
ubnt_cnx_rx_capacity_bps{index="0"} 1.79712e+08

# HELP ubnt_cnx_rx_packets Numbers of received packets on the current connection
# TYPE ubnt_cnx_rx_packets gauge
ubnt_cnx_rx_packets{index="0"} 68912

# HELP ubnt_cnx_tx_bytes Numbers of transmitted bytes on the current connection
# TYPE ubnt_cnx_tx_bytes gauge
ubnt_cnx_tx_bytes{index="0"} 1.15582397e+08

# HELP ubnt_cnx_tx_capacity_bps Airmax transmission capacity on the current connection in bit/s
# TYPE ubnt_cnx_tx_capacity_bps gauge
ubnt_cnx_tx_capacity_bps{index="0"} 2.101248e+08

# HELP ubnt_cnx_tx_latency_ms Transmission latency in ms for the current connection
# TYPE ubnt_cnx_tx_latency_ms gauge
ubnt_cnx_tx_latency_ms{index="0"} 0

# HELP ubnt_cnx_tx_packets Numbers of transmitted packets on the current connection
# TYPE ubnt_cnx_tx_packets gauge
ubnt_cnx_tx_packets{index="0"} 282762

# HELP ubnt_config_info Device's settings
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

# HELP ubnt_local_antenna_gain_dbi Antenna gain in dBi
# TYPE ubnt_local_antenna_gain_dbi gauge
ubnt_local_antenna_gain_dbi 16

# HELP ubnt_local_atpc Automatic Transmission Power Control status. 0: Disabled, 1: Adjusting, 2: Automatic, 3: Automatic failure, 4: Automatic limit reached
# TYPE ubnt_local_atpc gauge
ubnt_local_atpc 4

# HELP ubnt_local_channel_bw_mhz Channel band width in MHz
# TYPE ubnt_local_channel_bw_mhz gauge
ubnt_local_channel_bw_mhz 40

# HELP ubnt_local_channel_center_mhz Channel center frequency in MHz
# TYPE ubnt_local_channel_center_mhz gauge
ubnt_local_channel_center_mhz 5210

# HELP ubnt_local_channel_mhz Channel frequency in MHz
# TYPE ubnt_local_channel_mhz gauge
ubnt_local_channel_mhz 5200

# HELP ubnt_local_connections Number of connected peers
# TYPE ubnt_local_connections gauge
ubnt_local_connections 1

# HELP ubnt_local_dfs 1 when Dynamic Frequency Scanning is currently active
# TYPE ubnt_local_dfs gauge
ubnt_local_dfs 0

# HELP ubnt_local_info Information about local configuration
# TYPE ubnt_local_info gauge
ubnt_local_info{essid="b1r",ieeemode="11ACVHT40",mac="F4:E2:C6:86:9E:B0",mode="ap-ptp",security="WPA2"} 1

# HELP ubnt_local_noise_floor_dbm Antenna gain in dBi
# TYPE ubnt_local_noise_floor_dbm gauge
ubnt_local_noise_floor_dbm -86

# HELP ubnt_scrape_error 1 when error occurred while fetching device data
# TYPE ubnt_scrape_error gauge
ubnt_scrape_error 0

# HELP ubnt_sys_cpu_percent Device's CPU usage percentage
# TYPE ubnt_sys_cpu_percent gauge
ubnt_sys_cpu_percent 26

# HELP ubnt_sys_host_info Information about device
# TYPE ubnt_sys_host_info gauge
ubnt_sys_host_info{id="86721c212591f31cba1e70786d4c281e",model="NanoStation 5AC",name="B1R",role="bridge",version="v8.7.5"} 1

# HELP ubnt_sys_load Device's load average
# TYPE ubnt_sys_load gauge
ubnt_sys_load 0

# HELP ubnt_sys_localtime Device's internal time given as seconds since epoch
# TYPE ubnt_sys_localtime gauge
ubnt_sys_localtime 1.627069782e+09

# HELP ubnt_sys_powertime Number of seconds since the device's is powered on
# TYPE ubnt_sys_powertime gauge
ubnt_sys_powertime 303212

# HELP ubnt_sys_ram_free_mb Device's available memory in bytes
# TYPE ubnt_sys_ram_free_mb gauge
ubnt_sys_ram_free_mb 2.0697088e+07

# HELP ubnt_sys_ram_total_mb Device's total memory in bytes
# TYPE ubnt_sys_ram_total_mb gauge
ubnt_sys_ram_total_mb 6.344704e+07

# HELP ubnt_sys_ram_used_mb Device's available memory in bytes
# TYPE ubnt_sys_ram_used_mb gauge
ubnt_sys_ram_used_mb 4.2749952e+07

# HELP ubnt_sys_uptime Number of seconds since the device's main interface has been (re)initialized
# TYPE ubnt_sys_uptime gauge
ubnt_sys_uptime 294349
```
