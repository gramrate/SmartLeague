package main

import (
	"SmartLeague/internal/adapters/app"
	"SmartLeague/internal/adapters/controller/api/server"
	"log"
)

// @title           SmartLeague API
// @version         1.0
// @description     Backend service for SmartLeague. Cookie-based auth (HttpOnly access/refresh cookies).

// @contact.name    SmartLeague Backend

// @host            localhost:8000
// @schemes         http

// @securityDefinitions.apikey  CookieAuth
// @in                          cookie
// @name                        user_auth_access_token
// @description                 Authentication via HttpOnly cookies:\n- `user_auth_access_token` (short-lived)\n- `user_auth_refresh_token` (long-lived)\n\nProtected endpoints require valid cookies to be sent by the client.

// @tag.name        ping
// @tag.description Healthcheck

// @tag.name        user
// @tag.description Users and auth actions

// @tag.name        token
// @tag.description Access/refresh tokens

// @tag.name        profile
// @tag.description Profiles CRUD

// @tag.name        club
// @tag.description Clubs and membership

// @tag.name        series
// @tag.description Series created by clubs

// @tag.name        game
// @tag.description Games inside series

func main() {
	mainApp, err := app.New()
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}
	server.Setup(mainApp)
	mainApp.Start()
}
