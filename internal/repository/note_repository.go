package repository

import (
	"context"

	sqlc "github.com/meads/notes-api/internal/db/sqlc"
	"github.com/meads/notes-api/internal/domain"
)

// NoteSQLRepository wraps the standard sql.DB pool and the sqlc Querier.
// implements NoteRepository interface defined in domain.Note
//
//	type NoteSQLRepository struct {
//		db      *sql.DB
//		queries *sqlc.Queries
//	}
// func NewNoteRepository(db *sql.DB) *NoteSQLRepository {
// 	return &NoteSQLRepository{
// 		db:      db,
// 		queries: sqlc.New(db),
// 	}
// }

type NoteSQLRepository struct {
	queries sqlc.Querier
}

func NewNoteRepository(querier sqlc.Querier) *NoteSQLRepository {
	return &NoteSQLRepository{
		queries: querier,
	}
}

func (r *NoteSQLRepository) CreateNote(ctx context.Context, param domain.CreateNoteParams) (*domain.Note, error) {
	sqlcParam := sqlc.CreateNoteParams{
		Title:   param.Title,
		Content: param.Content,
		UserID:  param.UserID,
	}

	sqlcNote, err := r.queries.CreateNote(ctx, sqlcParam)
	if err != nil {
		return nil, err
	}

	return &domain.Note{
		ID:      sqlcNote.ID,
		Title:   sqlcNote.Title,
		Content: sqlcNote.Content,
		UserID:  sqlcNote.UserID,
	}, nil
}

func (r *NoteSQLRepository) GetNote(ctx context.Context, id int64) (*domain.Note, error) {
	sqlcNote, err := r.queries.GetNote(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Note{
		ID:      sqlcNote.ID,
		Title:   sqlcNote.Title,
		Content: sqlcNote.Content,
		UserID:  sqlcNote.UserID,
	}, nil
}

func (r *NoteSQLRepository) DeleteNote(ctx context.Context, param domain.DeleteNoteParams) error {
	sqlcParam := sqlc.DeleteNoteParams{
		ID:     param.ID,
		UserID: param.UserID,
	}

	return r.queries.DeleteNote(ctx, sqlcParam)
}

func (r *NoteSQLRepository) ListNotesByUserID(ctx context.Context, userID int64) ([]domain.Note, error) {
	sqlcNotes, err := r.queries.ListNotesByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	domainNotes := make([]domain.Note, 0, len(sqlcNotes))

	for _, n := range sqlcNotes {
		domainNotes = append(domainNotes, domain.Note{
			Content:   n.Content,
			CreatedAt: n.CreatedAt,
			ID:        n.ID,
			Title:     n.Title,
			UserID:    n.UserID,
		})
	}

	return domainNotes, nil
}

func (r *NoteSQLRepository) UpdateNote(ctx context.Context, param domain.UpdateNoteParams) error {
	sqlcParams := sqlc.UpdateNoteParams{
		Content: param.Content,
		ID:      param.ID,
		Title:   param.Title,
		UserID:  param.UserID,
	}

	return r.queries.UpdateNote(ctx, sqlcParams)
}
