package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *HTTPHandler) HandleSearchUsers(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		c.Logger().Error("invalid jwt token")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	term := c.QueryParam("search")

	users, err := h.store.User().Search(term, userID)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.JSON(http.StatusOK, users)
	return nil
}

func (h *HTTPHandler) HandleGetUserByID(c echo.Context) error {
	_, ok := c.Get("user_id").(string)
	if !ok {
		c.Logger().Error("invalid jwt token")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	userID := c.Param("id")

	user, err := h.store.User().FindByID(userID)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.JSON(http.StatusOK, user)
	return nil
}
