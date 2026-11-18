package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGroupDisplayPlatform(t *testing.T) {
	g := &Group{Platform: PlatformOpenCodeGo, DisplayPlatform: PlatformDeepseek}
	require.Equal(t, PlatformDeepseek, g.PresentationPlatform())
	require.Equal(t, PlatformOpenCodeGo, g.Platform)
	copy := cloneGroupForDuplicate(g, "display-test")
	require.Equal(t, g.DisplayPlatform, copy.DisplayPlatform)
	require.Equal(t, g.Platform, copy.Platform)
	g.DisplayPlatform = ""
	require.Equal(t, PlatformOpenCodeGo, g.PresentationPlatform())
	_, err := (&adminServiceImpl{}).CreateGroup(context.Background(), &CreateGroupInput{RateMultiplier: 1, DisplayPlatform: "invalid-provider"})
	require.ErrorContains(t, err, "Unknown display platform")
}
