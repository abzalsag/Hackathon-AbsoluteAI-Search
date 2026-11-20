package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"hackathon-absolutecinema/internal/service"
)

type CompanyHandler struct {
	companyService service.CompanyServiceInterface // Интерфейс
}

func NewCompanyHandler(companyService service.CompanyServiceInterface) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
	}
}

// GetCompanies GET /api/companies
func (h *CompanyHandler) GetCompanies(w http.ResponseWriter, r *http.Request) {
	// Middleware-тен userID алу
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"Қолданушы ID табылмады"}`, http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	companies, err := h.companyService.GetCompanies(ctx, userID)
	if err != nil {
		http.Error(w, `{"error":"Компанияларды алу қатесі: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"companies": companies,
	})
}

// CreateCompany POST /api/companies
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("user_id").(int64)
	if !ok {
		http.Error(w, `{"error":"Қолданушы ID табылмады"}`, http.StatusUnauthorized)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Website     string `json:"website"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON қатесі"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, `{"error":"Компания атын енгізіңіз"}`, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	companyID, err := h.companyService.CreateCompany(ctx, userID, req.Name, req.Website, req.Description)
	if err != nil {
		http.Error(w, `{"error":"Компания жасау қатесі: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Компания сәтті қосылды",
		"company_id": companyID,
	})
}

// GetCompany GET /api/companies/{id}
func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	// URL-ден ID алу (Go 1.22+ routing)
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Қате компания ID"}`, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	company, err := h.companyService.GetCompany(ctx, id)
	if err != nil {
		http.Error(w, `{"error":"Компания табылмады"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(company)
}

// UpdateCompany PUT /api/companies/{id}
func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Қате компания ID"}`, http.StatusBadRequest)
		return
	}

	var req struct {
		Name        string `json:"name"`
		Website     string `json:"website"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON қатесі"}`, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := h.companyService.UpdateCompany(ctx, id, req.Name, req.Website, req.Description); err != nil {
		http.Error(w, `{"error":"Компанияны жаңарту қатесі: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Компания сәтті жаңартылды",
	})
}

// DeleteCompany DELETE /api/companies/{id}
func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, `{"error":"Қате компания ID"}`, http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	if err := h.companyService.DeleteCompany(ctx, id); err != nil {
		http.Error(w, `{"error":"Компанияны өшіру қатесі: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Компания сәтті өшірілді",
	})
}
