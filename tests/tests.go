package tests

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"go-framework-guap/config"
	"go-framework-guap/core/errors"
	"go-framework-guap/fixtures"
)

var (
	flakyRetries   = flag.Int("flaky-retries", 1, "Number of retries for flaky tests")
	flakyThreshold = flag.Int("flaky-threshold", 2, "Number of failures before marking test as flaky")
	parallel       = flag.Int("parallel", 4, "Number of parallel tests")
	workers        = flag.Int("workers", 2, "Number of test workers")
	testEnv        = flag.String("env", "dev", "Test environment (dev, stage)")
)

type TestMetrics struct {
	mu           sync.Mutex
	passed       int
	failed       int
	skipped      int
	flaky        int
	total        int
	retryCount   int
	testsByName  map[string]*TestResult
}

type TestResult struct {
	Name         string
	Status       string
	Duration     time.Duration
	Attempts     int
	IsFlaky      bool
	ErrorMessage string
}

var metrics = &TestMetrics{
	testsByName: make(map[string]*TestResult),
}

func TestMain(m *testing.M) {
	flag.Parse()

	config.Load(*testEnv)

	_ = parallel
	_ = workers
	_ = flakyRetries
	_ = flakyThreshold

	code := m.Run()
	os.Exit(code)
}

func RunTestWithRetry(t *testing.T, fn func(t *testing.T)) {
	maxRetries := *flakyRetries + 1

	for attempt := 1; attempt <= maxRetries; attempt++ {
		metrics.mu.Lock()
		metrics.retryCount++
		metrics.mu.Unlock()

		t.Logf("Attempt %d/%d", attempt, maxRetries)

		fn(t)

		if !t.Failed() {
			return
		}

		if attempt < maxRetries {
			t.Logf("Test failed, retrying in 1 second...")
			time.Sleep(1 * time.Second)
		}
	}

	metrics.mu.Lock()
	if t.Failed() {
		metrics.failed++
	} else {
		metrics.passed++
	}
	metrics.total++
	metrics.mu.Unlock()
}

func ReportTestResult(t *testing.T, status string, duration time.Duration) {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	result := &TestResult{
		Name:     t.Name(),
		Status:   status,
		Duration: duration,
		Attempts: 1,
	}

	if status == "passed" {
		metrics.passed++
	} else if status == "failed" {
		metrics.failed++
	} else if status == "skipped" {
		metrics.skipped++
	}

	metrics.total++
	metrics.testsByName[t.Name()] = result
}

func PrintMetrics() {
	metrics.mu.Lock()
	defer metrics.mu.Unlock()

	fmt.Println("\n=== Test Metrics ===")
	fmt.Printf("Total:  %d\n", metrics.total)
	fmt.Printf("Passed: %d\n", metrics.passed)
	fmt.Printf("Failed: %d\n", metrics.failed)
	fmt.Printf("Skipped: %d\n", metrics.skipped)
	fmt.Printf("Flaky:  %d\n", metrics.flaky)
	fmt.Printf("Retries: %d\n", metrics.retryCount)

	if metrics.total > 0 {
		passRate := float64(metrics.passed) / float64(metrics.total) * 100
		fmt.Printf("Pass Rate: %.2f%%\n", passRate)
	}

	fmt.Println("\n=== Failed Tests ===")
	for name, result := range metrics.testsByName {
		if result.Status == "failed" {
			fmt.Printf("- %s\n", name)
			if result.ErrorMessage != "" {
				fmt.Printf("  Error: %s\n", result.ErrorMessage)
			}
		}
	}
}

type TestSuite struct {
	Name      string
	Tests     []testing.InternalTest
	Fragile   bool
	Retries   int
	Timeout   time.Duration
	SkipOnFail bool
}

func (ts *TestSuite) Run(ctx context.Context, t *testing.T) {
	for _, test := range ts.Tests {
		t.Run(test.Name, func(t *testing.T) {
			if ts.SkipOnFail && metrics.failed > 0 {
				t.Skip("Skipping due to previous failure (fail-fast)")
			}

			start := time.Now()
			test.F(t)

			duration := time.Since(start)
			if t.Failed() && ts.Fragile {
				for i := 0; i < ts.Retries; i++ {
					t.Logf("Retrying flaky test (attempt %d/%d)", i+1, ts.Retries)
					test.F(t)
					if !t.Failed() {
						break
					}
				}
			}

			if t.Failed() {
				ReportTestResult(t, "failed", duration)
			} else {
				ReportTestResult(t, "passed", duration)
			}
		})
	}
}

type WaitFunc func(ctx context.Context) error

func WaitForCondition(ctx context.Context, fn func() bool, timeout time.Duration, interval time.Duration) error {
	deadline := time.Now().Add(timeout)

	for {
		if fn() {
			return nil
		}

		if time.Now().After(deadline) {
			return errors.NewRetryableError(fmt.Errorf("condition timeout after %v", timeout))
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(interval):
		}
	}
}

func AssertWithRetry(t *testing.T, fn func() bool, retries int, interval time.Duration) {
	for i := 0; i < retries; i++ {
		if fn() {
			return
		}
		time.Sleep(interval)
	}
	if !fn() {
		t.Error("Assertion failed after retries")
	}
}
