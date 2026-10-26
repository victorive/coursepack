package middleware

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

type fakeValidator struct {
	err error
}

func (f fakeValidator) ValidateToken(ctx context.Context, token string) (string, error) {
	return "", f.err
}

func TestAuthDoesNotLeakTokenError(t *testing.T) {
	app := fiber.New()
	app.Use(Auth(fakeValidator{err: errors.New("signature is invalid: key 3 not found")}))
	app.Get("/", func(c *fiber.Ctx) error { return c.SendStatus(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer not-a-real-token")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("got status %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	if strings.Contains(string(body), "key 3 not found") {
		t.Fatalf("response leaked the validator error: %s", body)
	}
}
