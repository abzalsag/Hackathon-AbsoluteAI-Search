# Python негізгі бейнесі
FROM python:3.11-slim

# Жұмыс директориясын орнату
WORKDIR /app

# Тәуелділіктерді көшіру және орнату
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Барлық Python кодты көшіру
COPY . .

# Портты ашу (сенің .env файлыңда 8000 тұр)
EXPOSE 8000

# FastAPI қосымшасын іске қосу
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000", "--reload"]