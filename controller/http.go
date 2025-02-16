package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"porter"
)

func (h *HTTPController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/api/v1/ports":
		h.PortsFromStream(w, r)
	case "/api/v1/port":
		h.GetPort(w, r)
	default:
		http.Error(w, shouldEncode(ErrorResponse{Error: "Not found"}), http.StatusNotFound)
	}
}

func NewHTTPController(svc porter.Service) *HTTPController {
	return &HTTPController{
		svc: svc,
	}
}

type HTTPController struct {
	svc porter.Service
}

func (h *HTTPController) PortsFromStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, shouldEncode(ErrorResponse{Error: "Method not allowed."}), http.StatusMethodNotAllowed)

		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, shouldEncode(ErrorResponse{Error: "Content-Type must be application/json."}), http.StatusUnsupportedMediaType)

		return
	}

	ctx := r.Context()

	if err := h.svc.PortsFromStream(ctx, r.Body); err != nil {
		http.Error(w, shouldEncode(ErrorResponse{Error: err.Error()}), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, shouldEncode(SuccessResponse{Message: "Ports updated."}))
}

func (h *HTTPController) GetPort(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, shouldEncode(ErrorResponse{Error: "Method not allowed."}), http.StatusMethodNotAllowed)

		return
	}

	unloc := r.URL.Query().Get("unloc")
	if unloc == "" {
		http.Error(w, shouldEncode(ErrorResponse{Error: "Missing unloc parameter"}), http.StatusBadRequest)

		return
	}

	ctx := r.Context()

	port, err := h.svc.GetPort(ctx, unloc)
	if err != nil {
		http.Error(w, shouldEncode(ErrorResponse{Error: err.Error()}), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, shouldEncode(port))
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func shouldEncode(in any) string {
	jsonBytes, _ := json.Marshal(in)

	return string(jsonBytes)
}
