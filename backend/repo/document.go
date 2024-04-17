package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type baseDocument struct {
	Name string `db:"document_name"`
}

type CreateDocument struct {
	baseDocument
	OwnerId      int `db:"anon_id"`
	NumberOfRows int `db:"row_count"`
}

type Document struct {
	CreateDocument
	CreatedAt string `db:"document_created_at"`
}

func (pg *Postgres) CreateDocument(ctx context.Context, document CreateDocument) error {
	query := `INSERT INTO documents (document_name, anon_id) VALUES (@name, @ownerId)`
	args := pgx.NamedArgs{
		"name":    document.Name,
		"ownerId": document.OwnerId,
	}

	_, err := pg.DB.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("Failed to insert document: %w", err)
	}

	return nil
}

func (pg *Postgres) GetAllDocuments(ctx context.Context) ([]Document, error) {
	query := `SELECT document_name, document_created_at, count(users.document_id) as row_count
		FROM documents
		left join users
		on (documents.document_id = users.document_id)
		group by documents.document_id
`

	rows, err := pg.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get documents: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Document])
}
