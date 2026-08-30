package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meads/notes-api/internal/domain"
)

type UserServicer interface {
	DeleteUser(ctx context.Context, id int64) error
	ListUsers(ctx context.Context, params domain.ListUsersParams) ([]domain.User, error)
	ChangePassword(ctx context.Context, params domain.UpdateUserPasswordParams) error
}

type UserHandler struct {
	userService UserServicer
}

func NewUserHandler(userService UserServicer) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	idParam := ctx.Param("userid")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "id parameter must be a valid integer",
		})

		return
	}

	err = h.userService.DeleteUser(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)

		return
	}

	ctx.JSON(http.StatusOK, DeleteUserResponse{
		Message: "Resource successfully deleted",
	})

}

func (h *UserHandler) ListUsers(ctx *gin.Context) {
	limit := ctx.Query("limit")
	if limit == "0" || limit == "" {
		limit = "50"
	}

	offset := ctx.Query("offset")
	if offset == "" {
		offset = "0"
	}

	i, err := strconv.ParseInt(limit, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error parsing limit as int",
		})

		return
	}

	j, err := strconv.ParseInt(offset, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "error parsing offset as int",
		})

		return
	}

	users, err := h.userService.ListUsers(ctx, domain.ListUsersParams{Limit: int32(i), Offset: int32(j)})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})

		return
	}

	ctx.JSON(http.StatusOK, MapUsersToListUsersResponse(users))
}

func (h *UserHandler) PatchUser(ctx *gin.Context) {
	var req PatchUserRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := h.userService.ChangePassword(ctx, domain.UpdateUserPasswordParams{
		ID:              req.ID,
		Username:        req.Username,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("change password failed: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, PatchUserResponse{
		Message: "Resource successfully patched",
	})
}
