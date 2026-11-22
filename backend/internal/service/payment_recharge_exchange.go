package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const rechargeExchangeURL = "https://www.chinamoney.com.cn/r/cms/www/chinamoney/data/fx/ccpr.json"

type rechargeExchangeCache struct {
	mu         sync.Mutex
	multiplier float64
	published  time.Time
	nextCheck  time.Time
}

var rechargeExchangeClient = &http.Client{Timeout: 4 * time.Second}

// RechargeMultiplier returns the USD credit per unit of the payment currency.
func (s *PaymentConfigService) RechargeMultiplier(ctx context.Context, cfg *PaymentConfig, currency string) (float64, error) {
	switch currency {
	case "USD":
		return 1, nil
	case "CNY":
		if cfg != nil && !cfg.AutoRechargeExchange {
			multiplier := cfg.BalanceRechargeMultiplier
			if multiplier > 0 && multiplier < 1 {
				return multiplier, nil
			}
		} else if s != nil {
			if multiplier := s.rechargeExchange.current(ctx); multiplier > 0 {
				return multiplier, nil
			}
		}
	}
	return 0, infraerrors.ServiceUnavailable("RECHARGE_EXCHANGE_UNAVAILABLE", "Recharge exchange rate is unavailable for this currency")
}

func (cache *rechargeExchangeCache) current(ctx context.Context) float64 {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	now := time.Now()
	if !now.Before(cache.nextCheck) && ctx.Err() == nil {
		multiplier, published, err := fetchRechargeExchange(ctx)
		cache.nextCheck = now.Add(time.Minute)
		if err == nil && now.Sub(published) <= 7*24*time.Hour && !published.After(now) && !published.Before(cache.published) {
			cache.multiplier = multiplier
			cache.published = published
			cache.nextCheck = now.Add(time.Hour)
		}
	}
	if now.Sub(cache.published) <= 7*24*time.Hour {
		return cache.multiplier
	}
	return 0
}

func fetchRechargeExchange(ctx context.Context) (float64, time.Time, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rechargeExchangeURL, nil)
	if err != nil {
		return 0, time.Time{}, err
	}
	resp, err := rechargeExchangeClient.Do(req)
	if err != nil {
		return 0, time.Time{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, time.Time{}, fmt.Errorf("exchange rate feed returned HTTP %d", resp.StatusCode)
	}
	var payload struct {
		Data struct {
			LastDate string `json:"lastDate"`
		} `json:"data"`
		Records []struct {
			Pair  string `json:"vrtEName"`
			Price string `json:"price"`
		} `json:"records"`
	}
	if err := json.NewDecoder(io.LimitReader(resp.Body, 256<<10)).Decode(&payload); err != nil {
		return 0, time.Time{}, err
	}
	published, err := time.ParseInLocation("2006-01-02 15:04", payload.Data.LastDate, time.FixedZone("CST", 8*60*60))
	if err != nil {
		return 0, time.Time{}, err
	}
	for _, record := range payload.Records {
		if record.Pair == "USD/CNY" {
			rate, err := strconv.ParseFloat(record.Price, 64)
			if err == nil && rate > 1 && !math.IsInf(rate, 0) {
				return 1 / rate, published, nil
			}
			break
		}
	}
	return 0, time.Time{}, fmt.Errorf("exchange rate feed has no valid USD/CNY rate")
}
