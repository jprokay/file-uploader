package api

import (
	"log"
	"net/http"
	"strconv"

	"backend/repo"
	"backend/service"

	"github.com/labstack/echo/v4"
)

func (s Server) CreateNewDirectory(ctx echo.Context, params CreateNewDirectoryParams) error {
	form, err := ctx.MultipartForm()

	if err != nil {
		return err
	}

	id := params.UserId
	ds := service.NewDirectoryService(s.pg, repo.User{ID: id})

	log.Println(form.Value["excludeFirstRow"])
	files := form.File["filename"]
	exclude, err := strconv.ParseBool(form.Value["excludeFirstRow"][0])

	if err != nil {
		exclude = false
	}

	output, errors := ds.ProcessForm(ctx.Request().Context(), files, service.ProcessFormOpts{ExcludeFirstRow: exclude})

	if len(errors) > 0 {
		return ctx.JSON(http.StatusInternalServerError, output)
	}
	return ctx.JSON(http.StatusOK, output)
}

func (s Server) GetAllDirectories(ctx echo.Context, params GetAllDirectoriesParams) error {
	dirRepo := s.pg.NewDirectoryRepo()

	dirs, err := dirRepo.GetAllDirectoriesForUser(ctx.Request().Context(),
		repo.GetAllDirectoriesParams{Direction: "desc", UserId: params.UserId})

	if err != nil {
		return err
	}

	res, err := dirRepo.GetCountOfDirectoriesForUser(ctx.Request().Context(), params.UserId)

	if err != nil {
		return err
	}

	ds := make([]Directory, 0, len(dirs))

	for _, val := range dirs {
		status := (DirectoryDirectoryStatus)(val.Status)
		createdAt := val.CreatedAt.String()
		ds = append(ds, Directory{
			DirectoryId:        val.ID,
			DirectoryName:      val.Name,
			DirectoryStatus:    status,
			DirectoryCreatedAt: createdAt,
		})
	}

	resp := DirectoriesWithTotal{Total: res.Count, Items: ds}
	return ctx.JSON(http.StatusOK, resp)
}
