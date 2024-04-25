package api

import (
	"backend/repo"
)

type Server struct {
	pg repo.Postgres
}

func NewServer(pg repo.Postgres) Server {
	return Server{pg: pg}
}
