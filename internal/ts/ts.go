package ts

import (
	"os"
	"path/filepath"
	"time"

	"github.com/bull-cli/bull/internal/config"
	"github.com/nakabonne/tstorage"
)

func dbPath(name string) string {
	dir := filepath.Join(config.TSDir(), name)
	os.MkdirAll(dir, 0755)
	return dir
}

func openStorage(name string) (tstorage.Storage, error) {
	return tstorage.NewStorage(
		tstorage.WithDataPath(dbPath(name)),
		tstorage.WithTimestampPrecision(tstorage.Seconds),
	)
}

func Write(dbName, metric string, value float64, timestamp int64, labels map[string]string) error {
	s, err := openStorage(dbName)
	if err != nil {
		return err
	}
	defer s.Close()

	if timestamp == 0 {
		timestamp = time.Now().Unix()
	}

	var tsLabels []tstorage.Label
	for k, v := range labels {
		tsLabels = append(tsLabels, tstorage.Label{Name: k, Value: v})
	}

	return s.InsertRows([]tstorage.Row{
		{
			Metric:    metric,
			Labels:    tsLabels,
			DataPoint: tstorage.DataPoint{Timestamp: timestamp, Value: value},
		},
	})
}

type DataPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

func QueryRange(dbName, metric string, from, to int64, labels map[string]string) ([]DataPoint, error) {
	s, err := openStorage(dbName)
	if err != nil {
		return nil, err
	}
	defer s.Close()

	var tsLabels []tstorage.Label
	for k, v := range labels {
		tsLabels = append(tsLabels, tstorage.Label{Name: k, Value: v})
	}

	points, err := s.Select(metric, tsLabels, from, to)
	if err != nil {
		return nil, err
	}

	var result []DataPoint
	for _, p := range points {
		result = append(result, DataPoint{Timestamp: p.Timestamp, Value: p.Value})
	}
	return result, nil
}

func ListDBs() ([]string, error) {
	dir := config.TSDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}
