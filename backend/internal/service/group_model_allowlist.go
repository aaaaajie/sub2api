package service

import (
	"net/http"
	"slices"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/domain"
	"github.com/Wei-Shaw/sub2api/internal/pkg/claude"
	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/openai"
)

// GroupModelAllowlist 是 service 层的分组模型白名单（与 domain.GroupModelAllowlist
// 字段一致，ent 持久化用 domain 类型，边界处显式转换）。
type GroupModelAllowlist struct {
	Enabled      bool     `json:"enabled"`
	Models       []string `json:"models,omitempty"`
	HiddenModels []string `json:"hidden_models,omitempty"`
	DefaultModel string   `json:"default_model,omitempty"`
}

// DomainGroupModelAllowlist 把 service 白名单转换为 ent 持久化使用的 domain 类型。
func DomainGroupModelAllowlist(cfg GroupModelAllowlist) domain.GroupModelAllowlist {
	return domain.GroupModelAllowlist{Enabled: cfg.Enabled, Models: cfg.Models, HiddenModels: cfg.HiddenModels, DefaultModel: cfg.DefaultModel}
}

// GroupModelAllowlistFromDomain 把 ent 读出的 domain 白名单转换为 service 类型。
func GroupModelAllowlistFromDomain(cfg domain.GroupModelAllowlist) GroupModelAllowlist {
	return GroupModelAllowlist{Enabled: cfg.Enabled, Models: cfg.Models, HiddenModels: cfg.HiddenModels, DefaultModel: cfg.DefaultModel}
}

// supplementUnmappedOpenAIModels ensures a partial mapping catalog does not
// hide models from unmapped or unrestricted passthrough OpenAI accounts.
// An empty catalog is left unchanged so callers retain their existing
// discovery fallback.
func supplementUnmappedOpenAIModels(accounts []Account, models []string) []string {
	if len(models) == 0 {
		return models
	}
	for i := range accounts {
		account := &accounts[i]
		unrestrictedPassthrough := account.IsOpenAIPassthroughEnabled() && !account.IsOpenAIPassthroughModelFilterEnabled()
		if account.Platform == PlatformOpenAI && (unrestrictedPassthrough || len(account.GetModelMapping()) == 0) {
			return dedupeAndSortModelIDs(slices.Concat(models, openai.DefaultModelIDs()))
		}
	}
	return models
}

// normalizeGroupModelAllowlist 归一化管理端提交的分组模型白名单：
// 条目 TrimSpace、按小写去重保序；`*` 可出现在任意位置；
// enabled=true 且列表为空视为配置错误，返回 400 而不是运行时静默放行/拒绝。
func normalizeGroupModelAllowlist(cfg GroupModelAllowlist) (GroupModelAllowlist, error) {
	out := GroupModelAllowlist{
		Enabled:      cfg.Enabled,
		Models:       normalizeGroupModelPatterns(cfg.Models),
		HiddenModels: normalizeGroupModelPatterns(cfg.HiddenModels),
		DefaultModel: strings.TrimSpace(cfg.DefaultModel),
	}
	if out.Enabled && len(out.Models) == 0 && len(out.HiddenModels) == 0 {
		return out, infraerrors.New(http.StatusBadRequest, "INVALID_MODEL_ALLOWLIST", "model allowlist cannot be enabled with an empty model list")
	}
	if out.Enabled && out.DefaultModel != "" && !out.AllowsListing(out.DefaultModel) {
		return out, infraerrors.New(http.StatusBadRequest, "DEFAULT_MODEL_NOT_ALLOWLISTED", "default model must be visible in the enabled model allowlist")
	}
	return out, nil
}

func normalizeGroupModelPatterns(models []string) []string {
	seen := make(map[string]struct{}, len(models))
	var normalized []string
	for _, model := range models {
		model = strings.TrimSpace(model)
		if model == "" {
			continue
		}
		key := strings.ToLower(model)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		normalized = append(normalized, model)
	}
	return normalized
}

// ModelAllowlistEnabled 报告该分组是否启用了模型白名单。
// 开启后所有携带模型的网关请求与模型列表接口都受白名单约束。
func (g *Group) ModelAllowlistEnabled() bool {
	return g != nil && g.ModelAllowlist.Enabled
}

// Allows 判断客户端请求的模型是否命中白名单。
// 准入只看客户端书写的模型名，与账号映射、渠道映射、合成路由改写无关；
// 候选形式覆盖代码中已有的模型名等价规则（Gemini models/ 前缀、
// Antigravity/Claude -thinking 宽容规则、OpenAI 推理后缀），不做模糊匹配。
func (a GroupModelAllowlist) Allows(model string) bool {
	if !a.Enabled {
		return true
	}
	model = strings.TrimSpace(model)
	if model == "" {
		return true
	}
	candidates := groupModelAllowlistCandidates(model)
	return matchesGroupModelPatterns(a.Models, candidates) || matchesGroupModelPatterns(a.HiddenModels, candidates)
}

// AllowsListing only admits models selected for display when the allowlist is enabled.
func (a GroupModelAllowlist) AllowsListing(model string) bool {
	if !a.Enabled {
		return true
	}
	candidates := groupModelAllowlistCandidates(model)
	return matchesGroupModelPatterns(a.Models, candidates) && !matchesGroupModelPatterns(a.HiddenModels, candidates)
}

// IsHiddenFromListing applies the explicit visibility choice to descriptive catalogs.
func (a GroupModelAllowlist) IsHiddenFromListing(model string) bool {
	return a.Enabled && len(a.HiddenModels) > 0 &&
		matchesGroupModelPatterns(a.HiddenModels, groupModelAllowlistCandidates(model))
}

func matchesGroupModelPatterns(patterns, candidates []string) bool {
	for _, entry := range patterns {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		for _, candidate := range candidates {
			if groupAllowlistPatternMatches(entry, candidate) {
				return true
			}
		}
	}
	return false
}

// groupAllowlistPatternMatches treats only * as a wildcard, including in the middle
// or at either end. Matching is anchored and case-insensitive.
func groupAllowlistPatternMatches(pattern, model string) bool {
	pattern = strings.ToLower(pattern)
	model = strings.ToLower(model)
	parts := strings.Split(pattern, "*")
	if len(parts) == 1 {
		return pattern == model
	}
	if !strings.HasPrefix(model, parts[0]) {
		return false
	}
	model = model[len(parts[0]):]
	for _, part := range parts[1 : len(parts)-1] {
		index := strings.Index(model, part)
		if index < 0 {
			return false
		}
		model = model[index+len(part):]
	}
	return strings.HasSuffix(model, parts[len(parts)-1])
}

// groupModelAllowlistCandidates 返回客户端模型名在白名单匹配中的候选形式（均已小写）。
func groupModelAllowlistCandidates(model string) []string {
	candidates := make([]string, 0, 4)
	add := func(value string) {
		value = strings.ToLower(strings.TrimSpace(value))
		if value == "" {
			return
		}
		for _, existing := range candidates {
			if existing == value {
				return
			}
		}
		candidates = append(candidates, value)
	}

	add(model)
	add(strings.TrimPrefix(model, "models/"))
	add(claude.NormalizeModelID(strings.TrimSuffix(model, "-thinking")))
	add(NormalizeOpenAICompatRequestedModel(model))
	return candidates
}

// FilterForListing 按白名单条目顺序生成模型列表输出。
// 精确条目沿用既有规则：只有出现在 source（账号映射键 ∪ 平台默认列表）的模式
// 集合中才输出；通配条目展开为 source 中所有匹配项并保持 source 顺序；全局去重。
// 当 source 是前缀通配模式且白名单模式的固定前缀位于该前缀之内时，
// 可直接输出白名单模式：它的所有实例都属于 source。其他未处理的无限
// 交集不会被枚举为模型 ID，列表可能省略仍可请求的模型。
func (a GroupModelAllowlist) FilterForListing(source []string) []string {
	if !a.Enabled {
		return source
	}
	if len(a.Models) == 0 {
		// 可以只启用不公开列出的模型，此时目录为空。
		return nil
	}

	patterns := make([]string, 0, len(source))
	for _, pattern := range source {
		pattern = strings.TrimSpace(pattern)
		if pattern != "" {
			patterns = append(patterns, pattern)
		}
	}
	if len(patterns) == 0 {
		return nil
	}

	seen := make(map[string]struct{}, len(a.Models)+len(patterns))
	filtered := make([]string, 0, len(patterns))
	add := func(model string) {
		if a.IsHiddenFromListing(model) {
			return
		}
		key := strings.ToLower(model)
		if _, ok := seen[key]; ok {
			return
		}
		seen[key] = struct{}{}
		filtered = append(filtered, model)
	}

	for _, entry := range a.Models {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if strings.Contains(entry, "*") {
			for _, pattern := range patterns {
				if groupAllowlistPatternMatches(entry, pattern) {
					add(pattern)
				} else if strings.HasSuffix(pattern, "*") &&
					!strings.Contains(strings.TrimSuffix(pattern, "*"), "*") &&
					strings.HasPrefix(strings.ToLower(entry), strings.ToLower(strings.TrimSuffix(pattern, "*"))) {
					add(entry)
				}
			}
			continue
		}
		if allowlistSourcePatternAllowsModel(patterns, entry) {
			add(entry)
		}
	}
	return filtered
}

// OrderForListing applies the saved model order when the allowlist is disabled,
// and filters plus orders the catalog when it is enabled. The selected default
// model remains a request default and does not override the configured order.
func (a GroupModelAllowlist) OrderForListing(source []string) []string {
	if a.Enabled {
		return a.FilterForListing(source)
	}
	if len(a.Models) == 0 {
		return source
	}
	return a.OrderForDisplay(source)
}

// OrderForDisplay applies the group's preferred model order without filtering.
// The model plaza retains its pricing entries after applying explicit hiding.
func (a GroupModelAllowlist) OrderForDisplay(source []string) []string {
	ordered := make([]string, 0, len(source))
	seen := make(map[string]struct{}, len(source))
	add := func(model string) {
		key := strings.ToLower(strings.TrimSpace(model))
		if key == "" {
			return
		}
		if _, exists := seen[key]; exists {
			return
		}
		seen[key] = struct{}{}
		ordered = append(ordered, model)
	}
	for _, preferred := range a.Models {
		preferred = strings.TrimSpace(preferred)
		if preferred == "" {
			continue
		}
		for _, model := range source {
			if groupAllowlistPatternMatches(preferred, strings.TrimSpace(model)) {
				add(model)
			}
		}
	}
	for _, model := range source {
		add(model)
	}
	return ordered
}

// allowlistSourcePatternAllowsModel 沿用 filterModelsByCustomList 时代的匹配规则：
// 精确相等、source 通配模式的前缀匹配，以及 Claude 归一化（-thinking 后缀）
// 后的精确匹配。比较与 Allows 一律大小写不敏感。
func allowlistSourcePatternAllowsModel(patterns []string, model string) bool {
	for _, pattern := range patterns {
		if strings.EqualFold(pattern, model) {
			return true
		}
		if strings.HasSuffix(pattern, "*") && strings.HasPrefix(strings.ToLower(model), strings.ToLower(strings.TrimSuffix(pattern, "*"))) {
			return true
		}
	}
	normalizedClaudeModel := claude.NormalizeModelID(strings.TrimSuffix(model, "-thinking"))
	if !strings.EqualFold(normalizedClaudeModel, model) {
		for _, pattern := range patterns {
			if strings.EqualFold(pattern, normalizedClaudeModel) {
				return true
			}
		}
	}
	return false
}
