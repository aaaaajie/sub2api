package service

import (
	"errors"
	"io"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAdminUserMailTemplateEscapesContentAndChecksLinks(t *testing.T) {
	s := &AdminUserMailService{settings: &SettingService{
		settingRepo: &adminComplianceRepoStub{values: map[string]string{SettingKeySiteName: "SecAPI"}},
	}}
	draft := AdminUserMailDraft{
		Template: "maintenance", Subject: "<Account change>", Locale: "zh-CN",
		Body: "第一行\n第二行\n\n<script>alert(1)</script>", ActionLabel: "查看详情", ActionURL: "https://api.rzsec.cn/health",
	}
	subject, html, err := s.render(t.Context(), draft)
	require.NoError(t, err)
	require.Equal(t, draft.Subject, subject)
	require.Contains(t, html, "SecAPI")
	require.Contains(t, html, "维护通知")
	require.Contains(t, html, "&lt;Account change&gt;")
	require.Contains(t, html, "第一行<br>第二行")
	require.Contains(t, html, "&lt;script&gt;alert(1)&lt;/script&gt;")
	require.NotContains(t, html, "<script>")
	require.Contains(t, html, `href="https://api.rzsec.cn/health"`)
	draft.ActionURL = "javascript:alert(1)"
	_, _, err = s.render(t.Context(), draft)
	require.ErrorIs(t, err, errInvalidUserMailDraft)
}

func TestAdminUserMailSkipsDuplicateAndPlaceholderAddresses(t *testing.T) {
	got := eligibleUserMailRecipients([]AdminUserMailRecipient{
		{UserID: 1, Email: "alpha@example.test"},
		{UserID: 2, Email: " Alpha@example.test "},
		{UserID: 3, Email: "operator@linuxdo-connect.invalid"},
		{UserID: 4, Email: "Example <alice@example.test>"},
		{UserID: 5, Email: "beta@example.test"},
	})
	require.Equal(t, []AdminUserMailRecipient{
		{UserID: 1, Email: "alpha@example.test"},
		{UserID: 5, Email: "beta@example.test"},
	}, got)
}

func TestAdminUserMailDistinguishesRejectedAndUnconfirmedSMTP(t *testing.T) {
	status, reason := adminUserMailDeliveryResult(&textproto.Error{Code: 450, Msg: "rejected"})
	require.Equal(t, "failed", status)
	require.Equal(t, "smtp_rejected", reason)
	status, reason = adminUserMailDeliveryResult(errors.Join(errSMTPAcceptanceUnknown, io.EOF))
	require.Equal(t, "uncertain", status)
	require.Equal(t, "smtp_confirmation", reason)
}
