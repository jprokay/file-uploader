package api

import (
	"backend/repo"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) GetAllContacts(ctx echo.Context, params GetAllContactsParams) error {
	cRepo := s.pg.NewContactRepo()
	parms := repo.GetAllContactsParams{UserId: params.UserId, Limit: params.Limit, Offset: int(params.Offset)}

	dbContacts, err := cRepo.GetAllContacts(ctx.Request().Context(), parms)

	if err != nil {
		return err
	}

	res, err := cRepo.GetCountOfContacts(ctx.Request().Context(), parms)

	if err != nil {
		return err
	}

	contacts := make([]Contact, 0, len(dbContacts))

	for _, val := range dbContacts {
		contacts = append(contacts, Contact{
			ContactEmail:     val.Email,
			ContactFirstName: val.FirstName,
			ContactLastName:  val.LastName,
			ContactId:        val.ID,
			UserId:           val.OwnerId,
		})
	}

	resp := ContactsWithTotal{Total: res.Count, Items: contacts}

	return ctx.JSON(http.StatusOK, resp)
}
