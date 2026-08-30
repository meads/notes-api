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

func TestCreateNote_Post(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              NoteResponse
		setupExpectations func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener)
	}{
		{
			body:         bytes.NewBufferString("{\"title\":\"\",\"content\":\"\",\"userId\":\"\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "create note handler responds with status code 400 given fields are blank",
			responseCode: http.StatusBadRequest,
			route:        "/users/1/notes/",
			want:         NoteResponse{},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
			},
		},
		{
			body:         bytes.NewBufferString("{\"title\":\"title\",\"content\":\"content\",\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "create note handler responds with status code 500 given note service returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/",
			want:         NoteResponse{},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
				noteService.EXPECT().CreateNote(gomock.Any(), int64(1), "title", "content").
					Return(nil, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"title\":\"title\",\"content\":\"content\",\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "create note handler responds with status code 200 given note service succeeds",
			responseCode: http.StatusOK,
			route:        "/users/1/notes/",
			want: NoteResponse{
				ID:      int64(1),
				Title:   "title",
				Content: "content",
				UserID:  int64(1),
			},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
				noteService.EXPECT().CreateNote(gomock.Any(), int64(1), "title", "content").
					Return(&domain.Note{
						ID:        int64(1),
						Title:     "title",
						Content:   "content",
						UserID:    int64(1),
						CreatedAt: defaultDate,
					}, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockNoteService := NewMockNoteServicer(ctrl)
			noteHandler := NewNoteHandler(mockNoteService)
			mockTokener := security.NewMockTokener(ctrl)

			router := SetupRouter(nil, nil, noteHandler, mockTokener)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodPost, test.route, test.body)
			test.setupExpectations(request, mockNoteService, mockTokener)
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
			var got NoteResponse
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

func TestUpdateNote_Put(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              NoteResponse
		setupExpectations func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener)
	}{
		{
			body:         bytes.NewBufferString("{\"id\":0,\"title\":\"\",\"content\":\"\",\"userId\":0}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update note handler responds with status code 400 given fields are blank",
			responseCode: http.StatusBadRequest,
			route:        "/users/1/notes/",
			want:         NoteResponse{},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"title\":\"title\",\"content\":\"content\",\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update note handler responds with status code 500 given note service returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/",
			want:         NoteResponse{},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
				noteService.EXPECT().
					UpdateNote(gomock.Any(), int64(1), int64(1), "title", "content").
					Return(nil, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"title\":\"title\",\"content\":\"content\",\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update note handler responds with status code 200 given valid data supplied",
			responseCode: http.StatusOK,
			route:        "/users/1/notes/",
			want: NoteResponse{
				ID:      int64(1),
				Title:   "title",
				Content: "content",
				UserID:  int64(1),
			},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
				noteService.EXPECT().
					UpdateNote(gomock.Any(), int64(1), int64(1), "title", "content").
					Return(&domain.Note{
						ID:      int64(1),
						Title:   "title",
						Content: "content",
						UserID:  int64(1),
					}, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockNoteService := NewMockNoteServicer(ctrl)
			noteHandler := NewNoteHandler(mockNoteService)
			mockTokener := security.NewMockTokener(ctrl)

			router := SetupRouter(nil, nil, noteHandler, mockTokener)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodPut, test.route, test.body)
			test.setupExpectations(request, mockNoteService, mockTokener)
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
			var got NoteResponse
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

func TestListNotes_Get(t *testing.T) {
	tests := []struct {
		contentType       string
		name              string
		responseCode      int
		route             string
		want              ListNotesResponse
		setupExpectations func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener)
	}{
		{
			contentType:  "application/json; charset=utf-8",
			name:         "list notes handler responds with status code 400 given userid param is invalid int",
			responseCode: http.StatusBadRequest,
			route:        "/users/invalid/notes/",
			want:         ListNotesResponse{},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
			},
		},
		{
			contentType:  "application/json; charset=utf-8",
			name:         "list notes handler responds with status code 500 given note service returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/",
			want:         ListNotesResponse{},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
				noteService.EXPECT().
					ListNotes(gomock.Any(), int64(1)).
					Return([]domain.Note{}, errors.New("server error"))
			},
		},
		{
			contentType:  "application/json; charset=utf-8",
			name:         "list notes handler responds with status code 200 given valid data supplied",
			responseCode: http.StatusOK,
			route:        "/users/1/notes/",
			want: ListNotesResponse{
				Notes: []NoteResponse{{ID: int64(1), Title: "title", Content: "content", UserID: int64(1)}},
			},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
				noteService.EXPECT().
					ListNotes(gomock.Any(), int64(1)).
					Return(
						[]domain.Note{{ID: int64(1), Title: "title", Content: "content", UserID: int64(1)}},
						nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockNoteService := NewMockNoteServicer(ctrl)
			noteHandler := NewNoteHandler(mockNoteService)
			mockTokener := security.NewMockTokener(ctrl)

			router := SetupRouter(nil, nil, noteHandler, mockTokener)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodGet, test.route, nil)
			test.setupExpectations(request, mockNoteService, mockTokener)
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
			var got ListNotesResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&got)
			if err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			// Compare structural values
			if !reflect.DeepEqual(test.want, got) {
				t.Errorf("Response mismatch!\n Want: %+v\n Got:  %+v", test.want, got)
			}
		})
	}
}

func TestDeleteNote_Delete(t *testing.T) {
	tests := []struct {
		contentType       string
		name              string
		responseCode      int
		route             string
		want              DeleteNoteResponse
		setupExpectations func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener)
	}{
		{
			contentType:  "application/json; charset=utf-8",
			name:         "delete note handler responds with status code 400 given userid is an invalid int",
			responseCode: http.StatusBadRequest,
			route:        "/users/id/notes/1",
			want:         DeleteNoteResponse{},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
			},
		},
		{
			contentType:  "application/json; charset=utf-8",
			name:         "delete note handler responds with status code 400 given noteid is an invalid int",
			responseCode: http.StatusBadRequest,
			route:        "/users/1/notes/id",
			want:         DeleteNoteResponse{},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
			},
		},
		{
			contentType:  "application/json; charset=utf-8",
			name:         "delete note handler responds with status code 500 given note service returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/1",
			want:         DeleteNoteResponse{},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
				noteService.EXPECT().DeleteNote(gomock.Any(), int64(1), int64(1)).
					Return(errors.New("server error"))
			},
		},
		{
			contentType:  "application/json; charset=utf-8",
			name:         "delete note handler responds with status code 200 given note service succeeds",
			responseCode: http.StatusOK,
			route:        "/users/1/notes/1",
			want:         DeleteNoteResponse{Message: "Resource successfully deleted"},
			setupExpectations: func(r *http.Request, noteService *MockNoteServicer, tokener *security.MockTokener) {
				passClaimsMiddleware(r, tokener)
				noteService.EXPECT().DeleteNote(gomock.Any(), int64(1), int64(1)).
					Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockNoteService := NewMockNoteServicer(ctrl)
			noteHandler := NewNoteHandler(mockNoteService)
			mockTokener := security.NewMockTokener(ctrl)

			router := SetupRouter(nil, nil, noteHandler, mockTokener)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodDelete, test.route, nil)
			test.setupExpectations(request, mockNoteService, mockTokener)
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
			var got DeleteNoteResponse
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
