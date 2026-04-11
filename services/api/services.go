package api

import (
	"go-framework-guap/core/base"
)

type AuthService struct {
	client *base.Client
}

func NewAuthService(client *base.Client) *AuthService {
	return &AuthService{client: client}
}

func (s *AuthService) Login(ctx context.Context, login, password string) (*LoginResponse, error) {
	resp, err := s.client.Post(ctx, "/api/auth/login", map[string]string{
		"login":    login,
		"password": password,
	})
	if err != nil {
		return nil, err
	}

	var result LoginResponse
	if err := s.client.DecodeJSON(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	_, err := s.client.Post(ctx, "/api/auth/logout", nil,
		base.WithHeaders(map[string]string{"Authorization": "Bearer " + token}))
	return err
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	resp, err := s.client.Post(ctx, "/api/auth/refresh", map[string]string{
		"refresh_token": refreshToken,
	})
	if err != nil {
		return nil, err
	}

	var result TokenResponse
	if err := s.client.DecodeJSON(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type ScheduleService struct {
	client *base.Client
}

func NewScheduleService(client *base.Client) *ScheduleService {
	return &ScheduleService{client: client}
}

func (s *ScheduleService) GetSchedule(ctx context.Context, token, groupID string) ([]ScheduleItem, error) {
	resp, err := s.client.Get(ctx, "/api/schedule/"+groupID, nil,
		base.WithHeaders(map[string]string{"Authorization": "Bearer " + token}))
	if err != nil {
		return nil, err
	}

	var items []ScheduleItem
	if err := s.client.DecodeJSON(resp, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ScheduleService) GetScheduleByDate(ctx context.Context, token, groupID, date string) ([]ScheduleItem, error) {
	resp, err := s.client.Get(ctx, "/api/schedule/"+groupID, map[string]string{"date": date},
		base.WithHeaders(map[string]string{"Authorization": "Bearer " + token}))
	if err != nil {
		return nil, err
	}

	var items []ScheduleItem
	if err := s.client.DecodeJSON(resp, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ScheduleService) GetTeachers(ctx context.Context, token string) ([]Teacher, error) {
	resp, err := s.client.Get(ctx, "/api/teachers", nil,
		base.WithHeaders(map[string]string{"Authorization": "Bearer " + token}))
	if err != nil {
		return nil, err
	}

	var teachers []Teacher
	if err := s.client.DecodeJSON(resp, &teachers); err != nil {
		return nil, err
	}
	return teachers, nil
}

func (s *ScheduleService) GetTeacherSchedule(ctx context.Context, token, teacherID string) ([]ScheduleItem, error) {
	resp, err := s.client.Get(ctx, "/api/teachers/"+teacherID+"/schedule", nil,
		base.WithHeaders(map[string]string{"Authorization": "Bearer " + token}))
	if err != nil {
		return nil, err
	}

	var items []ScheduleItem
	if err := s.client.DecodeJSON(resp, &items); err != nil {
		return nil, err
	}
	return items, nil
}

type GradesService struct {
	client *base.Client
}

func NewGradesService(client *base.Client) *GradesService {
	return &GradesService{client: client}
}

func (s *GradesService) GetGrades(ctx context.Context, token, studentID string) ([]Grade, error) {
	resp, err := s.client.Get(ctx, "/api/students/"+studentID+"/grades", nil,
		base.WithHeaders(map[string]string{"Authorization": "Bearer " + token}))
	if err != nil {
		return nil, err
	}

	var grades []Grade
	if err := s.client.DecodeJSON(resp, &grades); err != nil {
		return nil, err
	}
	return grades, nil
}

func (s *GradesService) GetSubjectGrades(ctx context.Context, token, studentID, subjectID string) ([]Grade, error) {
	resp, err := s.client.Get(ctx, "/api/students/"+studentID+"/subjects/"+subjectID+"/grades", nil,
		base.WithHeaders(map[string]string{"Authorization": "Bearer " + token}))
	if err != nil {
		return nil, err
	}

	var grades []Grade
	if err := s.client.DecodeJSON(resp, &grades); err != nil {
		return nil, err
	}
	return grades, nil
}

func (s *GradesService) GetGPA(ctx context.Context, token, studentID string) (*GPA, error) {
	resp, err := s.client.Get(ctx, "/api/students/"+studentID+"/gpa", nil,
		base.WithHeaders(map[string]string{"Authorization": "Bearer " + token}))
	if err != nil {
		return nil, err
	}

	var gpa GPA
	if err := s.client.DecodeJSON(resp, &gpa); err != nil {
		return nil, err
	}
	return &gpa, nil
}

type ProfileService struct {
	client *base.Client
}

func NewProfileService(client *base.Client) *ProfileService {
	return &ProfileService{client: client}
}

func (s *ProfileService) GetProfile(ctx context.Context, token string) (*Student, error) {
	resp, err := s.client.Get(ctx, "/api/profile", nil,
		base.WithHeaders(map[string]string{"Authorization": "Bearer " + token}))
	if err != nil {
		return nil, err
	}

	var student Student
	if err := s.client.DecodeJSON(resp, &student); err != nil {
		return nil, err
	}
	return &student, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, token string, updates *ProfileUpdate) (*Student, error) {
	resp, err := s.client.Patch(ctx, "/api/profile", updates,
		base.WithHeaders(map[string]string{"Authorization": "Bearer " + token}))
	if err != nil {
		return nil, err
	}

	var student Student
	if err := s.client.DecodeJSON(resp, &student); err != nil {
		return nil, err
	}
	return &student, nil
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type ScheduleItem struct {
	ID        int    `json:"id"`
	Date      string `json:"date"`
	TimeStart string `json:"time_start"`
	TimeEnd   string `json:"time_end"`
	Subject   string `json:"subject"`
	Teacher   string `json:"teacher"`
	Room      string `json:"room"`
	GroupID   string `json:"group_id"`
	Type      string `json:"type"`
}

type Teacher struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Department string `json:"department"`
}

type Grade struct {
	ID        int     `json:"id"`
	SubjectID string  `json:"subject_id"`
	Subject   string  `json:"subject"`
	Value     float64 `json:"value"`
	Date      string  `json:"date"`
	Type      string  `json:"type"`
}

type GPA struct {
	Current float64 `json:"current"`
	Total   float64 `json:"total"`
	Credit  int     `json:"credits"`
}

type Student struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Group        string `json:"group"`
	Course       int    `json:"course"`
}

type ProfileUpdate struct {
	Email string `json:"email,omitempty"`
	Phone string `json:"phone,omitempty"`
}
