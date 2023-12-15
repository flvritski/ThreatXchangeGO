package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Get("/taxii2", app.getTaxii2)
	mux.Get("/{api-root}", app.getAPIRoot)
	mux.Get("/{api-root}/status/{status-id}", app.getAPIRootStatusID)
	mux.Get("/{api-root}/collections", app.getCollections)
	mux.Route("/{api-root}/collections/{collectionID}", func(r chi.Router) {
		r.Get("/", app.getCollectionByID)
		r.Post("/", app.getCollectionByID)
	})
	mux.Get("/{api-root}/collections/{collectionID}/manifest", app.getMetadataAboutCollectionID)
	mux.Route("/{api-root}/collections/{collectionID}/objects", func(r chi.Router) {
		r.Get("/", app.getObjects)
		r.Post("/", app.addObject)
		r.Delete("/", app.deleteObject)
	})
	mux.Get("/{api-root}/collections/{collectionID}/objects/{object-id}", app.getObjectByID)
	mux.Get("/{api-root}/collections/{collectionID}/objects/{object-id}/versions", app.getObjectIDVersions)
	return mux

}
