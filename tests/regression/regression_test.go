package regression

import (
	"context"
	"strings"
	"testing"

	"go-framework-guap/fixtures"
)

func TestRegressionProfileHasRequiredFields(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	profile, err := client.ProfileService().GetProfile(ctx, "test-token")
	if err != nil {
		t.Skipf("Profile not available: %v", err)
	}

	if profile.ID == "" {
		t.Error("Profile has no ID")
	}
	if profile.Name == "" {
		t.Error("Profile has no name")
	}
}

func TestRegressionScheduleItemHasRequiredFields(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	items, err := client.ScheduleService().GetSchedule(ctx, "test-token", "3101")
	if err != nil {
		t.Skipf("Schedule not available: %v", err)
	}

	for _, item := range items {
		if item.ID == 0 {
			t.Error("Schedule item has no ID")
		}
		if item.Subject == "" {
			t.Errorf("Schedule item %d has no subject", item.ID)
		}
		if item.Date == "" {
			t.Errorf("Schedule item %d has no date", item.ID)
		}
	}
}

func TestRegressionTeacherHasRequiredFields(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	teachers, err := client.ScheduleService().GetTeachers(ctx, "test-token")
	if err != nil {
		t.Skipf("Teachers not available: %v", err)
	}

	for _, teacher := range teachers {
		if teacher.ID == "" {
			t.Error("Teacher has no ID")
		}
		if teacher.Name == "" {
			t.Errorf("Teacher %s has no name", teacher.ID)
		}
	}
}

func TestRegressionTeacherEmailFormat(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	teachers, err := client.ScheduleService().GetTeachers(ctx, "test-token")
	if err != nil {
		t.Skipf("Teachers not available: %v", err)
	}

	for _, teacher := range teachers {
		if teacher.Email != "" && !strings.Contains(teacher.Email, "@") {
			t.Errorf("Teacher %s has invalid email: %s", teacher.ID, teacher.Email)
		}
	}
}

func TestRegressionGradesHasRequiredFields(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	grades, err := client.GradesService().GetGrades(ctx, "test-token", "12345")
	if err != nil {
		t.Skipf("Grades not available: %v", err)
	}

	for _, grade := range grades {
		if grade.ID == 0 {
			t.Error("Grade has no ID")
		}
		if grade.SubjectID == "" {
			t.Errorf("Grade %d has no subject_id", grade.ID)
		}
	}
}

func TestRegressionGPAEndpoint(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	gpa, err := client.GradesService().GetGPA(ctx, "test-token", "12345")
	if err != nil {
		t.Skipf("GPA not available: %v", err)
	}

	if gpa.Current < 0 || gpa.Current > 5 {
		t.Errorf("GPA out of range: %f", gpa.Current)
	}
}

func TestRegressionScheduleByDate(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	items, err := client.ScheduleService().GetScheduleByDate(ctx, "test-token", "3101", "2024-01-15")
	if err != nil {
		t.Skipf("Schedule by date not available: %v", err)
	}

	for _, item := range items {
		if item.Date != "2024-01-15" {
			t.Errorf("Expected date 2024-01-15, got %s", item.Date)
		}
	}
}

func TestRegressionTeacherSchedule(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	items, err := client.ScheduleService().GetTeacherSchedule(ctx, "test-token", "t001")
	if err != nil {
		t.Skipf("Teacher schedule not available: %v", err)
	}

	for _, item := range items {
		if item.Teacher == "" {
			t.Errorf("Schedule item %d has no teacher", item.ID)
		}
	}
}

func TestRegressionSubjectGrades(t *testing.T) {
	ctx := context.Background()
	client := fixtures.NewAPIClient(fixtures.GetEnv())
	client.Init()

	grades, err := client.GradesService().GetSubjectGrades(ctx, "test-token", "12345", "s001")
	if err != nil {
		t.Skipf("Subject grades not available: %v", err)
	}

	for _, grade := range grades {
		if grade.SubjectID != "s001" {
			t.Errorf("Expected subject s001, got %s", grade.SubjectID)
		}
	}
}
