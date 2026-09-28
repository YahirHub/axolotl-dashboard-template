package security

import (
	"crypto/tls"
	"fmt"
	"html"
	"net"
	netmail "net/mail"
	"net/smtp"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/YahirHub/axolotl-dashboard-template/app/facades"
)

func PasswordRecoveryEnabled() bool {
	host := strings.TrimSpace(facades.Config().GetString("mail.host"))
	from := strings.TrimSpace(facades.Config().GetString("mail.from.address"))
	appURL := strings.TrimSpace(facades.Config().GetString("app.url"))

	if host == "" || from == "" || appURL == "" || facades.Config().GetInt("mail.port") <= 0 {
		return false
	}
	if _, err := netmail.ParseAddress(from); err != nil {
		return false
	}

	parsed, err := url.Parse(appURL)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return false
	}
	if strings.EqualFold(facades.Config().GetString("app.env"), "production") && parsed.Scheme != "https" {
		return false
	}
	return true
}

func PasswordResetLink(token string) (string, error) {
	if !PasswordRecoveryEnabled() || !validRawResetToken(token) {
		return "", fmt.Errorf("recuperación por correo no configurada")
	}

	base := strings.TrimRight(facades.Config().GetString("app.url"), "/")
	return base + "/reset-password#token=" + url.QueryEscape(token), nil
}

func SendPasswordResetEmail(to, name, link string) error {
	if !PasswordRecoveryEnabled() {
		return fmt.Errorf("recuperación por correo no configurada")
	}

	toAddress, err := netmail.ParseAddress(strings.TrimSpace(to))
	if err != nil {
		return err
	}
	fromAddress, err := netmail.ParseAddress(strings.TrimSpace(facades.Config().GetString("mail.from.address")))
	if err != nil {
		return err
	}

	fromName := cleanHeader(facades.Config().GetString("mail.from.name"))
	if fromName == "" {
		fromName = "Axolotl"
	}
	fromAddress.Name = fromName

	subject := "Recupera tu cuenta de Axolotl"
	safeName := html.EscapeString(strings.TrimSpace(name))
	safeLink := html.EscapeString(link)
	body := fmt.Sprintf(`<!doctype html>
<html>
<body style="font-family:Arial,sans-serif;color:#24151d">
<h2>Recuperación de cuenta</h2>
<p>Hola %s.</p>
<p>Se solicitó restablecer la contraseña de tu cuenta Axolotl.</p>
<p><a href="%s">Restablecer contraseña</a></p>
<p>El enlace expira pronto y solo puede usarse una vez. Si no solicitaste este cambio, ignora este correo.</p>
</body>
</html>`, safeName, safeLink)

	message := strings.Join([]string{
		"From: " + fromAddress.String(),
		"To: " + toAddress.String(),
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/html; charset=UTF-8",
		"Content-Transfer-Encoding: 8bit",
		"",
		body,
	}, string([]byte{13, 10}))

	return sendSMTP(fromAddress.Address, toAddress.Address, []byte(message))
}

func sendSMTP(from, to string, message []byte) error {
	host := strings.TrimSpace(facades.Config().GetString("mail.host"))
	port := facades.Config().GetInt("mail.port")
	username := facades.Config().GetString("mail.username")
	password := facades.Config().GetString("mail.password")
	encryption := strings.ToLower(strings.TrimSpace(facades.Config().GetString("mail.encryption", "auto")))
	if encryption == "" {
		encryption = "auto"
	}

	address := net.JoinHostPort(host, strconv.Itoa(port))
	dialer := &net.Dialer{Timeout: 10 * time.Second}

	var client *smtp.Client
	var err error
	tlsActive := false

	switch encryption {
	case "tls":
		conn, dialErr := tls.DialWithDialer(dialer, "tcp", address, &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: host,
		})
		if dialErr != nil {
			return dialErr
		}
		client, err = smtp.NewClient(conn, host)
		tlsActive = true
	case "starttls":
		client, err = dialPlainSMTP(dialer, address, host)
		if err == nil {
			err = startTLS(client, host, true)
			tlsActive = err == nil
		}
	case "none":
		if !isLocalSMTPHost(host) {
			return fmt.Errorf("SMTP sin cifrado solo está permitido en localhost")
		}
		client, err = dialPlainSMTP(dialer, address, host)
	default:
		if port == 465 {
			conn, dialErr := tls.DialWithDialer(dialer, "tcp", address, &tls.Config{
				MinVersion: tls.VersionTLS12,
				ServerName: host,
			})
			if dialErr != nil {
				return dialErr
			}
			client, err = smtp.NewClient(conn, host)
			tlsActive = true
		} else {
			client, err = dialPlainSMTP(dialer, address, host)
			if err == nil {
				if ok, _ := client.Extension("STARTTLS"); ok {
					err = startTLS(client, host, true)
					tlsActive = err == nil
				} else if !isLocalSMTPHost(host) {
					_ = client.Close()
					return fmt.Errorf("servidor SMTP remoto sin STARTTLS")
				}
			}
		}
	}

	if err != nil {
		if client != nil {
			_ = client.Close()
		}
		return err
	}
	defer client.Close()

	if username != "" {
		if !tlsActive && !isLocalSMTPHost(host) {
			return fmt.Errorf("credenciales SMTP requieren TLS")
		}
		if ok, _ := client.Extension("AUTH"); !ok {
			return fmt.Errorf("servidor SMTP no anuncia AUTH")
		}
		if err := client.Auth(smtp.PlainAuth("", username, password, host)); err != nil {
			return err
		}
	}

	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	writer, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := writer.Write(message); err != nil {
		_ = writer.Close()
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return client.Quit()
}

func dialPlainSMTP(dialer *net.Dialer, address, host string) (*smtp.Client, error) {
	conn, err := dialer.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	return client, nil
}

func startTLS(client *smtp.Client, host string, required bool) error {
	ok, _ := client.Extension("STARTTLS")
	if !ok {
		if required {
			return fmt.Errorf("servidor SMTP no soporta STARTTLS")
		}
		return nil
	}
	return client.StartTLS(&tls.Config{
		MinVersion: tls.VersionTLS12,
		ServerName: host,
	})
}

func isLocalSMTPHost(host string) bool {
	host = strings.Trim(strings.ToLower(host), "[] ")
	if host == "localhost" || host == "::1" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

func cleanHeader(value string) string {
	value = strings.ReplaceAll(value, "", "")
	value = strings.Map(func(r rune) rune {
		if r == 13 || r == 10 {
			return -1
		}
		return r
	}, value)
	return strings.TrimSpace(value)
}
