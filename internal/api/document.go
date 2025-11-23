package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testProject/internal/api/dto"
	"testProject/internal/model"
	"testProject/internal/service"
	app_err "testProject/pkg/app-err"
	"testProject/pkg/utils"

	"github.com/gorilla/mux"
)

type DocumentApi struct {
	docService service.DocumentService
}

func NewDocumentApi(docService service.DocumentService) *DocumentApi {
	return &DocumentApi{docService: docService}
}

func (d *DocumentApi) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req dto.DocumentRequest

		URL := r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.RespondMessage(w, http.StatusBadRequest, app_err.WriteError(app_err.ParseErr.Error(), URL))
			return
		}

		if err := d.docService.Create(r.Context(), dto.ToDocument(&req)); err != nil {
			utils.RespondMessage(w, http.StatusInternalServerError, app_err.WriteError(app_err.ServerErr.Error(), URL))
			return
		}

		utils.RespondMessage(w, http.StatusCreated, nil)
	}
}

func (d *DocumentApi) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		URL := r.URL.Path

		limit, offset, err := getPaginationParams(query)
		if err != nil {
			utils.RespondMessage(w, http.StatusBadRequest, app_err.WriteError(app_err.ParseErr.Error(), URL))
			return
		}

		documents, err := d.docService.FindAll(r.Context(), limit, offset)
		if err != nil {
			utils.RespondMessage(w, http.StatusInternalServerError, app_err.WriteError(app_err.ServerErr.Error(), URL))
			return
		}

		utils.RespondMessage(w, http.StatusOK, documents)
	}
}

func (d *DocumentApi) FindByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		URL := r.URL.Path

		document, err := d.docService.FindByID(r.Context(), id)
		if err != nil {
			if errors.Is(err, app_err.DocumentNotFoundErr) {
				utils.RespondMessage(w, http.StatusNotFound, app_err.WriteError(app_err.DocumentNotFoundErr.Error(), URL))
				return
			}
			utils.RespondMessage(w, http.StatusInternalServerError, app_err.WriteError(app_err.ServerErr.Error(), URL))
			return
		}
		utils.RespondMessage(w, http.StatusOK, document)
	}
}

func (d *DocumentApi) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		URL := r.URL.Path

		if err := d.docService.Delete(r.Context(), id); err != nil {
			if errors.Is(err, app_err.DocumentNotFoundErr) {
				utils.RespondMessage(w, http.StatusNotFound, app_err.WriteError(app_err.DocumentNotFoundErr.Error(), URL))
				return
			}
			utils.RespondMessage(w, http.StatusInternalServerError, app_err.WriteError(app_err.ServerErr.Error(), URL))
			return
		}
		utils.RespondMessage(w, http.StatusOK, nil)
	}
}

func (d *DocumentApi) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var doc model.Document
		URL := r.URL.Path
		if err := json.NewDecoder(r.Body).Decode(&doc); err != nil {
			utils.RespondMessage(w, http.StatusBadRequest, app_err.WriteError(app_err.ParseErr.Error(), URL))
			return
		}

		if err := d.docService.Update(r.Context(), &doc); err != nil {
			if errors.Is(err, app_err.DocumentNotFoundErr) {
				utils.RespondMessage(w, http.StatusNotFound, app_err.WriteError(app_err.DocumentNotFoundErr.Error(), URL))
				return
			}
			utils.RespondMessage(w, http.StatusInternalServerError, app_err.WriteError(app_err.UpdateDocumentErr.Error(), URL))
			return
		}

		utils.RespondMessage(w, http.StatusOK, nil)
	}
}

func getPaginationParams(query url.Values) (limit, offset int, err error) {
	limit = 50
	offset = 0

	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return 0, 0, app_err.ParseErr
		}
	}

	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			return 0, 0, app_err.ParseErr
		}
	}

	if limit > 500 {
		limit = 500
	}

	return limit, offset, nil
}
