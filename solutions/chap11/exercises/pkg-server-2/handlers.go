package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
)

// This always returns the owner id as one of [1, 5] as the bootstrapping code only
// populates the users table with these records and since we have foreign key relationships
// the package owner must be one of those
func getOwnerId() int {
	return rand.Intn(4) + 1
}

func packageRegHandler(
	w http.ResponseWriter,
	r *http.Request,
	config appConfig,
) {
	d := pkgRegisterResponse{}
	err := r.ParseMultipartForm(5000)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}
	mForm := r.MultipartForm
	fHeader := mForm.File["filedata"][0]

	packageName := mForm.Value["name"][0]
	packageVersion := mForm.Value["version"][0]
	packageOwner := getOwnerId()

	d.ID = fmt.Sprintf(
		"%d/%s-%s-%s",
		packageOwner,
		packageName,
		packageVersion,
		fHeader.Filename,
	)

	nBytes, err := uploadData(config, d.ID, fHeader)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	err = updateDb(
		config,
		pkgRow{
			OwnerId:       packageOwner,
			Name:          packageName,
			Version:       packageVersion,
			ObjectStoreId: d.ID,
		},
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	config.logger.Printf(
		"Package uploaded: %s. Bytes written: %d\n",
		d.ID,
		nBytes,
	)
	jsonData, err := json.Marshal(d)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonData))
}

func packageQueryHandler(
	w http.ResponseWriter,
	r *http.Request,
	config appConfig,
) {
	queryParams := r.URL.Query()
	packageName := queryParams.Get("name")
	packageVersion := queryParams.Get("version")

	q := pkgQueryParams{
		packageVersion: packageVersion,
		packageName:    packageName,
	}

	pkgResults, err := queryDb(
		config, q,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	jsonData, err := json.Marshal(pkgResults)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(jsonData))
}

func packageGetHandler(w http.ResponseWriter, r *http.Request, config appConfig) {
	queryParams := r.URL.Query()
	packageName := queryParams.Get("name")
	packageVersion := queryParams.Get("version")

	q := pkgQueryParams{
		packageVersion: packageVersion,
		packageName:    packageName,
	}
	pkgResults, err := queryDb(
		config, q,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	if len(pkgResults) == 0 {
		http.Error(w, "No package found", http.StatusNotFound)
		return
	}

	url, err := config.packageBucket.SignedURL(
		r.Context(),
		pkgResults[0].ObjectStoreId,
		nil,
	)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func packageHandler(
	w http.ResponseWriter,
	r *http.Request,
	config appConfig,
) {
	switch r.Method {
	case "POST":
		packageRegHandler(w, r, config)
	case "GET":
		packageQueryHandler(w, r, config)
	default:
		http.Error(w,
			"Cannot process request",
			http.StatusMethodNotAllowed,
		)
	}
}
