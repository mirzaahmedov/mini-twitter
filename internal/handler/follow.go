package handler

import (
	"net/http"
	"twitter/internal/model"

	"github.com/labstack/echo/v4"
)

func (h *HTTPHandler) HandleFollow(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		c.Logger().Error("invalid jwt token")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	payload := new(struct {
		UserID string `json:"user_id"`
	})

	err := c.Bind(&payload)
	if err != nil {
		c.NoContent(http.StatusBadRequest)
		c.Logger().Error(err)
		return err
	}

	follow, err := h.store.Follow().Create(&model.Follow{
		FollowerID:  userID,
		FollowingID: payload.UserID,
	})
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.JSON(http.StatusCreated, follow)
	return nil
}

func (h *HTTPHandler) HandleUnfollow(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		c.Logger().Error("invalid jwt token")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	payload := new(struct {
		UserID string `json:"user_id"`
	})

	err := c.Bind(&payload)
	if err != nil {
		c.NoContent(http.StatusBadRequest)
		c.Logger().Error(err)
		return err
	}

	err = h.store.Follow().Delete(userID, payload.UserID)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.NoContent(http.StatusCreated)
	return nil
}

func (h *HTTPHandler) HandleGetFollows(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		c.Logger().Error("invalid jwt token")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	follows, err := h.store.Follow().FindAllWithFollowerID(userID)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.JSON(http.StatusCreated, follows)
	return nil
}
