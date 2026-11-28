package service

import (
	"bytes"
	"context"
	_ "embed"
	"html/template"
	"net/url"
	"strings"
	"unicode/utf8"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

type AdminUserMailDraft struct {
	Template    string `json:"template"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	ActionLabel string `json:"action_label"`
	ActionURL   string `json:"action_url"`
	Locale      string `json:"locale"`
}

var errInvalidUserMailDraft = infraerrors.BadRequest("INVALID_USER_MAIL_DRAFT", "Enter a subject and message, and provide a label and HTTPS link if adding a button.")

//go:embed admin_user_mail.html
var adminUserMailHTML string

var adminUserMailTemplate = template.Must(template.New("admin_user_mail").Parse(adminUserMailHTML))

type adminUserMailTemplateData struct {
	SiteName, Subject, Label, Footer, Language, ActionLabel, ActionURL string
	Paragraphs                                                         [][]string
}

func (s *AdminUserMailService) render(ctx context.Context, draft AdminUserMailDraft) (string, string, error) {
	subject := strings.TrimSpace(draft.Subject)
	body := strings.TrimSpace(strings.ReplaceAll(draft.Body, "\r\n", "\n"))
	label, actionURL := strings.TrimSpace(draft.ActionLabel), strings.TrimSpace(draft.ActionURL)
	if subject == "" || body == "" || strings.ContainsAny(subject, "\r\n") ||
		utf8.RuneCountInString(subject) > 160 || utf8.RuneCountInString(body) > 12000 ||
		utf8.RuneCountInString(label) > 40 || len(actionURL) > 2048 || (label == "") != (actionURL == "") {
		return "", "", errInvalidUserMailDraft
	}
	if actionURL != "" {
		link, err := url.Parse(actionURL)
		if err != nil || link.Scheme != "https" || link.Hostname() == "" || link.User != nil {
			return "", "", errInvalidUserMailDraft
		}
	}

	labels := map[string][2]string{
		"maintenance": {"维护通知", "Maintenance notice"},
		"service":     {"服务变更", "Service update"},
		"account":     {"账户消息", "Account message"},
	}
	names, ok := labels[draft.Template]
	if !ok {
		return "", "", errInvalidUserMailDraft
	}
	siteName := s.settings.GetSiteName(ctx)
	data := adminUserMailTemplateData{
		SiteName: siteName, Subject: subject, ActionLabel: label, ActionURL: actionURL,
		Label: names[0], Language: "zh-CN", Footer: "此邮件由 " + siteName + " 发送给你的平台账户。",
	}
	if normalizeNotificationLocale(draft.Locale) == notificationEmailDefaultLocale {
		data.Label, data.Language = names[1], "en"
		data.Footer = "This message was sent to your " + siteName + " account."
	}
	var paragraph []string
	for _, line := range strings.Split(body, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			paragraph = append(paragraph, line)
		} else if len(paragraph) > 0 {
			data.Paragraphs = append(data.Paragraphs, paragraph)
			paragraph = nil
		}
	}
	if len(paragraph) > 0 {
		data.Paragraphs = append(data.Paragraphs, paragraph)
	}
	var rendered bytes.Buffer
	if err := adminUserMailTemplate.Execute(&rendered, data); err != nil {
		return "", "", err
	}
	return subject, rendered.String(), nil
}
