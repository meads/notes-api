package repository

import (
	"context"
	"database/sql"
	"errors"
	reflect "reflect"
	"testing"
	"time"

	sqlc "github.com/meads/notes-api/internal/db/sqlc"
	"github.com/meads/notes-api/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

func TestCreateSession(t *testing.T) {
	tests := []struct {
		name              string
		expectedSession   *domain.Session
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:            "create session fails given querier create session returns a server error",
			expectedSession: nil,
			expectedError:   errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				sqlcParam := sqlc.CreateSessionParams{
					ID:           "sessionid",
					UserID:       int64(1),
					RefreshToken: "refreshtoken",
					IsRevoked:    false,
					ExpiresAt:    time.Time{},
				}
				querier.EXPECT().CreateSession(gomock.Any(), sqlcParam).
					Return(sqlc.Session{}, errors.New("server error"))
			},
		},
		{
			name: "create session returns a domain session given querier create session succeeds",
			expectedSession: &domain.Session{
				ID:           "sessionid",
				UserID:       int64(1),
				RefreshToken: "refreshtoken",
				IsRevoked:    false,
				ExpiresAt:    time.Time{},
			},
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				sqlcParam := sqlc.CreateSessionParams{
					ID:           "sessionid",
					UserID:       int64(1),
					RefreshToken: "refreshtoken",
					IsRevoked:    false,
					ExpiresAt:    time.Time{},
				}
				querier.EXPECT().CreateSession(gomock.Any(), sqlcParam).
					Return(sqlc.Session{
						ID:           "sessionid",
						UserID:       int64(1),
						RefreshToken: "refreshtoken",
						IsRevoked:    false,
						ExpiresAt:    time.Time{},
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
			domainCreateSessionParams := domain.CreateSessionParams{
				ID:           "sessionid",
				UserID:       int64(1),
				RefreshToken: "refreshtoken",
				IsRevoked:    false,
				ExpiresAt:    time.Time{},
			}

			// Act
			sessionRepository := NewSessionRepository(mockQuerier)
			actualSession, actualErr := sessionRepository.
				CreateSession(context.Background(), domainCreateSessionParams)

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}

			if !reflect.DeepEqual(test.expectedSession, actualSession) {
				t.Fatalf("expected: \n%v, got: \n%v", test.expectedSession, actualSession)
			}
		})
	}
}

func TestGetSession(t *testing.T) {
	tests := []struct {
		name              string
		expectedSession   *domain.Session
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:            "get session fails given querier get session returns sql no rows error",
			expectedSession: nil,
			expectedError:   errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().GetSession(gomock.Any(), "sessionid").
					Return(sqlc.Session{}, sql.ErrNoRows)
			},
		},
		{
			name:            "get session fails given querier get session returns a server error",
			expectedSession: nil,
			expectedError:   errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().GetSession(gomock.Any(), "sessionid").
					Return(sqlc.Session{}, errors.New("server error"))
			},
		},
		{
			name: "get session returns a domain session given querier get session succeeds",
			expectedSession: &domain.Session{
				ID:           "sessionid",
				UserID:       int64(1),
				RefreshToken: "refreshtoken",
				IsRevoked:    false,
				ExpiresAt:    time.Time{},
			},
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().GetSession(gomock.Any(), "sessionid").
					Return(sqlc.Session{
						ID:           "sessionid",
						UserID:       int64(1),
						RefreshToken: "refreshtoken",
						IsRevoked:    false,
						ExpiresAt:    time.Time{},
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

			// Act
			sessionRepository := NewSessionRepository(mockQuerier)
			actualSession, actualErr := sessionRepository.GetSession(context.Background(), "sessionid")

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}

			if !reflect.DeepEqual(test.expectedSession, actualSession) {
				t.Fatalf("session aren't equal, expected: \n%v, got: \n%v", test.expectedSession, actualSession)
			}
		})
	}
}

func TestDeleteSession(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "delete session fails given querier delete session returns a server error",
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().DeleteSession(gomock.Any(), "sessionid").
					Return(errors.New("server error"))
			},
		},
		{
			name:          "delete session returns nil given querier delete session succeeds",
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().DeleteSession(gomock.Any(), "sessionid").Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockQuerier := NewMockQuerier(ctrl)

			test.setupExpectations(mockQuerier)

			// Act
			sessionRepository := NewSessionRepository(mockQuerier)
			actualErr := sessionRepository.DeleteSession(context.Background(), "sessionid")

			// Assert
			// Check if one has an error and the other does not
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}
		})
	}
}

func TestRevokeSession(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "revoke session fails given querier revoke session returns a server error",
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().RevokeSession(gomock.Any(), "sessionid").
					Return(errors.New("server error"))
			},
		},
		{
			name:          "revoke session returns nil given querier revoke session succeeds",
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().RevokeSession(gomock.Any(), "sessionid").Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			mockQuerier := NewMockQuerier(ctrl)

			test.setupExpectations(mockQuerier)

			// Act
			sessionRepository := NewSessionRepository(mockQuerier)
			actualErr := sessionRepository.RevokeSession(context.Background(), "sessionid")

			// Assert
			// Check if one has an error and the other does not
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}
		})
	}
}
