package webserver

import (
	"strings"
	"testing"
)

func TestGenerateReverseProxyTemplate(t *testing.T) {
	domain := "example.com"
	rootDir := "/var/www/example.com"
	upstream := "127.0.0.1:8080"
	sslCert := "/etc/ssl/example.com.crt"
	sslKey := "/etc/ssl/example.com.key"

	tmpl := GenerateReverseProxyTemplate(domain, rootDir, upstream, sslCert, sslKey)

	if !strings.Contains(tmpl, "server_name example.com;") {
		t.Errorf("template missing server_name")
	}
	if !strings.Contains(tmpl, "proxy_pass http://127.0.0.1:8080;") {
		t.Errorf("template missing proxy_pass upstream")
	}
	if !strings.Contains(tmpl, "ssl_certificate /etc/ssl/example.com.crt;") {
		t.Errorf("template missing ssl_certificate directive")
	}
}
