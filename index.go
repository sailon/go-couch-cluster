package main

import (
	"encoding/json"
	"fmt"
	"github.com/couchbaselabs/gocb"
	"net/http"
	"regexp"
)

var bucket *gocb.Bucket

type Asset struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
	Note string `json:"note,omitempty"`
}

func isValidUri(uri string) bool {
	r := regexp.MustCompile("\\w{3,}\\:\\/\\/\\w{2,}\\/\\w*")
	return r.MatchString(uri)
}

func sanitize(query string) string {
	r := regexp.MustCompile(" ")
	return r.ReplaceAllString(query, "")
}

func docQuery(uri string) *gocb.N1qlQuery {
	return gocb.NewN1qlQuery("SELECT uri, name, note FROM `assets` WHERE assets.uri=" + sanitize(uri))
}

// Handle all asset requests /api/assets
func assetHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":
		retrieveAssets(w, r)
	case "POST":
		createAssset(w, r)
	case "PUT":
	   	updateAsset(w, r)
	case "DELETE":
        removeAsset(w, r)
	}
}

// Retrieve the list of assets GET /api/assets
func retrieveAssets(w http.ResponseWriter, r *http.Request) {
	var assets []Asset
	var row Asset

	query := gocb.NewN1qlQuery("SELECT name, uri, note FROM `assets` ORDER BY name ASC")
	rows, err := bucket.ExecuteN1qlQuery(query, nil)

	if err != nil {
		http.Error(w, "N1QL query error: "+err.Error(), 500)
		return
	}

	for rows.Next(&row) {
		assets = append(assets, row)
	}

	if err := rows.Close(); err != nil {
		http.Error(w, "N1QL query error: "+err.Error(), 500)
		return
	}

	b, _ := json.Marshal(assets)
	w.Write(b)
}

// Create a new asset POST /api/assets
func createAssset(w http.ResponseWriter, r *http.Request) {
	var a Asset

	json.NewDecoder(r.Body).Decode(&a)

	if !isValidUri(a.Uri) {
		http.Error(w, "Must provide a valid URI.", 400)
		return
	}

	_, err := bucket.Upsert(a.Uri, &a, 0)

	if err != nil {
		http.Error(w, "Error: "+err.Error(), 500)
		return
	}

	b, _ := json.Marshal(a)
    w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

// Add a note to an asset PUT /api/assets
func updateAsset(w http.ResponseWriter, r *http.Request) {
	var req Asset
	var newAsset Asset

	json.NewDecoder(r.Body).Decode(&req)

	if !isValidUri(req.Uri) {
		http.Error(w, "Must provide a valid URI.", 400)
		return
	}

	if req.Note == "" {
		http.Error(w, "Must provide a note.", 400)
		return
	}

	_, err := bucket.Get(req.Uri, &newAsset)

	if err != nil {
		http.Error(w, "Error: "+err.Error(), 500)
		return
	}

    if newAsset.Uri == "" {
        http.Error(w, "Asset does not exist.", 404)
        return
    }

	newAsset.Note = req.Note

	_, err = bucket.Upsert(req.Uri, &newAsset, 0)

	if err != nil {
		http.Error(w, "Error: "+err.Error(), 500)
		return
	}

	b, _ := json.Marshal(newAsset)
    w.WriteHeader(http.StatusAccepted)
	w.Write(b)
}

// Delete an asset DELETE /api/assets
func removeAsset(w http.ResponseWriter, r *http.Request) {
    var req Asset
    var newAsset Asset

    json.NewDecoder(r.Body).Decode(&req)

    if !isValidUri(req.Uri) {
        http.Error(w, "Must provide a valid URI.", 400)
        return
    }

    _, err := bucket.Get(req.Uri, &newAsset)

    if err != nil {
        http.Error(w, "Error: "+err.Error(), 500)
        return
    }

    if newAsset.Uri == "" {
        http.Error(w, "Asset does not exist.", 404)
        return
    }

    _, err = bucket.Remove(req.Uri, 0)

    if err != nil {
        http.Error(w, "Error: "+err.Error(), 500)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func main() {
	cluster, err := gocb.Connect("couchbase://127.0.0.1")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	bucket, err = cluster.OpenBucket("assets", "")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	http.HandleFunc("/api/assets", assetHandler)

	fmt.Printf("Starting server on :9980\n")
	http.ListenAndServe(":9990", nil)

}
