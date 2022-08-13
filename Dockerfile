ARG ENVIRONMENT
ARG SERVICE

FROM python:3.10.5-slim-buster as base

ENV PYTHONUNBUFFERED 1
ENV PYTHONDONTWRITEBYTECODE 1

WORKDIR /code

FROM base AS branch-shared-env-dev

RUN apt update && apt install -y \ 
    git \
    && rm -rf /var/lib/apt/lists/*

FROM base AS branch-shared-env-prod

FROM branch-shared-env-${ENVIRONMENT} as branch-app-scrappy

FROM branch-shared-env-${ENVIRONMENT} as branch-app-listener

FROM branch-shared-env-${ENVIRONMENT} as branch-app-discorder

FROM branch-shared-env-${ENVIRONMENT} as branch-app-configurator

FROM branch-app-${SERVICE} AS final
ARG SERVICE

COPY ${SERVICE}/requirements.txt ${SERVICE}/constraints.txt ${SERVICE}/
RUN pip install -r ${SERVICE}/requirements.txt -c ${SERVICE}/constraints.txt 

COPY ${SERVICE} /code/${SERVICE}
COPY utils /code/utils
COPY listener /code/listener
COPY pytest.ini conftest.py make.py ./