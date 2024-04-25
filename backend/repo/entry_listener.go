package repo

import (
	"context"
	"encoding/json"
	"log"
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
		log.Println("Error acquiring connection:", err)
	}
	defer conn.Release()

	_, err = conn.Exec(context.Background(), "listen new_entry")

	if err != nil {
		log.Println("Error listening to channel:", err)
	}

	for {
		notification, err := conn.Conn().WaitForNotification(context.Background())
		if err != nil {
			log.Println("Error waiting for notification:", err)
		}

		var entry DirectoryEntryNotification
		json.Unmarshal([]byte(notification.Payload), &entry)

		repo := l.pg.NewContactRepo()
		_, err = repo.CreateContact(context.Background(), CreateContact{FirstName: entry.FirstName, LastName: entry.LastName, Email: entry.Email, OwnerId: entry.UserID})

		if err != nil {
			log.Println("Error creating contact: ", err)
		}
	}
}
