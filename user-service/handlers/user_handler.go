package handlers

import (
	"net/http"
	"real-estate-system/user-service/models"
	repository "real-estate-system/user-service/repository/interfaces"
	"strconv"

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
	name := c.FormValue("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "name is required"})
	}

	user := models.User{
		Name: name,
	}

	err := h.Repo.CreateUser(&user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"result": true,
		"user":   user,
	})
}
