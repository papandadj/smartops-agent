// +build !windows

package system

import (
	"fmt"
	"github.com/anchnet/smartops-agent/pkg/metric"
	"time"

	log "github.com/cihub/seelog"
	"github.com/shirou/gopsutil/disk"
)

var (
	diskTs    int64
	diskStats map[string]disk.IOCountersStat
)

const (
	diskIOMetric = "system.disk.io.%s"
)

func runIOStatsCheck(t time.Time) ([]metric.MetricSample, error) {
	var samples []metric.MetricSample
	ioMap, err := disk.IOCounters()
	if err != nil {
		log.Errorf("system.IOCheck: could not retrieve io diskStats: %s", err)
		return samples, err
	}
	// timestamp
	now := time.Now().Unix()
	delta := float64(now - diskTs)

	for device, ioStats := range ioMap {
		if diskTs == 0 {
			continue
		}
		lastIOStats, ok := diskStats[device]
		if !ok {
			log.Debug("New device diskStats (possible hotplug) - full diskStats unavailable this iteration.")
			continue
		}
		if delta == 0 {
			log.Debug("No delta to compute - skipping.")
			continue
		}
		tag := make(map[string]string, 1)
		tag["device"] = device
		rBytes := float64(ioStats.ReadBytes - lastIOStats.ReadBytes)
		wBytes := float64(ioStats.WriteBytes - lastIOStats.WriteBytes)
		rCount := float64(ioStats.ReadCount - lastIOStats.ReadCount)
		wCount := float64(ioStats.WriteCount - lastIOStats.WriteCount)
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(diskIOMetric, "byte.read"), rBytes, metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(diskIOMetric, "byte.write"), wBytes, metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(diskIOMetric, "byte.read.sec"), rBytes/delta, metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(diskIOMetric, "byte.write.sec"), wBytes/delta, metric.UnitByte, t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(diskIOMetric, "read.count"), rCount, "", t, tag))
		samples = append(samples, metric.NewServerMetricSample(fmt.Sprintf(diskIOMetric, "write.count"), wCount, "", t, tag))
	}
	diskStats = ioMap
	diskTs = now
	return samples, nil
}
