package config

import "testing"

func TestValidateAllowsDefaultsOutsideProduction(t *testing.T) {
	cfg := &Config{
		App:  AppConfig{Env: "development"},
		SMTP: SMTPConfig{FromEmail: "noreply@example.com"},
		CORS: CORSConfig{AllowedOrigins: "*"},
	}

	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateRejectsPlaceholderSMTPInProduction(t *testing.T) {
	cfg := &Config{
		App:  AppConfig{Env: "production"},
		SMTP: SMTPConfig{FromEmail: "noreply@example.com"},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected an error for the example.com sender")
	}
}

func TestValidateRejectsWildcardCORSInProduction(t *testing.T) {
	cfg := &Config{
		App:  AppConfig{Env: "production"},
		SMTP: SMTPConfig{FromEmail: "noreply@coursepack.app"},
		CORS: CORSConfig{AllowedOrigins: "*"},
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected an error for wildcard CORS")
	}
}
