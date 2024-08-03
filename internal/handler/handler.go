package handler

import (
	"twitter/internal/config"
	"twitter/internal/store"
)

type HTTPHandler struct {
	store  store.Store
	config *config.Config
}

func New(s store.Store, c *config.Config) *HTTPHandler {
	return &HTTPHandler{
		store:  s,
		config: c,
	}
}
