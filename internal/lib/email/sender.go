package email

import (
	"fmt"
	"net/smtp"
	"strings"
)

const (
	//todo update template and url
	uiUrl                 = "https://example.com"
	resetPasswordTemplate = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<style>
			body {
				font-family: Arial, sans-serif;
				background: #f4f6f9;
				padding: 20px;
			}
	
			.container {
				background-color: #ffffff;
				max-width: 600px;
				margin: auto;
				padding: 30px;
				border-radius: 8px;
				box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
			}
			h2 {
				color: #333333;
			}
			p {
				color: #555555;
				font-size: 16px;
				line-height: 1.5;
			}
			.button {
				display: inline-block;
				margin-top: 20px;
				padding: 12px 24px;
				background-color: #28a745;
				color: white;
				text-decoration: none;
				border-radius: 5px;
				font-weight: bold;
				box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
				transition: background-color 0.3s ease, box-shadow 0.3s ease;
			}
			.button:hover {
				background-color: #218838;
				box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
			}
			.footer {
				font-size: 12px;
				color: #999999;
				margin-top: 30px;
				text-align: center;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h2>Password Reset Request</h2>
			<p>Hello,</p>
			<p>We received a request to reset your password. Click the button below to choose a new one:</p>
			<a href="%s" class="button">Reset Password</a>
			<p>If you didn't request a password reset, you can safely ignore this email.</p>
			<div class="footer">
				&copy; 2025 Your Company. All rights reserved.
			</div>
		</div>
	</body>
	</html>
	`
)

type Sender interface {
	SendResetPasswordEmail(toEmail, resetToken string) error
}

type SmtpSender struct {
	auth smtp.Auth

	username string
	from     string
	smtpAddr string
}

func NewSmtpSender(host, port, username, password, from string) *SmtpSender {
	auth := smtp.PlainAuth("", username, password, host)

	return &SmtpSender{
		username: username,
		from:     from,
		auth:     auth,
		smtpAddr: fmt.Sprintf("%s:%s", host, port),
	}
}

func (s *SmtpSender) SendResetPasswordEmail(toEmail, resetToken string) error {
	link := fmt.Sprintf("%s/reset-password?token=%s", uiUrl, resetToken)
	subject := "Password Reset Instructions"

	body := fmt.Sprintf(resetPasswordTemplate, link)

	msg := s.buildMessage(toEmail, subject, body)
	return smtp.SendMail(s.smtpAddr, s.auth, s.username, []string{toEmail}, msg)
}

func (s *SmtpSender) buildMessage(to, subject, htmlBody string) []byte {
	headers := make(map[string]string)
	headers["From"] = s.username
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n" + htmlBody)
	return []byte(msg.String())
}
