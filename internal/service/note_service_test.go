package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	domain "github.com/meads/notes-api/internal/domain"
	"go.uber.org/mock/gomock"
)

func TestNoteService_CreateNote(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		userID            int64
		title             string
		content           string
		note              *domain.Note
		setupExpectations func(nrepo *MockNoteRepository, urepo *MockUserRepository)
	}{
		{
			name:          "note service create note fails given user repo error",
			expectedError: errors.New("user repo error"),
			userID:        int64(1),
			title:         "title",
			content:       "content",
			note:          nil,
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).
					Return(false, errors.New("user repo error"))
				nrepo.EXPECT().
					CreateNote(gomock.Any(), domain.CreateNoteParams{
						Content: "content", Title: "title", UserID: int64(1)}).
					Times(0)
			},
		},
		{
			name:          "note service create note fails given user not found",
			expectedError: errors.New("user not found"),
			userID:        int64(1),
			title:         "title",
			content:       "content",
			note:          nil,
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(false, nil)
				nrepo.EXPECT().
					CreateNote(gomock.Any(), domain.CreateNoteParams{
						Content: "content", Title: "title", UserID: int64(1)}).
					Times(0)
			},
		},
		{
			name:          "note service create note fails given note repo fails create note",
			expectedError: errors.New("note repo error"),
			userID:        int64(1),
			title:         "title",
			content:       "content",
			note:          nil,
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(true, nil)
				nrepo.EXPECT().
					CreateNote(gomock.Any(), domain.CreateNoteParams{
						Content: "content", Title: "title", UserID: int64(1)}).
					Return(nil, errors.New("note repo error"))
			},
		},
		{
			name:          "note service create note succeeds given valid data",
			expectedError: nil,
			userID:        int64(1),
			title:         "title",
			content:       "content",
			note: &domain.Note{
				ID: int64(1), UserID: int64(1), Title: "title", Content: "content", CreatedAt: time.Time{}},
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(true, nil)
				nrepo.EXPECT().
					CreateNote(gomock.Any(), domain.CreateNoteParams{
						Content: "content", Title: "title", UserID: int64(1)}).
					Return(&domain.Note{
						ID: int64(1), UserID: int64(1), Title: "title", Content: "content", CreatedAt: time.Time{}}, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockNoteRepo := NewMockNoteRepository(ctrl)
			mockUserRepo := NewMockUserRepository(ctrl)

			noteService := NewNoteService(mockNoteRepo, mockUserRepo)
			test.setupExpectations(mockNoteRepo, mockUserRepo)

			// Act
			note, err := noteService.
				CreateNote(context.Background(), test.userID, test.title, test.content)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}

			if !reflect.DeepEqual(test.note, note) {
				t.Fatalf("note does not match, expected: \n%v, got: \n%v", test.note, note)
			}
		})
	}
}

func TestNoteService_UpdateNote(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		noteID            int64
		userID            int64
		title             string
		content           string
		note              *domain.Note
		setupExpectations func(nrepo *MockNoteRepository, urepo *MockUserRepository)
	}{
		{
			name:          "note service update note fails given user repo error",
			expectedError: errors.New("user repo error"),
			noteID:        int64(1),
			userID:        int64(1),
			title:         "title",
			content:       "new content",
			note:          nil,
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).
					Return(false, errors.New("user repo error"))
				nrepo.EXPECT().UpdateNote(gomock.Any(), domain.UpdateNoteParams{
					Content: "new content", ID: int64(1), Title: "title", UserID: int64(1),
				}).Times(0)
			},
		},
		{
			name:          "note service update note fails given user not found error",
			expectedError: errors.New("user not found"),
			noteID:        int64(1),
			userID:        int64(1),
			title:         "title",
			content:       "new content",
			note:          nil,
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(false, nil)
				nrepo.EXPECT().UpdateNote(gomock.Any(), domain.UpdateNoteParams{
					Content: "new content", ID: int64(1), Title: "title", UserID: int64(1),
				}).Times(0)
			},
		},
		{
			name:          "note service update note fails given note repo update note returns an error",
			expectedError: errors.New("note repo error"),
			noteID:        int64(1),
			userID:        int64(1),
			title:         "title",
			content:       "new content",
			note:          nil,
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(true, nil)
				nrepo.EXPECT().UpdateNote(gomock.Any(), domain.UpdateNoteParams{
					Content: "new content",
					ID:      int64(1),
					Title:   "title",
					UserID:  int64(1),
				}).Return(errors.New("note repo error"))
			},
		},
		{
			name:          "note service update note succeeds given valid data",
			expectedError: nil,
			noteID:        int64(1),
			userID:        int64(1),
			title:         "title",
			content:       "new content",
			note:          &domain.Note{ID: int64(1), UserID: int64(1), Title: "title", Content: "new content"},
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(true, nil)
				nrepo.EXPECT().UpdateNote(gomock.Any(), domain.UpdateNoteParams{
					Content: "new content",
					ID:      int64(1),
					Title:   "title",
					UserID:  int64(1),
				}).Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockNoteRepo := NewMockNoteRepository(ctrl)
			mockUserRepo := NewMockUserRepository(ctrl)

			noteService := NewNoteService(mockNoteRepo, mockUserRepo)
			test.setupExpectations(mockNoteRepo, mockUserRepo)

			// Act
			note, err := noteService.
				UpdateNote(context.Background(), test.noteID, test.userID, test.title, test.content)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}

			if !reflect.DeepEqual(test.note, note) {
				t.Fatalf("note does not match, expected: \n%v, got: \n%v", test.note, note)
			}
		})
	}
}

func TestNoteService_ListNotes(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		userID            int64
		notes             []domain.Note
		setupExpectations func(nrepo *MockNoteRepository, urepo *MockUserRepository)
	}{
		{
			name:          "note service list notes fails given user repo error",
			expectedError: errors.New("user repo error"),
			userID:        int64(1),
			notes:         nil,
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(false, errors.New("user repo error"))
				nrepo.EXPECT().ListNotesByUserID(gomock.Any(), int64(1)).Times(0)
			},
		},
		{
			name:          "note service list notes fails given user not found",
			expectedError: errors.New("user not found error"),
			userID:        int64(1),
			notes:         nil,
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(false, nil)
				nrepo.EXPECT().ListNotesByUserID(gomock.Any(), int64(1)).Times(0)
			},
		},
		{
			name:          "note service list notes fails given note repo error",
			expectedError: errors.New("note repo error"),
			userID:        int64(1),
			notes:         nil,
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(true, nil)
				nrepo.EXPECT().ListNotesByUserID(gomock.Any(), int64(1)).
					Return(nil, errors.New("note repo error"))
			},
		},
		{
			name:          "note service list notes succeeds given valid data supplied",
			expectedError: nil,
			userID:        int64(1),
			notes: []domain.Note{
				{ID: int64(1), UserID: int64(1), Title: "note1", Content: "content1"},
				{ID: int64(2), UserID: int64(1), Title: "note2", Content: "content2"},
				{ID: int64(3), UserID: int64(1), Title: "note3", Content: "content3"},
			},
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(true, nil)
				nrepo.EXPECT().ListNotesByUserID(gomock.Any(), int64(1)).
					Return([]domain.Note{
						{ID: int64(1), UserID: int64(1), Title: "note1", Content: "content1"},
						{ID: int64(2), UserID: int64(1), Title: "note2", Content: "content2"},
						{ID: int64(3), UserID: int64(1), Title: "note3", Content: "content3"},
					}, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockNoteRepo := NewMockNoteRepository(ctrl)
			mockUserRepo := NewMockUserRepository(ctrl)

			noteService := NewNoteService(mockNoteRepo, mockUserRepo)
			test.setupExpectations(mockNoteRepo, mockUserRepo)

			// Act
			notes, err := noteService.ListNotes(context.Background(), test.userID)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}

			if !reflect.DeepEqual(test.notes, notes) {
				t.Fatalf("note does not match, expected: \n%v, got: \n%v", test.notes, notes)
			}
		})
	}
}

func TestNoteService_DeleteNote(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		noteID            int64
		userID            int64
		setupExpectations func(nrepo *MockNoteRepository, urepo *MockUserRepository)
	}{
		{
			name:          "note service delete note fails given user repo error",
			expectedError: errors.New("user repo error"),
			noteID:        int64(1),
			userID:        int64(1),
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).
					Return(false, errors.New("user repo error"))
				nrepo.EXPECT().DeleteNote(gomock.Any(),
					domain.DeleteNoteParams{ID: int64(1), UserID: int64(1)}).
					Times(0)
			},
		},
		{
			name:          "note service delete note fails given user not found error",
			expectedError: errors.New("user not found error"),
			noteID:        int64(1),
			userID:        int64(1),
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(false, nil)
				nrepo.EXPECT().DeleteNote(gomock.Any(),
					domain.DeleteNoteParams{ID: int64(1), UserID: int64(1)}).
					Times(0)
			},
		},
		{
			name:          "note service delete note fails given note repo error",
			expectedError: errors.New("note repo error"),
			noteID:        int64(1),
			userID:        int64(1),
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(true, nil)
				nrepo.EXPECT().DeleteNote(gomock.Any(),
					domain.DeleteNoteParams{ID: int64(1), UserID: int64(1)}).
					Return(errors.New("note repo error"))
			},
		},
		{
			name:          "note service delete note succeeds given valid data supplied",
			expectedError: nil,
			noteID:        int64(1),
			userID:        int64(1),
			setupExpectations: func(nrepo *MockNoteRepository, urepo *MockUserRepository) {
				urepo.EXPECT().UserExists(gomock.Any(), int64(1)).Return(true, nil)
				nrepo.EXPECT().DeleteNote(gomock.Any(),
					domain.DeleteNoteParams{ID: int64(1), UserID: int64(1)}).Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockNoteRepo := NewMockNoteRepository(ctrl)
			mockUserRepo := NewMockUserRepository(ctrl)

			noteService := NewNoteService(mockNoteRepo, mockUserRepo)
			test.setupExpectations(mockNoteRepo, mockUserRepo)

			// Act
			err := noteService.DeleteNote(context.Background(), test.noteID, test.userID)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}
		})
	}
}
