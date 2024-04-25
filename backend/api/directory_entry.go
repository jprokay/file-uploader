package api

import (
	"backend/repo"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) GetEntriesForDirectory(ctx echo.Context, id int, params GetEntriesForDirectoryParams) error {
	eRepo := s.pg.NewEntryRepo()

	parms := repo.GetAllEntriesParams{DirectoryId: id, Offset: params.Offset, Limit: params.Limit, UserId: params.UserId}

	entries, err := eRepo.GetAllDirectoryEntries(ctx.Request().Context(), parms)

	if err != nil {
		return err
	}

	res, err := eRepo.GetCountOfDirectoryEntries(ctx.Request().Context(), parms)

	if err != nil {
		return err
	}

	des := make([]DirectoryEntry, 0, len(entries))

	for _, val := range entries {
		des = append(des, DirectoryEntry{
			DirectoryId:     id,
			EntryEmail:      val.Email,
			EntryEmailValid: val.EmailValid,
			EntryFirstName:  val.FirstName,
			EntryLastName:   val.LastName,
			OrderId:         val.OrderID,
			EntryId:         val.ID,
			UserId:          params.UserId,
		})
	}

	resp := EntriesWithTotal{Total: res.Count, Items: des}
	return ctx.JSON(http.StatusOK, resp)
}
