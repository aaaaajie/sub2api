package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/adminaccess"
	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserMailHandler struct {
	email *service.AdminUserMailService
}

func NewUserMailHandler(email *service.AdminUserMailService) *UserMailHandler {
	return &UserMailHandler{email: email}
}

func (h *UserMailHandler) RequireFullAccess(c *gin.Context) {
	if adminaccess.RestrictedFromContext(c.Request.Context()) != nil {
		response.Error(c, http.StatusForbidden, "User email requires full administrator access.")
		c.Abort()
		return
	}
	c.Next()
}

type userMailFilters struct {
	Role          string           `json:"role" binding:"omitempty,oneof=admin user"`
	Status        string           `json:"status" binding:"omitempty,oneof=active disabled"`
	Search        string           `json:"search" binding:"max=100"`
	GroupName     string           `json:"group_name" binding:"max=160"`
	APIKeyGroupID int64            `json:"api_key_group_id" binding:"min=0"`
	Attributes    map[int64]string `json:"attributes" binding:"max=100"`
}

func (h *UserMailHandler) Preview(c *gin.Context) {
	var req struct {
		Draft   service.AdminUserMailDraft `json:"draft"`
		IDs     []int64                    `json:"user_ids" binding:"max=500,dive,gt=0"`
		Filters *userMailFilters           `json:"filters"`
	}
	if c.Request.Body != nil {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 128<<10)
	}
	if err := c.ShouldBindJSON(&req); err != nil || (len(req.IDs) > 0 && req.Filters != nil) {
		response.BadRequest(c, "Invalid preview parameters.")
		return
	}
	var filters *service.UserListFilters
	if in := req.Filters; in != nil {
		filters = &service.UserListFilters{
			Role: in.Role, Status: in.Status, Search: strings.TrimSpace(in.Search), GroupName: strings.TrimSpace(in.GroupName),
			APIKeyGroupID: in.APIKeyGroupID, Attributes: in.Attributes,
		}
	}
	preview, err := h.email.Preview(c.Request.Context(), req.Draft, req.IDs, filters)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, preview)
}

func (h *UserMailHandler) Queue(c *gin.Context) {
	actor, ok := middleware.GetAuthSubjectFromContext(c)
	if !ok || actor.UserID <= 0 {
		response.Error(c, http.StatusUnauthorized, "Administrator session required.")
		return
	}
	var req struct {
		ID     string                     `json:"id" binding:"required"`
		Draft  service.AdminUserMailDraft `json:"draft"`
		IDs    []int64                    `json:"user_ids" binding:"max=500,dive,gt=0"`
		IsTest bool                       `json:"is_test"`
	}
	if c.Request.Body != nil {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 128<<10)
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid email parameters.")
		return
	}
	id, err := uuid.Parse(req.ID)
	if err != nil || id == uuid.Nil {
		response.BadRequest(c, "Invalid email job ID.")
		return
	}
	job, err := h.email.Queue(c.Request.Context(), actor.UserID, id, req.Draft, req.IDs, req.IsTest)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Accepted(c, job)
}

func (h *UserMailHandler) List(c *gin.Context) {
	page, limit := response.ParsePagination(c)
	jobs, total, err := h.email.List(c.Request.Context(), pagination.PaginationParams{Page: page, PageSize: limit})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Paginated(c, jobs, total, page, limit)
}

func (h *UserMailHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "Invalid email job ID.")
		return
	}
	job, err := h.email.Get(c.Request.Context(), id)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, job)
}

func (h *UserMailHandler) Retry(c *gin.Context) {
	id, idErr := uuid.Parse(c.Param("id"))
	userID, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if idErr != nil || userErr != nil || userID <= 0 {
		response.BadRequest(c, "Invalid email job or recipient ID.")
		return
	}
	job, err := h.email.Retry(c.Request.Context(), id, userID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, job)
}
