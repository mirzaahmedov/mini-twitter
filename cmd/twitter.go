package main

import (
	"log"

	"twitter/internal/config"
	"twitter/internal/handler"
	"twitter/internal/middleware"
	"twitter/internal/store/postgres"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Static("/", "/web")

	c, err := config.Load()
	if err != nil {
		log.Fatal("failed to start the application: ", err.Error())
	}

	s := postgres.NewPostgresStore(c)
	err = s.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	h := handler.New(s, c)

	publicRoutes := e.Group("/api/v1/auth")

	publicRoutes.POST("/register", h.HandleRegister)
	publicRoutes.POST("/login", h.HandleLogin)
	publicRoutes.GET("/logout", h.HandleLogout)

	protectedRoutes := e.Group("/api/v1")
	protectedRoutes.Use(middleware.Authentication(c.JWT.Secret))

	protectedRoutes.POST("/tweet", h.HandleCreateTweet)
	protectedRoutes.GET("/tweet/feed", h.HandleGetFeed)
	protectedRoutes.GET("/users/:id/tweet", h.HandleGetUserTweets)
	protectedRoutes.DELETE("/tweet/:id", h.HandleDeleteTweet)
	protectedRoutes.PATCH("/tweet/:id", h.HandleUpdateTweet)

	protectedRoutes.POST("/follow", h.HandleFollow)
	protectedRoutes.POST("/unfollow", h.HandleUnfollow)
	protectedRoutes.GET("/follows", h.HandleGetFollows)

	protectedRoutes.GET("/users", h.HandleSearchUsers)
	protectedRoutes.GET("/users/:id", h.HandleGetUserByID)

	log.Fatal(e.Start(c.HTTP.Addr))
}
