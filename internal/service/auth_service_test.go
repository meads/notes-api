package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	domain "github.com/meads/notes-api/internal/domain"
	"github.com/meads/notes-api/internal/security"
	"go.uber.org/mock/gomock"
)

func TestAuthService_Register(t *testing.T) {
	var testAuthSessionResult *domain.AuthSessionResult
	tests := []struct {
		name              string
		username          string
		password          string
		expectedError     error
		setupExpectations func(
			tokener *security.MockTokener, hasher *security.MockHasher,
			srepo *MockSessionRepository, urepo *MockUserRepository,
		)
	}{
		{
			name:          "auth service register fails given username exists returns an error",
			username:      "username",
			password:      "password",
			expectedError: errors.New("server error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, errors.New("server error"))

				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service register fails given username not available",
			username:      "username",
			password:      "password",
			expectedError: errors.New("please choose another username"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(true, nil)

				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service register fails given error hashing password",
			username:      "username",
			password:      "password",
			expectedError: errors.New("error hashing password"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("", errors.New("error hashing password"))

				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service register fails given error creating user",
			username:      "username",
			password:      "password",
			expectedError: errors.New("user repo create error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("hashpassword", nil)
				urepo.EXPECT().CreateUser(gomock.Any(), "username", "hashpassword").
					Return(nil, errors.New("user repo create error"))

				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service register fails given error generating access token",
			username:      "username",
			password:      "password",
			expectedError: errors.New("tokener error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("hashpassword", nil)
				urepo.EXPECT().CreateUser(gomock.Any(), "username", "hashpassword").
					Return(&domain.User{ID: int64(1), Username: "username", Password: "hashpassword"}, nil)

				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("", nil, errors.New("tokener error"))

				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service register fails given error generating refresh token",
			username:      "username",
			password:      "password",
			expectedError: errors.New("tokener error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("hashpassword", nil)
				urepo.EXPECT().CreateUser(gomock.Any(), "username", "hashpassword").
					Return(&domain.User{ID: int64(1), Username: "username", Password: "hashpassword"}, nil)

				accessClaims, _ := security.NewUserClaims(int64(1), "username", "access", 15*time.Minute)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("accesstoken", accessClaims, nil)

				tokener.EXPECT().
					GenerateToken(int64(1), "username", "refresh", 24*time.Hour).
					Return("", nil, errors.New("tokener error"))

				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service register fails given error creating session",
			username:      "username",
			password:      "password",
			expectedError: errors.New("server error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("hashpassword", nil)
				urepo.EXPECT().CreateUser(gomock.Any(), "username", "hashpassword").
					Return(&domain.User{ID: int64(1), Username: "username", Password: "hashpassword"}, nil)

				accessClaims, _ := security.NewUserClaims(int64(1), "username", "access", 15*time.Minute)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("accesstoken", accessClaims, nil)

				refreshClaims, _ := security.NewUserClaims(int64(1), "username", "access", 24*time.Hour)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "refresh", 24*time.Hour).
					Return("refreshtoken", refreshClaims, nil)

				createSessionParams := domain.CreateSessionParams{
					ID:           refreshClaims.RegisteredClaims.ID,
					UserID:       int64(1),
					RefreshToken: "refreshtoken",
					IsRevoked:    false,
					ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
				}
				srepo.EXPECT().
					CreateSession(gomock.Any(), createSessionParams).
					Return(nil, errors.New("server error"))
				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service register succeeds given valid data supplied",
			username:      "username",
			password:      "password",
			expectedError: nil,
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().UsernameExists(gomock.Any(), "username").Return(false, nil)
				hasher.EXPECT().HashPassword("password").Return("hashpassword", nil)
				urepo.EXPECT().CreateUser(gomock.Any(), "username", "hashpassword").
					Return(&domain.User{ID: int64(1), Username: "username", Password: "hashpassword"}, nil)

				accessClaims, _ := security.NewUserClaims(int64(1), "username", "access", 15*time.Minute)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("accesstoken", accessClaims, nil)

				refreshClaims, _ := security.NewUserClaims(int64(1), "username", "access", 24*time.Hour)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "refresh", 24*time.Hour).
					Return("refreshtoken", refreshClaims, nil)

				createSessionParams := domain.CreateSessionParams{
					ID:           refreshClaims.RegisteredClaims.ID,
					UserID:       int64(1),
					RefreshToken: "refreshtoken",
					IsRevoked:    false,
					ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
				}
				srepo.EXPECT().
					CreateSession(gomock.Any(), createSessionParams).
					Return(&domain.Session{
						ID:           refreshClaims.RegisteredClaims.ID,
						UserID:       int64(1),
						RefreshToken: "refreshtoken",
						IsRevoked:    false,
						ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
					}, nil)
				testAuthSessionResult = &domain.AuthSessionResult{
					RefreshToken:          "refreshtoken",
					RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
					AccessToken:           "accesstoken",
					AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
					SessionID:             refreshClaims.RegisteredClaims.ID,
					UserID:                int64(1),
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)

			tokener := security.NewMockTokener(ctrl)
			hasher := security.NewMockHasher(ctrl)
			sessionRepo := NewMockSessionRepository(ctrl)
			userRepo := NewMockUserRepository(ctrl)
			test.setupExpectations(tokener, hasher, sessionRepo, userRepo)

			authService := NewAuthService(userRepo, sessionRepo, tokener, hasher, 15*time.Minute, 24*time.Hour)

			// Act
			authSessionResult, err := authService.Register(context.Background(), test.username, test.password)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}

			if !reflect.DeepEqual(testAuthSessionResult, authSessionResult) {
				t.Fatalf("register result does not match, expected: \n%v, got: \n%v", testAuthSessionResult, authSessionResult)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	var testAuthSessionResult *domain.AuthSessionResult
	tests := []struct {
		name              string
		username          string
		password          string
		expectedError     error
		setupExpectations func(
			tokener *security.MockTokener, hasher *security.MockHasher,
			srepo *MockSessionRepository, urepo *MockUserRepository,
		)
	}{
		{
			name:          "auth service login fails given error user not found returned from get user by username",
			username:      "username",
			password:      "password",
			expectedError: domain.ErrInvalidCredentials,
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().GetUserByUsername(gomock.Any(), "username").
					Return(nil, domain.ErrUserNotFound)

				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service login fails given error returned from get user by username",
			username:      "username",
			password:      "password",
			expectedError: errors.New("user repo error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				urepo.EXPECT().GetUserByUsername(gomock.Any(), "username").
					Return(nil, errors.New("user repo error"))

				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service login fails given error returned from compare password",
			username:      "username",
			password:      "password",
			expectedError: errors.New("hasher error compare password"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				user := &domain.User{ID: int64(1), Username: "username", Password: "password"}
				urepo.EXPECT().GetUserByUsername(gomock.Any(), "username").Return(user, nil)
				hasher.EXPECT().ComparePassword(user.Password, "password").
					Return(errors.New("hasher error compare password"))

				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service login fails given error returned from generate token for access token",
			username:      "username",
			password:      "password",
			expectedError: errors.New("generate access token error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				user := &domain.User{ID: int64(1), Username: "username", Password: "password"}
				urepo.EXPECT().GetUserByUsername(gomock.Any(), "username").Return(user, nil)
				hasher.EXPECT().ComparePassword(user.Password, "password").Return(nil)

				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("", nil, errors.New("generate access token error"))

				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service login fails given error returned from generate token for refresh token",
			username:      "username",
			password:      "password",
			expectedError: errors.New("generate refresh token error"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				user := &domain.User{ID: int64(1), Username: "username", Password: "password"}
				urepo.EXPECT().GetUserByUsername(gomock.Any(), "username").Return(user, nil)
				hasher.EXPECT().ComparePassword(user.Password, "password").Return(nil)

				accessClaims, _ := security.NewUserClaims(int64(1), "username", "access", 15*time.Minute)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("accesstoken", accessClaims, nil)

				tokener.EXPECT().
					GenerateToken(int64(1), "username", "refresh", 24*time.Hour).
					Return("", nil, errors.New("generate refresh token error"))

				testAuthSessionResult = nil
			},
		},
		{
			name:          "auth service login fails given error creating session",
			username:      "username",
			password:      "password",
			expectedError: errors.New("server error creating session"),
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				user := &domain.User{ID: int64(1), Username: "username", Password: "password"}
				urepo.EXPECT().GetUserByUsername(gomock.Any(), "username").Return(user, nil)
				hasher.EXPECT().ComparePassword(user.Password, "password").Return(nil)

				accessClaims, _ := security.NewUserClaims(int64(1), "username", "access", 15*time.Minute)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("accesstoken", accessClaims, nil)

				refreshClaims, _ := security.NewUserClaims(int64(1), "username", "access", 24*time.Hour)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "refresh", 24*time.Hour).
					Return("refreshtoken", refreshClaims, nil)

				createSessionParams := domain.CreateSessionParams{
					ID:           refreshClaims.RegisteredClaims.ID,
					UserID:       int64(1),
					RefreshToken: "refreshtoken",
					IsRevoked:    false,
					ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
				}
				srepo.EXPECT().
					CreateSession(gomock.Any(), createSessionParams).
					Return(nil, errors.New("server error creating session"))
			},
		},
		{
			name:          "auth service login succeeds given valid data supplied",
			username:      "username",
			password:      "password",
			expectedError: nil,
			setupExpectations: func(
				tokener *security.MockTokener, hasher *security.MockHasher,
				srepo *MockSessionRepository, urepo *MockUserRepository,
			) {
				user := &domain.User{ID: int64(1), Username: "username", Password: "password"}
				urepo.EXPECT().GetUserByUsername(gomock.Any(), "username").Return(user, nil)
				hasher.EXPECT().ComparePassword(user.Password, "password").Return(nil)

				accessClaims, _ := security.NewUserClaims(int64(1), "username", "access", 15*time.Minute)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "access", 15*time.Minute).
					Return("accesstoken", accessClaims, nil)

				refreshClaims, _ := security.NewUserClaims(int64(1), "username", "access", 24*time.Hour)
				tokener.EXPECT().
					GenerateToken(int64(1), "username", "refresh", 24*time.Hour).
					Return("refreshtoken", refreshClaims, nil)

				createSessionParams := domain.CreateSessionParams{
					ID:           refreshClaims.RegisteredClaims.ID,
					UserID:       int64(1),
					RefreshToken: "refreshtoken",
					IsRevoked:    false,
					ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
				}
				srepo.EXPECT().
					CreateSession(gomock.Any(), createSessionParams).
					Return(&domain.Session{
						ID:           refreshClaims.RegisteredClaims.ID,
						UserID:       int64(1),
						RefreshToken: "refreshtoken",
						IsRevoked:    false,
						ExpiresAt:    refreshClaims.RegisteredClaims.ExpiresAt.Time,
					}, nil)

				testAuthSessionResult = &domain.AuthSessionResult{
					RefreshToken:          "refreshtoken",
					RefreshTokenExpiresAt: refreshClaims.ExpiresAt.Time,
					AccessToken:           "accesstoken",
					AccessTokenExpiresAt:  accessClaims.ExpiresAt.Time,
					SessionID:             refreshClaims.RegisteredClaims.ID,
					UserID:                int64(1),
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)

			tokener := security.NewMockTokener(ctrl)
			hasher := security.NewMockHasher(ctrl)
			sessionRepo := NewMockSessionRepository(ctrl)
			userRepo := NewMockUserRepository(ctrl)
			test.setupExpectations(tokener, hasher, sessionRepo, userRepo)

			authService := NewAuthService(userRepo, sessionRepo, tokener, hasher, 15*time.Minute, 24*time.Hour)

			// Act
			authSessionResult, err := authService.Login(context.Background(), test.username, test.password)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}

			if !reflect.DeepEqual(testAuthSessionResult, authSessionResult) {
				t.Fatalf("login result does not match, expected: \n%v, got: \n%v", testAuthSessionResult, authSessionResult)
			}
		})
	}
}

func TestAuthService_Logout(t *testing.T) {
	tests := []struct {
		name              string
		sessionID         string
		expectedError     error
		setupExpectations func(srepo *MockSessionRepository)
	}{
		{
			name:          "auth service logout fails given delete session returns an error",
			sessionID:     "session-id",
			expectedError: errors.New("delete session error"),
			setupExpectations: func(srepo *MockSessionRepository) {
				srepo.EXPECT().DeleteSession(gomock.Any(), "session-id").
					Return(errors.New("delete session error"))
			},
		},
		{
			name:          "auth service logout succeeds given valid session id",
			sessionID:     "valid-session-id",
			expectedError: nil,
			setupExpectations: func(srepo *MockSessionRepository) {
				srepo.EXPECT().DeleteSession(gomock.Any(), "valid-session-id").Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)

			tokener := security.NewMockTokener(ctrl)
			hasher := security.NewMockHasher(ctrl)
			sessionRepo := NewMockSessionRepository(ctrl)
			userRepo := NewMockUserRepository(ctrl)
			test.setupExpectations(sessionRepo)

			authService := NewAuthService(userRepo, sessionRepo, tokener, hasher, 15*time.Minute, 24*time.Hour)

			// Act
			err := authService.Logout(context.Background(), test.sessionID)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}
		})
	}
}

func TestAuthService_RenewAccessToken(t *testing.T) {
	var testRenewAccessTokenResult *domain.RenewAccessTokenResult
	tests := []struct {
		name              string
		refreshToken      string
		expectedError     error
		setupExpectations func(tokener *security.MockTokener, srepo *MockSessionRepository)
	}{
		{
			name:          "auth service renew access token fails given error returned from verify token",
			refreshToken:  "refresh-token",
			expectedError: errors.New("invalid token or claims"),
			setupExpectations: func(tokener *security.MockTokener, srepo *MockSessionRepository) {
				tokener.EXPECT().VerifyToken("refresh-token").
					Return(nil, errors.New("invalid token or claims"))

				testRenewAccessTokenResult = nil
			},
		},
		{
			name:          "auth service renew access token fails given error returned from get session",
			refreshToken:  "refresh-token",
			expectedError: errors.New("error getting session"),
			setupExpectations: func(tokener *security.MockTokener, srepo *MockSessionRepository) {
				refreshClaims, _ := security.NewUserClaims(int64(1), "username", "refresh", 24*time.Hour)
				tokener.EXPECT().VerifyToken("refresh-token").
					Return(refreshClaims, nil)
				srepo.EXPECT().GetSession(gomock.Any(), refreshClaims.RegisteredClaims.ID).
					Return(nil, errors.New("error getting session"))

				testRenewAccessTokenResult = nil
			},
		},
		{
			name:          "auth service renew access token fails given session is revoked",
			refreshToken:  "refresh-token",
			expectedError: errors.New("error generating token"),
			setupExpectations: func(tokener *security.MockTokener, srepo *MockSessionRepository) {
				refreshClaims, _ := security.NewUserClaims(int64(1), "username", "refresh", 24*time.Hour)
				tokener.EXPECT().VerifyToken("refresh-token").
					Return(refreshClaims, nil)
				srepo.EXPECT().GetSession(gomock.Any(), refreshClaims.RegisteredClaims.ID).
					Return(&domain.Session{IsRevoked: true}, nil)
				testRenewAccessTokenResult = nil
			},
		},
		{
			name:          "auth service renew access token fails given error returned from generate token",
			refreshToken:  "refresh-token",
			expectedError: errors.New("error generating token"),
			setupExpectations: func(tokener *security.MockTokener, srepo *MockSessionRepository) {
				refreshClaims, _ := security.NewUserClaims(int64(1), "username", "refresh", 24*time.Hour)
				tokener.EXPECT().VerifyToken("refresh-token").
					Return(refreshClaims, nil)
				srepo.EXPECT().GetSession(gomock.Any(), refreshClaims.RegisteredClaims.ID).
					Return(&domain.Session{IsRevoked: false}, nil)
				tokener.EXPECT().GenerateToken(refreshClaims.UserID, refreshClaims.Username, "access", 15*time.Minute).
					Return("", nil, errors.New("error generating token"))
				testRenewAccessTokenResult = nil
			},
		},
		{
			name:          "auth service renew access token succeeds given valid refresh token and non revoked session",
			refreshToken:  "valid-refresh-token",
			expectedError: nil,
			setupExpectations: func(tokener *security.MockTokener, srepo *MockSessionRepository) {
				refreshClaims, _ := security.NewUserClaims(int64(1), "username", "refresh", 24*time.Hour)
				tokener.EXPECT().VerifyToken("valid-refresh-token").
					Return(refreshClaims, nil)
				srepo.EXPECT().GetSession(gomock.Any(), refreshClaims.RegisteredClaims.ID).
					Return(&domain.Session{IsRevoked: false}, nil)
				accessClaims, _ := security.NewUserClaims(int64(1), "username", "access", 15*time.Minute)
				tokener.EXPECT().GenerateToken(refreshClaims.UserID, refreshClaims.Username, "access", 15*time.Minute).
					Return("accesstoken", accessClaims, nil)
				testRenewAccessTokenResult = &domain.RenewAccessTokenResult{
					AccessToken:          "accesstoken",
					AccessTokenExpiresAt: accessClaims.ExpiresAt.Time,
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)

			tokener := security.NewMockTokener(ctrl)
			hasher := security.NewMockHasher(ctrl)
			sessionRepo := NewMockSessionRepository(ctrl)
			userRepo := NewMockUserRepository(ctrl)
			test.setupExpectations(tokener, sessionRepo)

			authService := NewAuthService(userRepo, sessionRepo, tokener, hasher, 15*time.Minute, 24*time.Hour)

			// Act
			renewAccessTokenResult, err := authService.RenewAccessToken(context.Background(), test.refreshToken)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}

			if !reflect.DeepEqual(testRenewAccessTokenResult, renewAccessTokenResult) {
				t.Fatalf("login result does not match, expected: \n%v, got: \n%v",
					testRenewAccessTokenResult, renewAccessTokenResult)
			}
		})
	}
}

func TestAuthService_RevokeSession(t *testing.T) {
	tests := []struct {
		name              string
		sessionID         string
		expectedError     error
		setupExpectations func(srepo *MockSessionRepository)
	}{
		{
			name:          "auth service revoke session succeeds given a valid session id",
			sessionID:     "session-id",
			expectedError: nil,
			setupExpectations: func(srepo *MockSessionRepository) {
				srepo.EXPECT().RevokeSession(gomock.Any(), "session-id").
					Return(nil)
			},
		},
		{
			name:          "auth service revoke session fails given a session repo returns an error",
			sessionID:     "session-id",
			expectedError: errors.New("session repo error"),
			setupExpectations: func(srepo *MockSessionRepository) {
				srepo.EXPECT().RevokeSession(gomock.Any(), "session-id").
					Return(errors.New("session repo error"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)

			tokener := security.NewMockTokener(ctrl)
			hasher := security.NewMockHasher(ctrl)
			sessionRepo := NewMockSessionRepository(ctrl)
			userRepo := NewMockUserRepository(ctrl)
			test.setupExpectations(sessionRepo)

			authService := NewAuthService(userRepo, sessionRepo, tokener, hasher, 15*time.Minute, 24*time.Hour)

			// Act
			err := authService.RevokeSession(context.Background(), test.sessionID)

			// Assert
			if (err != nil) != (test.expectedError != nil) {
				t.Fatalf("expected error presence: %v, got: %v", test.expectedError != nil, err)
			}
		})
	}
}
