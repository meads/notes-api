package domain

import (
	"context"
	"time"
)

type Note struct {
	ID        int64
	Title     string
	Content   string
	UserID    int64
	CreatedAt time.Time
}

type CreateNoteParams struct {
	Title   string
	Content string
	UserID  int64
}

type CreateNoteResult struct {
	ID      int64
	Title   string
	Content string
	UserID  int64
}

type DeleteNoteParams struct {
	ID     int64
	UserID int64
}

type UpdateNoteParams struct {
	ID      int64
	Title   string
	Content string
	UserID  int64
}

type NoteRepository interface {
	CreateNote(ctx context.Context, param CreateNoteParams) (*Note, error)
	GetNote(ctx context.Context, id int64) (*Note, error)
	DeleteNote(ctx context.Context, param DeleteNoteParams) error
	ListNotesByUserID(ctx context.Context, userID int64) ([]Note, error)
	UpdateNote(ctx context.Context, param UpdateNoteParams) error
}
