package email

import "context"

type Mailer interface {
	Send(ctx context.Context, to, subject, tmpl string, data map[string]any) error
}
