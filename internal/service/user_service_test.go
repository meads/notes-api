package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/meads/notes-api/internal/domain"
	"github.com/meads/notes-api/internal/security"
	"go.uber.org/mock/gomock"
)

func TestUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		userID            int64
		setupExpectations func(urepo *MockUserRepository)
	}{
		{
			name:          "user service delete user succeeds given a valid user id supplied for active user",
			expectedError: nil,
			userID:        1,
			setupExpectations: func(urepo *MockUserRepository) {
				urepo.EXPECT().DeleteUser(gomock.Any(), int64(1)).Return(nil)
			},
		},
		{
			name:          "user service delete user fails given an error is returned from user repo",
			expectedError: errors.New("user repo error"),
			userID:        1,
			setupExpectations: func(urepo *MockUserRepository) {
				urepo.EXPECT().DeleteUser(gomock.Any(), int64(1)).
					Return(errors.New("user repo error"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)

			mockUserRepo := NewMockUserRepository(ctrl)
			mockHasher := security.NewMockHasher(ctrl)
			userService := NewUserService(mockUserRepo, mockHasher)

			test.setupExpectations(mockUserRepo)

			// Act
			err := userService.DeleteUser(context.Background(), test.userID)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}
		})
	}
}

func TestUserService_ListUsers(t *testing.T) {
	tests := []struct {
		name              string
		expectedError     error
		params            domain.ListUsersParams
		users             []domain.User
		setupExpectations func(urepo *MockUserRepository)
	}{
		{
			name:          "user service list users succeeds given valid parameters",
			expectedError: nil,
			params:        domain.ListUsersParams{Offset: 0, Limit: 50},
			users: []domain.User{
				{ID: 1, Username: "userone", Password: "password"},
				{ID: 2, Username: "usertwo", Password: "password"},
				{ID: 3, Username: "userthree", Password: "password"},
			},
			setupExpectations: func(urepo *MockUserRepository) {
				params := domain.ListUsersParams{Offset: 0, Limit: 50}
				urepo.EXPECT().ListUsers(gomock.Any(), params).Return([]domain.User{
					{ID: 1, Username: "userone", Password: "password"},
					{ID: 2, Username: "usertwo", Password: "password"},
					{ID: 3, Username: "userthree", Password: "password"},
				}, nil)
			},
		},
		{
			name:          "user service list users fails given user repo returns an error",
			expectedError: errors.New("user repo error"),
			params:        domain.ListUsersParams{Offset: 0, Limit: 50},
			users:         nil,
			setupExpectations: func(urepo *MockUserRepository) {
				params := domain.ListUsersParams{Offset: 0, Limit: 50}
				urepo.EXPECT().ListUsers(gomock.Any(), params).
					Return(nil, errors.New("user repo error"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			mockUserRepo := NewMockUserRepository(ctrl)
			mockHasher := security.NewMockHasher(ctrl)
			userService := NewUserService(mockUserRepo, mockHasher)

			test.setupExpectations(mockUserRepo)

			users, err := userService.ListUsers(
				context.Background(), test.params)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}

			if !reflect.DeepEqual(test.users, users) {
				t.Fatalf("users do not match, expected: \n%v, got: \n%v", test.users, users)
			}
		})
	}
}

func TestUserService_ChangePassword(t *testing.T) {
	tests := []struct {
		name                     string
		expectedError            error
		updateUserPasswordParams domain.UpdateUserPasswordParams
		setupExpectations        func(urepo *MockUserRepository, hasher *security.MockHasher)
	}{
		{
			name:          "user service change password fails given get user fails",
			expectedError: errors.New("user repo get user error"),
			updateUserPasswordParams: domain.UpdateUserPasswordParams{
				ID:              int64(1),
				Username:        "username",
				CurrentPassword: "currenthashed",
				NewPassword:     "secret",
			},
			setupExpectations: func(urepo *MockUserRepository, hasher *security.MockHasher) {
				urepo.EXPECT().GetUser(gomock.Any(), int64(1)).
					Return(nil, errors.New("user repo get user error"))
				hasher.EXPECT().ComparePassword("currenthashed", "currenthashed").Times(0)
			},
		},
		{
			name:          "user service change password fails given compare password fails",
			expectedError: errors.New("hasher compare error"),
			updateUserPasswordParams: domain.UpdateUserPasswordParams{
				ID:              int64(1),
				Username:        "username",
				CurrentPassword: "currenthashed",
				NewPassword:     "secret",
			},
			setupExpectations: func(urepo *MockUserRepository, hasher *security.MockHasher) {
				user := &domain.User{
					ID:       int64(1),
					Username: "username",
					Password: "currenthashed",
				}
				urepo.EXPECT().GetUser(gomock.Any(), int64(1)).Return(user, nil)
				hasher.EXPECT().ComparePassword(user.Password, "currenthashed").
					Return(errors.New("hasher compare error"))
				hasher.EXPECT().HashPassword("secret").Times(0)
			},
		},
		{
			name:          "user service change password fails given hash password fails",
			expectedError: errors.New("hasher hash password error"),
			updateUserPasswordParams: domain.UpdateUserPasswordParams{
				ID:              int64(1),
				Username:        "username",
				CurrentPassword: "currenthashed",
				NewPassword:     "secret",
			},
			setupExpectations: func(urepo *MockUserRepository, hasher *security.MockHasher) {
				user := &domain.User{
					ID:       int64(1),
					Username: "username",
					Password: "currenthashed",
				}
				urepo.EXPECT().GetUser(gomock.Any(), int64(1)).Return(user, nil)
				hasher.EXPECT().ComparePassword(user.Password, "currenthashed").Return(nil)
				hasher.EXPECT().HashPassword("secret").Return("", errors.New("hasher hash password error"))
				urepo.EXPECT().UpdateUserPassword(gomock.Any(), int64(1), "secrethashed").
					Times(0)
			},
		},
		{
			name:          "user service change password fails given update user password fails",
			expectedError: errors.New("server error"),
			updateUserPasswordParams: domain.UpdateUserPasswordParams{
				ID:              int64(1),
				Username:        "username",
				CurrentPassword: "currenthashed",
				NewPassword:     "secret",
			},
			setupExpectations: func(urepo *MockUserRepository, hasher *security.MockHasher) {
				user := &domain.User{
					ID:       int64(1),
					Username: "username",
					Password: "currenthashed",
				}
				urepo.EXPECT().GetUser(gomock.Any(), int64(1)).Return(user, nil)
				hasher.EXPECT().ComparePassword(user.Password, "currenthashed").Return(nil)
				hasher.EXPECT().HashPassword("secret").Return("secrethashed", nil)
				urepo.EXPECT().UpdateUserPassword(gomock.Any(), int64(1), "secrethashed").
					Return(errors.New("server error"))
			},
		},
		{
			name:          "user service change password succeeds given valid data",
			expectedError: nil,
			updateUserPasswordParams: domain.UpdateUserPasswordParams{
				ID:              int64(1),
				Username:        "username",
				CurrentPassword: "currenthashed",
				NewPassword:     "secret",
			},
			setupExpectations: func(urepo *MockUserRepository, hasher *security.MockHasher) {
				user := &domain.User{
					ID:       int64(1),
					Username: "username",
					Password: "currenthashed",
				}
				urepo.EXPECT().GetUser(gomock.Any(), int64(1)).Return(user, nil)
				hasher.EXPECT().ComparePassword(user.Password, "currenthashed").Return(nil)
				hasher.EXPECT().HashPassword("secret").Return("secrethashed", nil)
				urepo.EXPECT().UpdateUserPassword(gomock.Any(), int64(1), "secrethashed").Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)

			mockUserRepo := NewMockUserRepository(ctrl)
			mockHasher := security.NewMockHasher(ctrl)

			userService := NewUserService(mockUserRepo, mockHasher)

			test.setupExpectations(mockUserRepo, mockHasher)

			// Act
			err := userService.ChangePassword(context.Background(), test.updateUserPasswordParams)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}

		})
	}
}
