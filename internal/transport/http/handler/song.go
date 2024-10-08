package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/DSorbon/effective-mobile-task/internal/models"
	"github.com/DSorbon/effective-mobile-task/internal/transport/http/request"
	"github.com/DSorbon/effective-mobile-task/internal/transport/http/response"
	"github.com/DSorbon/effective-mobile-task/internal/transport/http/validation"
	"github.com/DSorbon/effective-mobile-task/pkg/logger"
	"github.com/go-chi/chi/v5"
)

// @Summary  Get All Songs
// @Tags songs
// @Description  get all songs
// @ModuleID List
// @Accept  json
// @Produce  json
// @Param artist query string false "search by artist"
// @Param group query string false "search by group"
// @Param title query string false "search by title"
// @Param releaseDate query string false "search by releaseDate"
// @Param page query int false "paginated by page"
// @Success 200 {object} models.SongPagination
// @Failure 400 {object} response.ResponseMessage
// @Failure 500 {object} response.ResponseMessage
// @Failure default {object} response.ResponseMessage
// @Router /songs [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	var page int64 = 1
	if queryParams.Has("page") {
		pageParam := queryParams.Get("page")

		parsedPage, err := strconv.ParseInt(pageParam, 10, 64)
		if err != nil {
			logger.Errorf("parsing page param: %v", err)
			response.WithBody(w, http.StatusBadRequest, response.ResponseMessage{Message: "invalid param page"})
			return
		}
		page = parsedPage
	}

	releaseDateParam := queryParams.Get("releaseDate")
	var releaseDate *time.Time
	if releaseDateParam != "" {
		date, err := time.Parse(time.DateOnly, releaseDateParam)
		if err != nil {
			logger.Errorf("parsing releaseDate: %v", err)
			response.WithBody(w, http.StatusInternalServerError, response.ResponseMessage{Message: err.Error()})
			return
		}
		releaseDate = &date
	}

	artistParam := queryParams.Get("artist")
	groupParam := queryParams.Get("group")
	titleParam := queryParams.Get("title")

	filter := &models.SongFilter{
		Artist:      artistParam,
		Group:       groupParam,
		Title:       titleParam,
		ReleaseDate: releaseDate,
		Page:        int(page),
	}
	res, err := h.songService.List(r.Context(), filter)
	if err != nil {
		logger.Errorf("get list of song: %v", err)
		response.WithBody(w, http.StatusInternalServerError, response.ResponseMessage{Message: err.Error()})
		return
	}

	response.WithBody(w, http.StatusOK, response.ResponsePagintion{Data: res.Data, Page: res.Page})
}

// @Summary  Create Song
// @Tags songs
// @Description  create song
// @ModuleID Create
// @Accept  json
// @Produce  json
// @Param input body request.SongCreate true "create song"
// @Success 200 {object} response.ResponseMessage
// @Failure 400 {object} response.ResponseMessage
// @Failure 422 {object} response.ResponseValidationErrors
// @Failure 500 {object} response.ResponseMessage
// @Failure default {object} response.ResponseMessage
// @Router /songs [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input request.SongCreate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Errorf("decoding SongCreate: %v", err)
		response.WithBody(w, http.StatusBadRequest, response.ResponseMessage{Message: "invalid request body"})
		return
	}

	// Check validation
	if errBytes := validation.ValidateStruct(input); errBytes != nil {
		logger.Error("validating SongCreate")
		response.ValidateErrors(w, http.StatusUnprocessableEntity, errBytes)
		return
	}

	releaseDate, err := time.Parse(time.DateOnly, input.ReleaseDate)
	if err != nil {
		logger.Errorf("parsing releaseDate: %v", err)
		response.WithBody(w, http.StatusInternalServerError, response.ResponseMessage{Message: err.Error()})
		return
	}

	model := &models.SongCreate{
		Artist:      input.Artist,
		Group:       input.Group,
		Title:       input.Title,
		Lyrics:      input.Lyrics,
		ReleaseDate: releaseDate,
	}

	if err = h.songService.Create(r.Context(), model); err != nil {
		logger.Errorf("creating Song: %v", err)
		response.WithBody(w, http.StatusInternalServerError, response.ResponseMessage{Message: err.Error()})
		return
	}

	response.WithBody(w, http.StatusOK, response.ResponseMessage{Message: "song successfuly created"})
}

// @Summary  Get Song By ID
// @Tags songs
// @Description  get song by id
// @ModuleID Get
// @Accept  json
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} models.Song
// @Failure 400 {object} response.ResponseMessage
// @Failure 500 {object} response.ResponseMessage
// @Failure default {object} response.ResponseMessage
// @Router /songs/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil {
		logger.Errorf("parsing param id: %v", err)
		response.WithBody(w, http.StatusBadRequest, response.ResponseMessage{Message: "invalid param id"})
		return
	}

	res, err := h.songService.Get(r.Context(), ID)
	if err != nil {
		logger.Errorf("get song by ID: %v %v", ID, err)
		response.WithBody(w, http.StatusInternalServerError, response.ResponseMessage{Message: err.Error()})
		return
	}

	response.WithBody(w, http.StatusOK, response.ResponseData{Data: res})
}

// @Summary  Update Song By ID
// @Tags songs
// @Description  update song by id
// @ModuleID Update
// @Accept  json
// @Produce  json
// @Param id path int true "update by id"
// @Param input body request.SongUpdate false "update song"
// @Success 200 {object} response.ResponseMessage
// @Failure 400 {object} response.ResponseMessage
// @Failure 422 {object} response.ResponseValidationErrors
// @Failure 500 {object} response.ResponseMessage
// @Failure default {object} response.ResponseMessage
// @Router /songs/{id} [patch]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil {
		logger.Errorf("parsing param id: %v", err)
		response.WithBody(w, http.StatusBadRequest, response.ResponseMessage{Message: "invalid param id"})
		return
	}

	var input request.SongUpdate
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Errorf("decoding SongUpdate: %v", err)
		response.WithBody(w, http.StatusBadRequest, response.ResponseMessage{Message: "invalid request body"})
		return
	}

	// Check validation
	if errBytes := validation.ValidateStruct(input); errBytes != nil {
		logger.Error("validating SongUpdate")
		response.ValidateErrors(w, http.StatusUnprocessableEntity, errBytes)
		return
	}

	var releaseDate *time.Time
	if input.ReleaseDate != nil {
		date, err := time.Parse(time.DateOnly, *input.ReleaseDate)
		if err != nil {
			logger.Errorf("parsing releaseDate: %v", err)
			response.WithBody(w, http.StatusInternalServerError, response.ResponseMessage{Message: err.Error()})
			return
		}
		releaseDate = &date
	}

	model := &models.SongUpdate{
		Artist:      input.Artist,
		Group:       input.Group,
		Title:       input.Title,
		Lyrics:      input.Lyrics,
		ReleaseDate: releaseDate,
	}

	if err := h.songService.Update(r.Context(), ID, model); err != nil {
		logger.Errorf("updating Song by ID: %v %v", ID, err)
		response.WithBody(w, http.StatusInternalServerError, response.ResponseMessage{Message: err.Error()})
		return
	}

	response.WithBody(w, http.StatusOK, response.ResponseMessage{Message: "song successfuly updated"})
}

// @Summary  Delete Song By ID
// @Tags songs
// @Description  delete song by id
// @ModuleID Delete
// @Accept  json
// @Produce  json
// @Param id path int true "update by id"
// @Success 200 {object} response.ResponseMessage
// @Failure 400 {object} response.ResponseMessage
// @Failure 500 {object} response.ResponseMessage
// @Failure default {object} response.ResponseMessage
// @Router /songs/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	IDParam := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(IDParam, 10, 64)
	if err != nil {
		logger.Errorf("parsing param id: %v", err)
		response.WithBody(w, http.StatusBadRequest, response.ResponseMessage{Message: "invalid param id"})
		return
	}

	if err := h.songService.Delete(r.Context(), ID); err != nil {
		logger.Errorf("deleting song by ID: %v %v", ID, err)
		response.WithBody(w, http.StatusInternalServerError, response.ResponseMessage{Message: err.Error()})
		return
	}

	response.WithBody(w, http.StatusOK, response.ResponseData{Data: "song successfuly deleted"})
}
