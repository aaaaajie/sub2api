//go:build unit

package service

import (
	"context"
	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestPlazaOfficialTimePricingMatchesBilling(t *testing.T) {
	pricing := newStubPricingServiceFromMap(map[string]*LiteLLMModelPricing{})
	svc := &ModelPlazaService{billingService: NewBillingService(&config.Config{}, pricing)}
	p := svc.lookupOfficialPricing(context.Background(), "DeepSeek-v4-flash", map[string]*PlazaOfficialPricing{})
	require.NotNil(t, p)
	require.NotNil(t, p.TimePricing)
	require.InDelta(t, 0.15e-6, *p.InputPrice, 1e-15)
	cfg := &ChannelTimePricing{Timezone: p.TimePricing.Timezone, WeekdaysOnly: p.TimePricing.WeekdaysOnly}
	for _, period := range p.TimePricing.Periods {
		cfg.Periods = append(cfg.Periods, ChannelTimePricingPeriod{StartTime: period.StartTime, EndTime: period.EndTime, Multiplier: period.Multiplier})
	}
	// Every minute of a week includes both boundaries, weekdays and weekends.
	start := time.Date(2026, 11, 16, 0, 0, 0, 0, time.UTC)
	for minute := 0; minute < 7*24*60; minute++ {
		at := start.Add(time.Duration(minute) * time.Minute)
		require.Equal(t, deepseekPeakMultiplierAt(at), cfg.MultiplierAt(at), at.String())
	}
}
