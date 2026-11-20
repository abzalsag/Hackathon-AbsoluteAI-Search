package ai_client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Analyze(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/internal/analyze", r.URL.Path)
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		response := AnalysisResponse{
			VisibilityScore: 75.5,
			Competitors: []CompetitorResult{
				{Name: "Competitor 1", Website: "https://comp1.com", MentionCount: 10},
				{Name: "Competitor 2", Website: "https://comp2.com", MentionCount: 5},
			},
			Responses: []AIResponseResult{
				{
					Question:       "What is Test Corp?",
					ResponseText:   "Test Corp is a leading company...",
					BrandMentioned: true,
				},
			},
			Recommendations: []string{
				"Improve website SEO",
				"Add more AI-readable content",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := NewClient(server.URL)

	req := AnalysisRequest{
		CompanyID:   1,
		Name:        "Test Corp",
		Website:     "https://test.com",
		Description: "Test description",
	}

	result, err := client.Analyze(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 75.5, result.VisibilityScore)
	assert.Len(t, result.Competitors, 2)
	assert.Len(t, result.Responses, 1)
	assert.Len(t, result.Recommendations, 2)
	assert.True(t, result.Responses[0].BrandMentioned)
}

func TestClient_Analyze_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := NewClient(server.URL)

	req := AnalysisRequest{
		CompanyID: 1,
		Name:      "Test Corp",
	}

	result, err := client.Analyze(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "500")
}
