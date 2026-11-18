package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPlazaDisplayPlatformPreservesPricingAndSource(t *testing.T) {
	price := 0.123
	g := &service.PlazaGroup{Platform: service.PlatformOpenCodeGo, RateMultiplier: 0.9,
		Models: []service.PlazaModel{{Name: "deepseek-flash", Platform: service.PlatformOpenCodeGo, Pricing: &service.ChannelModelPricing{InputPrice: &price}}}}
	before := toModelPlazaGroupDTO(g, nil)
	g.DisplayPlatform = service.PlatformDeepseek
	after := toModelPlazaGroupDTO(g, nil)
	require.Equal(t, service.PlatformDeepseek, after.Platform)
	require.Equal(t, service.PlatformDeepseek, after.Models[0].Platform)
	require.Equal(t, before.Models[0].Pricing, after.Models[0].Pricing)
	require.Equal(t, before.RateMultiplier, after.RateMultiplier)
	require.Equal(t, service.PlatformOpenCodeGo, g.Models[0].Platform)
	require.Equal(t, service.PlatformOpenCodeGo, g.Platform)
	g.DisplayPlatform = ""
	require.Equal(t, before, toModelPlazaGroupDTO(g, nil))
}

func TestChannelDisplayPlatformKeepsModels(t *testing.T) {
	sections := []userChannelPlatformSection{{Platform: "opencode_go", Groups: []userAvailableGroup{
		{ID: 1, Platform: "opencode_go", displayPlatform: "deepseek"},
		{ID: 2, Platform: "opencode_go"},
	}, SupportedModels: []userSupportedModel{{Name: "deepseek-flash", Platform: "opencode_go"}}}}
	result := presentChannelSections(sections)
	require.Len(t, result, 2)
	require.Equal(t, "deepseek", result[0].Platform)
	require.Equal(t, "deepseek-flash", result[0].SupportedModels[0].Name)
	require.Equal(t, "deepseek", result[0].SupportedModels[0].Platform)
	require.Equal(t, "opencode_go", result[1].Platform)
	require.Equal(t, "opencode_go", sections[0].SupportedModels[0].Platform)
}
