package ai_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client interface {
	Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResponse, error)
}

type httpClient struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) Client {
	return &httpClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *httpClient) Analyze(ctx context.Context, req AnalysisRequest) (*AnalysisResponse, error) {
	url := fmt.Sprintf("%s/internal/analyze", c.baseURL)

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("JSON маршалинг қатесі: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("Сұраныс жасау қатесі: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Python сервисіне қосылу қатесі: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Python сервисі қате қайтарды: %d", resp.StatusCode)
	}

	var result AnalysisResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Жауапты оқу қатесі: %w", err)
	}

	return &result, nil
}

type AnalysisRequest struct {
	CompanyID   int64  `json:"company_id"`
	Name        string `json:"name"`
	Website     string `json:"website"`
	Description string `json:"description"`
}

type AnalysisResponse struct {
	VisibilityScore float64            `json:"visibility_score"`
	Competitors     []CompetitorResult `json:"competitors"`
	Responses       []AIResponseResult `json:"responses"`
	Recommendations []string           `json:"recommendations"`
}

type CompetitorResult struct {
	Name         string `json:"name"`
	Website      string `json:"website"`
	MentionCount int    `json:"mention_count"`
}

type AIResponseResult struct {
	Question       string `json:"question"`
	ResponseText   string `json:"response_text"`
	BrandMentioned bool   `json:"brand_mentioned"`
}
