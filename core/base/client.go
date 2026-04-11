package base

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"go-framework-guap/core/errors"
	"go-framework-guap/core/utils"
)

type Config struct {
	BaseURL        string
	Timeout        time.Duration
	MaxRetries     int
	RetryDelay     time.Duration
	RetryCondition func(*http.Response, error) bool
}

type Client struct {
	baseURL string
	client  *http.Client
	config  Config
}

func NewClient(cfg Config) *Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	if cfg.MaxRetries == 0 {
		cfg.MaxRetries = 3
	}
	if cfg.RetryDelay == 0 {
		cfg.RetryDelay = 500 * time.Millisecond
	}
	if cfg.RetryCondition == nil {
		cfg.RetryCondition = defaultRetryCondition
	}

	return &Client{
		baseURL: cfg.BaseURL,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
		config: cfg,
	}
}

func defaultRetryCondition(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}
	return resp.StatusCode >= 500 || resp.StatusCode == 429
}

func (c *Client) buildURL(path string, query map[string]string) string {
	u, _ := url.Parse(c.baseURL + path)
	if len(query) > 0 {
		q := u.Query()
		for k, v := range query {
			q.Add(k, v)
		}
		u.RawQuery = q.Encode()
	}
	return u.String()
}

type RequestOption func(*http.Request)

func WithHeaders(headers map[string]string) RequestOption {
	return func(req *http.Request) {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
}

func WithBody(body interface{}) RequestOption {
	return func(req *http.Request) {
		if body == nil {
			return
		}
		var buf bytes.Buffer
		if s, ok := body.(string); ok {
			buf.WriteString(s)
			req.Header.Set("Content-Type", "application/json")
		} else {
			json.NewEncoder(&buf).Encode(body)
			req.Header.Set("Content-Type", "application/json")
		}
		req.Body = io.NopCloser(&buf)
	}
}

type Response struct {
	StatusCode int
	Headers    http.Header
	Body       []byte
	Duration   time.Duration
	RequestURL string
}

func (c *Client) doRequest(ctx context.Context, method, path string, query map[string]string, opts ...RequestOption) (*Response, error) {
	url := c.buildURL(path, query)

	var body []byte
	for _, opt := range opts {
		if r, ok := opt.(*requestWithBody); ok {
			body = r.body
			break
		}
	}

	reqLog := &utils.RequestLog{
		Method:  method,
		URL:     url,
		Headers: make(map[string]string),
		Body:    string(body),
		Time:    time.Now(),
	}

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for _, opt := range opts {
		opt(req)
	}

	for k := range req.Header {
		reqLog.Headers[k] = req.Header.Get(k)
	}

	utils.LogRequest(reqLog)

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		start := time.Now()

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < c.config.MaxRetries && c.config.RetryCondition(resp, err) {
				utils.Warn("Request failed (attempt %d/%d): %v", attempt+1, c.config.MaxRetries, err)
				time.Sleep(c.config.RetryDelay * time.Duration(attempt+1))
				continue
			}
			return nil, errors.NewRetryableError(err)
		}

		respBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		respLog := &utils.ResponseLog{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Body:       string(respBody),
			Duration:   time.Since(start),
		}

		utils.LogResponse(respLog, url)

		if resp.StatusCode >= 400 && attempt < c.config.MaxRetries && c.config.RetryCondition(resp, nil) {
			utils.Warn("Got %d (attempt %d/%d)", resp.StatusCode, attempt+1, c.config.MaxRetries)
			time.Sleep(c.config.RetryDelay * time.Duration(attempt+1))
			continue
		}

		return &Response{
			StatusCode: resp.StatusCode,
			Headers:    resp.Header,
			Body:       respBody,
			Duration:   time.Since(start),
			RequestURL: url,
		}, nil
	}

	return nil, lastErr
}

type requestWithBody struct {
	body []byte
}

func (r *requestWithBody) Apply(req *http.Request) {
	req.Body = io.NopCloser(bytes.NewReader(r.body))
}

func WithRequestBody(body interface{}) RequestOption {
	var buf bytes.Buffer
	if s, ok := body.(string); ok {
		buf.WriteString(s)
	} else {
		json.NewEncoder(&buf).Encode(body)
	}
	return &requestWithBody{body: buf.Bytes()}
}

func (c *Client) Get(ctx context.Context, path string, query map[string]string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodGet, path, query, opts...)
}

func (c *Client) Post(ctx context.Context, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodPost, path, nil, append(opts, WithRequestBody(body))...)
}

func (c *Client) Put(ctx context.Context, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodPut, path, nil, append(opts, WithRequestBody(body))...)
}

func (c *Client) Patch(ctx context.Context, path string, body interface{}, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodPatch, path, nil, append(opts, WithRequestBody(body))...)
}

func (c *Client) Delete(ctx context.Context, path string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(ctx, http.MethodDelete, path, nil, opts...)
}

func (c *Client) DecodeJSON(resp *Response, target interface{}) error {
	return json.Unmarshal(resp.Body, target)
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ValidateRequired(value, fieldName string) error {
	if strings.TrimSpace(value) == "" {
		return errors.NewValidationError(fieldName, "is required")
	}
	return nil
}

func ValidateRange(value, fieldName string, min, max int) error {
	if len(value) < min || len(value) > max {
		return fmt.Errorf("field %s must be between %d and %d characters", fieldName, min, max)
	}
	return nil
}
