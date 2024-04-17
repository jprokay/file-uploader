package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"backend/repo"
)

type Hello struct {
	Name string `json:"name"`
}

type ProcessedFile struct {
	FileName string `json:"fileName"`
	Rows     int64  `json:"rows"`
	Error    string `json:"error"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type FileResponse struct {
	Files []ProcessedFile `json:"files"`
}

type UsersResponse struct {
	Users []repo.User `json:"users"`
}

func upload(pg repo.Postgres, c echo.Context) ([]ProcessedFile, error) {
	form, err := c.MultipartForm()
	if err != nil {
		return nil, err
	}

	files := form.File["files"]

	processedFiles := make([]ProcessedFile, 0, len(files))

	for _, file := range files {
		processed := ProcessedFile{FileName: file.Filename}
		src, err := file.Open()
		if err != nil {
			processed.Error = fmt.Sprintf("Failed to open file: %v", err)
			processedFiles = append(processedFiles, processed)
			continue
		}

		defer src.Close()

		asCsv := csv.NewReader(src)
		rows, err := asCsv.ReadAll()
		// rows, err := pg.CopyFromCSV(c.Request().Context(), asCsv)

		if err != nil {
			processed.Error = fmt.Sprintf("Failed to copy rows: %v", err)
		}

		users := make([]repo.User, 0, len(rows))

		for _, row := range rows {
			users = append(users, repo.User{FirstName: row[0], LastName: row[1], Email: row[2]})
		}

		copied, err := pg.UsersCopyFrom(c.Request().Context(), users)
		processed.Rows = copied

		if err != nil {
			processed.Error = fmt.Sprintf("Failed to copy rows: %v", err)
		}

		processedFiles = append(processedFiles, processed)
	}

	return processedFiles, nil
}

func main() {
	e := echo.New()

	pg, err := repo.NewPool(context.Background(), os.Getenv("DATABASE_URL"))

	if err != nil {
		log.Fatalf("Failed to open DB connection: %v", err)
	}

	defer pg.Close()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, Hello{Name: "Julian"})
	})

	e.POST("/upload", func(c echo.Context) error {
		files, err := upload(pg, c)

		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("%v", err)})
		}

		return c.JSON(http.StatusOK, FileResponse{Files: files})
	})

	e.GET("/files", func(c echo.Context) error {
		res, err := pg.GetAllUsers(context.Background())

		if err != nil {
			return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: fmt.Sprintf("%v", err)})
		}

		return c.JSON(http.StatusOK, UsersResponse{Users: res})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
