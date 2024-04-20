package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"backend/repo"
	"backend/service"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type FileResponse struct {
	Files []service.ProcessedFile `json:"files"`
}

type EntriesResponse struct {
	Entries []repo.DirectoryEntry `json:"entries"`
}

type DirectoriesResponse struct {
	Directories []repo.Directory `json:"directories"`
}

func upload(pg repo.Postgres, c echo.Context) ([]service.ProcessedFile, error) {
	form, err := c.MultipartForm()

	if err != nil {
		return nil, err
	}

	id := form.Value["ownerId"][0]

	ds := service.NewDirectoryService(pg, repo.User{ID: id})

	return ds.ProcessForm(c.Request().Context(), form)
}

func listen(pg repo.Postgres) {
	conn, err := pg.DB.Acquire(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error acquiring connection:", err)
		os.Exit(1)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "listen new_entry")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error listening to chat channel:", err)
		os.Exit(1)
	}

	for {
		notification, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error waiting for notification:", err)
			os.Exit(1)
		}

		var entry repo.DirectoryEntryNotification
		json.Unmarshal([]byte(notification.Payload), &entry)

		contact, err := pg.CreateContact(context.Background(), repo.CreateContact{FirstName: entry.FirstName, LastName: entry.LastName, Email: entry.Email, OwnerId: entry.UserID})

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating contact: ", err)
		} else {
			fmt.Sprintln(contact)
		}
	}
}

func main() {
	e := echo.New()

	pg, err := repo.NewPool(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}

	go listen(pg)
	defer pg.Close()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/upload", func(c echo.Context) error {
		files, err := upload(pg, c)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("%v", err)})
		}

		return c.JSON(http.StatusOK, FileResponse{Files: files})
	})

	e.GET("/entries", func(c echo.Context) error {
		res, err := pg.GetAllDirectoryEntries(context.Background())

		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("%v", err)})
		}

		return c.JSON(http.StatusOK, EntriesResponse{Entries: res})
	})

	e.GET("/directories", func(c echo.Context) error {
		res, err := pg.GetAllDirectories(context.Background())

		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("%v", err)})
		}

		return c.JSON(http.StatusOK, DirectoriesResponse{Directories: res})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
