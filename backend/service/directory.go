package service

import (
	"backend/repo"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"mime/multipart"
)

type Row struct {
	Id    int    `json:"rowId"`
	Error string `json:"error"`
}

type ProcessedFile struct {
	FileName string `json:"fileName"`
	Rows     int64  `json:"rows"`
	Error    string `json:"error"`
}

type DirectoryService struct {
	pg    repo.Postgres
	owner repo.User
}

func NewDirectoryService(pg repo.Postgres, owner repo.User) DirectoryService {
	return DirectoryService{pg: pg, owner: owner}
}

type errorResponse struct {
	FileName string
	Error    string
}

type ProcessFormOpts struct {
	ExcludeFirstRow bool
}

func (ds *DirectoryService) ProcessForm(ctx context.Context, files []*multipart.FileHeader, opts ProcessFormOpts) ([]repo.Directory, []errorResponse) {

	processedFiles := make([]repo.Directory, 0, len(files))
	errors := make([]errorResponse, 0)

	for _, file := range files {
		doc := repo.NewCreateDirectory(
			repo.BaseDirectory{Name: file.Filename, OwnerId: ds.owner.ID})

		dirRepo := ds.pg.NewDirectoryRepo()
		dir, err := dirRepo.CreateDirectory(ctx, doc)

		if err != nil {
			errors = append(errors, errorResponse{FileName: file.Filename, Error: fmt.Sprintf("Failed to create directory: %v", err)})
			continue
		}

		updateDir := repo.UpdateDirectoryParams{ID: dir.ID, Name: dir.Name, Status: dir.Status}

		src, err := file.Open()

		if err != nil {
			errors = append(errors, errorResponse{FileName: file.Filename, Error: fmt.Sprintf("Failed to open file: %v", err)})

			updateDir.Status = "error"
			_, _ = dirRepo.UpdateDirectory(ctx, updateDir)
			continue
		}

		defer src.Close()

		asCsv := csv.NewReader(src)
		rows, err := asCsv.ReadAll()
		// rows, err := pg.CopyFromCSV(c.Request().Context(), asCsv)

		if err != nil {
			updateDir.Status = "error"

			errors = append(errors, errorResponse{FileName: file.Filename, Error: fmt.Sprintf("Failed to read rows: %v", err)})
		}

		users := make([]repo.DirectoryEntry, 0, len(rows))

		doc.Entries = len(rows)

		for i, row := range rows {
			if opts.ExcludeFirstRow && i == 0 {
				continue
			}
			users = append(users, repo.NewDirectoryEntry(repo.BaseDirectoryEntry{DirectoryID: dir.ID, OrderID: i, FirstName: row[0], LastName: row[1], Email: row[2], UserID: ds.owner.ID}))
		}

		entryRepo := ds.pg.NewEntryRepo()
		_, err = entryRepo.DirectoryEntriesCopyFrom(ctx, users)

		if err != nil {
			updateDir.Status = "error"
			errors = append(errors, errorResponse{FileName: file.Filename, Error: fmt.Sprintf("Failed to copy rows to DB: %v", err)})
		} else {
			updateDir.Status = "completed"
		}

		updated, err := dirRepo.UpdateDirectory(ctx, updateDir)

		if err != nil {
			log.Printf("Failed to update status: %v", err)
		}

		processedFiles = append(processedFiles, updated)
	}

	return processedFiles, errors

}
