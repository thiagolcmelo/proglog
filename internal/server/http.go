package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// NewHTTPServer creates a Log server
func NewHTTPServer(addr string) *http.Server {
	httpsvr := newHTTPServer()
	r := mux.NewRouter()
	r.HandleFunc("/", httpsvr.handleProduce).Methods("POST")
	r.HandleFunc("/", httpsvr.handleConsume).Methods("GET")
	return &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
}

type httpServer struct {
	Log *Log
}

func newHTTPServer() *httpServer {
	return &httpServer{
		Log: NewLog(),
	}
}

func (s httpServer) handleProduce(w http.ResponseWriter, r *http.Request) {
	// unmarshal request
	var req ProduceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("error decoding: %s\n", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// run the logic
	off, err := s.Log.Append(req.Record)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// marshal response
	res := ProduceResponse{Offset: off}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s httpServer) handleConsume(w http.ResponseWriter, r *http.Request) {
	// unmarshal request
	var req ConsumeRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// run the logic
	record, err := s.Log.Read(req.Offset)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// marshal response
	res := ConsumeResponse{Record: record}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// ProduceRequest encapsulates a produce request
type ProduceRequest struct {
	Record Record `json:"record"`
}

// ProduceResponse encapsulates a produce response
type ProduceResponse struct {
	Offset uint64 `json:"offset"`
}

// ConsumeRequest encapsulates a consume request
type ConsumeRequest struct {
	Offset uint64 `json:"offset"`
}

// ConsumeResponse encapsulates a consume response
type ConsumeResponse struct {
	Record Record `json:"record"`
}
