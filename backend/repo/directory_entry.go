package repo

import (
	"context"
	"fmt"
	"net/mail"

	"github.com/jackc/pgx/v5"
)

type BaseDirectoryEntry struct {
	FirstName   string `db:"entry_first_name"`
	LastName    string `db:"entry_last_name"`
	Email       string `db:"entry_email"`
	DirectoryID int    `db:"directory_id"`
	UserID      string `db:"user_id"`
	OrderID     int    `db:"order_id"`
}

type DirectoryEntry struct {
	BaseDirectoryEntry
	ID         int  `db:"entry_id"`
	EmailValid bool `db:"entry_email_valid"`
}

type DirectoryEntryNotification struct {
	FirstName   string `json:"entry_first_name"`
	LastName    string `json:"entry_last_name"`
	Email       string `json:"entry_email"`
	DirectoryID int    `json:"directory_id"`
	UserID      string `json:"user_id"`
	OrderID     int    `json:"order_id"`

	ID         int  `json:"entry_id"`
	EmailValid bool `json:"entry_email_valid"`
}

func NewDirectoryEntry(b BaseDirectoryEntry) DirectoryEntry {
	u := DirectoryEntry{BaseDirectoryEntry: b, EmailValid: true}
	u.validate()

	return u
}

// validate checks the input for any issues and sets appropriate errors
func (u *DirectoryEntry) validate() {
	_, err := mail.ParseAddress(u.Email)

	if err != nil {
		u.EmailValid = false
	}

}

// CreateDirectoryEntry inserts the validated input data into the database
func (pg *Postgres) CreateDirectoryEntry(ctx context.Context, entry DirectoryEntry) (DirectoryEntry, error) {
	entry.validate()

	query := `INSERT INTO directory_entries
						(order_id, entry_first_name, entry_last_name, entry_email, entry_email_valid, directory_id, user_id)
						VALUES
						(@orderId, @firstName, @lastName, @email, @emailValid, @directoryID, @userID)
						RETURNING *`

	args := pgx.NamedArgs{
		"orderID":     entry.OrderID,
		"firstName":   entry.FirstName,
		"lastName":    entry.LastName,
		"email":       entry.Email,
		"emailValid":  entry.EmailValid,
		"directoryID": entry.DirectoryID,
		"userID":      entry.UserID,
	}

	rows, err := pg.DB.Query(ctx, query, args)
	if err != nil {
		return DirectoryEntry{}, fmt.Errorf("Failed to insert row: %w", err)
	}
	defer rows.Close()

	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[DirectoryEntry])
}

// GetAllDirectoryEntries queries the database for all DirectoryEntries
func (pg *Postgres) GetAllDirectoryEntries(ctx context.Context) ([]DirectoryEntry, error) {
	query := `SELECT *
						FROM directory_entries`

	rows, err := pg.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("Failed to get directory_entries: %w", err)
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[DirectoryEntry])
}

// DirectoryEntriesCopyFrom takes a slice of directory_entries and puts all in the database using the
// Postgres COPY FROM command
func (pg *Postgres) DirectoryEntriesCopyFrom(ctx context.Context, es []DirectoryEntry) (int64, error) {
	return pg.DB.CopyFrom(
		ctx,
		pgx.Identifier{"directory_entries"},
		[]string{"order_id", "entry_first_name", "entry_last_name", "entry_email", "entry_email_valid", "directory_id", "user_id"},
		pgx.CopyFromSlice(len(es), func(i int) ([]any, error) {
			return []any{es[i].OrderID, es[i].FirstName, es[i].LastName,
				es[i].Email, es[i].EmailValid,
				es[i].DirectoryID, es[i].UserID}, nil
		}),
	)
}
