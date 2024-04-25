package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
)

type EntryNotificationListener struct {
	pg Postgres
}

func NewEntryNotificationListener(pg Postgres) EntryNotificationListener {
	return EntryNotificationListener{pg: pg}
}

func (l EntryNotificationListener) Listen() {
	go l.listen()
}

func (l EntryNotificationListener) listen() {
	conn, err := l.pg.DB.Acquire(context.Background())
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

		var entry DirectoryEntryNotification
		json.Unmarshal([]byte(notification.Payload), &entry)

		repo := l.pg.NewContactRepo()
		contact, err := repo.CreateContact(context.Background(), CreateContact{FirstName: entry.FirstName, LastName: entry.LastName, Email: entry.Email, OwnerId: entry.UserID})

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating contact: ", err)
		} else {
			fmt.Sprintln(contact)
		}
	}
}
