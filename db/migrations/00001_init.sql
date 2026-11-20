-- +goose Up
CREATE TABLE users (
                       id BIGSERIAL PRIMARY KEY,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       name VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE companies (
                           id BIGSERIAL PRIMARY KEY,
                           user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
                           name VARCHAR(255) NOT NULL,
                           website VARCHAR(500) NOT NULL,
                           description TEXT,
                           created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE analyses (
                          id BIGSERIAL PRIMARY KEY,
                          user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
                          company_id BIGINT REFERENCES companies(id) ON DELETE CASCADE,
                          status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, processing, completed, failed
                          result JSONB,
                          created_at TIMESTAMP DEFAULT NOW(),
                          updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE competitors (
                             id BIGSERIAL PRIMARY KEY,
                             analysis_id BIGINT REFERENCES analyses(id) ON DELETE CASCADE,
                             name VARCHAR(255) NOT NULL,
                             website VARCHAR(500),
                             mention_count INT DEFAULT 0,
                             created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE ai_responses (
                              id BIGSERIAL PRIMARY KEY,
                              analysis_id BIGINT REFERENCES analyses(id) ON DELETE CASCADE,
                              question TEXT NOT NULL,
                              response_text TEXT NOT NULL,
                              brand_mentioned BOOLEAN DEFAULT FALSE,
                              created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_companies_user_id ON companies(user_id);
CREATE INDEX idx_analyses_user_id ON analyses(user_id);
CREATE INDEX idx_analyses_company_id ON analyses(company_id);
CREATE INDEX idx_competitors_analysis_id ON competitors(analysis_id);
CREATE INDEX idx_ai_responses_analysis_id ON ai_responses(analysis_id);

-- +goose Down
DROP TABLE IF EXISTS ai_responses;
DROP TABLE IF EXISTS competitors;
DROP TABLE IF EXISTS analyses;
DROP TABLE IF EXISTS companies;
DROP TABLE IF EXISTS users;