package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/minio/minio-go"
)

type ModelHandler struct {
	minioClient *minio.Client
	expiry      time.Duration
}

func (h *ModelHandler) GetModels(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	buckets, err := minioClient.ListBuckets()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	bucketNames := make([]string, 0)
	for _, bucket := range buckets {
		if bucket.Name != "parser" {
			bucketNames = append(bucketNames, bucket.Name)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bucketNames)
}

func (h *ModelHandler) GetModelById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// Check if bucket exists
	exists, err := minioClient.BucketExists(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Send Presigned URL
	reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject(id, "model", h.expiry, reqParams)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}

func (h *ModelHandler) UploadModel(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	bucketName := p.ByName("id")
	if bucketName == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Create bucket for parsers if it doesn't exist
	exists, err := minioClient.BucketExists(bucketName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	err = minioClient.MakeBucket(bucketName, "us-east-1")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	presignedURL, err := minioClient.PresignedPutObject(bucketName, "model", h.expiry)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}
