package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/minio/minio-go"
)

type ModelDataHandler struct {
	minioClient *minio.Client
	expiry      time.Duration
}

func (h *ModelDataHandler) GetModelData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")

	// Check if bucket exists
	exists, err := minioClient.BucketExists(model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	modelData := make([]string, 0)
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectsCh := minioClient.ListObjectsV2(model, "data:", true, doneCh)
	for object := range objectsCh {
		if object.Err == nil {
			modelData = append(modelData, object.Key)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modelData)
}

func (h *ModelDataHandler) GetModelDataById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")
	id := p.ByName("id")

	// Check if bucket exists
	exists, err := minioClient.BucketExists(model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Send Presigned URL
	reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject(model, "data:"+id, h.expiry, reqParams)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}

func (h *ModelDataHandler) HeadModelDataById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")
	id := p.ByName("id")

	// Check if bucket exists
	exists, err := minioClient.BucketExists(model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Get object HEAD
	modelDataInfo, err := minioClient.StatObject(model, "data:"+id, minio.StatObjectOptions{})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Send headers
	w.Header().Set("Content-Type", modelDataInfo.ContentType)
	for key, value := range modelDataInfo.Metadata {
		for _, v := range value {
			w.Header().Set(key, v)
		}
	}

	w.WriteHeader(200)
}

func (h *ModelDataHandler) UploadModelData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")
	id := p.ByName("id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	presignedURL, err := minioClient.PresignedPutObject(model, "data:"+id, h.expiry)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}

func (h *ModelDataHandler) GetModelLabels(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")

	// Check if bucket exists
	exists, err := minioClient.BucketExists(model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	modelLabels := make([]string, 0)
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectsCh := minioClient.ListObjectsV2(model, "labels:", true, doneCh)
	for object := range objectsCh {
		if object.Err == nil {
			modelLabels = append(modelLabels, object.Key)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(modelLabels)
}

func (h *ModelDataHandler) GetModelLabelsById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")
	id := p.ByName("id")

	// Check if bucket exists
	exists, err := minioClient.BucketExists(model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Send Presigned URL
	reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject(model, "labels:"+id, h.expiry, reqParams)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}

func (h *ModelDataHandler) HeadModelLabelsById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")
	id := p.ByName("id")

	// Check if bucket exists
	exists, err := minioClient.BucketExists(model)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Get object HEAD
	modelLabelsInfo, err := minioClient.StatObject(model, "labels:"+id, minio.StatObjectOptions{})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Send headers
	w.Header().Set("Content-Type", modelLabelsInfo.ContentType)
	for key, value := range modelLabelsInfo.Metadata {
		for _, v := range value {
			w.Header().Set(key, v)
		}
	}

	w.WriteHeader(200)
}

func (h *ModelDataHandler) UploadModelLabels(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")
	id := p.ByName("id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	presignedURL, err := minioClient.PresignedPutObject(model, "labels:"+id, h.expiry)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}
