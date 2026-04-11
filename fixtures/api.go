package fixtures

import (
	"context"
	"os"
	"sync"
	"time"

	"go-framework-guap/config"
	"go-framework-guap/core/base"
	"go-framework-guap/services/api"
)

type APIClient struct {
	client      *base.Client
	authSvc     *api.AuthService
	scheduleSvc *api.ScheduleService
	gradesSvc   *api.GradesService
	profileSvc  *api.ProfileService
	baseURL     string
	env         string
	token       string
	initOnce    sync.Once
	initErr     error
}

var (
	defaultClient *APIClient
	defaultOnce  sync.Once
)

func NewAPIClient(env string) *APIClient {
	return &APIClient{
		env:     env,
		baseURL: config.GetBaseURL(env),
	}
}

func (ac *APIClient) Init() error {
	ac.initOnce.Do(func() {
		cfg := base.Config{
			BaseURL:    ac.baseURL,
			Timeout:    30 * time.Second,
			MaxRetries: 3,
			RetryDelay: 500 * time.Millisecond,
		}
		ac.client = base.NewClient(cfg)
		ac.authSvc = api.NewAuthService(ac.client)
		ac.scheduleSvc = api.NewScheduleService(ac.client)
		ac.gradesSvc = api.NewGradesService(ac.client)
		ac.profileSvc = api.NewProfileService(ac.client)
	})
	return ac.initErr
}

func (ac *APIClient) Client() *base.Client {
	if ac.client == nil {
		ac.Init()
	}
	return ac.client
}

func (ac *APIClient) AuthService() *api.AuthService {
	if ac.authSvc == nil {
		ac.Init()
	}
	return ac.authSvc
}

func (ac *APIClient) ScheduleService() *api.ScheduleService {
	if ac.scheduleSvc == nil {
		ac.Init()
	}
	return ac.scheduleSvc
}

func (ac *APIClient) GradesService() *api.GradesService {
	if ac.gradesSvc == nil {
		ac.Init()
	}
	return ac.gradesSvc
}

func (ac *APIClient) ProfileService() *api.ProfileService {
	if ac.profileSvc == nil {
		ac.Init()
	}
	return ac.profileSvc
}

func (ac *APIClient) SetToken(token string) {
	ac.token = token
}

func (ac *APIClient) GetToken() string {
	return ac.token
}

func (ac *APIClient) BaseURL() string {
	return ac.baseURL
}

func (ac *APIClient) Env() string {
	return ac.env
}

func GetDefaultClient() *APIClient {
	defaultOnce.Do(func() {
		env := os.Getenv("TEST_ENV")
		if env == "" {
			env = config.GetEnv()
		}
		defaultClient = NewAPIClient(env)
	})
	return defaultClient
}

func GetEnv() string {
	env := os.Getenv("TEST_ENV")
	if env != "" {
		return env
	}
	return config.GetEnv()
}

type AuthFixture struct {
	client       *APIClient
	token        string
	refreshToken string
}

func NewAuthFixture(client *APIClient) *AuthFixture {
	return &AuthFixture{client: client}
}

func (f *AuthFixture) Setup(ctx context.Context, login, password string) (string, error) {
	f.client.Init()
	resp, err := f.client.AuthService().Login(ctx, login, password)
	if err != nil {
		return "", err
	}
	f.token = resp.AccessToken
	f.refreshToken = resp.RefreshToken
	f.client.SetToken(f.token)
	return f.token, nil
}

func (f *AuthFixture) Teardown(ctx context.Context) error {
	if f.token != "" {
		return f.client.AuthService().Logout(ctx, f.token)
	}
	return nil
}

func (f *AuthFixture) GetToken() string {
	return f.token
}

type ScheduleFixture struct {
	client   *APIClient
	token    string
	groupID  string
	schedule []api.ScheduleItem
}

func NewScheduleFixture(client *APIClient) *ScheduleFixture {
	return &ScheduleFixture{client: client}
}

func (f *ScheduleFixture) Setup(ctx context.Context, token, groupID string) ([]api.ScheduleItem, error) {
	f.client.Init()
	f.token = token
	f.groupID = groupID

	schedule, err := f.client.ScheduleService().GetSchedule(ctx, token, groupID)
	if err != nil {
		return nil, err
	}
	f.schedule = schedule
	return schedule, nil
}

func (f *ScheduleFixture) Get() []api.ScheduleItem {
	return f.schedule
}

func (f *ScheduleFixture) GetByDate(date string) []api.ScheduleItem {
	var result []api.ScheduleItem
	for _, item := range f.schedule {
		if item.Date == date {
			result = append(result, item)
		}
	}
	return result
}

type GradesFixture struct {
	client    *APIClient
	token     string
	studentID string
	grades    []api.Grade
}

func NewGradesFixture(client *APIClient) *GradesFixture {
	return &GradesFixture{client: client}
}

func (f *GradesFixture) Setup(ctx context.Context, token, studentID string) ([]api.Grade, error) {
	f.client.Init()
	f.token = token
	f.studentID = studentID

	grades, err := f.client.GradesService().GetGrades(ctx, token, studentID)
	if err != nil {
		return nil, err
	}
	f.grades = grades
	return grades, nil
}

func (f *GradesFixture) Get() []api.Grade {
	return f.grades
}

func (f *GradesFixture) GetBySubject(subjectID string) []api.Grade {
	var result []api.Grade
	for _, g := range f.grades {
		if g.SubjectID == subjectID {
			result = append(result, g)
		}
	}
	return result
}

type APIFixture struct {
	client   *APIClient
	ctx      context.Context
	cancel   context.CancelFunc
	fixtures []Fixture
}

type Fixture interface {
	Setup(ctx context.Context) error
	Teardown(ctx context.Context) error
}

func NewAPIFixture(client *APIClient) *APIFixture {
	return &APIFixture{client: client}
}

func (f *APIFixture) WithFixtures(fixtures ...Fixture) *APIFixture {
	f.fixtures = append(f.fixtures, fixtures...)
	return f
}

func (f *APIFixture) Setup(ctx context.Context) error {
	f.client.Init()
	f.ctx, f.cancel = context.WithCancel(ctx)

	for _, fix := range f.fixtures {
		if err := fix.Setup(f.ctx); err != nil {
			return err
		}
	}
	return nil
}

func (f *APIFixture) Teardown(ctx context.Context) error {
	for i := len(f.fixtures) - 1; i >= 0; i-- {
		if err := f.fixtures[i].Teardown(ctx); err != nil {
			return err
		}
	}
	if f.cancel != nil {
		f.cancel()
	}
	return nil
}

func (f *APIFixture) Context() context.Context {
	return f.ctx
}

func (f *APIFixture) Client() *APIClient {
	return f.client
}
