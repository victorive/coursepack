package smtp

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"
	"time"

	gomail "gopkg.in/gomail.v2"
)

type Client struct {
	dialer    *gomail.Dialer
	fromAddr  string
	fromName  string
	appURL    string
	templates *template.Template
}

func New(host string, port int, username, password, fromAddr, fromName, appURL string, templateFS fs.FS) (*Client, error) {
	tmpl, err := template.ParseFS(templateFS, "*.html")
	if err != nil {
		return nil, fmt.Errorf("smtp: parse templates: %w", err)
	}

	return &Client{
		dialer:    gomail.NewDialer(host, port, username, password),
		fromAddr:  fromAddr,
		fromName:  fromName,
		appURL:    appURL,
		templates: tmpl,
	}, nil
}

func (c *Client) Send(_ context.Context, to, subject, tmpl string, data map[string]any) error {
	if data == nil {
		data = make(map[string]any)
	}
	data["AppURL"] = c.appURL
	data["Year"] = time.Now().Year()

	var buf bytes.Buffer
	if err := c.templates.ExecuteTemplate(&buf, tmpl, data); err != nil {
		return fmt.Errorf("smtp: render %s: %w", tmpl, err)
	}

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", c.fromAddr, c.fromName)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", buf.String())

	return c.dialer.DialAndSend(msg)
}
