package main

import (
	"flag"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var (
	allocCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "zfs_pool",
			Name:      "allocated",
			Help:      "The allocated space for this pool",
		},
		[]string{"pool"},
	)

	sizeCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "zfs_pool",
			Name:      "size",
			Help:      "The amount of bytes used in this pool",
		},
		[]string{"pool"},
	)

	freeCounter = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "zfs_pool",
			Name:      "free",
			Help:      "The amount of bytes free in this pool",
		},
		[]string{"pool"},
	)

	capacity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "zfs_pool",
			Name:      "capacity",
			Help:      "The capacity reported by this pool, out of 100",
		},
		[]string{"pool"},
	)
)

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9312", "Address on which to expose metrics and web interface.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		zpoolPath     = flag.String("zfs.zpool-path", "/sbin/zpool", "Path to execute the zpool binary.")
	)
	flag.Parse()

	// Ensure zpool exists.
	findZpoolBinary(*zpoolPath)

	// Begin collection.
	prometheus.MustRegister(allocCounter, sizeCounter, freeCounter, capacity)
	recordMetrics()

	http.Handle(*metricsPath, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, http.DefaultServeMux))
}

func recordMetrics() {
	go func() {
		for {
			// Query metrics.
			pools := queryPoolMetrics()

			// Convert to usable metrics.
			for _, pool := range pools {
				allocCounter.With(prometheus.Labels{"pool": pool.Name}).Set(pool.Alloc)
				sizeCounter.With(prometheus.Labels{"pool": pool.Name}).Set(pool.Size)
				freeCounter.With(prometheus.Labels{"pool": pool.Name}).Set(pool.Free)
				capacity.With(prometheus.Labels{"pool": pool.Name}).Set(pool.Capacity)
			}

			// Query every 10 seconds.
			time.Sleep(10 * time.Second)
		}
	}()
}
