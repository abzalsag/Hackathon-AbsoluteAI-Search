package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"hackathon-absolutecinema/internal/service"
)

type AnalysisHandler struct {
	analysisService service.AnalysisServiceInterface // Интерфейс
}

func NewAnalysisHandler(analysisService service.AnalysisServiceInterface) *AnalysisHandler {
	return &AnalysisHandler{
		analysisService: analysisService,
	}
}

func (h *AnalysisHandler) CreateAnalysis(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"Қолданушы ID табылмады"}`, http.StatusUnauthorized)
		return
	}

	var req struct {
		CompanyID int64 `json:"company_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON қатесі"}`, http.StatusBadRequest)
		return
	}

	if req.CompanyID == 0 {
		http.Error(w, `{"error":"company_id міндетті"}`, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	analysisID, err := h.analysisService.CreateAnalysis(ctx, userID, req.CompanyID)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":     "Анализ басталды",
		"analysis_id": analysisID,
		"status":      "pending",
	})
}

func (h *AnalysisHandler) GetAnalysis(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Қате анализ ID"}`, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	analysis, err := h.analysisService.GetAnalysis(ctx, id)
	if err != nil {
		http.Error(w, `{"error":"Анализ табылмады"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analysis)
}
