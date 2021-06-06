FROM python:3.8-slim

ENV PYTHONUNBUFFERED 1
ENV PYTHONDONTWRITEBYTECODE 1

COPY ./requirements.txt ./
RUN pip install -r requirements.txt

COPY . .

CMD python3 app.py
