package handlers

import (
	"net/http"
	"real-estate-system/user-service/models"
	repository "real-estate-system/user-service/repository/interfaces"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	Repo repository.UserRepository
}

func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	pageNum, _ := strconv.Atoi(c.QueryParam("page_num"))
	if pageNum < 1 {
		pageNum = 1
	}
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}

	users, err := h.Repo.GetUsers(pageNum, pageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"result": true,
		"users":  users,
	})
}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.Repo.GetUser(id)
	if err != nil || user == nil {
		return echo.NewHTTPError(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"result": true,
		"user":   user,
	})
}

func (h *UserHandler) CreateUser(c echo.Context) error {
	var user models.User
	if err := c.Bind(&user); err != nil || strings.TrimSpace(user.Name) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid user input")
	}

	now := time.Now().UnixMicro()
	user.CreatedAt = now
	user.UpdatedAt = now

	if err := h.Repo.CreateUser(&user); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"result": true,
		"user":   user,
	})
}
