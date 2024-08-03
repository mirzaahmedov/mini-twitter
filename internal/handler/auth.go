package handler

import (
	"fmt"
	"net/http"
	"twitter/internal/model"
	"twitter/internal/utils"

	"github.com/labstack/echo/v4"
)

func (h *HTTPHandler) HandleRegister(c echo.Context) error {
	payload := new(model.UserCreatePayload)

	err := c.Bind(&payload)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	hash, err := utils.HashPassword(payload.Password)
	if err != nil {
		return err
	}
	payload.Password = hash

	user, err := h.store.User().Create(payload)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	token, err := utils.GenerateToken(h.config.JWT.Secret, user.ID)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	})
	c.JSON(http.StatusCreated, user)
	return nil
}

func (h *HTTPHandler) HandleLogin(c echo.Context) error {
	payload := new(model.UserLoginPayload)

	err := c.Bind(&payload)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	user, err := h.store.User().FindByUsername(payload.Username)
	if err != nil {
		return err
	}

	if !utils.CheckPasswordHash(payload.Password, user.Password) {
		c.Logger().Error("incorrect password")
		return fmt.Errorf("unauthorized user")
	}

	token, err := utils.GenerateToken(h.config.JWT.Secret, user.ID)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
		Path:     "/",
	})

	c.JSON(http.StatusOK, user)
	return nil
}

func (h *HTTPHandler) HandleLogout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:   "token",
		MaxAge: -1,
	})
	c.NoContent(http.StatusOK)
	return nil
}
