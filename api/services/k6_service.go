package services

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"prahara-api/models"
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
		cmd := exec.Command("k6", "run", "--out", "influxdb", scriptPath)
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

func (s *K6Service) GetMetrics(ctx context.Context, timeRange string) ([]map[string]interface{}, error) {
	queryAPI := s.InfluxClient.QueryAPI(s.Org)
	query := fmt.Sprintf(`from(bucket: "%s") 
		|> range(start: %s) 
		|> filter(fn: (r) => r["_measurement"] == "http_req_duration" or r["_measurement"] == "http_reqs")
		|> filter(fn: (r) => r["_field"] == "value")
		|> aggregateWindow(every: 10s, fn: mean, createEmpty: false)
		|> yield(name: "mean")`, s.Bucket, timeRange)

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

func (s *K6Service) RunDynamicTest(targetURL, method string, vus int, duration string) (*models.TestRun, error) {
	testRun := models.TestRun{
		Name:         fmt.Sprintf("Quick Storm: %s %s", method, targetURL),
		TargetURL:    targetURL,
		Method:       method,
		VUs:          vus,
		Duration:     duration,
		Status:       "running",
		InfluxBucket: s.Bucket,
		StartedAt:    time.Now(),
	}
	s.DB.Create(&testRun)

	// Generate dynamic k6 script
	scriptContent := fmt.Sprintf(`
import http from 'k6/http';
import { sleep } from 'k6';

export const options = {
  vus: %d,
  duration: '%s',
};

export default function () {
  http.request('%s', '%s');
  sleep(1);
}`, vus, duration, method, targetURL)

	scriptPath := fmt.Sprintf("./dynamic_script_%d.js", testRun.ID)
	os.WriteFile(scriptPath, []byte(scriptContent), 0644)

	go func() {
		defer os.Remove(scriptPath)
		cmd := exec.Command("k6", "run", "--out", "influxdb", scriptPath)
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
