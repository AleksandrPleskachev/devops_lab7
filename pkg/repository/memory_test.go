package repository

import (
	"context"
	"testing"

	"DevOps/pkg/course"
)

func TestCoursePersistence(t *testing.T) {
	ctx := context.Background()
	svc := course.NewCourseService(NewMemoryRepository())
	if _, err := svc.Create(ctx, course.CreateRequest{Title: "Go", Instructor: "Teacher", Duration: 0}); err == nil {
		t.Fatal("accepted non-positive duration")
	}
	id, err := svc.Create(ctx, course.CreateRequest{Title: "Go", Instructor: "Teacher", Duration: 24})
	if err != nil {
		t.Fatal(err)
	}
	before, err := svc.GetByID(ctx, id)
	if err != nil || before.Published {
		t.Fatalf("initial course: %+v, %v", before, err)
	}
	if err := svc.Publish(ctx, id); err != nil {
		t.Fatal(err)
	}
	after, err := svc.GetByID(ctx, id)
	if err != nil || !after.Published || before.CreatedAt != after.CreatedAt {
		t.Fatalf("published course: %+v, %v", after, err)
	}
	if _, err := svc.GetByID(ctx, "missing"); err == nil {
		t.Fatal("found an unknown course")
	}
}
