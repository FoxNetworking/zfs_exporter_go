package main

import (
	"errors"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// zpoolPath is what we need to execute to invoke zpool.
var zpoolPath string

// findZpoolBinary determines the path to invoke the "zpool" binary via several methods.
func findZpoolBinary(flagPath string) {
	var err error

	// Ensure the passed zpool location is valid.
	_, err = os.Stat(flagPath)
	if !errors.Is(err, fs.ErrNotExist) {
		// Phew, close call.
		zpoolPath = flagPath
		return
	}

	// It appears the passed value was incorrect.
	// Determine whether zpool is in our path.
	zpoolPath, err = exec.LookPath("zpool")
	if err == nil {
		// We found it! Phew.
		return
	}

	log.Fatal("unable to find zpool binary")
}

// poolStats represents scraped statistics for a pool.
type poolStats struct {
	Name     string
	Alloc    float64
	Size     float64
	Free     float64
	Capacity float64
}

// queryPoolMetrics queries statistics for all pools.
func queryPoolMetrics() []poolStats {
	var pools []poolStats

	// Invoke "zpool list".
	// -H prints output with tab delimiters and no tab names.
	// -p allows exact output in bytes.
	// -o specifies what columns we care about.
	cmd, err := exec.Command(zpoolPath, "list", "-H", "-p", "-o", "name,size,alloc,free,capacity").Output()
	if err != nil {
		log.Fatalf("error while executing zpool binary: %v", err)
	}
	// Trim our trailing newline.
	cmdOutput := strings.TrimSuffix(string(cmd), "\n")

	// Each newline will have output for a different pool.
	poolOutputs := strings.Split(cmdOutput, "\n")
	for _, poolOutput := range poolOutputs {
		// We expect five columns: name, size, allocation, free, and capacity.
		columns := strings.Split(poolOutput, "\t")
		if len(columns) != 5 {
			log.Fatalf("error parsing zpool output: expected 5 columns, got %d", len(columns))
		}

		// Synthesize our pool statistics.
		pools = append(pools, poolStats{
			Name:     columns[0],
			Size:     atof(columns[1]),
			Alloc:    atof(columns[2]),
			Free:     atof(columns[3]),
			Capacity: atof(columns[4]),
		})
	}

	return pools
}

// atof converts a ZFS property value into a float64, usable for Prometheus.
func atof(value string) float64 {
	f, err := strconv.Atoi(value)
	if err != nil {
		panic(err)
	}
	return float64(f)
}
