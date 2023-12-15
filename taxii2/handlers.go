package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"opentaxii/internal/models"

	"github.com/go-chi/chi/v5"
)

func (app *application) getTaxii2(w http.ResponseWriter, r *http.Request) {
	var apiroots []string
	for _, apiRoot := range app.config.TAXII.ApiRoots {
		apiroots = append(apiroots, apiRoot.URL)
	}

	discovery := models.Discovery{
		Title:       app.config.TAXII.Title,
		Description: app.config.TAXII.Description,
		Contact:     app.config.TAXII.Contact,
		Default:     app.config.TAXII.Default,
		API_ROOTS:   apiroots,
	}

	fmt.Println(app.config.TAXII.Default)

	// Respond with the discovery JSON
	app.respondJSON(w, discovery)
}

func (app *application) getAPIRoot(w http.ResponseWriter, r *http.Request) {
	api_root := chi.URLParam(r, "api-root")
	apiRoot, err := app.DB.GetApiRoot(api_root)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	out, err := json.MarshalIndent(apiRoot, "", "   ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/taxii+json;version=2.1")
	w.Write(out)
}

func (app *application) getAPIRootStatusID(w http.ResponseWriter, r *http.Request) {
	api_root := chi.URLParam(r, "api-root")
	id := chi.URLParam(r, "status-id")

	status := map[string]interface{}{
		"id":                "2d086da7-4bdc-4f91-900e-d77486753710",
		"status":            "pending",
		"request_timestamp": "2016-11-02T12:34:34.12345Z",
		"total_count":       4,
		"success_count":     1,
		"successes": models.StatusDetails{
			ID:      "indicator--c410e480-e42b-47d1-9476-85307c12bcbf",
			Version: "2018-05-27T12:02:41.312Z",
		},
		"failure_count": 1,
		"failures": models.StatusDetails{
			ID:      "malware--664fa29d-bf65-4f28-a667-bdb76f29ec98",
			Version: "2018-05-28T14:03:42.543Z",
			Message: "Unable to process object",
		},
		"pending_count": 2,
		"pendings": []models.StatusDetails{
			{
				ID:      "indicator--252c7c11-daf2-42bd-843b-be65edca9f61",
				Version: "2018-05-18T20:16:21.148Z",
			},
			{
				ID:      "relationship--045585ad-a22f-4333-af33-bfd503a683b5",
				Version: "2018-05-15T10:13:32.579Z",
			},
		},
	}
	fmt.Println(api_root, id)
	app.respondJSON(w, status)
}

func (app *application) getCollections(w http.ResponseWriter, r *http.Request) {
	api_root := chi.URLParam(r, "api-root")
	allCollections, err := app.DB.GetAllCollections()
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	fmt.Println(api_root)
	app.respondJSON(w, allCollections)
}

func (app *application) getCollectionByID(w http.ResponseWriter, r *http.Request) {
	// api_root := chi.URLParam(r, "api-root")
	id := chi.URLParam(r, "collectionID")

	collectionID, err := strconv.Atoi(id)
	collection, err := app.DB.GetCollectionID(collectionID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	app.respondJSON(w, collection)
}

func (app *application) getMetadataAboutCollectionID(w http.ResponseWriter, r *http.Request) {
	api_root := chi.URLParam(r, "api-root")
	id := chi.URLParam(r, "collectionID")

	collectionID, err := strconv.Atoi(id)
	manifest, err := app.DB.GetManifestRecord(collectionID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	fmt.Printf(api_root, collectionID)
	app.respondJSON(w, manifest)

}

func (app *application) getObjects(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getObjectByID(w http.ResponseWriter, r *http.Request) {

}

func (app *application) getObjectIDVersions(w http.ResponseWriter, r *http.Request) {

}

func (app *application) addObject(w http.ResponseWriter, r *http.Request) {

}

func (app *application) deleteObject(w http.ResponseWriter, r *http.Request) {

}

func (app *application) postCollection(w http.ResponseWriter, r *http.Request, collectionID string) {
	response := map[string]string{
		"status": "success",
	}

	app.respondJSON(w, response)
}
