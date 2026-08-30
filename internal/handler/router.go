package handler

import (
	"errors"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/meads/notes-api/internal/security"
)

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func claimsMiddleware(h gin.HandlerFunc, tokener security.Tokener) gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context) {

		// Read the access token from the cookie
		tokenString, err := ctx.Cookie("access_token")
		if err != nil {
			if err == http.ErrNoCookie {
				ctx.JSON(http.StatusUnauthorized, errorResponse(err))
				return
			}
			ctx.JSON(http.StatusBadRequest, errorResponse(err))
			return
		}

		userClaims, err := tokener.VerifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, errorResponse(err))
			return
		}

		if userClaims.Usage != "access" {
			ctx.JSON(http.StatusUnauthorized, errorResponse(errors.New("invalid token type")))
			return
		}

		h(ctx)
	})
}

func SetupRouter(
	authHandler *AuthHandler,
	userHandler *UserHandler,
	noteHandler *NoteHandler,
	tokener security.Tokener,
) *gin.Engine {
	config := cors.DefaultConfig()
	origins := os.Getenv("ALLOW_ORIGINS")
	if origins == "" {
		origins = "http://localhost:3000"
	}
	config.AllowOrigins = []string{origins}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"}
	config.ExposeHeaders = []string{"Authorization"}

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(cors.New(config))

	if authHandler != nil {
		router.POST("/register/", authHandler.Register)
		router.POST("/login/", authHandler.Login)
		router.POST("/refresh/", authHandler.RenewAccessToken)
		router.POST("/logout/:sessionid", authHandler.Logout)
		router.POST("/revoke/:sessionid", authHandler.RevokeSession)
	}
	if userHandler != nil {
		router.GET("/users/", claimsMiddleware(userHandler.ListUsers, tokener))
		router.PATCH("/users/", claimsMiddleware(userHandler.PatchUser, tokener))
		router.DELETE("/users/:userid/", claimsMiddleware(userHandler.DeleteUser, tokener))
	}
	if noteHandler != nil {
		router.POST("/users/:userid/notes/", claimsMiddleware(noteHandler.CreateNote, tokener))
		router.PUT("/users/:userid/notes/", claimsMiddleware(noteHandler.UpdateNote, tokener))
		router.GET("/users/:userid/notes/", claimsMiddleware(noteHandler.ListNotes, tokener))
		router.DELETE("/users/:userid/notes/:noteid", claimsMiddleware(noteHandler.DeleteNote, tokener))
	}

	return router
}
