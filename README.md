# zfs_exporter
Monitors a pool, returning various ZFS properties such as the pool's size and allocation.

### Usage
Ensure libzfs is installed, and `go get -u github.com/FoxNetworking/zfs_exporter`. This will produce a binary at `~/go/bin/zfs_exporter` (unless configured otherwise.)

Use the following flags to configure to your situation:

```
  -web.listen-address string
    	Address on which to expose metrics and web interface. (default ":9312")
  -web.telemetry-path string
    	Path under which to expose metrics. (default "/metrics")
  -zfs.pool-name string
    	Pool to monitor metrics with. (default "zpool")
```

Start the service.

You can then add this endpoint as an exporter to Prometheus:

```yaml
scrape_configs:
  - job_name: 'zfs'
    static_configs:
    - targets: ['[::1]:9312']
```

This will give you three new gauge metrics: `poolname_alloc_size` (the pool's used space, in bytes), `poolname_total_size` (the pool's combined size, in bytes), and `poolname_capacity` (the capacity reported by ZFS as a whole number, out of 100).
You can apply these to a graph in various systems and use to whatever fits your situation best.

![Grafana graph containing these two keys](https://owo.whats-th.is/8wB5Aiy.png)

### Known limitations
Currently, this implementation does not monitor more than one zpool. For its internal needs, this works perfectly - however, in situations involving multiple, this may be desired. PRs are always welcome to expand functionality!

Additionally, there is no handling or testing if the pool is abruptly removed.

This has only been tested on Linux with OpenZFS 0.8.4 and may not work on other platforms supporting ZFS.