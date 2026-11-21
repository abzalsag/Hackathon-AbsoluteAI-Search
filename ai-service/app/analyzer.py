
import os

from dotenv import load_dotenv
from google import genai

load_dotenv()

gemini = genai.Client(
    api_key=os.getenv("GEMINI_API_KEY")
)


def analyze_with_ai(question, sources):
    if not sources:
        return "Недостаточно подходящих источников для ответа."

    context = ""

    for i, source in enumerate(sources, 1):
        context += (
            f"\nИсточник [{i}]\n"
            f"Название: {source['title']}\n"
            f"URL: {source['url']}\n"
            f"Текст: {source['content']}\n"
        )

    prompt = f"""
Ты — помощник по выбору магазинов.

Ответь на вопрос пользователя, используя
только предоставленные источники.

Правила:
1. Не придумывай магазины и факты.
2. Не используй нерелевантные источники.
3. Указывай номера источников: [1], [2].
4. Не рекомендуй магазин, если недостаточно данных.
5. Не добавляй информацию из памяти модели.

Вопрос пользователя:
{question}

Источники:
{context}

Сформируй краткий список подходящих магазинов
с объяснением выбора и ссылками на номера источников.
"""

    response = gemini.models.generate_content(
        model="gemini-3.5-flash-lite",
        contents=prompt
    )

    return response.text or ""
