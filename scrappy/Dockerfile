FROM python:3.10.5-slim-buster

ENV PYTHONUNBUFFERED 1
ENV PYTHONDONTWRITEBYTECODE 1

WORKDIR /code

COPY scrappy/requirements.txt scrappy/constraints.txt scrappy/
RUN pip install -r scrappy/requirements.txt -c scrappy/constraints.txt 

COPY scrappy /code/scrappy
COPY utils /code/utils
COPY pytest.ini conftest.py make.py ./