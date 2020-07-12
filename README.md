# zfs_exporter
Monitors a pool, returning various ZFS properties such as the pool's size and allocation.

### Known limitations
Currently, this implementation does not monitor more than one zpool. For its internal needs, this works perfectly - however, in situations involving multiple, this may be desired. PRs are always welcome to expand functionality!

Additionally, there is no handling or testing if the pool is abruptly removed.

### Usage
```
Usage of zfs_exporter:
  -web.listen-address string
    	Address on which to expose metrics and web interface. (default ":9312")
  -web.telemetry-path string
    	Path under which to expose metrics. (default "/metrics")
  -zfs.pool-name string
    	Pool to monitor metrics with. (default "zpool")
```