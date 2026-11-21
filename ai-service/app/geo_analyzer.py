import json
import re
from datetime import datetime, timezone
from urllib.parse import urlparse

from app.search import search_sources
from app.analyzer import analyze_with_ai


def normalize_text(text):
    text = text.lower().replace("ё", "е")
    text = re.sub(r"[^\w\s]", " ", text)
    return " ".join(text.split())


def brand_is_mentioned(answer, aliases):
    normalized_answer = normalize_text(answer)

    return any(
        normalize_text(alias) in normalized_answer
        for alias in aliases
        if alias.strip()
    )


def get_domain(url):
    domain = urlparse(url).hostname or ""
    return domain.lower().removeprefix("www.")


def analyze_brand(brand, website, city, category, aliases=None):
    aliases = aliases or [brand]

    question = f"Где купить {category} в городе {city}?"

    sources = search_sources(question)
    answer = analyze_with_ai(question, sources)

    mentioned = brand_is_mentioned(answer, aliases)

    official_domain = get_domain(website)

    website_found = any(
        get_domain(source["url"]) == official_domain
        for source in sources
    )

    result = {
        "brand": brand,
        "website": website,
        "city": city,
        "category": category,
        "question": question,
        "answer": answer,
        "visibility": {
            "mentioned": mentioned,
            "website_found": website_found
        },
        "sources_count": len(sources),
        "sources": sources,
        "analyzed_at": datetime.now(timezone.utc).isoformat()
    }

    return result


if __name__ == "__main__":
    result = analyze_brand(
        brand="Центр Красок #1",
        website="https://centr-krasok.kz",
        city="Алматы",
        category="качественную краску для стен",
        aliases=[
            "Центр Красок #1",
            "Центр Красок"
        ]
    )

    print(json.dumps(
        result,
        ensure_ascii=False,
        indent=2
    ))
