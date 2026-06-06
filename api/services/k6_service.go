package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"prahara-api/models"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"gorm.io/gorm"
)

type K6Service struct {
	DB           *gorm.DB
	InfluxClient influxdb2.Client
	Token        string
	Org          string
	Bucket       string
}

func NewK6Service(db *gorm.DB, url, token, org, bucket string) *K6Service {
	client := influxdb2.NewClient(url, token)
	return &K6Service{
		DB:           db,
		InfluxClient: client,
		Token:        token,
		Org:          org,
		Bucket:       bucket,
	}
}

func (s *K6Service) RunTest(scriptID uint) (*models.TestRun, error) {
	var script models.TestScript
	if err := s.DB.First(&script, scriptID).Error; err != nil {
		return nil, err
	}

	testRun := models.TestRun{
		TestScriptID: scriptID,
		Name:         script.Name,
		Method:       "SCRIPT",
		Category:     "CUSTOM",
		Status:       "running",
		InfluxBucket: s.Bucket,
		StartedAt:    time.Now(),
	}
	s.DB.Create(&testRun)

	// Save script to temp file
	scriptPath := fmt.Sprintf("./temp_script_%d.js", scriptID)
	os.WriteFile(scriptPath, []byte(script.Content), 0644)

	// Execute k6
	go func() {
		defer os.Remove(scriptPath)
		cmd := exec.Command("k6", "run", "--out", "influxdb", "--tag", fmt.Sprintf("run_id=%d", testRun.ID), "--tag", fmt.Sprintf("category=%s", testRun.Category), scriptPath)
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("K6_INFLUXDB_URL=http://influxdb:8086"),
			fmt.Sprintf("K6_INFLUXDB_ORGANIZATION=%s", s.Org),
			fmt.Sprintf("K6_INFLUXDB_BUCKET=%s", s.Bucket),
			fmt.Sprintf("K6_INFLUXDB_TOKEN=%s", s.Token),
		)
		err := cmd.Run()

		status := "completed"
		if err != nil {
			status = "failed"
		}

		now := time.Now()
		s.DB.Model(&testRun).Updates(models.TestRun{
			Status:     status,
			FinishedAt: &now,
		})
	}()

	return &testRun, nil
}

func (s *K6Service) GetMetrics(ctx context.Context, timeRange string, category string) ([]map[string]interface{}, error) {
	queryAPI := s.InfluxClient.QueryAPI(s.Org)

	filterCategory := ""
	if category != "" && category != "ALL" {
		filterCategory = fmt.Sprintf(`|> filter(fn: (r) => r["category"] == "%s")`, category)
	}

	query := fmt.Sprintf(`from(bucket: "%s") 
		|> range(start: %s) 
		|> filter(fn: (r) => r["_measurement"] == "http_req_duration" or r["_measurement"] == "http_reqs")
		|> filter(fn: (r) => r["_field"] == "value")
		%s
		|> aggregateWindow(every: 2s, fn: mean, createEmpty: false)
		|> yield(name: "mean")`, s.Bucket, timeRange, filterCategory)

	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	var metrics []map[string]interface{}
	for result.Next() {
		metrics = append(metrics, map[string]interface{}{
			"time":  result.Record().Time(),
			"value": result.Record().Value(),
			"field": result.Record().Field(),
			"meas":  result.Record().Measurement(),
		})
	}

	return metrics, nil
}

func (s *K6Service) RunDynamicTest(db *gorm.DB, targetURL, method string, vus int, duration string, category string, scriptID uint, customScript string) (*models.TestRun, error) {
	testName := fmt.Sprintf("Quick Storm: %s %s", method, targetURL)
	scriptToRun := customScript

	if scriptToRun == "" && scriptID > 0 {
		var script models.TestScript
		if err := db.First(&script, scriptID).Error; err != nil {
			return nil, err
		}
		scriptToRun = script.Content
		testName = fmt.Sprintf("Script Storm [%s] -> %s", script.Name, targetURL)
	} else if scriptToRun != "" {
		testName = fmt.Sprintf("In-Context Storm -> %s", targetURL)
	}

	testRun := models.TestRun{
		TestScriptID: scriptID,
		Name:         testName,
		TargetURL:    targetURL,
		Method:       method,
		VUs:          vus,
		Duration:     duration,
		Category:     category,
		Status:       "running",
		InfluxBucket: s.Bucket,
		StartedAt:    time.Now(),
	}
	s.DB.Create(&testRun)

	// Create temp script file
	scriptPath := fmt.Sprintf("./temp_%d.js", testRun.ID)
	if scriptToRun == "" {
		scriptTemplate := `
import http from 'k6/http';
import { sleep, check } from 'k6';

export const options = {
  vus: %d,
  duration: '%s',
};

export default function () {
  const res = http.%s('%s');
  check(res, {
    'status is 200': (r) => r.status === 200,
  });
  sleep(1);
}`
		scriptToRun = fmt.Sprintf(scriptTemplate, vus, duration, strings.ToLower(method), targetURL)
	}

	if err := os.WriteFile(scriptPath, []byte(scriptToRun), 0644); err != nil {
		return nil, err
	}

	go func() {
		defer os.Remove(scriptPath)
		cmd := exec.Command("k6", "run", "--out", "influxdb", "--tag", fmt.Sprintf("run_id=%d", testRun.ID), "--tag", fmt.Sprintf("category=%s", testRun.Category), scriptPath)
		cmd.Env = append(os.Environ(),
			fmt.Sprintf("K6_INFLUXDB_URL=http://influxdb:8086"),
			fmt.Sprintf("K6_INFLUXDB_ORGANIZATION=%s", s.Org),
			fmt.Sprintf("K6_INFLUXDB_BUCKET=%s", s.Bucket),
			fmt.Sprintf("K6_INFLUXDB_TOKEN=%s", s.Token),
		)
		err := cmd.Run()

		status := "completed"
		if err != nil {
			status = "failed"
		}

		now := time.Now()
		s.DB.Model(&testRun).Updates(models.TestRun{
			Status:     status,
			FinishedAt: &now,
		})
	}()

	return &testRun, nil
}

func (s *K6Service) GetRunMetrics(ctx context.Context, runID uint) (map[string]interface{}, error) {
	queryAPI := s.InfluxClient.QueryAPI(s.Org)

	// Query for summary statistics for this specific run
	query := fmt.Sprintf(`from(bucket: "%s") 
		|> range(start: -24h) 
		|> filter(fn: (r) => r["run_id"] == "%d")
		|> filter(fn: (r) => r["_measurement"] == "http_req_duration" or r["_measurement"] == "http_reqs" or r["_measurement"] == "checks")
		|> filter(fn: (r) => r["_field"] == "value")`, s.Bucket, runID)

	result, err := queryAPI.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	summary := map[string]interface{}{
		"avg_latency": 0.0,
		"max_latency": 0.0,
		"requests":    0,
		"success":     0,
		"failed":      0,
	}

	var latencies []float64
	for result.Next() {
		val := result.Record().Value()
		meas := result.Record().Measurement()

		if meas == "http_req_duration" {
			v := val.(float64)
			latencies = append(latencies, v)
			if v > summary["max_latency"].(float64) {
				summary["max_latency"] = v
			}
		} else if meas == "http_reqs" {
			summary["requests"] = summary["requests"].(int) + 1
		} else if meas == "checks" {
			if v, ok := val.(float64); ok && v == 1.0 {
				summary["success"] = summary["success"].(int) + 1
			} else {
				summary["failed"] = summary["failed"].(int) + 1
			}
		}
	}

	if len(latencies) > 0 {
		var sum float64
		for _, l := range latencies {
			sum += l
		}
		summary["avg_latency"] = sum / float64(len(latencies))
	}

	return summary, nil
}
