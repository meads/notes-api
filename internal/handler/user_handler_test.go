package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meads/notes-api/internal/domain"
	"github.com/meads/notes-api/internal/security"
	"go.uber.org/mock/gomock"
)

func TestRegisterUserHandler_Post(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              AuthResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer, authService *MockAuthServicer)
	}{
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/register/",
			want: AuthResponse{
				SessionID: "uuid",
				UserID:    1,
			},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer, authService *MockAuthServicer) {
				requestUsername, requestPassword := "newuser", "message"
				result := &domain.AuthSessionResult{
					RefreshToken:          "mockrefreshtoken",
					RefreshTokenExpiresAt: defaultDate,
					AccessToken:           "mockaccesstoken",
					AccessTokenExpiresAt:  defaultDate,
					SessionID:             "uuid",
					UserID:                1,
				}
				authService.EXPECT().Register(gomock.Any(), requestUsername, requestPassword).
					Return(result, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 500 given register returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			want:         AuthResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer, authService *MockAuthServicer) {
				requestUsername, requestPassword := "newuser", "message"
				authService.EXPECT().Register(gomock.Any(), requestUsername, requestPassword).
					Return(&domain.AuthSessionResult{}, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"\",\"password\":\"\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 400 given invalid params are supplied",
			responseCode: http.StatusBadRequest,
			route:        "/register/",
			want:         AuthResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer, authService *MockAuthServicer) {
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockUserService := NewMockUserServicer(ctrl)
			userHandler := NewUserHandler(mockUserService)
			mockTokener := security.NewMockTokener(ctrl)

			mockAuthService := NewMockAuthServicer(ctrl)
			authHandler := NewAuthHandler(mockAuthService, "localhost")

			router := SetupRouter(authHandler, userHandler, nil, mockTokener)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodPost, test.route, test.body)
			test.setupExpectations(request, mockTokener, mockUserService, mockAuthService)
			router.ServeHTTP(responseRecorder, request)

			// Assert

			// validate the response code
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			// validate the content type
			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			var got AuthResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&got)
			if err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			// Compare structural values
			if got != test.want {
				t.Errorf("Response mismatch!\n Want: %+v\n Got:  %+v", test.want, got)
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              DeleteUserResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer)
	}{
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "delete handler responds with Status Code 200 given valid request",
			responseCode: http.StatusOK,
			route:        "/users/1/",
			want:         DeleteUserResponse{Message: "Resource successfully deleted"},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer) {
				passClaimsMiddleware(r, tokener)
				userService.EXPECT().DeleteUser(gomock.Any(), int64(1)).Return(nil)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "delete handler responds with Status Code 400 given param id not supplied",
			responseCode: http.StatusBadRequest,
			route:        "/users//",
			want:         DeleteUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer) {
				passClaimsMiddleware(r, tokener)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "delete handler responds with Status Code 400 given param id is not a valid integer",
			responseCode: http.StatusBadRequest,
			route:        "/users/invalid/",
			want:         DeleteUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer) {
				passClaimsMiddleware(r, tokener)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "delete handler responds with Status Code 500 given there is a server error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/",
			want:         DeleteUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer) {
				passClaimsMiddleware(r, tokener)
				userService.EXPECT().DeleteUser(gomock.Any(), int64(1)).Return(errors.New("oops"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockUserService := NewMockUserServicer(ctrl)
			userHandler := NewUserHandler(mockUserService)
			mockTokener := security.NewMockTokener(ctrl)

			router := SetupRouter(nil, userHandler, nil, mockTokener)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodDelete, test.route, test.body)
			test.setupExpectations(request, mockTokener, mockUserService)
			router.ServeHTTP(responseRecorder, request)

			// Assert

			// validate the response code
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			// validate the content type
			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			var got DeleteUserResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&got)
			if err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			// Compare structural values
			if got != test.want {
				t.Errorf("Response mismatch!\n Want: %+v\n Got:  %+v", test.want, got)
			}
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              ListUsersResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer)
	}{
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "list handler responds with Status Code 400 given limit param is invalid int",
			responseCode: http.StatusBadRequest,
			route:        "/users/?limit=invalid",
			want:         ListUsersResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer) {
				passClaimsMiddleware(r, tokener)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "list handler responds with Status Code 400 given offset param is invalid int",
			responseCode: http.StatusBadRequest,
			route:        "/users/?offset=invalid",
			want:         ListUsersResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer) {
				passClaimsMiddleware(r, tokener)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "list handler responds with Status Code 500 given there is a server error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/",
			want:         ListUsersResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer) {
				passClaimsMiddleware(r, tokener)
				params := domain.ListUsersParams{Limit: 50, Offset: 0}
				userService.EXPECT().ListUsers(gomock.Any(), params).Return([]domain.User{}, errors.New("oops."))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "list handler responds with Status Code 200 given a valid request",
			responseCode: http.StatusOK,
			route:        "/users/",
			want:         ListUsersResponse{Users: []UserResponse{{ID: 69, Username: "foo"}}},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer) {
				passClaimsMiddleware(r, tokener)
				params := domain.ListUsersParams{Limit: 50, Offset: 0}
				userService.EXPECT().ListUsers(gomock.Any(), params).Return([]domain.User{
					{ID: 69, Username: "foo"},
				}, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockUserService := NewMockUserServicer(ctrl)
			userHandler := NewUserHandler(mockUserService)
			mockTokener := security.NewMockTokener(ctrl)

			router := SetupRouter(nil, userHandler, nil, mockTokener)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodGet, test.route, test.body)
			test.setupExpectations(request, mockTokener, mockUserService)
			router.ServeHTTP(responseRecorder, request)

			// Assert

			// validate the response code
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			// validate the content type
			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			var got ListUsersResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&got)
			if err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			// Compare structural values

			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Response mismatch!\n Want: %+v\n Got:  %+v", test.want, got)
			}

		})
	}
}

func TestUserHandler_Patch(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              PatchUserResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer)
	}{
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"new\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/users/",
			want:         PatchUserResponse{Message: "Resource successfully patched"},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer) {
				passClaimsMiddleware(r, tokener)
				reqID, reqUsername, reqCurrentPassword, reqNewPassword := int64(1), "user", "current", "new"
				userService.EXPECT().ChangePassword(gomock.Any(), domain.UpdateUserPasswordParams{
					ID:              reqID,
					Username:        reqUsername,
					CurrentPassword: reqCurrentPassword,
					NewPassword:     reqNewPassword,
				})
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":69,\"username\":\"user\",\"wrong\":\"newpass\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update handler responds with Status Code 400 when invalid data supplied",
			responseCode: http.StatusBadRequest,
			route:        "/users/",
			want:         PatchUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer) {
				passClaimsMiddleware(r, tokener)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"new\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update handler responds with Status Code 500 when server error on update",
			responseCode: http.StatusInternalServerError,
			route:        "/users/",
			want:         PatchUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, userService *MockUserServicer) {
				passClaimsMiddleware(r, tokener)
				reqID, reqUsername, reqCurrentPassword, reqNewPassword := int64(1), "user", "current", "new"
				userService.EXPECT().ChangePassword(gomock.Any(), domain.UpdateUserPasswordParams{
					ID:              reqID,
					Username:        reqUsername,
					CurrentPassword: reqCurrentPassword,
					NewPassword:     reqNewPassword,
				}).Return(errors.New("server error"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockUserService := NewMockUserServicer(ctrl)
			userHandler := NewUserHandler(mockUserService)
			mockTokener := security.NewMockTokener(ctrl)

			router := SetupRouter(nil, userHandler, nil, mockTokener)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodPatch, test.route, test.body)
			test.setupExpectations(request, mockTokener, mockUserService)
			router.ServeHTTP(responseRecorder, request)

			// Assert

			// validate the response code
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			// validate the content type
			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			var got PatchUserResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&got)
			if err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			// Compare structural values
			if got != test.want {
				t.Errorf("Response mismatch!\n Want: %+v\n Got:  %+v", test.want, got)
			}
		})
	}
}
