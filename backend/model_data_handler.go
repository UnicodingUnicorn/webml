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
  expiry time.Duration
};

func (h *ModelDataHandler) GetModelData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  model := p.ByName("model")

  modelData :=  make([]string, 0)
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

  // Send Presigned URL
  reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject(model, "data:" + id, h.expiry, reqParams)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

  http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}

func (h *ModelDataHandler) UploadModelData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  model := p.ByName("model")

  id := r.FormValue("id")
  if id == "" {
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return
  }

  presignedURL, err := minioClient.PresignedPutObject(model, "data:" + id, h.expiry)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}

func (h *ModelDataHandler) GetModelLabels(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  model := p.ByName("model")

  modelLabels :=  make([]string, 0)
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

  // Send Presigned URL
  reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject(model, "labels:" + id, h.expiry, reqParams)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

  http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}

func (h *ModelDataHandler) UploadModelLabels(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  model := p.ByName("model")

  id := r.FormValue("id")
  if id == "" {
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return
  }

  presignedURL, err := minioClient.PresignedPutObject(model, "labels:" + id, h.expiry)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
    return
  }

  http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}
