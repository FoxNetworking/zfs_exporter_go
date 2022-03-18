# zfs_exporter
Monitors a pool, returning various ZFS properties such as the pool's size and allocation.

### Usage
Simply `go get -u github.com/FoxNetworking/zfs_exporter`. This will produce a binary at `~/go/bin/zfs_exporter` (unless configured otherwise.)

Use the following flags to configure to your situation:

```
  -web.listen-address string
    	Address on which to expose metrics and web interface. (default ":9312")
  -web.telemetry-path string
    	Path under which to expose metrics. (default "/metrics")
  -zfs.zpool-path string
    	Path to execute the zpool binary. (default "/sbin/zpool")
```
Note that if /sbin/zpool does not exist, the zpool binary will be search for in $PATH.
Ensure this is properly configured if your zpool binary is in a non-standard location.

Start the service.

You can then add this endpoint as an exporter to Prometheus:

```yaml
scrape_configs:
  - job_name: 'zfs'
    static_configs:
    - targets: ['[::1]:9312']
```

This will give you four new gauge metrics:
 - `zfs_pool_allocated`, the pool's allocated space, in bytes
 - `zfs_pool_size`, the pool's used space, in bytes
 - `zfs_pool_free`, the pool's free space, in bytes
 - `zfs_pool_capacity`, the capacity reported by ZFS as a whole number, out of 100

You can apply these to a graph in various systems and use to whatever fits your situation best.

![Grafana graph containing these two keys](https://owo.whats-th.is/8wB5Aiy.png)

### Wishlist
 - Find a nicer way to scrape command output
 - Handle removing pools nicely
 - Query `/dev/zfs` directly via ioctls in order to avoid libzfs
 - Expose more metrics - PRs/feature requests via issues welcome!
