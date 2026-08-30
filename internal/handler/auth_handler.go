package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meads/notes-api/internal/domain"
)

type AuthServicer interface {
	Register(ctx context.Context, username, password string) (*domain.AuthSessionResult, error)
	Login(ctx context.Context, username, password string) (*domain.AuthSessionResult, error)
	Logout(ctx context.Context, sessionID string) error
	RenewAccessToken(ctx context.Context, refreshToken string) (*domain.RenewAccessTokenResult, error)
	RevokeSession(ctx context.Context, sessionID string) error
}

type AuthHandler struct {
	authService  AuthServicer
	cookieDomain string
}

func NewAuthHandler(authService AuthServicer, cookieDomain string) *AuthHandler {
	return &AuthHandler{authService: authService, cookieDomain: cookieDomain}
}

func (h *AuthHandler) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authSessionResult, err := h.authService.Register(ctx, req.Username, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("Registration failed")))
		return
	}

	ctx.SetCookieData(&http.Cookie{
		Domain:   h.cookieDomain,
		Name:     "access_token",
		Value:    authSessionResult.AccessToken,
		Path:     "/",
		Expires:  authSessionResult.AccessTokenExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	ctx.SetCookieData(&http.Cookie{
		Domain:   h.cookieDomain,
		Name:     "refresh_token",
		Value:    authSessionResult.RefreshToken,
		Path:     "/refresh",
		Expires:  authSessionResult.RefreshTokenExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	ctx.JSON(http.StatusOK, MapToAuthResponse(authSessionResult))
}

func (h *AuthHandler) Login(ctx *gin.Context) {
	var req LoginRequest

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	authSessionResult, err := h.authService.Login(ctx, req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrInvalidCredentials):
			ctx.JSON(http.StatusUnauthorized, errorResponse(domain.ErrInvalidCredentials))

		default:
			ctx.JSON(http.StatusInternalServerError, errorResponse(errors.New("internal server error")))
		}
		return
	}

	ctx.SetCookieData(&http.Cookie{
		Domain:   h.cookieDomain,
		Name:     "access_token",
		Value:    authSessionResult.AccessToken,
		Path:     "/",
		Expires:  authSessionResult.AccessTokenExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	ctx.SetCookieData(&http.Cookie{
		Domain:   h.cookieDomain,
		Name:     "refresh_token",
		Value:    authSessionResult.RefreshToken,
		Path:     "/refresh",
		Expires:  authSessionResult.RefreshTokenExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	ctx.JSON(http.StatusOK, MapToAuthResponse(authSessionResult))
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	idParam := ctx.Param("sessionid")
	err := h.authService.Logout(ctx, idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("logout failed %w", err)))
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) RenewAccessToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	renewAccessTokenResult, err := h.authService.RenewAccessToken(ctx, refreshToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error renewing access token: %w", err)))
		return
	}

	ctx.SetCookieData(&http.Cookie{
		Domain:   h.cookieDomain,
		Name:     "access_token",
		Value:    renewAccessTokenResult.AccessToken,
		Path:     "/",
		Expires:  renewAccessTokenResult.AccessTokenExpiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode,
	})

	ctx.JSON(http.StatusOK, nil)
}

func (h *AuthHandler) RevokeSession(ctx *gin.Context) {
	idParam := ctx.Param("sessionid")
	err := h.authService.RevokeSession(ctx, idParam)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("error revoking session: %w", err))
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}
