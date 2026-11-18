package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

type publicActivityTrendResponse struct {
	Code int `json:"code"`
	Data struct {
		Trend []publicActivityTrendPoint `json:"trend"`
	} `json:"data"`
}

func TestPublicActivityTrendReturnsOnlyRecentHourlyRelativeActivityValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := timezone.Now()
	usageRepo := &dailyUsageRepoStub{recentActivity: []service.RecentUsageActivityHour{
		{BucketStart: now.Truncate(time.Hour), TotalTokens: 12_000},
		{BucketStart: now.Add(-2 * time.Hour).Truncate(time.Hour), TotalTokens: 6_000},
		{BucketStart: now.Add(-4 * time.Hour).Truncate(time.Hour), TotalTokens: 2_000},
	}}
	usageService := service.NewUsageService(usageRepo, nil, nil, nil)
	usageHandler := NewUsageHandler(usageService, nil, nil, nil)
	router := gin.New()
	router.GET("/api/v1/public/recent-usage-activity", usageHandler.PublicActivityTrend)

	request := httptest.NewRequest(http.MethodGet, "/api/v1/public/recent-usage-activity", nil)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	require.Equal(t, http.StatusOK, responseRecorder.Code)
	require.Equal(t, "public, max-age=60, stale-while-revalidate=240", responseRecorder.Header().Get("Cache-Control"))
	var payload publicActivityTrendResponse
	require.NoError(t, json.Unmarshal(responseRecorder.Body.Bytes(), &payload))
	require.Equal(t, 0, payload.Code)
	require.Len(t, payload.Data.Trend, 3)
	require.InDelta(t, 16.67, payload.Data.Trend[0].ActivityPercent, 0.01)
	require.InDelta(t, 50, payload.Data.Trend[1].ActivityPercent, 0.01)
	require.InDelta(t, 100, payload.Data.Trend[2].ActivityPercent, 0.01)
	require.Contains(t, payload.Data.Trend[0].Date, "T")
	require.NotContains(t, responseRecorder.Body.String(), "total_tokens")
	require.NotContains(t, responseRecorder.Body.String(), "requests")
	require.Equal(t, 1, usageRepo.recentActivityCalls)
	require.Equal(t, 10, usageRepo.recentActivityLimit)

	secondResponse := httptest.NewRecorder()
	router.ServeHTTP(secondResponse, httptest.NewRequest(http.MethodGet, "/api/v1/public/recent-usage-activity", nil))
	require.Equal(t, http.StatusOK, secondResponse.Code)
	require.Equal(t, 1, usageRepo.recentActivityCalls)
	require.Equal(t, "public, max-age=60, stale-while-revalidate=240", secondResponse.Header().Get("Cache-Control"))
}

func TestPublicActivityTrendReturnsNoPointsWhenThereIsNoUsage(t *testing.T) {
	gin.SetMode(gin.TestMode)
	usageRepo := &dailyUsageRepoStub{}
	usageHandler := NewUsageHandler(service.NewUsageService(usageRepo, nil, nil, nil), nil, nil, nil)
	router := gin.New()
	router.GET("/recent-usage-activity", usageHandler.PublicActivityTrend)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, httptest.NewRequest(http.MethodGet, "/recent-usage-activity", nil))

	var payload publicActivityTrendResponse
	require.NoError(t, json.Unmarshal(responseRecorder.Body.Bytes(), &payload))
	require.Empty(t, payload.Data.Trend)
}
