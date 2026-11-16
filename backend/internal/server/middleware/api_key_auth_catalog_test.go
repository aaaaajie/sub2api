//go:build unit

package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestModelCatalogBalanceExemption(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, balance := range []float64{0, -1, 1} {
		for _, tc := range []struct {
			method, path string
			catalog      bool
		}{
			{"GET", "/v1/models", true}, {"GET", "/models?client_version=0.158.0", true},
			{"GET", "/backend-api/codex/models?client_version=0.158.0", true},
			{"POST", "/v1/models", false}, {"GET", "/v1/models/other", false},
			{"POST", "/v1/responses", false}, {"GET", "/v1/responses", false},
			{"POST", "/v1/chat/completions", false}, {"POST", "/v1/messages/count_tokens", false},
		} {
			t.Run(tc.method+tc.path, func(t *testing.T) {
				key := &service.APIKey{ID: 1, Key: "catalog-test", Status: service.StatusActive, UserID: 1,
					User: &service.User{ID: 1, Status: service.StatusActive, Balance: balance}}
				cfg := &config.Config{RunMode: config.RunModeStandard}
				repo := &stubApiKeyRepo{getByKey: func(context.Context, string) (*service.APIKey, error) { return key, nil }}
				svc := service.NewAPIKeyService(repo, nil, nil, nil, nil, nil, cfg)
				router := gin.New()
				router.Use(gin.HandlerFunc(NewAPIKeyAuthMiddleware(svc, nil, cfg)))
				reached := false
				router.NoRoute(func(c *gin.Context) { reached = true; c.Status(200) })
				req := httptest.NewRequest(tc.method, tc.path, nil)
				req.Header.Set("Authorization", "Bearer catalog-test")
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)
				if tc.catalog || balance > 0 {
					require.Equal(t, 200, w.Code)
					require.True(t, reached)
				} else {
					require.Equal(t, 403, w.Code)
					require.False(t, reached)
					require.Contains(t, w.Body.String(), "INSUFFICIENT_BALANCE")
				}
			})
		}
	}
}

func TestModelCatalogRetainsAccessChecks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name   string
		change func(*service.APIKey)
		code   string
	}{
		{"disabled key", func(k *service.APIKey) { k.Status = "disabled" }, "API_KEY_DISABLED"},
		{"expired status", func(k *service.APIKey) { k.Status = service.StatusAPIKeyExpired }, "API_KEY_EXPIRED"},
		{"expired time", func(k *service.APIKey) { v := time.Now().Add(-time.Hour); k.ExpiresAt = &v }, "API_KEY_EXPIRED"},
		{"quota status", func(k *service.APIKey) { k.Status = service.StatusAPIKeyQuotaExhausted }, "QUOTA_EXHAUSTED"},
		{"quota runtime", func(k *service.APIKey) { k.Quota = 1; k.QuotaUsed = 1 }, "QUOTA_EXHAUSTED"},
		{"inactive user", func(k *service.APIKey) { k.User.Status = "disabled" }, "USER_INACTIVE"},
		{"ip restriction", func(k *service.APIKey) { k.IPWhitelist = []string{"203.0.113.10"} }, "ACCESS_DENIED"},
		{"disabled group", func(k *service.APIKey) { k.Group.Status = "disabled" }, "GROUP_DISABLED"},
		{"exclusive group", func(k *service.APIKey) { k.Group.IsExclusive = true }, "GROUP_NOT_ALLOWED"},
		{"missing subscription", func(k *service.APIKey) { k.Group.SubscriptionType = service.SubscriptionTypeSubscription }, "SUBSCRIPTION_NOT_FOUND"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			group := &service.Group{ID: 3, Status: service.StatusActive, Hydrated: true}
			key := &service.APIKey{ID: 1, Key: "catalog-test", Status: service.StatusActive, UserID: 1, Group: group, GroupID: &group.ID, User: &service.User{ID: 1, Status: service.StatusActive, Balance: 0}}
			tc.change(key)
			cfg := &config.Config{RunMode: config.RunModeStandard}
			repo := &stubApiKeyRepo{getByKey: func(context.Context, string) (*service.APIKey, error) { return key, nil }}
			svc := service.NewAPIKeyService(repo, nil, nil, nil, nil, nil, cfg)
			subRepo := &stubUserSubscriptionRepo{getActive: func(context.Context, int64, int64) (*service.UserSubscription, error) {
				return nil, service.ErrSubscriptionNotFound
			}}
			subs := service.NewSubscriptionService(nil, subRepo, nil, nil, cfg)
			t.Cleanup(subs.Stop)
			router := gin.New()
			router.Use(gin.HandlerFunc(NewAPIKeyAuthMiddleware(svc, subs, cfg)))
			router.GET("/v1/models", func(c *gin.Context) { t.Error("invalid key reached catalog"); c.Status(200) })
			req := httptest.NewRequest(http.MethodGet, "/v1/models", nil)
			req.Header.Set("Authorization", "Bearer catalog-test")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			require.Contains(t, []int{401, 403, 429}, w.Code)
			require.Contains(t, w.Body.String(), tc.code)
		})
	}
}
