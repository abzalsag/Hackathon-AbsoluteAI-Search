
import os
from urllib.parse import urlparse

from dotenv import load_dotenv
from tavily import TavilyClient

load_dotenv()

tavily = TavilyClient(
    api_key=os.getenv("TAVILY_API_KEY")
)

PAINT_KEYWORDS = [
    "краск",
    "лакокрас",
    "эмал",
    "грунтов",
    "покраск",
    "отделочн",
    "paint",
    "color studio",
    "декоративн"
]

BLOCKED_KEYWORDS = [
    "ваканси",
    "криптовалют",
    "автозапчаст",
    "фармацевтик",
    "препарат",
    "linkedin.com",
    "cryptorank.io",
    "kraken.com"
]


def is_relevant(source):
    title = source.get("title", "")
    content = source.get("content", "")
    url = source.get("url", "")

    text = f"{title} {content} {url}".lower()

    if any(word in text for word in BLOCKED_KEYWORDS):
        return False

    return any(word in text for word in PAINT_KEYWORDS)


def search_sources(question, max_results=10):
    result = tavily.search(
        query=question + " Казахстан магазины краски",
        search_depth="basic",
        max_results=max_results
    )

    sources = []
    seen_urls = set()

    for item in result.get("results", []):
        url = item.get("url", "")

        if not url:
            continue

        parsed = urlparse(url)

        # Убираем дубликаты одинаковых страниц
        normalized_url = url.rstrip("/").lower()

        if normalized_url in seen_urls:
            continue

        source = {
            "title": item.get("title", ""),
            "url": url,
            "domain": parsed.netloc.lower(),
            "content": item.get("content", "")[:1200]
        }

        if not is_relevant(source):
            continue

        seen_urls.add(normalized_url)
        sources.append(source)

    return sources
