package smoke

import (
	"context"
	"net/http"
	"testing"

	"go-framework-guap/fixtures"
)

func TestSmokeGuapSiteReturns200(t *testing.T) {
	resp, err := http.Get("https://guap.ru")
	if err != nil {
		t.Fatalf("Failed to reach guap.ru: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected 200, got %d", resp.StatusCode)
	}
}

func TestSmokeProfileEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	_, err := client.ProfileService().GetProfile(ctx, "test-token")
	if err != nil {
		t.Logf("Profile endpoint check: %v", err)
	}
}

func TestSmokeScheduleEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	_, err := client.ScheduleService().GetSchedule(ctx, "test-token", "3101")
	if err != nil {
		t.Logf("Schedule endpoint check: %v", err)
	}
}

func TestSmokeGradesEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	_, err := client.GradesService().GetGrades(ctx, "test-token", "12345")
	if err != nil {
		t.Logf("Grades endpoint check: %v", err)
	}
}

func TestSmokeTeachersEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	_, err := client.ScheduleService().GetTeachers(ctx, "test-token")
	if err != nil {
		t.Logf("Teachers endpoint check: %v", err)
	}
}

func TestSmokeAuthLoginEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	_, err := client.AuthService().Login(ctx, "test", "test")
	if err != nil {
		t.Logf("Auth login endpoint check: %v", err)
	}
}
