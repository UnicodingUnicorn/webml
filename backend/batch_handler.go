package main

import (
	"bytes"
	"crypto/rand"
	"io"
	"log"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/minio/minio-go"
	"github.com/yuin/gopher-lua"
)

type BatchHandler struct {
	minioClient *minio.Client
	expiry      time.Duration
}

func (h *BatchHandler) GetBatch(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")

	ids := make([]string, 0)
	doneCh := make(chan struct{})
	defer close(doneCh)
	objectsCh := h.minioClient.ListObjectsV2(model, "batch:data:", true, doneCh)
	for object := range objectsCh {
		if object.Err == nil {
			ids = append(ids, object.Key)
		}
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(ids))))
	w.Write([]byte(ids[n.Int64()]))
}

func (h *BatchHandler) GetBatchData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")
	id := p.ByName("id")

	reqParams := make(url.Values)
	presignedURL, err := h.minioClient.PresignedGetObject(model, "batch:data:"+id, h.expiry, reqParams)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}

func (h *BatchHandler) GetBatchLabels(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")
	id := p.ByName("id")

	reqParams := make(url.Values)
	presignedURL, err := h.minioClient.PresignedGetObject(model, "batch:labels:"+id, h.expiry, reqParams)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, presignedURL.String(), http.StatusTemporaryRedirect)
}

// Parse and split a dataset into batches
func (h *BatchHandler) BatchData(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	model := p.ByName("model")
	dataId := p.ByName("id")

	dataParserId := r.FormValue("data_parser")
	if dataParserId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	labelParserId := r.FormValue("label_parser")
	if labelParserId == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	batchSizeText := r.FormValue("batch_size")
	if batchSizeText == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	batchSize, err := strconv.Atoi(batchSizeText)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	dataParserObject, err := h.minioClient.GetObject("parser", dataParserId, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	dataParserBytes := new(bytes.Buffer)
	dataParserBytes.ReadFrom(dataParserObject)
	dataParser := dataParserBytes.String()

	labelParserObject, err := h.minioClient.GetObject("parser", labelParserId, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	labelParserBytes := new(bytes.Buffer)
	labelParserBytes.ReadFrom(labelParserObject)
	labelParser := labelParserBytes.String()

	dataL := lua.NewState()
	defer dataL.Close()
	err = dataL.DoString(dataParser)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	dataObject, err := h.minioClient.GetObject(model, "data:"+dataId, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	batchIds := make([]string, 0)

	buf := make([]byte, 512)
	batch := make([][]byte, 0)
	for {
		n, err := dataObject.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = dataL.CallByParam(lua.P{
			Fn:      dataL.GetGlobal("parse"),
			NRet:    1,
			Protect: true,
		}, lua.LString(buf), lua.LNumber(n))
		if err != nil {
			log.Printf("%s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		lv := dataL.Get(-1)
		dataL.Pop(1)
		if table, ok := lv.(*lua.LTable); ok {
			table.ForEach(func(_ lua.LValue, v lua.LValue) {
				val := []byte(v.(lua.LString).String())
				batch = append(batch, val)
				if len(batch) >= batchSize {
					batchId := RandomHex()
					data := make([]byte, 0)
					for _, datum := range batch {
						data = append(data, datum...)
					}
					_, err := h.minioClient.PutObject(model, "batch:data:"+batchId, bytes.NewReader(data), -1, minio.PutObjectOptions{})
					if err != nil {
						log.Printf("%s", err)
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						return
					}
					batchIds = append(batchIds, batchId)
				}
			})
		}
	}

	labelL := lua.NewState()
	defer labelL.Close()
	err = labelL.DoString(labelParser)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	labelObject, err := h.minioClient.GetObject(model, "labels:"+dataId, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	buf = make([]byte, 512)
	i := 0
	for {
		n, err := labelObject.Read(buf)
		if err == io.EOF {
			break
		} else if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		err = labelL.CallByParam(lua.P{
			Fn:      labelL.GetGlobal("parse"),
			NRet:    1,
			Protect: true,
		}, lua.LString(buf), lua.LNumber(n))
		if err != nil {
			log.Printf("%s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		lv := labelL.Get(-1)
		labelL.Pop(1)
		if table, ok := lv.(*lua.LTable); ok {
			table.ForEach(func(_ lua.LValue, v lua.LValue) {
				val := []byte(v.(lua.LString).String())
				batch = append(batch, val)
				if len(batch) >= batchSize {
					batchId := batchIds[i]
					data := make([]byte, 0)
					for _, datum := range batch {
						data = append(data, datum...)
					}
					_, err := h.minioClient.PutObject(model, "batch:labels:"+batchId, bytes.NewReader(data), -1, minio.PutObjectOptions{})
					if err != nil {
						log.Printf("%s", err)
						http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
						return
					}
					i += 1
					if i >= len(batchIds) {
						w.WriteHeader(200)
						return
					}
				}
			})
		}
	}

	w.WriteHeader(200)
}
