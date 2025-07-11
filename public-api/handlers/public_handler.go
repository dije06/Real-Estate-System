package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

var (
	UserServiceURL    = os.Getenv("USER_SERVICE_URL")
	ListingServiceURL = os.Getenv("LISTING_SERVICE_URL")
)

// CreateUser forwards request to user-service
func CreateUser(c echo.Context) error {
	reqBody, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body")
	}

	resp, err := http.Post(UserServiceURL+"/users", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, "User service unavailable")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return c.Blob(resp.StatusCode, "application/json", body)
}

// CreateListing forwards request to listing-service
func CreateListing(c echo.Context) error {
	var data map[string]interface{}
	if err := c.Bind(&data); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON body")
	}

	// Convert JSON to form-urlencoded
	form := make(map[string][]string)
	for k, v := range data {
		form[k] = []string{ToString(v)}
	}
	resp, err := http.PostForm(ListingServiceURL+"/listings", form)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, "Listing service unavailable")
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return c.Blob(resp.StatusCode, "application/json", body)
}

// GetListings fetches from listing-service and enriches with user-service
func GetListings(c echo.Context) error {
	// Forward query params
	query := c.Request().URL.RawQuery
	listingResp, err := http.Get(ListingServiceURL + "/listings?" + query)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadGateway, "Listing service unavailable")
	}
	defer listingResp.Body.Close()

	var listingPayload struct {
		Result   bool `json:"result"`
		Listings []struct {
			ID          int    `json:"id"`
			UserID      int    `json:"user_id"`
			ListingType string `json:"listing_type"`
			Price       int    `json:"price"`
			CreatedAt   int64  `json:"created_at"`
			UpdatedAt   int64  `json:"updated_at"`
			User        any    `json:"user,omitempty"` // Will be filled later
		} `json:"listings"`
	}

	body, _ := io.ReadAll(listingResp.Body)
	if err := json.Unmarshal(body, &listingPayload); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to decode listings")
	}

	// Fetch and embed user data per listing
	for i, listing := range listingPayload.Listings {
		userResp, err := http.Get(UserServiceURL + "/users/" + ToString(listing.UserID))
		if err != nil {
			continue
		}
		defer userResp.Body.Close()

		var userPayload struct {
			Result bool        `json:"result"`
			User   interface{} `json:"user"`
		}
		uBody, _ := io.ReadAll(userResp.Body)
		if err := json.Unmarshal(uBody, &userPayload); err == nil && userPayload.Result {
			listingPayload.Listings[i].User = userPayload.User
		}
	}

	// Return enriched listing result
	return c.JSON(http.StatusOK, map[string]interface{}{
		"result":   true,
		"listings": listingPayload.Listings,
	})
}

// Converts any number/string/float to string
func ToString(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		return strconv.Itoa(int(v))
	default:
		return ""
	}
}
