package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"real-estate-system/listing-service/handlers"
	"real-estate-system/listing-service/models"
	"real-estate-system/listing-service/repository/mocks"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateListing_Success(t *testing.T) {
	mockRepo := new(mocks.ListingRepositoryMock)
	handler := handlers.NewListingHandler(mockRepo)

	form := "user_id=1&listing_type=rent&price=200000"
	req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader(form))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	mockRepo.On("CreateListing", mock.AnythingOfType("*models.Listing")).Return(nil)

	err := handler.CreateListing(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetListings_Success(t *testing.T) {
	mockRepo := new(mocks.ListingRepositoryMock)
	handler := handlers.NewListingHandler(mockRepo)

	expected := []models.Listing{
		{ID: 1, ListingType: "rent", Price: 100000},
		{ID: 2, ListingType: "sale", Price: 200000},
	}
	mockRepo.On("GetListings", 1, 10).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/listings?page_num=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := handler.GetListings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateListing_InvalidUserID(t *testing.T) {
	mockRepo := new(mocks.ListingRepositoryMock)
	handler := handlers.NewListingHandler(mockRepo)

	req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader("user_id=abc&listing_type=rent&price=100000"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := handler.CreateListing(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestCreateListing_InvalidPrice(t *testing.T) {
	mockRepo := new(mocks.ListingRepositoryMock)
	handler := handlers.NewListingHandler(mockRepo)

	req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader("user_id=1&listing_type=rent&price=-100"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := handler.CreateListing(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestCreateListing_InvalidListingType(t *testing.T) {
	mockRepo := new(mocks.ListingRepositoryMock)
	handler := handlers.NewListingHandler(mockRepo)

	req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader("user_id=1&listing_type=other&price=100000"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := handler.CreateListing(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestCreateListing_RepoError(t *testing.T) {
	mockRepo := new(mocks.ListingRepositoryMock)
	handler := handlers.NewListingHandler(mockRepo)

	mockRepo.On("CreateListing", mock.Anything).Return(errors.New("db error"))

	req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader("user_id=1&listing_type=rent&price=100000"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := handler.CreateListing(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

func TestGetListings_RepoError(t *testing.T) {
	mockRepo := new(mocks.ListingRepositoryMock)
	handler := handlers.NewListingHandler(mockRepo)

	mockRepo.On("GetListings", 1, 10).Return([]models.Listing{}, errors.New("db error"))

	req := httptest.NewRequest(http.MethodGet, "/listings?page_num=1&page_size=10", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := handler.GetListings(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

func TestGetListings_DefaultPagination(t *testing.T) {
	mockRepo := new(mocks.ListingRepositoryMock)
	handler := handlers.NewListingHandler(mockRepo)

	mockRepo.On("GetListings", 1, 10).Return([]models.Listing{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/listings", nil)
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)

	err := handler.GetListings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
