package handler

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/meads/notes-api/internal/domain"
)

type NoteServicer interface {
	CreateNote(ctx context.Context, userID int64, title, content string) (*domain.Note, error)
	UpdateNote(ctx context.Context, noteID, userID int64, title, content string) (*domain.Note, error)
	ListNotes(ctx context.Context, userID int64) ([]domain.Note, error)
	DeleteNote(ctx context.Context, noteID, userID int64) error
}

type NoteHandler struct {
	noteService NoteServicer
}

func NewNoteHandler(noteService NoteServicer) *NoteHandler {
	return &NoteHandler{noteService: noteService}
}

func (h *NoteHandler) CreateNote(ctx *gin.Context) {
	var req CreateNoteRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	note, err := h.noteService.CreateNote(ctx, req.UserID, req.Title, req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, NoteResponse{
		ID:      note.ID,
		Title:   note.Title,
		Content: note.Content,
		UserID:  note.UserID,
	})
}

func (h *NoteHandler) UpdateNote(ctx *gin.Context) {
	var req UpdateNoteRequest
	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	note, err := h.noteService.UpdateNote(ctx, req.ID, req.UserID, req.Title, req.Content)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, NoteResponse{
		ID:      note.ID,
		Title:   note.Title,
		Content: note.Content,
		UserID:  note.UserID,
	})
}

func (h *NoteHandler) ListNotes(ctx *gin.Context) {
	userIDParam := ctx.Param("userid")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid user id param")))
		return
	}

	notes, err := h.noteService.ListNotes(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, MapNotesToListNotesResponse(notes))
}

func (h *NoteHandler) DeleteNote(ctx *gin.Context) {
	userIDParam := ctx.Param("userid")
	userID, err := strconv.ParseInt(userIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid user id param")))
		return
	}

	noteIDParam := ctx.Param("noteid")
	noteID, err := strconv.ParseInt(noteIDParam, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(errors.New("invalid note id param")))
		return
	}

	err = h.noteService.DeleteNote(ctx, noteID, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Errorf("delete note failed: %w", err))
		return
	}

	ctx.JSON(http.StatusOK, DeleteNoteResponse{
		Message: "Resource successfully deleted",
	})
}
