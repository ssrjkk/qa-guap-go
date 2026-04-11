package critical

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"go-framework-guap/fixtures"
)

func TestCriticalAuthLogin(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	resp, err := client.AuthService().Login(ctx, "testuser", "testpass")
	if err != nil {
		t.Skipf("Auth not available: %v", err)
	}

	if resp.AccessToken == "" {
		t.Error("Access token is empty")
	}
}

func TestCriticalAuthRefresh(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	resp, err := client.AuthService().Refresh(ctx, "test-refresh-token")
	if err != nil {
		t.Skipf("Auth refresh not available: %v", err)
	}

	if resp.AccessToken == "" {
		t.Error("Access token is empty after refresh")
	}
}

func TestCriticalProfileRetrieval(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	profile, err := client.ProfileService().GetProfile(ctx, "test-token")
	if err != nil {
		t.Skipf("Profile not available: %v", err)
	}

	if profile.ID == "" {
		t.Error("Profile ID is empty")
	}
}

func TestCriticalScheduleRetrieval(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	items, err := client.ScheduleService().GetSchedule(ctx, "test-token", "3101")
	if err != nil {
		t.Skipf("Schedule not available: %v", err)
	}

	if len(items) == 0 {
		t.Error("Schedule is empty")
	}
}

func TestCriticalGradesRetrieval(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	grades, err := client.GradesService().GetGrades(ctx, "test-token", "12345")
	if err != nil {
		t.Skipf("Grades not available: %v", err)
	}

	for _, g := range grades {
		if g.Value < 0 || g.Value > 5 {
			t.Errorf("Grade value out of range: %f", g.Value)
		}
	}
}

func TestNegativeInvalidCredentials(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	_, err := client.AuthService().Login(ctx, "invalid", "invalid")
	if err == nil {
		t.Error("Expected error for invalid credentials")
	}
}

func TestNegativeUnauthorizedAccess(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	_, err := client.ProfileService().GetProfile(ctx, "invalid-token")
	if err == nil {
		t.Error("Expected error for unauthorized access")
	}
}

func TestNegativeInvalidGroupID(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	items, err := client.ScheduleService().GetSchedule(ctx, "test-token", "invalid-group")
	if err == nil && len(items) == 0 {
		t.Log("Empty schedule returned for invalid group")
	}
}

func TestNegativeInvalidStudentID(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	grades, err := client.GradesService().GetGrades(ctx, "test-token", "999999")
	if err == nil && len(grades) == 0 {
		t.Log("Empty grades returned for invalid student")
	}
}

func TestNegativeInvalidEmailFormat(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	teachers, err := client.ScheduleService().GetTeachers(ctx, "test-token")
	if err != nil {
		t.Skipf("Teachers not available: %v", err)
	}

	for _, teacher := range teachers {
		if teacher.Email != "" && !isValidEmail(teacher.Email) {
			t.Errorf("Teacher %s has invalid email: %s", teacher.ID, teacher.Email)
		}
	}
}

func TestNegativeEmptyRequiredFields(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	items, err := client.ScheduleService().GetSchedule(ctx, "test-token", "3101")
	if err != nil {
		t.Skipf("Schedule not available: %v", err)
	}

	for _, item := range items {
		if item.Subject == "" {
			t.Errorf("Schedule item %d has empty subject", item.ID)
		}
		if item.Teacher == "" {
			t.Errorf("Schedule item %d has empty teacher", item.ID)
		}
	}
}

func TestNegativeSiteAvailability(t *testing.T) {
	resp, err := http.Get("https://guap.ru")
	if err != nil {
		t.Fatalf("Site is not available: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestNegativeAPIEndpointNotFound(t *testing.T) {
	resp, err := http.Get("https://guap.ru/api/invalid")
	if err != nil {
		t.Logf("Expected error for invalid endpoint: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNotFound && resp.StatusCode != http.StatusInternalServerError {
		t.Logf("Got status %d for invalid endpoint", resp.StatusCode)
	}
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}
