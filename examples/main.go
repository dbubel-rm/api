package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dbubel/api"
	"github.com/julienschmidt/httprouter"
	"time"
	"github.com/satori/go.uuid"
	"net/http"
)

type NullTime sql.NullTime
type NullString sql.NullString

type Foo struct {
	ID    string `json:"id"  validate:"required"`
	Index int    `json:"index" validate:"required"`
	CreatedAt NullTime
	ReportedAt NullTime
	Status NullString
}
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}

func postit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var f Foo
	err := api.UnmarshalJSON(r.Body, &f)
	if err != nil {
		api.RespondError(w, r, err, http.StatusOK)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, map[string]interface{}{"status": "OK"})
}

func getit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	orgId := orgIdHelper(r.Context())
	api.RespondJSON(w, r, http.StatusOK, &Foo{
		ID:    orgId.String(),
		Index: 1337,
		CreatedAt:NullTime{Time:time.Now(), Valid:true},
		ReportedAt:NullTime{Valid:false},
	})
}

func middlethis(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := context.WithValue(r.Context(), "middlethis", "1")
		next(w, r.WithContext(ctx), params)
	}
}

func middlethat(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		uid, err := uuid.FromString(params.ByName("orgId"))
		if err != nil {
			api.RespondError(w, r, err, http.StatusBadRequest, "invalid orgId")
			return
		}

		ctx := context.WithValue(r.Context(), "orgId", uid)
		next(w, r.WithContext(ctx), params)
	}
}

func orgIdHelper(ctx context.Context) uuid.UUID {
	return ctx.Value("orgId").(uuid.UUID)
}

func globalmiddle(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		next(w, r, params)
	}
}

func main() {
	endpoints := api.Endpoints{
		api.NewEndpoint(http.MethodGet, "/:orgId/test", getit),
		api.NewEndpoint(http.MethodGet, "/:orgId", getit),
	}

	moreEndpoints := api.Endpoints{
		api.NewEndpoint(http.MethodPost, "/test", postit),
	}

	endpoints.Use(middlethat)
	moreEndpoints.Use(middlethat)

	app := api.NewDefault()
	app.GlobalMiddleware(globalmiddle)
	app.Endpoints(endpoints, moreEndpoints)

	app.Run(&http.Server{
		Addr:           ":8000",
		Handler:        app.Router,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	})
}
