package ts

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
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

type BatchRow struct {
	Metric    string            `json:"metric"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp,omitempty"`
	Labels    map[string]string `json:"labels,omitempty"`
}

func WriteBatch(dbName string, rows []BatchRow) (int, error) {
	s, err := openStorage(dbName)
	if err != nil {
		return 0, err
	}
	defer s.Close()

	now := time.Now().Unix()
	var tsRows []tstorage.Row
	for _, r := range rows {
		ts := r.Timestamp
		if ts == 0 {
			ts = now
		}
		var labels []tstorage.Label
		for k, v := range r.Labels {
			labels = append(labels, tstorage.Label{Name: k, Value: v})
		}
		tsRows = append(tsRows, tstorage.Row{
			Metric:    r.Metric,
			Labels:    labels,
			DataPoint: tstorage.DataPoint{Timestamp: ts, Value: r.Value},
		})
	}
	return len(tsRows), s.InsertRows(tsRows)
}

func WriteBatchFromNDJSON(dbName, ndjsonFile string) (int, error) {
	f, err := os.Open(ndjsonFile)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	var rows []BatchRow
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var row BatchRow
		if err := json.Unmarshal(line, &row); err != nil {
			return len(rows), err
		}
		rows = append(rows, row)
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	return WriteBatch(dbName, rows)
}

func Latest(dbName, metric string, labels map[string]string) (*DataPoint, error) {
	now := time.Now().Unix()
	points, err := QueryRange(dbName, metric, now-86400*30, now, labels)
	if err != nil {
		return nil, err
	}
	if len(points) == 0 {
		return nil, fmt.Errorf("no data points found for metric %q", metric)
	}
	return &points[len(points)-1], nil
}

func Count(dbName, metric string, from, to int64, labels map[string]string) (int, error) {
	points, err := QueryRange(dbName, metric, from, to, labels)
	if err != nil {
		return 0, err
	}
	return len(points), nil
}

func ExportCSV(dbName, metric string, from, to int64, labels map[string]string, w io.Writer) error {
	points, err := QueryRange(dbName, metric, from, to, labels)
	if err != nil {
		return err
	}
	fmt.Fprintln(w, "timestamp,value")
	for _, p := range points {
		fmt.Fprintf(w, "%d,%f\n", p.Timestamp, p.Value)
	}
	return nil
}

func DeleteDB(name string) error {
	return os.RemoveAll(filepath.Join(config.TSDir(), name))
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
