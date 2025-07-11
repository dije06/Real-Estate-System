package handlers

import (
	"net/http"
	"real-estate-system/listing-service/models"
	"real-estate-system/listing-service/repository/interfaces"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type ListingHandler struct {
	Repo interfaces.ListingRepository
}

func NewListingHandler(repo interfaces.ListingRepository) *ListingHandler {
	return &ListingHandler{Repo: repo}
}

func (h *ListingHandler) CreateListing(c echo.Context) error {
	userIDStr := c.FormValue("user_id")
	listingType := c.FormValue("listing_type")
	priceStr := c.FormValue("price")

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user_id")
	}

	price, err := strconv.Atoi(priceStr)
	if err != nil || price <= 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid price")
	}

	if listingType != "rent" && listingType != "sale" {
		return echo.NewHTTPError(http.StatusBadRequest, "listing_type must be 'rent' or 'sale'")
	}

	timestamp := time.Now().UnixMicro()

	listing := models.Listing{
		UserID:      userID,
		Price:       price,
		ListingType: listingType,
		CreatedAt:   timestamp,
		UpdatedAt:   timestamp,
	}

	if err := h.Repo.CreateListing(&listing); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"result":  true,
		"listing": listing,
	})
}

func (h *ListingHandler) GetListings(c echo.Context) error {
	pageNum, _ := strconv.Atoi(c.QueryParam("page_num"))
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	listings, err := h.Repo.GetListings(pageNum, pageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"result":   true,
		"listings": listings,
	})
}
