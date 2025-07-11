package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"real-estate-system/public-api/handlers"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Success(t *testing.T) {
	e := echo.New()
	payload := `{"name": "John Doe"}`

	mockUserService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/users", r.URL.Path)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"result":true,"user":{"id":1,"name":"John Doe"}}`))
	}))
	defer mockUserService.Close()

	handlers.UserServiceURL = mockUserService.URL

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.CreateUser(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "John Doe")
}

func TestCreateUser_InvalidBody(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/users", &badReader{})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.CreateUser(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestCreateUser_ServiceUnavailable(t *testing.T) {
	e := echo.New()
	payload := `{"name": "John Doe"}`

	handlers.UserServiceURL = "http://localhost:9999" // unreachable

	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.CreateUser(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadGateway, err.(*echo.HTTPError).Code)
}

func TestCreateListing_Success(t *testing.T) {
	e := echo.New()
	payload := `{"user_id": 1, "listing_type": "rent", "price": 100000}`

	mockListingService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/listings", r.URL.Path)
		err := r.ParseForm()
		assert.NoError(t, err)
		assert.Equal(t, "1", r.FormValue("user_id"))
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"result":true,"listing":{"id":1}}`))
	}))
	defer mockListingService.Close()

	handlers.ListingServiceURL = mockListingService.URL

	req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.CreateListing(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), "listing")
}

func TestCreateListing_ServiceUnavailable(t *testing.T) {
	e := echo.New()
	payload := `{"user_id": 1, "listing_type": "rent", "price": 100000}`
	handlers.ListingServiceURL = "http://localhost:9999"

	req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader(payload))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.CreateListing(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadGateway, err.(*echo.HTTPError).Code)
}

func TestGetListings_Success(t *testing.T) {
	e := echo.New()

	listingPayload := map[string]interface{}{
		"result": true,
		"listings": []map[string]interface{}{
			{"id": 1, "user_id": 2, "listing_type": "rent", "price": 100000, "created_at": 12345678, "updated_at": 12345678},
		},
	}
	listingBody, _ := json.Marshal(listingPayload)

	userPayload := map[string]interface{}{
		"result": true,
		"user":   map[string]interface{}{"id": 2, "name": "Alice"},
	}
	userBody, _ := json.Marshal(userPayload)

	mockListingService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(listingBody)
	}))
	mockUserService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(userBody)
	}))
	defer mockListingService.Close()
	defer mockUserService.Close()

	handlers.ListingServiceURL = mockListingService.URL
	handlers.UserServiceURL = mockUserService.URL

	req := httptest.NewRequest(http.MethodGet, "/listings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetListings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Alice")
}

func TestGetListings_UserEnrichFailSafe(t *testing.T) {
	e := echo.New()

	listingPayload := map[string]interface{}{
		"result": true,
		"listings": []map[string]interface{}{
			{"id": 1, "user_id": 2, "listing_type": "rent", "price": 100000, "created_at": 12345678, "updated_at": 12345678},
		},
	}
	listingBody, _ := json.Marshal(listingPayload)

	mockListingService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(listingBody)
	}))
	mockUserService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockListingService.Close()
	defer mockUserService.Close()

	handlers.ListingServiceURL = mockListingService.URL
	handlers.UserServiceURL = mockUserService.URL

	req := httptest.NewRequest(http.MethodGet, "/listings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetListings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "listings")
}

func TestGetListings_BadListingResponse(t *testing.T) {
	e := echo.New()

	mockListingService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer mockListingService.Close()

	handlers.ListingServiceURL = mockListingService.URL
	handlers.UserServiceURL = "http://dummy"

	req := httptest.NewRequest(http.MethodGet, "/listings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetListings(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusInternalServerError, err.(*echo.HTTPError).Code)
}

func TestGetListings_ListingServiceUnavailable(t *testing.T) {
	e := echo.New()

	handlers.ListingServiceURL = "http://localhost:9999"
	handlers.UserServiceURL = "http://dummy"

	req := httptest.NewRequest(http.MethodGet, "/listings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetListings(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadGateway, err.(*echo.HTTPError).Code)
}

func TestGetListings_InvalidUserJSON(t *testing.T) {
	e := echo.New()

	listingPayload := map[string]interface{}{
		"result": true,
		"listings": []map[string]interface{}{
			{"id": 1, "user_id": 2, "listing_type": "rent", "price": 100000, "created_at": 12345678, "updated_at": 12345678},
		},
	}
	listingBody, _ := json.Marshal(listingPayload)

	mockListingService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(listingBody)
	}))
	mockUserService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("invalid-json"))
	}))
	defer mockListingService.Close()
	defer mockUserService.Close()

	handlers.ListingServiceURL = mockListingService.URL
	handlers.UserServiceURL = mockUserService.URL

	req := httptest.NewRequest(http.MethodGet, "/listings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetListings(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

type badReader struct{}

func (b *badReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func TestCreateListing_InvalidJSON(t *testing.T) {
	e := echo.New()
	badJSON := `{"user_id": 1, "listing_type": "rent", "price": ` // broken JSON

	req := httptest.NewRequest(http.MethodPost, "/listings", strings.NewReader(badJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.CreateListing(c)
	assert.Error(t, err)
	assert.Equal(t, http.StatusBadRequest, err.(*echo.HTTPError).Code)
}

func TestGetListings_UserGetFails_Completely(t *testing.T) {
	e := echo.New()

	listingPayload := map[string]interface{}{
		"result": true,
		"listings": []map[string]interface{}{
			{"id": 1, "user_id": 123, "listing_type": "rent", "price": 100000, "created_at": 12345678, "updated_at": 12345678},
		},
	}
	listingBody, _ := json.Marshal(listingPayload)

	// Listing service returns valid data
	mockListingService := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(listingBody)
	}))
	defer mockListingService.Close()

	handlers.ListingServiceURL = mockListingService.URL
	handlers.UserServiceURL = "http://localhost:9999" // unreachable

	req := httptest.NewRequest(http.MethodGet, "/listings", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handlers.GetListings(c)
	assert.NoError(t, err) // still handled safely
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "listings")
}

func TestToString_UnsupportedType(t *testing.T) {
	result := handlers.ToString(true) // bool is not handled in switch
	assert.Equal(t, "", result)
}
