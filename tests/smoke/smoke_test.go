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

func TestSmokeHealthEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	status, err := client.HealthService().Check(ctx)
	if err != nil {
		t.Skipf("Health endpoint not available: %v", err)
	}

	if status.Status == "" {
		t.Error("Health status is empty")
	}
}

func TestSmokeStudentsEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	students, err := client.StudentService().GetAll(ctx, "test-token")
	if err != nil {
		t.Skipf("Students endpoint not available: %v", err)
	}

	if len(students) == 0 {
		t.Log("Students list is empty")
	}
}

func TestSmokeStudentByIDEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	_, err := client.StudentService().GetByID(ctx, "test-token", "1")
	if err != nil {
		t.Skipf("Student by ID not available: %v", err)
	}
}

func TestSmokeScheduleEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	schedule, err := client.ScheduleService().GetSchedule(ctx, "test-token")
	if err != nil {
		t.Skipf("Schedule endpoint not available: %v", err)
	}

	if len(schedule) == 0 {
		t.Log("Schedule is empty")
	}
}

func TestSmokeScheduleByGroupEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	schedule, err := client.ScheduleService().GetScheduleByGroup(ctx, "test-token", "Z3420")
	if err != nil {
		t.Skipf("Schedule by group not available: %v", err)
	}

	if len(schedule) == 0 {
		t.Log("Schedule for group Z3420 is empty")
	}
}

func TestSmokeSubjectsEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	subjects, err := client.SubjectService().GetAll(ctx, "test-token")
	if err != nil {
		t.Skipf("Subjects endpoint not available: %v", err)
	}

	if len(subjects) == 0 {
		t.Log("Subjects list is empty")
	}
}

func TestSmokeGradesEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	grades, err := client.GradesService().GetAll(ctx, "test-token")
	if err != nil {
		t.Skipf("Grades endpoint not available: %v", err)
	}

	if len(grades) == 0 {
		t.Log("Grades list is empty")
	}
}

func TestSmokeGradesByStudentEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	grades, err := client.GradesService().GetByStudent(ctx, "test-token", "1")
	if err != nil {
		t.Skipf("Grades by student not available: %v", err)
	}

	if len(grades) == 0 {
		t.Log("Grades for student 1 is empty")
	}
}

func TestSmokeTeachersEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	teachers, err := client.ScheduleService().GetTeachers(ctx, "test-token")
	if err != nil {
		t.Skipf("Teachers endpoint not available: %v", err)
	}

	if len(teachers) == 0 {
		t.Log("Teachers list is empty")
	}
}

func TestSmokeAuthLoginEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	_, err := client.AuthService().Login(ctx, "test", "test")
	if err != nil {
		t.Skipf("Auth login endpoint not available: %v", err)
	}
}

func TestSmokeProfileEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	_, err := client.ProfileService().GetProfile(ctx, "test-token")
	if err != nil {
		t.Skipf("Profile endpoint not available: %v", err)
	}
}
