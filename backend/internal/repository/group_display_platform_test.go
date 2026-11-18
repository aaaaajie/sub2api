package repository

import (
	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGroupDisplayPlatformEntityMapping(t *testing.T) {
	row := &dbent.Group{ID: 1, Platform: "opencode_go", DisplayPlatform: "deepseek"}
	group := groupEntityToService(row)
	require.Equal(t, "opencode_go", group.Platform)
	require.Equal(t, "deepseek", group.DisplayPlatform)
}
