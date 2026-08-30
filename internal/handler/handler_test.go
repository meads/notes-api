package handler

import (
	"net/http"
	"time"

	"github.com/meads/notes-api/internal/security"
)

func passClaimsMiddleware(r *http.Request, tokener *security.MockTokener) {
	tokenString := "mocktoken"
	r.AddCookie(&http.Cookie{
		Name:  "access_token",
		Value: tokenString,
	})

	userClaims := &security.UserClaims{Usage: "access"}
	tokener.EXPECT().VerifyToken(tokenString).Return(userClaims, nil)
}

var defaultDate time.Time = time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)
