package main

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	badger "github.com/dgraph-io/badger"
)

const ALPHA = 0.9

type Tensor struct {
	Shape []int     `json:"shape"`
	Data  []float64 `json:"data"`
}

type Session struct {
	Loss    float64 `json:"loss"`
	Weights Tensor  `json:"weights"`
	Alpha   float64 `json:"alpha"`
	Model		string	`json:"model"`
}

type ValuesHandler struct {
	Badger *badger.DB
}

func (h *ValuesHandler) RetrieveSession(id string) (*Session, error) {
	var session Session

	err := h.Badger.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}

		rawSession, err := item.ValueCopy(nil)
		if err != nil {
			return err
		}

		err = json.Unmarshal(rawSession, &session)
		return err
	})

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (h *ValuesHandler) GetSession(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	session, err := h.RetrieveSession(id)
	if err == badger.ErrKeyNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

type LossRequest struct {
	Loss float64 `json:"loss"`
}

func (h *ValuesHandler) PostLoss(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// Retrieve Session
	session, err := h.RetrieveSession(id)
	if err == badger.ErrKeyNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Decode request
	var loss LossRequest
	err = json.NewDecoder(r.Body).Decode(&loss)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Update loss
	session.Loss = loss.Loss

	// Update session
	sessionData, err := json.Marshal(&session)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_ = h.Badger.Update(func(txn *badger.Txn) error { // Error should only be ErrReadOnlyTxn
		return txn.Set([]byte(id), sessionData)
	})

	w.WriteHeader(200)
}

func (h *ValuesHandler) PostWeights(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// Retrieve Session
	session, err := h.RetrieveSession(id)
	if err == badger.ErrKeyNotFound {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Decode request
	var newWeights Tensor
	err = json.NewDecoder(r.Body).Decode(&newWeights)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check shape is the same
	same := true
	shapeLength := 0
	if len(newWeights.Shape) != len(session.Weights.Shape) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	for i, v := range newWeights.Shape {
		if v != session.Weights.Shape[i] {
			same = false
			break
		} else {
			shapeLength *= v
		}
	}
	if !same {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Check data length matches shape
	if len(newWeights.Data) != shapeLength {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Merge two fields
	for i, newWeight := range newWeights.Data {
		session.Weights.Data[i] = session.Weights.Data[i] * session.Alpha + newWeight * (1.0 - session.Alpha)
	}

	// Update session
	sessionData, err := json.Marshal(&session)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_ = h.Badger.Update(func(txn *badger.Txn) error { // Error should only be ErrReadOnlyTxn
		return txn.Set([]byte(id), sessionData)
	})

	w.WriteHeader(200)
}

type NewSessionReq struct {
	Shape []int   `json:"shape"`
	Loss  float64 `json:"loss,omitempty"`
	Alpha float64 `json:"alpha,omitempty"`
}

func (h *ValuesHandler) NewSession(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")

	// Check model exists
	exists, err := minioClient.BucketExists(id)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Decode body
	var req NewSessionReq
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	// Default values
	if req.Loss == 0.0 {
		req.Loss = 1.0
	}

	if req.Alpha == 0.0 {
		req.Alpha = ALPHA
	}

	// Make weights
	weightsLength := 1
	for _, n := range req.Shape {
		weightsLength *= n
	}
	weights := make([]float64, weightsLength)

	// Generate ID
	sid := RandomHex()

	session := Session {
		Loss:  req.Loss,
		Alpha: req.Alpha,
		Weights: Tensor{
			Shape: req.Shape,
			Data:  weights,
		},
		Model: id,
	}

	sessionData, err := json.Marshal(&session)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	_ = h.Badger.Update(func(txn *badger.Txn) error { // Error should only be ErrReadOnlyTxn
		return txn.Set([]byte(sid), sessionData)
	})

	w.WriteHeader(200)
}
