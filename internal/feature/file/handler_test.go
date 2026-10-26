package file

import (
	"net/http"
	"testing"
)

func TestIsUnsafeUploadType(t *testing.T) {
	unsafe := []string{
		"text/html; charset=utf-8",
		"text/xml; charset=utf-8",
		"application/xhtml+xml",
		"image/svg+xml",
	}
	for _, mimeType := range unsafe {
		if !isUnsafeUploadType(mimeType) {
			t.Errorf("%q should be rejected", mimeType)
		}
	}

	safe := []string{
		"application/pdf",
		"image/png",
		"application/zip",
		"application/octet-stream",
		"text/plain; charset=utf-8",
	}
	for _, mimeType := range safe {
		if isUnsafeUploadType(mimeType) {
			t.Errorf("%q should be allowed", mimeType)
		}
	}
}

func TestSniffedHTMLIsUnsafe(t *testing.T) {
	payload := []byte("<!DOCTYPE html><html><body><script>alert(1)</script></body></html>")

	if !isUnsafeUploadType(http.DetectContentType(payload)) {
		t.Fatal("an uploaded HTML page was not rejected")
	}
}
