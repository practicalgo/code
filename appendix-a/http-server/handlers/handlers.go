package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	users "github.com/practicalgo/code/appendix-a/grpc-server/service"
	"github.com/practicalgo/code/appendix-a/http-server/config"
	"github.com/practicalgo/code/appendix-a/http-server/storage"
	"github.com/practicalgo/code/appendix-a/http-server/types"
	"google.golang.org/grpc"
)

// This always returns the owner id as one of [1, 5] as the bootstrapping code only
// populates the users table with these records and since we have foreign key relationships
// the package owner must be one of those
func validateOwner(config *config.AppConfig, r *http.Request) (int, error) {
	authToken := r.Header.Get("X-Auth-Token")
	if len(authToken) == 0 {
		return 0, errors.New("specify a header X-AuthToken")
	}

	req := users.UserGetRequest{Auth: authToken}
	u, err := config.UsersSvc.GetUser(
		config.SpanCtx,
		&req,
		grpc.WaitForReady(true),
	)
	if err != nil {
		return 0, err
	}

	if u.User == nil {
		return 0, errors.New("invalid user")
	}
	return int(u.User.Id), nil
}

func packageRegHandler(
	w http.ResponseWriter,
	r *http.Request,
	config *config.AppConfig,
) {
	config.Span.AddEvent("form_data_read")

	d := types.PkgRegisterResponse{}
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
	packageOwner, err := validateOwner(config, r)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusUnauthorized,
		)
		return
	}

	q := types.PkgQueryParams{
		OwnerId: packageOwner,
		Version: packageVersion,
		Name:    packageName,
	}
	pkgResults, err := storage.QueryDb(
		config, q,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(pkgResults) != 0 {
		http.Error(w, "Package version for the owner exists", http.StatusBadRequest)
		return
	}

	d.ID = fmt.Sprintf(
		"%d/%s-%s-%s",
		packageOwner,
		packageName,
		packageVersion,
		fHeader.Filename,
	)
	config.Span.AddEvent("upload_package")

	nBytes, err := storage.UploadData(config, d.ID, fHeader)
	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	config.Span.AddEvent("update_package_db")

	err = storage.UpdateDb(
		r.Context(),
		config,
		types.PkgRow{
			OwnerId:       int(packageOwner), // TODO: Do something nicer here
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

	config.Logger.Printf(
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

func packageGetHandler(w http.ResponseWriter, r *http.Request, config *config.AppConfig) {
	queryParams := r.URL.Query()
	owner := queryParams.Get("owner_id")
	name := queryParams.Get("name")
	version := queryParams.Get("version")

	if len(owner) == 0 || len(name) == 0 || len(version) == 0 {
		http.Error(
			w,
			"Must specify package owner, name and version",
			http.StatusBadRequest,
		)
		return
	}
	ownerId, err := strconv.Atoi(owner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	q := types.PkgQueryParams{
		OwnerId: ownerId,
		Version: version,
		Name:    name,
	}
	pkgResults, err := storage.QueryDb(
		config, q,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(pkgResults) == 0 {
		http.Error(w, "No package found", http.StatusNotFound)
		return
	}

	url, err := config.PackageBucket.SignedURL(
		r.Context(),
		pkgResults[0].ObjectStoreId,
		nil,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func packageHandler(
	w http.ResponseWriter,
	r *http.Request,
	config *config.AppConfig,
) {
	switch r.Method {
	case "POST":
		packageRegHandler(w, r, config)
	case "GET":
		packageGetHandler(w, r, config)
	default:
		http.Error(w,
			"Cannot process request",
			http.StatusMethodNotAllowed,
		)
	}
}
