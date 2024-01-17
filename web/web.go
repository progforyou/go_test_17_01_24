package web

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type PaginationOrderDirection string

const (
	DESC PaginationOrderDirection = "DESC"
	ASC                           = "ASC"
)

type PaginationSort struct {
	Field string                   `json:"field"`
	Order PaginationOrderDirection `json:"order"`
}

type Pagination struct {
	Page    int                    `json:"page"`
	PerPage int                    `json:"perPage"`
	Filter  map[string]interface{} `json:"filter"`
}

type CrudRouterController struct {
	Name   string
	GetAll func(p Pagination) interface{}
	Create func(r []byte) (interface{}, error)
	Update func(id uint64, r []byte) (interface{}, error)
	Delete func(id uint64) (error, int64)
	Log    zerolog.Logger
}

var ZeroPagination = Pagination{Page: 0, PerPage: 20}

func NewPagination(v url.Values) Pagination {
	res := ZeroPagination
	p := v.Get("page")
	if p != "" {
		page, err := strconv.Atoi(p)
		if err != nil {
			log.Error().Err(err).Str("p", p).Msg("parse page")
		}
		res.Page = page - 1
	}
	p = v.Get("per_page")
	if p != "" {
		perPage, err := strconv.Atoi(p)
		if err != nil {
			log.Error().Err(err).Str("p", p).Msg("parse per_page")
		}
		res.PerPage = perPage
	}
	p = v.Get("filter")
	if p != "" {
		var filter map[string]interface{}
		err := json.Unmarshal([]byte(p), &filter)
		if err != nil {
			log.Error().Err(err).Str("p", p).Msg("unmarshal filter")
			return ZeroPagination
		}
		res.Filter = filter
	}
	log.Info().Interface("pagination", res).Msg("p")
	return res
}

func NewCrudRouter(controller CrudRouterController) func(r chi.Router) {
	return func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			if controller.GetAll == nil {
				w.WriteHeader(501)
				return
			}
			p := NewPagination(r.URL.Query())
			log.Info().Interface("p", p).Msg("Pagination")
			data := controller.GetAll(p)
			render.JSON(w, r, data)
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			if controller.Update == nil {
				w.WriteHeader(501)
				controller.Log.Error().Msg("not implemented")
				return
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				controller.Log.Error().Err(err).Msg("read body")
				w.WriteHeader(403)
				return
			}
			obj, err := controller.Create(body)
			if err != nil {
				controller.Log.Error().Err(err).Msg("create")
				w.WriteHeader(403)
				return
			}
			render.JSON(w, r, obj)
		})
		r.Put("/{id:\\d+}/", func(w http.ResponseWriter, r *http.Request) {
			if controller.Update == nil {
				w.WriteHeader(501)
				controller.Log.Error().Msg("not implemented")
				return
			}
			id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				w.WriteHeader(403)
				controller.Log.Error().Err(err).Msg("parse id")
				return
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				w.WriteHeader(403)
				controller.Log.Error().Err(err).Msg("read body")
				return
			}
			obj, err := controller.Update(id, body)
			if err != nil {
				w.WriteHeader(500)
				controller.Log.Error().Err(err).Msg("update")
				return
			}
			render.JSON(w, r, obj)
		})
		r.Delete("/{id:\\d+}/", func(w http.ResponseWriter, r *http.Request) {
			if controller.Delete == nil {
				w.WriteHeader(501)
				controller.Log.Error().Msg("not implemented")
				return
			}
			id, err := strconv.ParseUint(chi.URLParam(r, "id"), 10, 64)
			if err != nil {
				w.WriteHeader(403)
				controller.Log.Error().Err(err).Msg("parse id")
				return
			}
			err, rowsAffected := controller.Delete(id)

			if rowsAffected == 0 {
				w.WriteHeader(404)
				controller.Log.Error().Err(err).Msg("id not found")
				return
			}

			if err != nil {
				w.WriteHeader(500)
				controller.Log.Error().Err(err).Msg("delete")
				return
			}
			w.WriteHeader(202)
		})
	}
}
