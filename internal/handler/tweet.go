package handler

import (
	"net/http"
	"twitter/internal/model"

	"github.com/labstack/echo/v4"
)

func (h *HTTPHandler) HandleCreateTweet(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		c.Logger().Error("invalid jwt token")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	payload := new(model.TweetCreatePayload)

	err := c.Bind(&payload)
	if err != nil {
		c.NoContent(http.StatusBadRequest)
		c.Logger().Error(err)
		return err
	}

	payload.AuthorID = userID

	tweet, err := h.store.Tweet().Create(payload)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.JSON(http.StatusCreated, tweet)
	return nil
}

func (h *HTTPHandler) HandleGetUserTweets(c echo.Context) error {
	_, ok := c.Get("user_id").(string)
	if !ok {
		c.Logger().Error("invalid jwt token")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	userID := c.Param("id")

	tweets, err := h.store.Tweet().GetTweetsByUser(userID)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.JSON(http.StatusOK, tweets)
	return nil
}

func (h *HTTPHandler) HandleGetFeed(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		c.Logger().Error("invalid jwt token")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	tweets, err := h.store.Tweet().GetTweetsFromFollowedUsers(userID)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.JSON(http.StatusOK, tweets)
	return nil
}

func (h *HTTPHandler) HandleDeleteTweet(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		c.Logger().Error("invalid jwt token")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	id := c.Param("id")

	err := h.store.Tweet().Delete(id, userID)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.NoContent(http.StatusOK)
	return nil
}

func (h *HTTPHandler) HandleUpdateTweet(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok {
		c.Logger().Error("invalid jwt token")
		c.NoContent(http.StatusBadRequest)
		return nil
	}

	id := c.Param("id")

	payload := new(model.Tweet)

	err := c.Bind(&payload)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	tweet, err := h.store.Tweet().Update(id, userID, payload)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	c.JSON(http.StatusOK, tweet)
	return nil
}
