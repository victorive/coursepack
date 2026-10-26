package file

import (
	"encoding/json"
	"testing"

	"github.com/victorive/coursepack/internal/domain"
)

// The typed category and status must keep the JSON contract unchanged.
func TestFileResponseJSONIsUnchanged(t *testing.T) {
	resp := toFileResponse(&domain.File{
		UUID:          "file-1",
		UploadedBy:    1,
		ClassCourseID: 2,
		Category:      string(domain.FileCategoryLectureNote),
		Title:         "Week 1",
		SizeBytes:     3,
		MimeType:      "application/pdf",
		Status:        string(domain.FileStatusApproved),
	})

	got, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	want := `{"uuid":"file-1","uploaded_by":1,"class_course_id":2,"category":"lecture_note","title":"Week 1","size_bytes":3,"mime_type":"application/pdf","status":"approved","created_at":"0001-01-01T00:00:00Z"}`
	if string(got) != want {
		t.Fatalf("json changed:\n got: %s\nwant: %s", got, want)
	}
}
