package handler

import (
	"github.com/meads/notes-api/internal/domain"
)

// Auth

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	SessionID string `json:"sessionId"`
	UserID    int64  `json:"userId"`
}

func MapToAuthResponse(authSessionResult *domain.AuthSessionResult) *AuthResponse {
	return &AuthResponse{
		SessionID: authSessionResult.SessionID,
		UserID:    authSessionResult.UserID,
	}
}

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type LoginResponse struct {
	RefreshToken string `json:"refreshToken"`
	AccessToken  string `json:"accessToken"`
	SessionID    string `json:"sessionId"`
	UserID       int64  `json:"userId"`
}

type RenewAccessTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type RenewAccessTokenResponse struct {
	AccessToken string `json:"accessToken"`
}

func MapToRenewAccessTokenResponse(result *domain.RenewAccessTokenResult) *RenewAccessTokenResponse {
	return &RenewAccessTokenResponse{
		AccessToken: result.AccessToken,
	}
}

// Users

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	// CreatedAt sql.NullTime `json:"createdAt"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

type ListUsersResponse struct {
	Users []UserResponse `json:"users"`
}

func MapUsersToListUsersResponse(users []domain.User) ListUsersResponse {
	userResponses := make([]UserResponse, 0, len(users))
	for _, u := range users {
		userResponses = append(userResponses, UserResponse{
			ID:       u.ID,
			Username: u.Username,
			Password: u.Password,
			// CreatedAt: u.CreatedAt,
		})
	}
	return ListUsersResponse{Users: userResponses}
}

type PatchUserRequest struct {
	ID              int64  `json:"id" binding:"required"`
	Username        string `json:"username" binding:"required"`
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required"`
}

type PatchUserResponse struct {
	Message string `json:"message"`
}

// Notes

type NoteResponse struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	UserID  int64  `json:"userId"`
}

type CreateNoteRequest struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
	UserID  int64  `json:"userId" binding:"required"`
}

type UpdateNoteRequest struct {
	ID      int64  `json:"id" binding:"required"`
	UserID  int64  `json:"userId" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type ListNotesResponse struct {
	Notes []NoteResponse `json:"notes"`
}

type DeleteNoteResponse struct {
	Message string `json:"message"`
}

func MapNotesToListNotesResponse(notes []domain.Note) ListNotesResponse {
	noteResponses := make([]NoteResponse, 0, len(notes))
	for _, n := range notes {
		noteResponses = append(noteResponses, NoteResponse{
			ID:      n.ID,
			Title:   n.Title,
			Content: n.Content,
			UserID:  n.UserID,
		})
	}
	return ListNotesResponse{Notes: noteResponses}
}
