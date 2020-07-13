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

var allocCounter prometheus.Gauge
var capacity prometheus.Gauge
var sizeCounter prometheus.Gauge

func main() {
	var (
		listenAddress = flag.String("web.listen-address", ":9312", "Address on which to expose metrics and web interface.")
		metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		poolNameFlag  = flag.String("zfs.pool-name", "zpool", "Pool to monitor metrics with.")
	)
	flag.Parse()

	poolName = *poolNameFlag

	allocCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: poolName + "_alloc_size",
		Help: "The size allocated in this pool",
	})

	capacity = promauto.NewGauge(prometheus.GaugeOpts{
		Name: poolName + "_capacity",
		Help: "The capacity reported by this pool, out of 100",
	})

	sizeCounter = promauto.NewGauge(prometheus.GaugeOpts{
		Name: poolName + "_total_size",
		Help: "The size of this pool",
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

			alloc, err := pool.GetProperty(zfs.PoolPropAllocated)
			if err != nil {
				panic(err)
			}
			allocCounter.Set(atof(alloc))

			poolCapacity, err := pool.GetProperty(zfs.PoolPropCapacity)
			if err != nil {
				panic(err)
			}
			capacity.Set(atof(poolCapacity))

			size, err := pool.GetProperty(zfs.PoolPropSize)
			if err != nil {
				panic(err)
			}
			sizeCounter.Set(atof(size))

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
