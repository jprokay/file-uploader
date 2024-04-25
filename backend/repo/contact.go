package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type CreateContact struct {
	FirstName string `db:"contact_first_name"`
	LastName  string `db:"contact_last_name"`
	Email     string `db:"contact_email"`
	OwnerId   string `db:"user_id"`
	SearchCol string `db:"contact_textsearchable_index_col"`
}
type Contact struct {
	CreateContact
	ID string `db:"contact_id"`
}

type ContactRepo struct {
	Postgres
}

func (pg Postgres) NewContactRepo() ContactRepo {
	return ContactRepo{pg}
}

type GetAllContactsParams struct {
	UserId string
	Limit  int
	Offset int
}

// GetAllContacts queries the database for all Contacts
func (pg *ContactRepo) GetAllContacts(ctx context.Context, params GetAllContactsParams) ([]Contact, error) {
	query := `SELECT * FROM contacts
						WHERE user_id = @userId
						ORDER BY contact_id
						LIMIT @limit
						OFFSET @offset
						`

	args := pgx.NamedArgs{
		"userId": params.UserId,
		"offset": params.Offset * params.Limit,
		"limit":  params.Limit,
	}

	rows, err := pg.DB.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("Failed to get contacts: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[Contact])
}

// GetAllContacts queries the database for all Contacts
func (pg *ContactRepo) GetCountOfContacts(ctx context.Context, params GetAllContactsParams) (CountResult, error) {
	query := `SELECT count(*) FROM contacts
						WHERE user_id = @userId
						`

	args := pgx.NamedArgs{
		"userId": params.UserId,
	}

	rows, err := pg.DB.Query(ctx, query, args)
	if err != nil {
		return CountResult{}, fmt.Errorf("Failed to get contacts: %w", err)
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[CountResult])
}

// CreateDirectoryEntry inserts the validated input data into the database
func (pg *ContactRepo) CreateContact(ctx context.Context, contact CreateContact) (Contact, error) {
	query := `INSERT INTO contacts
						(contact_first_name, contact_last_name, contact_email, user_id)
						VALUES
						(@firstName, @lastName, @email, @userID)
						ON CONFLICT (user_id, contact_email) DO UPDATE SET
							contact_first_name = excluded.contact_first_name,
							contact_last_name = excluded.contact_last_name
						RETURNING *`

	args := pgx.NamedArgs{
		"firstName": contact.FirstName,
		"lastName":  contact.LastName,
		"email":     contact.Email,
		"userID":    contact.OwnerId,
	}

	rows, err := pg.DB.Query(ctx, query, args)
	if err != nil {
		return Contact{}, fmt.Errorf("Failed to insert row: %w", err)
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByNameLax[Contact])
}
