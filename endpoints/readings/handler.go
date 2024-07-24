package readings

import (
	"encoding/json"
	"fmt"
	"io"
	"joi-energy-golang/api"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"joi-energy-golang/domain"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) StoreReadings(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		api.Error(w, r, fmt.Errorf("read request body failed: %w", err), http.StatusBadRequest)
		return
	}
	var req domain.StoreReadings
	if err := json.Unmarshal(body, &req); err != nil {
		api.Error(w, r, fmt.Errorf("unmarshal request body failed: %w", err), http.StatusBadRequest)
		return
	}
	err = validateSmartMeterId(req.SmartMeterId)
	if err != nil {
		api.Error(w, r, err, http.StatusBadRequest)
		return
	}
	h.service.StoreReadings(req.SmartMeterId, req.ElectricityReadings)
	api.Success(w, r, nil)
}

func (h *Handler) GetReadings(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	smartMeterId := urlParams.ByName("smartMeterId")
	err := validateSmartMeterId(smartMeterId)
	if err != nil {
		api.Error(w, r, err, http.StatusBadRequest)
		return
	}
	readings := h.service.GetReadings(smartMeterId)
	result := domain.StoreReadings{
		SmartMeterId:        smartMeterId,
		ElectricityReadings: readings,
	}
	api.SuccessJson(w, r, result)
}

func (h *Handler) GetLastWeekUsage(w http.ResponseWriter, r *http.Request, urlParams httprouter.Params) {
	smartMeterId := urlParams.ByName("smartMeterId")
	err := validateSmartMeterId(smartMeterId)
	if err != nil {
		api.Error(w, r, err, http.StatusBadRequest)
		return
	}
	cost, err := h.service.CalculateLastWeekUsageCost(smartMeterId, 1)
	if err != nil {
		api.Error(w, r, err, http.StatusInternalServerError)
		return
	}
	result := domain.LastWeekUsage{
		SmartMeterId: smartMeterId,
		Cost:         cost,
	}
	api.SuccessJson(w, r, result)
}
