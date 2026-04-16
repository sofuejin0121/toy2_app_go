package mailer

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/smtp"
	"net/url"
	texttemplate "text/template"

	"github.com/sofuejin0121/toy_app_go/internal/mailer/components"
	"github.com/sofuejin0121/toy_app_go/internal/model"
)

//go:embed templates/*.txt
var templateFS embed.FS

func buildAccountActivation(host, from string, user *model.User) (subject, to, fromAddr, textBody, htmlBody string, err error) {
	activationURL := fmt.Sprintf("http://%s/account_activations/%s/edit?email=%s",
		host, user.ActivationToken, url.QueryEscape(user.Email))

	data := EmailData{
		User:          user,
		ActivationURL: activationURL,
	}

	txtTmpl, err := texttemplate.ParseFS(templateFS, "templates/account_activation.txt")
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("parse text template: %w", err)
	}
	var txtBuf bytes.Buffer
	if err := txtTmpl.Execute(&txtBuf, data); err != nil {
		return "", "", "", "", "", fmt.Errorf("execute text template: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err := components.AccountActivationHTML(user.Name, activationURL).Render(context.Background(), &htmlBuf); err != nil {
		return "", "", "", "", "", fmt.Errorf("render html component: %w", err)
	}

	return "Account activation", user.Email, from, txtBuf.String(), htmlBuf.String(), nil
}

func (m *LogMailer) BuildAccountActivation(user *model.User) (subject, to, from, textBody, htmlBody string, err error) {
	return buildAccountActivation(m.Host, m.From, user)
}

func (m *LogMailer) SendAccountActivation(user *model.User) error {
	subject, to, from, textBody, htmlBody, err := m.BuildAccountActivation(user)
	if err != nil {
		return fmt.Errorf("build account activation email: %w", err)
	}

	log.Printf("=== Account Activation Email ===\n"+
		"To: %s\nFrom: %s\nSubject: %s\n\n"+
		"テキスト版:\n%s\n\nHTML版:\n%s\n"+
		"================================\n",
		to, from, subject, textBody, htmlBody)
	return nil
}

func (m *SMTPMailer) BuildAccountActivation(user *model.User) (subject, to, from, textBody, htmlBody string, err error) {
	return buildAccountActivation(m.AppHost, m.From, user)
}

func (m *SMTPMailer) SendAccountActivation(user *model.User) error {
	subject, to, from, textBody, htmlBody, err := m.BuildAccountActivation(user)
	if err != nil {
		return fmt.Errorf("build account activation email: %w", err)
	}

	boundary := "sample-app-boundary"
	var msg bytes.Buffer
	fmt.Fprintf(&msg, "From: %s\r\n", from)
	fmt.Fprintf(&msg, "To: %s\r\n", to)
	fmt.Fprintf(&msg, "Subject: %s\r\n", subject)
	fmt.Fprint(&msg, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&msg, "Content-Type: multipart/alternative; boundary=%q\r\n\r\n", boundary)

	fmt.Fprintf(&msg, "--%s\r\n", boundary)
	fmt.Fprint(&msg, "Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	fmt.Fprint(&msg, textBody)
	fmt.Fprint(&msg, "\r\n")

	fmt.Fprintf(&msg, "--%s\r\n", boundary)
	fmt.Fprint(&msg, "Content-Type: text/html; charset=UTF-8\r\n\r\n")
	fmt.Fprint(&msg, htmlBody)
	fmt.Fprint(&msg, "\r\n")

	fmt.Fprintf(&msg, "--%s--\r\n", boundary)

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	if err := smtp.SendMail(addr, auth, from, []string{to}, msg.Bytes()); err != nil {
		return fmt.Errorf("send account activation email: %w", err)
	}
	return nil
}

func buildPasswordReset(host, from string, user *model.User) (subject, to, fromAddr, textBody, htmlBody string, err error) {
	resetURL := fmt.Sprintf("http://%s/password_resets/%s/edit?email=%s",
		host, user.ResetToken, url.QueryEscape(user.Email))

	data := EmailData{
		User:     user,
		ResetURL: resetURL,
	}

	txtTmpl, err := texttemplate.ParseFS(templateFS, "templates/password_reset.txt")
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("parse text template: %w", err)
	}
	var txtBuf bytes.Buffer
	if err := txtTmpl.Execute(&txtBuf, data); err != nil {
		return "", "", "", "", "", fmt.Errorf("execute text template: %w", err)
	}

	var htmlBuf bytes.Buffer
	if err := components.PasswordResetHTML(resetURL).Render(context.Background(), &htmlBuf); err != nil {
		return "", "", "", "", "", fmt.Errorf("render html component: %w", err)
	}

	return "Password reset", user.Email, from, txtBuf.String(), htmlBuf.String(), nil
}

func (m *LogMailer) BuildPasswordReset(user *model.User) (subject, to, from, textBody, htmlBody string, err error) {
	return buildPasswordReset(m.Host, m.From, user)
}

func (m *LogMailer) SendPasswordReset(user *model.User) error {
	subject, to, from, textBody, htmlBody, err := m.BuildPasswordReset(user)
	if err != nil {
		return fmt.Errorf("build password reset email: %w", err)
	}

	log.Printf("=== Password Reset Email ===\n"+
		"To: %s\nFrom: %s\nSubject: %s\n\n"+
		"テキスト版:\n%s\n\nHTML版:\n%s\n"+
		"============================\n",
		to, from, subject, textBody, htmlBody)
	return nil
}

func (m *SMTPMailer) BuildPasswordReset(user *model.User) (subject, to, from, textBody, htmlBody string, err error) {
	return buildPasswordReset(m.AppHost, m.From, user)
}

// SendFollowNotification はフォロー通知メールをログに出力する（開発用）
func (m *LogMailer) SendFollowNotification(to *model.User, follower *model.User) error {
	log.Printf("=== Follow Notification Email ===\n"+
		"To: %s\nFrom: %s\nSubject: %s がフォローしました\n\n"+
		"%s さんから新しいフォローがありました。\n"+
		"================================\n",
		to.Email, m.From, follower.Name, follower.Name)
	return nil
}

// SendLikeNotification はいいね通知メールをログに出力する（開発用）
func (m *LogMailer) SendLikeNotification(to *model.User, liker *model.User, content string) error {
	log.Printf("=== Like Notification Email ===\n"+
		"To: %s\nFrom: %s\nSubject: %s があなたの投稿をいいねしました\n\n"+
		"%s さんが「%s」をいいねしました。\n"+
		"==============================\n",
		to.Email, m.From, liker.Name, liker.Name, content)
	return nil
}

// SendFollowNotification はフォロー通知メールをSMTP送信する
func (m *SMTPMailer) SendFollowNotification(to *model.User, follower *model.User) error {
	subject := follower.Name + " があなたをフォローしました"
	body := follower.Name + " さんから新しいフォローがありました。\r\n"

	var msg bytes.Buffer
	fmt.Fprintf(&msg, "From: %s\r\n", m.From)
	fmt.Fprintf(&msg, "To: %s\r\n", to.Email)
	fmt.Fprintf(&msg, "Subject: %s\r\n", subject)
	fmt.Fprint(&msg, "Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	fmt.Fprint(&msg, body)

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	return smtp.SendMail(addr, auth, m.From, []string{to.Email}, msg.Bytes())
}

// SendLikeNotification はいいね通知メールをSMTP送信する
func (m *SMTPMailer) SendLikeNotification(to *model.User, liker *model.User, content string) error {
	subject := liker.Name + " があなたの投稿をいいねしました"
	body := liker.Name + " さんが「" + content + "」をいいねしました。\r\n"

	var msg bytes.Buffer
	fmt.Fprintf(&msg, "From: %s\r\n", m.From)
	fmt.Fprintf(&msg, "To: %s\r\n", to.Email)
	fmt.Fprintf(&msg, "Subject: %s\r\n", subject)
	fmt.Fprint(&msg, "Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	fmt.Fprint(&msg, body)

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	return smtp.SendMail(addr, auth, m.From, []string{to.Email}, msg.Bytes())
}

func (m *SMTPMailer) SendPasswordReset(user *model.User) error {
	subject, to, from, textBody, htmlBody, err := m.BuildPasswordReset(user)
	if err != nil {
		return fmt.Errorf("build password reset email: %w", err)
	}

	boundary := "sample-app-boundary"
	var msg bytes.Buffer
	fmt.Fprintf(&msg, "From: %s\r\n", from)
	fmt.Fprintf(&msg, "To: %s\r\n", to)
	fmt.Fprintf(&msg, "Subject: %s\r\n", subject)
	fmt.Fprint(&msg, "MIME-Version: 1.0\r\n")
	fmt.Fprintf(&msg, "Content-Type: multipart/alternative; boundary=%q\r\n\r\n", boundary)

	fmt.Fprintf(&msg, "--%s\r\n", boundary)
	fmt.Fprint(&msg, "Content-Type: text/plain; charset=UTF-8\r\n\r\n")
	fmt.Fprint(&msg, textBody)
	fmt.Fprint(&msg, "\r\n")

	fmt.Fprintf(&msg, "--%s\r\n", boundary)
	fmt.Fprint(&msg, "Content-Type: text/html; charset=UTF-8\r\n\r\n")
	fmt.Fprint(&msg, htmlBody)
	fmt.Fprint(&msg, "\r\n")

	fmt.Fprintf(&msg, "--%s--\r\n", boundary)

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)
	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)
	if err := smtp.SendMail(addr, auth, from, []string{to}, msg.Bytes()); err != nil {
		return fmt.Errorf("send password reset email: %w", err)
	}
	return nil
}

// brevoSend は Brevo HTTP API でメールを送信する共通ヘルパー
func brevoSend(apiKey, fromAddr, toAddr, toName, subject, htmlBody, textBody string) error {
	payload := map[string]any{
		"sender":      map[string]string{"email": fromAddr},
		"to":          []map[string]string{{"email": toAddr, "name": toName}},
		"subject":     subject,
		"htmlContent": htmlBody,
		"textContent": textBody,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("brevo marshal: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.brevo.com/v3/smtp/email", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("brevo new request: %w", err)
	}
	req.Header.Set("api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("brevo send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		log.Printf("brevo API error: status %d body: %s", resp.StatusCode, string(respBody))
		return fmt.Errorf("brevo API error: status %d body: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (m *BrevoMailer) SendAccountActivation(user *model.User) error {
	activationURL := fmt.Sprintf("https://%s/account_activations/%s/edit?email=%s",
		m.AppHost, user.ActivationToken, url.QueryEscape(user.Email))

	var htmlBuf bytes.Buffer
	if err := components.AccountActivationHTML(user.Name, activationURL).Render(context.Background(), &htmlBuf); err != nil {
		return fmt.Errorf("render html: %w", err)
	}

	return brevoSend(m.APIKey, m.From, user.Email, user.Name,
		"Account activation", htmlBuf.String(),
		fmt.Sprintf("Hi %s,\nActivate your account: %s", user.Name, activationURL))
}

func (m *BrevoMailer) SendPasswordReset(user *model.User) error {
	resetURL := fmt.Sprintf("https://%s/password_resets/%s/edit?email=%s",
		m.AppHost, user.ResetToken, url.QueryEscape(user.Email))

	var htmlBuf bytes.Buffer
	if err := components.PasswordResetHTML(resetURL).Render(context.Background(), &htmlBuf); err != nil {
		return fmt.Errorf("render html: %w", err)
	}

	return brevoSend(m.APIKey, m.From, user.Email, user.Name,
		"Password reset", htmlBuf.String(),
		fmt.Sprintf("Hi %s,\nReset your password: %s", user.Name, resetURL))
}

func (m *BrevoMailer) SendFollowNotification(to *model.User, follower *model.User) error {
	subject := follower.Name + " があなたをフォローしました"
	body := follower.Name + " さんから新しいフォローがありました。"
	return brevoSend(m.APIKey, m.From, to.Email, to.Name, subject, body, body)
}

func (m *BrevoMailer) SendLikeNotification(to *model.User, liker *model.User, content string) error {
	subject := liker.Name + " があなたの投稿をいいねしました"
	body := liker.Name + " さんが「" + content + "」をいいねしました。"
	return brevoSend(m.APIKey, m.From, to.Email, to.Name, subject, body, body)
}

// resendSend は Resend HTTP API でメールを送信する共通ヘルパー
func resendSend(apiKey, fromAddr, toAddr, subject, htmlBody, textBody string) error {
	payload := map[string]any{
		"from":    fromAddr,
		"to":      []string{toAddr},
		"subject": subject,
		"html":    htmlBody,
		"text":    textBody,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("resend marshal: %w", err)
	}

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("resend new request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("resend send: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		log.Printf("resend API error: status %d body: %s", resp.StatusCode, string(respBody))
		return fmt.Errorf("resend API error: status %d body: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (m *ResendMailer) SendAccountActivation(user *model.User) error {
	activationURL := fmt.Sprintf("https://%s/account_activations/%s/edit?email=%s",
		m.AppHost, user.ActivationToken, url.QueryEscape(user.Email))

	var htmlBuf bytes.Buffer
	if err := components.AccountActivationHTML(user.Name, activationURL).Render(context.Background(), &htmlBuf); err != nil {
		return fmt.Errorf("render html: %w", err)
	}

	return resendSend(m.APIKey, m.From, user.Email,
		"Account activation", htmlBuf.String(),
		fmt.Sprintf("Hi %s,\nActivate your account: %s", user.Name, activationURL))
}

func (m *ResendMailer) SendPasswordReset(user *model.User) error {
	resetURL := fmt.Sprintf("https://%s/password_resets/%s/edit?email=%s",
		m.AppHost, user.ResetToken, url.QueryEscape(user.Email))

	var htmlBuf bytes.Buffer
	if err := components.PasswordResetHTML(resetURL).Render(context.Background(), &htmlBuf); err != nil {
		return fmt.Errorf("render html: %w", err)
	}

	return resendSend(m.APIKey, m.From, user.Email,
		"Password reset", htmlBuf.String(),
		fmt.Sprintf("Hi %s,\nReset your password: %s", user.Name, resetURL))
}

func (m *ResendMailer) SendFollowNotification(to *model.User, follower *model.User) error {
	subject := follower.Name + " があなたをフォローしました"
	body := follower.Name + " さんから新しいフォローがありました。"
	return resendSend(m.APIKey, m.From, to.Email, subject, body, body)
}

func (m *ResendMailer) SendLikeNotification(to *model.User, liker *model.User, content string) error {
	subject := liker.Name + " があなたの投稿をいいねしました"
	body := liker.Name + " さんが「" + content + "」をいいねしました。"
	return resendSend(m.APIKey, m.From, to.Email, subject, body, body)
}
