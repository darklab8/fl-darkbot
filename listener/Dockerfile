FROM python:3.10.5-slim-buster

ENV PYTHONUNBUFFERED 1
ENV PYTHONDONTWRITEBYTECODE 1

WORKDIR /code

COPY listener/requirements.txt listener/constraints.txt listener/
RUN pip install -r listener/requirements.txt -c listener/constraints.txt 

COPY listener /code/listener
COPY utils /code/utils
COPY pytest.ini conftest.py make.py ./