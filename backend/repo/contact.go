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
}
type Contact struct {
	CreateContact
	ID string `db:"contact_id"`
}

// GetAllContacts queries the database for all Contacts
func (pg *Postgres) GetAllContacts(ctx context.Context) ([]Contact, error) {
	query := `SELECT * FROM contacts`

	rows, err := pg.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get contacts: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Contact])
}

// CreateDirectoryEntry inserts the validated input data into the database
func (pg *Postgres) CreateContact(ctx context.Context, contact CreateContact) (Contact, error) {
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

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[Contact])
}
