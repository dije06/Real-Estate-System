package tests

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"real-estate-system/user-service/handlers"
	"real-estate-system/user-service/models"
	"real-estate-system/user-service/repository/mocks"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	h := handlers.NewUserHandler(mockRepo)

	body := strings.NewReader("name=Alice")
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.On("CreateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Name == "Alice"
	})).Return(nil)

	err := h.CreateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_BindError(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	h := handlers.NewUserHandler(mockRepo)

	body := strings.NewReader("name=")
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := h.CreateUser(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestCreateUser_RepoError(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	h := handlers.NewUserHandler(mockRepo)

	body := strings.NewReader("name=RepoFail")
	req := httptest.NewRequest(http.MethodPost, "/users", body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.On("CreateUser", mock.MatchedBy(func(u *models.User) bool {
		return u.Name == "RepoFail"
	})).Return(errors.New("mock error"))

	err := h.CreateUser(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	mockRepo.AssertExpectations(t)
}

func TestGetUsers_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	h := handlers.NewUserHandler(mockRepo)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users?page_num=1&page_size=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockUsers := []models.User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
	}
	mockRepo.On("GetUsers", 1, 2).Return(mockUsers, nil)

	err := h.GetUsers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string][]models.User
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response["users"], 2)
	assert.Equal(t, "Alice", response["users"][0].Name)
	assert.Equal(t, "Bob", response["users"][1].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetUsers_RepoError(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	h := handlers.NewUserHandler(mockRepo)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users?page_num=1&page_size=2", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mockRepo.On("GetUsers", 1, 2).Return([]models.User(nil), errors.New("repo error"))

	err := h.GetUsers(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	h := handlers.NewUserHandler(mockRepo)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")

	mockUser := &models.User{ID: 1, Name: "Charlie"}
	mockRepo.On("GetUser", 1).Return(mockUser, nil)

	err := h.GetUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response map[string]models.User
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), response["user"].ID)
	assert.Equal(t, "Charlie", response["user"].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetUser_InvalidID(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	h := handlers.NewUserHandler(mockRepo)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("abc")

	err := h.GetUser(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	h := handlers.NewUserHandler(mockRepo)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("99")

	mockRepo.On("GetUser", 99).Return((*models.User)(nil), errors.New("not found"))

	err := h.GetUser(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusNotFound, err.(*echo.HTTPError).Code)
	mockRepo.AssertExpectations(t)
}

func TestGetUsers_InvalidQueryParams(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	h := handlers.NewUserHandler(mockRepo)

	// Invalid params default to 1 and 10
	mockRepo.On("GetUsers", 1, 10).Return([]models.User{}, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users?page_num=abc&page_size=xyz", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetUsers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetUsers_ZeroPageParams(t *testing.T) {
	mockRepo := new(mocks.UserRepositoryMock)
	h := handlers.NewUserHandler(mockRepo)

	// 0 and negative should default to 1 and 10
	mockRepo.On("GetUsers", 1, 10).Return([]models.User{}, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/users?page_num=0&page_size=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := h.GetUsers(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}
