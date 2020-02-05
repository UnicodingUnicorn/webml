package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/minio/minio-go"
)

type ParserHandler struct {
	minioClient *minio.Client
	expiry      time.Duration
}

func (h *ParserHandler) GetParsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	parsers := make([]string, 0)

	doneCh := make(chan struct{})
	defer close(doneCh)
	objectsCh := minioClient.ListObjectsV2("parser", "", true, doneCh)
	for object := range objectsCh {
		if object.Err == nil {
			parsers = append(parsers, object.Key)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parsers)
}

func (h *ParserHandler) GetParserById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// Send Presigned URL
	reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject("parser", id, h.expiry, reqParams)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}

func (h *ParserHandler) HeadParserById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// Get object HEAD
	parserInfo, err := minioClient.StatObject("parser", id, minio.StatObjectOptions{})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Send headers
	w.Header().Set("Content-Type", parserInfo.ContentType)
	for key, value := range parserInfo.Metadata {
		for _, v := range value {
			w.Header().Set(key, v)
		}
	}

	w.WriteHeader(200)
}

func (h *ParserHandler) UploadParser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	id := RandomHex()

	// Validate metadata
	if r.Header.Get("x-amz-meta-name") == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}


	presignedURL, err := minioClient.PresignedPutObject("parser", id, h.expiry)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}
