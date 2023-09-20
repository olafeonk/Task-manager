package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

const (
	authorizationHeader   = "Authorization"
	wwwAuthenticateHeader = "WWW-Authenticate"
	userCtx               = "userId"
	authBasic             = "Basic"
)

func (h *Handler) userIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.Header(wwwAuthenticateHeader, authBasic)
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	user, password, ok := c.Request.BasicAuth()
	if !ok {
		c.Header(wwwAuthenticateHeader, authBasic)
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	slog.Info("Auth user",
		"username", user,
	)
	userId, err := h.services.Authorization.GetUserId(user, password)
	if err != nil {
		c.Header(wwwAuthenticateHeader, authBasic)
		newErrorResponse(c, http.StatusUnauthorized, "empty auth handler")
		return
	}
	c.Set(userCtx, userId)
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("user id is of invalid type")
	}

	return idInt, nil
}
