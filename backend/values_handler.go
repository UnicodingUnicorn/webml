package main

import (
  "encoding/json"
  "fmt"
  "net/http"

  "github.com/julienschmidt/httprouter"
)

const ALPHA = 0.9

type Tensor struct {
  Shape []int `json:"shape"`
  Data []float64 `json:"data"`
}

type Session struct {
  Loss float64 `json:"loss"`
  Weights Tensor `json:"weights"`
  Alpha float64 `json:"alpha"`
}

type ValuesHandler struct {
  Sessions map[string]Session
}

func (h *ValuesHandler) GetLoss(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  id := p.ByName("id")
  if loss, ok := h.Sessions[id]; ok {
    lossString := fmt.Sprintf("%f", loss)
    w.Write([]byte(lossString))
  } else {
    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
    return
  }
}

type LossRequest struct {
  Loss float64 `json:"loss"`
}
func (h *ValuesHandler) PostLoss(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  id := p.ByName("id")

  var loss LossRequest
  decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&loss)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
  }

  if session, ok := h.Sessions[id]; ok {
    session.Loss = loss.Loss
    w.WriteHeader(200)
  } else {
    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
    return
  }
}

func (h *ValuesHandler) PostWeights(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  id := p.ByName("id")

  var newWeights Tensor
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&newWeights)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return
  }

  // Check session exists
  var session Session
  if s, ok := h.Sessions[id]; ok {
    session = s
  } else {
    http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
    return
  }

  // Check shape is the same
  same := true
  shapeLength := 0
  if len(newWeights.Shape) != len(session.Weights.Shape) {
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return
  }
  for i, v := range(newWeights.Shape) {
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
  for i, newWeight := range(newWeights.Data) {
    session.Weights.Data[i] = session.Weights.Data[i] * session.Alpha + newWeight * (1.0 - session.Alpha)
  }

  w.WriteHeader(200)
}

type NewSessionReq struct {
  Shape  []int `json:"shape"`
  Loss  float64 `json:"loss, omitempty"`
  Alpha float64 `json:"alpha, omitempty"`
}
func (h *ValuesHandler) NewSession(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
  id := p.ByName("id")

  var req NewSessionReq
  decoder := json.NewDecoder(r.Body)
  err := decoder.Decode(&req)
  if err != nil {
    http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
    return
  }

  if _, ok := h.Sessions[id]; !ok {
    if req.Loss == 0.0 {
      req.Loss = 1.0
    }

    if req.Alpha == 0.0 {
      req.Alpha = ALPHA
    }

    weightsLength := 1
    for _, n := range(req.Shape) {
      weightsLength *= n
    }
    weights := make([]float64, weightsLength)

    h.Sessions[id] = Session {
      Loss: req.Loss,
      Alpha: req.Alpha,
      Weights: Tensor {
        Shape: req.Shape,
        Data: weights,
      },
    }

    w.WriteHeader(200)
  } else {
    http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
    return
  }
}
