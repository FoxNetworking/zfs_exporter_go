package main

import (
	"flag"
	"github.com/bicomsystems/go-libzfs"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"strconv"
	"time"
)

var poolName string
var sizeCounter prometheus.Gauge
var allocCounter prometheus.Gauge

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9312", "Address on which to expose metrics and web interface.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		poolNameFlag  = flag.String("zfs.pool-name", "zpool", "Pool to monitor metrics with.")
	)
	flag.Parse()

	poolName = *poolNameFlag

	sizeCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: poolName + "_total_size",
		Help: "The size of this pool",
	})

	allocCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: poolName + "_alloc_size",
		Help: "The size allocated in this pool",
	})

	recordMetrics()

	http.Handle(*metricsPath, promhttp.Handler())
	log.Fatal(http.ListenAndServe(*listenAddress, http.DefaultServeMux))
}

func recordMetrics() {
	go func() {
		for {
			pool, err := zfs.PoolOpen(poolName)
			if err != nil {
				panic(err)
			}

			size, err := pool.GetProperty(zfs.PoolPropSize)
			if err != nil {
				panic(err)
			}
			sizeCounter.Set(atof(size))

			alloc, err := pool.GetProperty(zfs.PoolPropAllocated)
			if err != nil {
				panic(err)
			}
			allocCounter.Set(atof(alloc))

			time.Sleep(2 * time.Second)
		}
	}()
}

func atof(p zfs.Property) float64 {
	size := p.Value
	f, err := strconv.Atoi(size)
	if err != nil {
		panic(err)
	}
	return float64(f)
}
