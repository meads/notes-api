package repository

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	sqlc "github.com/meads/notes-api/internal/db/sqlc"
	"github.com/meads/notes-api/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name              string
		expectedUser      *domain.User
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "create user fails given querier create user returns a server error",
			expectedUser:  nil,
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				sqlcParam := sqlc.CreateUserParams{
					Username: "username",
					Password: "password",
				}
				querier.EXPECT().CreateUser(gomock.Any(), sqlcParam).
					Return(sqlc.User{}, errors.New("server error"))
			},
		},
		{
			name: "create user returns a domain user given querier create user returns no error",
			expectedUser: &domain.User{
				ID:       int64(1),
				Username: "username",
				Password: "password",
			},
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				sqlcParam := sqlc.CreateUserParams{
					Username: "username",
					Password: "password",
				}
				querier.EXPECT().CreateUser(gomock.Any(), sqlcParam).
					Return(sqlc.User{ID: int64(1), Username: "username", Password: "password"}, nil)
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
			userRepository := NewUserRepository(mockQuerier)
			actualUser, actualErr := userRepository.
				CreateUser(context.Background(), "username", "password")

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}

			if test.expectedUser == nil && actualUser != nil {
				t.Fatalf("expected nil note, got: %v", actualUser)
			}

		})
	}
}

func TestUsernameExists(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "username exists fails given querier username exists returns a server error",
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().UsernameExists(gomock.Any(), "username").
					Return(false, errors.New("server error"))
			},
		},
		{
			name:          "username exists succeeds given querier username exists returns no error",
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().UsernameExists(gomock.Any(), "username").
					Return(false, nil)
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
			userRepository := NewUserRepository(mockQuerier)
			_, actualErr := userRepository.
				UsernameExists(context.Background(), "username")

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "delete user fails given querier delete user returns a server error",
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().DeleteUser(gomock.Any(), int64(1)).
					Return(errors.New("server error"))
			},
		},
		{
			name:          "delete user succeeds given querier delete user returns no error",
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().DeleteUser(gomock.Any(), int64(1)).
					Return(nil)
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
			userRepository := NewUserRepository(mockQuerier)
			actualErr := userRepository.DeleteUser(context.Background(), int64(1))

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name              string
		expectedUser      *domain.User
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "get user fails given querier get user returns sql no rows error",
			expectedUser:  nil,
			expectedError: errors.New("user not found"),
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().GetUser(gomock.Any(), int64(1)).
					Return(sqlc.User{}, sql.ErrNoRows)
			},
		},
		{
			name:          "get user fails given querier get user returns a server error",
			expectedUser:  nil,
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().GetUser(gomock.Any(), int64(1)).
					Return(sqlc.User{}, errors.New("server error"))
			},
		},
		{
			name: "get user returns a domain user given querier get user returns no error",
			expectedUser: &domain.User{
				ID:       int64(1),
				Username: "username",
				Password: "password",
			},
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().GetUser(gomock.Any(), int64(1)).
					Return(sqlc.User{ID: int64(1), Username: "username", Password: "password"}, nil)
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
			userRepository := NewUserRepository(mockQuerier)
			actualUser, actualErr := userRepository.GetUser(context.Background(), int64(1))

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}

			if !reflect.DeepEqual(test.expectedUser, actualUser) {
				t.Fatalf("expected: \n%v, got: \n%v", test.expectedUser, actualUser)
			}
		})
	}
}

func TestListUsers(t *testing.T) {
	tests := []struct {
		name              string
		expectedUsers     []domain.User
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "list users fails given querier list users returns a server error",
			expectedUsers: nil,
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				sqlcParams := sqlc.ListUsersParams{Limit: 0, Offset: 0}
				querier.EXPECT().ListUsers(gomock.Any(), sqlcParams).
					Return(nil, errors.New("server error"))
			},
		},
		{
			name: "list users returns a slice of domain user given querier list users returns no error",
			expectedUsers: []domain.User{
				{ID: int64(1), Username: "username", Password: "password"},
				{ID: int64(2), Username: "username", Password: "password"},
			},
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				sqlcParams := sqlc.ListUsersParams{Limit: 0, Offset: 0}
				querier.EXPECT().ListUsers(gomock.Any(), sqlcParams).
					Return([]sqlc.User{
						{ID: int64(1), Username: "username", Password: "password"},
						{ID: int64(2), Username: "username", Password: "password"},
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
			userRepository := NewUserRepository(mockQuerier)

			domainListUsersParams := domain.ListUsersParams{Limit: 0, Offset: 0}
			actualUsers, actualErr := userRepository.ListUsers(context.Background(), domainListUsersParams)

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}

			if !reflect.DeepEqual(test.expectedUsers, actualUsers) {
				t.Fatalf("expected: \n%v, got: \n%v", test.expectedUsers, actualUsers)
			}
		})
	}
}

func TestUpdateUserPassword(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "update user password fails given querier update user password returns a server error",
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				sqlcParam := sqlc.UpdateUserPasswordParams{
					ID:       int64(1),
					Password: "newpassword",
				}
				querier.EXPECT().UpdateUserPassword(gomock.Any(), sqlcParam).
					Return(errors.New("server error"))
			},
		},
		{
			name:          "update user password succeeds given querier update user password returns no error",
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				sqlcParam := sqlc.UpdateUserPasswordParams{
					ID:       int64(1),
					Password: "newpassword",
				}
				querier.EXPECT().UpdateUserPassword(gomock.Any(), sqlcParam).Return(nil)
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
			userRepository := NewUserRepository(mockQuerier)
			actualErr := userRepository.
				UpdateUserPassword(context.Background(), int64(1), "newpassword")

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}
		})
	}
}

func TestUserExists(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "user exists fails given querier user exists returns a server error",
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().UserExists(gomock.Any(), int64(1)).
					Return(false, errors.New("server error"))
			},
		},
		{
			name:          "user exists succeeds given querier user exists returns no error",
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().UserExists(gomock.Any(), int64(1)).
					Return(false, nil)
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
			userRepository := NewUserRepository(mockQuerier)
			_, actualErr := userRepository.
				UserExists(context.Background(), int64(1))

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	tests := []struct {
		name              string
		expectedUser      *domain.User
		expectedError     error
		setupExpectations func(querier *MockQuerier)
	}{
		{
			name:          "get user by username fails given querier get user by username returns sql no rows error",
			expectedUser:  nil,
			expectedError: domain.ErrUserNotFound,
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().GetUserByUsername(gomock.Any(), "username").
					Return(sqlc.User{}, sql.ErrNoRows)
			},
		},
		{
			name:          "get user by username fails given querier get user by username returns a server error",
			expectedUser:  nil,
			expectedError: errors.New("server error"),
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().GetUserByUsername(gomock.Any(), "username").
					Return(sqlc.User{}, errors.New("server error"))
			},
		},
		{
			name: "get user by username returns a domain user given querier get user by username returns no error",
			expectedUser: &domain.User{
				ID:       int64(1),
				Username: "username",
				Password: "password",
			},
			expectedError: nil,
			setupExpectations: func(querier *MockQuerier) {
				querier.EXPECT().GetUserByUsername(gomock.Any(), "username").
					Return(sqlc.User{ID: int64(1), Username: "username", Password: "password"}, nil)
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
			userRepository := NewUserRepository(mockQuerier)
			actualUser, actualErr := userRepository.GetUserByUsername(context.Background(), "username")

			// Assert
			if (actualErr != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, actualErr)
			}

			if !reflect.DeepEqual(test.expectedUser, actualUser) {
				t.Fatalf("expected: \n%v, got: \n%v", test.expectedUser, actualUser)
			}
		})
	}
}
