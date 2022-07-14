FROM python:3.10.5-slim-buster

ENV PYTHONUNBUFFERED 1
ENV PYTHONDONTWRITEBYTECODE 1

WORKDIR /code

COPY scrappy/requirements.txt scrappy/constraints.txt ./
RUN pip install -r requirements.txt -c constraints.txt

COPY scrappy /code/scrappy
COPY utils /code/utils