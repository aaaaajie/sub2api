package middleware

import (
	"encoding/json"
	"mime"
	"net/http"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/httputil"
	"github.com/Wei-Shaw/sub2api/internal/pkg/requestmodel"
	"github.com/gin-gonic/gin"
)

// GroupSystemPrompt prepends the configured group prompt to supported JSON
// generation requests. It runs after API-key authentication so the prompt comes
// from the authenticated key's group snapshot, never from client input.
func GroupSystemPrompt() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey, ok := GetAPIKeyFromContext(c)
		if !ok || apiKey == nil || apiKey.Group == nil || strings.TrimSpace(apiKey.Group.SystemPrompt) == "" {
			c.Next()
			return
		}
		protocol := groupSystemPromptProtocol(c)
		if protocol == "" || !isJSONContentType(c.GetHeader("Content-Type")) {
			c.Next()
			return
		}

		body, err := httputil.ReadRequestBodyWithPrealloc(c.Request)
		if err != nil {
			if RejectRequestBodyReadError(c, "middleware.group_system_prompt", err) {
				return
			}
			AbortWithRequestBodyReadFailure(c, "middleware.group_system_prompt", err)
			return
		}
		modified, changed := prependGroupSystemPrompt(body, protocol, apiKey.Group.SystemPrompt)
		if changed {
			requestmodel.ResetRequestBody(c.Request, modified)
		}
		c.Next()
	}
}

func groupSystemPromptProtocol(c *gin.Context) string {
	if c.Request == nil || c.Request.Method != http.MethodPost {
		return ""
	}
	path := strings.TrimRight(c.Request.URL.Path, "/")
	switch path {
	case "/v1/messages", "/messages", "/antigravity/v1/messages":
		return "messages"
	case "/v1/chat/completions", "/chat/completions":
		return "chat"
	case "/v1/responses", "/responses", "/backend-api/codex/responses":
		return "responses"
	}

	if c.FullPath() != "/v1beta/models/*modelAction" && c.FullPath() != "/antigravity/v1beta/models/*modelAction" {
		return ""
	}
	action := strings.TrimPrefix(c.Param("modelAction"), "/")
	if strings.Contains(action, ":generateContent") || strings.Contains(action, ":streamGenerateContent") {
		return "gemini"
	}
	return ""
}

func isJSONContentType(value string) bool {
	mediaType, _, err := mime.ParseMediaType(value)
	return err == nil && (mediaType == "application/json" || strings.HasSuffix(mediaType, "+json"))
}

func prependGroupSystemPrompt(body []byte, protocol, prompt string) ([]byte, bool) {
	if strings.TrimSpace(prompt) == "" {
		return body, false
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(body, &payload); err != nil || payload == nil {
		return body, false
	}

	switch protocol {
	case "chat":
		var messages []json.RawMessage
		if current := payload["messages"]; len(current) > 0 && string(current) != "null" {
			if err := json.Unmarshal(current, &messages); err != nil || messages == nil {
				return body, false
			}
		}
		message, _ := json.Marshal(map[string]string{"role": "system", "content": prompt})
		messages = append([]json.RawMessage{message}, messages...)
		payload["messages"], _ = json.Marshal(messages)
	case "responses":
		var existing string
		if current := payload["instructions"]; len(current) > 0 && string(current) != "null" {
			if err := json.Unmarshal(current, &existing); err != nil {
				return body, false
			}
		}
		payload["instructions"], _ = json.Marshal(prefixPrompt(prompt, existing))
	case "messages":
		current := payload["system"]
		if len(current) == 0 || string(current) == "null" {
			payload["system"], _ = json.Marshal(prompt)
		} else if len(current) > 0 && current[0] == '"' {
			var existing string
			if err := json.Unmarshal(current, &existing); err != nil {
				return body, false
			}
			payload["system"], _ = json.Marshal(prefixPrompt(prompt, existing))
		} else if len(current) > 0 && current[0] == '[' {
			var blocks []json.RawMessage
			if err := json.Unmarshal(current, &blocks); err != nil {
				return body, false
			}
			block, _ := json.Marshal(map[string]string{"type": "text", "text": prompt})
			blocks = append([]json.RawMessage{block}, blocks...)
			payload["system"], _ = json.Marshal(blocks)
		} else {
			return body, false
		}
	case "gemini":
		var instruction map[string]json.RawMessage
		current := payload["systemInstruction"]
		if len(current) > 0 && string(current) != "null" {
			if err := json.Unmarshal(current, &instruction); err != nil || instruction == nil {
				return body, false
			}
		} else {
			instruction = make(map[string]json.RawMessage)
		}
		var parts []json.RawMessage
		if currentParts := instruction["parts"]; len(currentParts) > 0 && string(currentParts) != "null" {
			if err := json.Unmarshal(currentParts, &parts); err != nil || parts == nil {
				return body, false
			}
		}
		part, _ := json.Marshal(map[string]string{"text": prompt})
		parts = append([]json.RawMessage{part}, parts...)
		instruction["parts"], _ = json.Marshal(parts)
		payload["systemInstruction"], _ = json.Marshal(instruction)
	default:
		return body, false
	}

	modified, err := json.Marshal(payload)
	if err != nil {
		return body, false
	}
	return modified, true
}

func prefixPrompt(prompt, existing string) string {
	if strings.TrimSpace(existing) == "" {
		return prompt
	}
	return prompt + "\n\n" + existing
}
