package config

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	baseURL    string
	timeout    int
	maxRetries int
	retryDelay int
	enableLogs bool
	logLevel  string
	parallel  int
	workers   int
}

var (
	cfg     *Config
	cfgOnce sync.Once
	envMap  = map[string]map[string]string{
		"dev": {
			"base_url":    "https://guap.ru",
			"timeout":     "30",
			"max_retries": "3",
			"retry_delay": "500",
			"enable_logs": "true",
			"log_level":   "info",
			"parallel":    "4",
			"workers":     "2",
		},
		"stage": {
			"base_url":    "https://guap.ru",
			"timeout":     "60",
			"max_retries": "5",
			"retry_delay": "1000",
			"enable_logs": "true",
			"log_level":   "debug",
			"parallel":    "8",
			"workers":     "4",
		},
	}
)

func GetBaseURL(env string) string {
	if url, ok := envMap[env]["base_url"]; ok {
		if envURL := os.Getenv("API_BASE_URL"); envURL != "" {
			return envURL
		}
		return url
	}
	return "https://guap.ru"
}

func GetTimeout(env string) int {
	if timeout, ok := envMap[env]["timeout"]; ok {
		if t := os.Getenv("API_TIMEOUT"); t != "" {
			if parsed, err := strconv.Atoi(t); err == nil {
				return parsed
			}
		}
		if parsed, err := strconv.Atoi(timeout); err == nil {
			return parsed
		}
	}
	return 30
}

func GetMaxRetries(env string) int {
	if retries, ok := envMap[env]["max_retries"]; ok {
		if r := os.Getenv("API_MAX_RETRIES"); r != "" {
			if parsed, err := strconv.Atoi(r); err == nil {
				return parsed
			}
		}
		if parsed, err := strconv.Atoi(retries); err == nil {
			return parsed
		}
	}
	return 3
}

func GetRetryDelay(env string) int {
	if delay, ok := envMap[env]["retry_delay"]; ok {
		if d := os.Getenv("API_RETRY_DELAY"); d != "" {
			if parsed, err := strconv.Atoi(d); err == nil {
				return parsed
			}
		}
		if parsed, err := strconv.Atoi(delay); err == nil {
			return parsed
		}
	}
	return 500
}

func GetEnv() string {
	if env := os.Getenv("TEST_ENV"); env != "" {
		return env
	}
	return "dev"
}

func GetLogLevel() string {
	if level := os.Getenv("LOG_LEVEL"); level != "" {
		return level
	}
	return "info"
}

func GetParallel() int {
	if p := os.Getenv("TEST_PARALLEL"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			return parsed
		}
	}
	return 4
}

func GetWorkers() int {
	if w := os.Getenv("TEST_WORKERS"); w != "" {
		if parsed, err := strconv.Atoi(w); err == nil {
			return parsed
		}
	}
	return 2
}

func Load(env string) *Config {
	cfgOnce.Do(func() {
		cfg = &Config{
			baseURL:    GetBaseURL(env),
			timeout:    GetTimeout(env),
			maxRetries: GetMaxRetries(env),
			retryDelay: GetRetryDelay(env),
			enableLogs: os.Getenv("ENABLE_LOGS") == "true",
			logLevel:  GetLogLevel(),
			parallel:  GetParallel(),
			workers:   GetWorkers(),
		}
	})
	return cfg
}

func Get() *Config {
	if cfg == nil {
		return Load(GetEnv())
	}
	return cfg
}
