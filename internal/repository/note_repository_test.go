package repository

import (
	"context"
	"errors"
	reflect "reflect"
	"testing"

	sqlc "github.com/meads/notes-api/internal/db/sqlc"
	"github.com/meads/notes-api/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

func TestCreateNote(t *testing.T) {
	tests := []struct {
		name              string
		expectedNote      *domain.Note
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "create note fails given querier create note returns a server error",
			expectedNote:  nil,
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				sqlcParam := sqlc.CreateNoteParams{Title: "title", Content: "content", UserID: int64(1)}
				querier.EXPECT().CreateNote(gomock.Any(), sqlcParam).
					Return(sqlc.Note{}, errors.New("server error"))
			},
		},
		{
			name:          "create note returns a domain note given querier create note returns sqlc note and no error",
			expectedNote:  &domain.Note{ID: int64(1), Title: "title", Content: "content", UserID: int64(1)},
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				sqlcParam := sqlc.CreateNoteParams{Title: "title", Content: "content", UserID: int64(1)}
				querier.EXPECT().CreateNote(gomock.Any(), sqlcParam).
					Return(sqlc.Note{ID: int64(1), Title: "title", Content: "content", UserID: int64(1)}, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockQuerier := NewMockQuerier(ctrl)

			test.setupExpectations(mockQuerier)
			domainCreateNoteParams := domain.CreateNoteParams{Title: "title", Content: "content", UserID: int64(1)}

			// Act
			noteRepository := NewNoteRepository(mockQuerier)
			actualNote, actualErr := noteRepository.
				CreateNote(context.Background(), domainCreateNoteParams)

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}

			if test.expectedNote == nil && actualNote != nil {
				t.Fatalf("expected nil note, got: %v", actualNote)
			}

		})
	}
}

func TestGetNote(t *testing.T) {
	tests := []struct {
		name              string
		expectedNote      *domain.Note
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "get note fails given querier get note returns a server error",
			expectedNote:  nil,
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				noteID := int64(1)
				querier.EXPECT().GetNote(gomock.Any(), noteID).Return(sqlc.Note{}, errors.New("server error"))
			},
		},
		{
			name:          "get note returns a domain note given querier get note returns sqlc note and no error",
			expectedNote:  &domain.Note{ID: int64(1), Title: "title", Content: "content", UserID: int64(1)},
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				noteID := int64(1)
				querier.EXPECT().GetNote(gomock.Any(), noteID).
					Return(sqlc.Note{ID: int64(1), Title: "title", Content: "content", UserID: int64(1)}, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockQuerier := NewMockQuerier(ctrl)

			test.setupExpectations(mockQuerier)
			noteID := int64(1)

			// Act
			noteRepository := NewNoteRepository(mockQuerier)
			actualNote, actualErr := noteRepository.GetNote(context.Background(), noteID)

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}

			if test.expectedNote == nil && actualNote != nil {
				t.Fatalf("expected nil note, got: %v", actualNote)
			}

		})
	}
}

func TestDeleteNote(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "delete note fails given querier delete note returns a server error",
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				deleteNoteParams := sqlc.DeleteNoteParams{
					ID:     int64(1),
					UserID: int64(1),
				}
				querier.EXPECT().DeleteNote(gomock.Any(), deleteNoteParams).Return(errors.New("server error"))
			},
		},
		{
			name:          "delete note returns nil given querier delete note succeeds",
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				deleteNoteParams := sqlc.DeleteNoteParams{
					ID:     int64(1),
					UserID: int64(1),
				}
				querier.EXPECT().DeleteNote(gomock.Any(), deleteNoteParams).Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockQuerier := NewMockQuerier(ctrl)

			test.setupExpectations(mockQuerier)

			deleteNoteParams := domain.DeleteNoteParams{
				ID:     int64(1),
				UserID: int64(1),
			}

			// Act
			noteRepository := NewNoteRepository(mockQuerier)
			actualErr := noteRepository.DeleteNote(context.Background(), deleteNoteParams)

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}
		})
	}
}

func TestListNotesByUserID(t *testing.T) {
	tests := []struct {
		name              string
		expectedNotes     []domain.Note
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "list notes fails given querier list notes returns a server error",
			expectedNotes: nil,
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				userID := int64(1)
				querier.EXPECT().ListNotesByUserID(gomock.Any(), userID).Return(nil, errors.New("server error"))
			},
		},
		{
			name: "list notes returns a slice of domain note given querier list notes returns sqlc notes",
			expectedNotes: []domain.Note{
				{ID: int64(1), Title: "title1", Content: "content1", UserID: int64(1)},
				{ID: int64(2), Title: "title2", Content: "content2", UserID: int64(1)},
				{ID: int64(3), Title: "title3", Content: "content3", UserID: int64(1)},
			},
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				userID := int64(1)
				querier.EXPECT().ListNotesByUserID(gomock.Any(), userID).
					Return(
						[]sqlc.Note{
							{ID: int64(1), Title: "title1", Content: "content1", UserID: int64(1)},
							{ID: int64(2), Title: "title2", Content: "content2", UserID: int64(1)},
							{ID: int64(3), Title: "title3", Content: "content3", UserID: int64(1)},
						}, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockQuerier := NewMockQuerier(ctrl)

			test.setupExpectations(mockQuerier)
			userID := int64(1)

			// Act
			noteRepository := NewNoteRepository(mockQuerier)
			actualNotes, actualErr := noteRepository.ListNotesByUserID(context.Background(), userID)

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}

			if !reflect.DeepEqual(test.expectedNotes, actualNotes) {
				t.Fatalf("expected: %v, got: %v", test.expectedNotes, actualNotes)
			}

		})
	}
}

func TestUpdateNote(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "update note fails given querier update note returns a server error",
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				sqlcParam := sqlc.UpdateNoteParams{Title: "title", Content: "content", ID: int64(1), UserID: int64(1)}
				querier.EXPECT().UpdateNote(gomock.Any(), sqlcParam).Return(errors.New("server error"))
			},
		},
		{
			name:          "update note returns nil given querier update note returns no error",
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				sqlcParam := sqlc.UpdateNoteParams{Title: "title", Content: "content", ID: int64(1), UserID: int64(1)}
				querier.EXPECT().UpdateNote(gomock.Any(), sqlcParam).Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockQuerier := NewMockQuerier(ctrl)

			test.setupExpectations(mockQuerier)
			domainUpdateNoteParams := domain.UpdateNoteParams{
				Title: "title", Content: "content", ID: int64(1), UserID: int64(1),
			}

			// Act
			noteRepository := NewNoteRepository(mockQuerier)
			actualErr := noteRepository.UpdateNote(context.Background(), domainUpdateNoteParams)

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}
		})
	}
}
