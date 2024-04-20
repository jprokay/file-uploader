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

func (ds *DirectoryService) ProcessForm(ctx context.Context, f *multipart.Form) ([]ProcessedFile, error) {
	files := f.File["files"]

	processedFiles := make([]ProcessedFile, 0, len(files))

	for _, file := range files {
		processed := ProcessedFile{FileName: file.Filename}
		doc := repo.NewCreateDirectory(
			repo.BaseDirectory{Name: file.Filename, OwnerId: ds.owner.ID})

		dir, err := ds.pg.CreateDirectory(ctx, doc)
		if err != nil {
			processed.Error = fmt.Sprintf("Failed to create directory: %v", err)
			processedFiles = append(processedFiles, processed)
			continue
		}

		updateDir := repo.UpdateDirectoryParams{ID: dir.ID, Name: dir.Name, Status: dir.Status}

		src, err := file.Open()
		if err != nil {
			updateDir.Status = "error"
			processed.Error = fmt.Sprintf("Failed to open file: %v", err)
			processedFiles = append(processedFiles, processed)
			continue
		}

		defer src.Close()

		asCsv := csv.NewReader(src)
		rows, err := asCsv.ReadAll()
		// rows, err := pg.CopyFromCSV(c.Request().Context(), asCsv)

		if err != nil {
			updateDir.Status = "error"
			processed.Error = fmt.Sprintf("Failed to read rows: %v", err)
		}

		users := make([]repo.DirectoryEntry, 0, len(rows))

		doc.Entries = len(rows)

		for i, row := range rows {
			users = append(users, repo.NewDirectoryEntry(repo.BaseDirectoryEntry{DirectoryID: dir.ID, OrderID: i, FirstName: row[0], LastName: row[1], Email: row[2], UserID: ds.owner.ID}))
		}

		copied, err := ds.pg.DirectoryEntriesCopyFrom(ctx, users)
		processed.Rows = copied

		if err != nil {
			updateDir.Status = "error"
			processed.Error = fmt.Sprintf("Failed to copy rows: %v", err)
		} else {
			updateDir.Status = "completed"
		}

		err = ds.pg.UpdateDirectory(ctx, updateDir)

		if err != nil {
			log.Printf("Failed to update status: %v", err)
		}

		processedFiles = append(processedFiles, processed)
	}

	return processedFiles, nil

}
