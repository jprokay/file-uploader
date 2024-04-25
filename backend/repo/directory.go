package repo

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

type BaseDirectory struct {
	Name    string `db:"directory_name"`
	OwnerId string `db:"user_id"`
	Entries int    `db:"directory_entries"`
}

type CreateDirectory struct {
	BaseDirectory
	Status string `db:"directory_status"`
}

type Directory struct {
	CreateDirectory
	ID        int       `db:"directory_id"`
	CreatedAt time.Time `db:"directory_created_at"`
}

type DirectoryRepo struct {
	Postgres
}

func (pg Postgres) NewDirectoryRepo() DirectoryRepo {
	return DirectoryRepo{pg}
}

func NewCreateDirectory(base BaseDirectory) CreateDirectory {
	return CreateDirectory{BaseDirectory: base, Status: "processing"}
}

func (doc *CreateDirectory) Errored() {
	doc.Status = "error"
}

func (doc *CreateDirectory) Completed() {
	doc.Status = "completed"
}

func (pg *DirectoryRepo) CreateDirectory(ctx context.Context, directory CreateDirectory) (Directory, error) {
	query := `INSERT INTO directories (directory_name, user_id)
						VALUES (@name, @ownerId)
						RETURNING *						
`
	args := pgx.NamedArgs{
		"name":    directory.Name,
		"ownerId": directory.OwnerId,
	}

	rows, err := pg.DB.Query(ctx, query, args)

	if err != nil {
		return Directory{}, fmt.Errorf("Failed to get directories: %w", err)
	}

	defer rows.Close()
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[Directory])
}

func (pg *DirectoryRepo) GetAllDirectories(ctx context.Context) ([]Directory, error) {
	query := `SELECT * FROM directories
`

	rows, err := pg.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get directories: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Directory])
}

type GetAllDirectoriesParams struct {
	UserId    string
	Direction string
}

func (pg *DirectoryRepo) GetAllDirectoriesForUser(ctx context.Context, params GetAllDirectoriesParams) ([]Directory, error) {
	query := `SELECT * FROM directories
		WHERE user_id = @ownerId
		ORDER BY directory_created_at desc 
`
	args := pgx.NamedArgs{
		"ownerId": params.UserId,
	}

	rows, err := pg.DB.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("Failed to get directories: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Directory])
}

type CountResult struct {
	Count int `db:"count"`
}

func (pg *DirectoryRepo) GetCountOfDirectoriesForUser(ctx context.Context, uid string) (CountResult, error) {
	query := `SELECT count(*) as count FROM directories
		WHERE user_id = @ownerId
`
	args := pgx.NamedArgs{
		"ownerId": uid,
	}

	rows, err := pg.DB.Query(ctx, query, args)
	if err != nil {
		return CountResult{}, fmt.Errorf("Failed to get directories: %w", err)
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[CountResult])
}

type UpdateDirectoryParams struct {
	ID     int
	Name   string
	Status string
}

func (pg *DirectoryRepo) UpdateDirectory(ctx context.Context, dir UpdateDirectoryParams) (Directory, error) {
	query := `UPDATE directories
						SET directory_name = @name, directory_status = @status
						WHERE directory_id = @id
						RETURNING *`

	args := pgx.NamedArgs{
		"id":     dir.ID,
		"name":   dir.Name,
		"status": dir.Status,
	}

	rows, err := pg.DB.Query(ctx, query, args)
	if err != nil {
		return Directory{}, fmt.Errorf("Failed to update directory: %w", err)
	}

	defer rows.Close()
	return pgx.CollectOneRow(rows, pgx.RowToStructByName[Directory])
}
