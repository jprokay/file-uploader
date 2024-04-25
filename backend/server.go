package main

import (
	"backend/api"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"backend/repo"
)

type UserMiddleware struct {
	repo repo.UserRepo
}

func (s *UserMiddleware) GetOrSetCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		_, err := c.Cookie("userId")

		if err == nil {
			return next(c)
		}

		user, err := s.repo.Create(c.Request().Context())

		if err != nil {
			log.Printf("Failed to create a cookie: %v\n", err)
			return err
		}
		cookie := new(http.Cookie)
		cookie.Name = "userId"
		cookie.Value = user.ID
		c.SetCookie(cookie)

		return next(c)
	}
}

func main() {
	log.SetOutput(os.Stdout)

	pg, err := repo.NewPool(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}

	defer pg.Close()

	listener := repo.NewEntryNotificationListener(pg)
	listener.Listen()

	server := api.NewServer(pg)
	e := echo.New()

	um := UserMiddleware{repo: pg.NewUserRepo()}
	// TODO: Update with production domain & use env vars
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{os.Getenv("BACKEND_ALLOW_ORIGIN")},
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(um.GetOrSetCookie)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	api.RegisterHandlers(e, server)

	log.Fatal(e.Start("0.0.0.0:" + os.Getenv("PORT")))
}
