
from fastapi import FastAPI
from pydantic import BaseModel

from app.geo_analyzer import analyze_brand

app = FastAPI(title="GEO Intelligence AI Service")


class AnalyzeRequest(BaseModel):
    brand: str
    website: str
    city: str
    category: str


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/analyze")
def analyze(request: AnalyzeRequest):
    return analyze_brand(
        brand=request.brand,
        website=request.website,
        city=request.city,
        category=request.category,
        aliases=[
            request.brand,
            "Центр Красок"
        ]
    )
