package handler

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/pkg/usagestats"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type recentHoursRepo struct {
	service.UsageLogRepository
	userID   int64
	timezone string
	err      error
}

func (r *recentHoursRepo) GetRecentUserUsageHours(_ context.Context, id int64, tz string) ([]usagestats.TrendDataPoint, error) {
	r.userID, r.timezone = id, tz
	return []usagestats.TrendDataPoint{{Date: "2000-01-01 10:00", Requests: 1}}, r.err
}
func TestDashboardRecentHours(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		name, query   string
		authenticated bool
		fail          bool
		status        int
	}{
		{"own history", "?timezone=Asia%2FShanghai&user_id=999", true, false, 200},
		{"default UTC", "", true, false, 200},
		{"invalid timezone", "?timezone=wrong", true, false, 400},
		{"local timezone", "?timezone=Local", true, false, 400},
		{"unauthenticated", "", false, false, 401},
		{"repository failure", "", true, true, 500},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &recentHoursRepo{}
			if tc.fail {
				repo.err = errors.New("query failed")
			}
			h := NewUsageHandler(service.NewUsageService(repo, nil, nil, nil), nil, nil, nil)
			r := gin.New()
			r.GET("/recent", func(c *gin.Context) {
				if tc.authenticated {
					c.Set(string(middleware2.ContextKeyUser), middleware2.AuthSubject{UserID: 42})
				}
				h.DashboardRecentHours(c)
			})
			w := httptest.NewRecorder()
			r.ServeHTTP(w, httptest.NewRequest("GET", "/recent"+tc.query, nil))
			require.Equal(t, tc.status, w.Code)
			if tc.status == 200 {
				require.Equal(t, int64(42), repo.userID)
				require.Contains(t, w.Body.String(), `"requests":1`)
				if tc.query == "" {
					require.Equal(t, "UTC", repo.timezone)
				} else {
					require.Equal(t, "Asia/Shanghai", repo.timezone)
				}
			}
		})
	}
}
