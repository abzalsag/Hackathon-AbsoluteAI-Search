-- name: CreateAnalysis :one
INSERT INTO analyses (user_id, company_id, status)
VALUES ($1, $2, $3) RETURNING id, user_id, company_id, status, result, created_at, updated_at;

-- name: GetAnalysisByID :one
SELECT id, user_id, company_id, status, result, created_at, updated_at
FROM analyses WHERE id = $1;

-- name: UpdateAnalysisStatus :exec
UPDATE analyses SET status = $2, updated_at = NOW() WHERE id = $1;

-- name: UpdateAnalysisResult :exec
UPDATE analyses SET status = $2, result = $3, updated_at = NOW() WHERE id = $1;

-- name: GetCompetitorsByAnalysisID :many
SELECT id, analysis_id, name, website, mention_count, created_at
FROM competitors WHERE analysis_id = $1;

-- name: GetAIResponsesByAnalysisID :many
SELECT id, analysis_id, question, response_text, brand_mentioned, created_at
FROM ai_responses WHERE analysis_id = $1;