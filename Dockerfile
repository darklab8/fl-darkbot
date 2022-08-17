ARG ENVIRONMENT
ARG SERVICE

FROM python:3.10.5-slim-buster as base

ENV PYTHONUNBUFFERED 1
ENV PYTHONDONTWRITEBYTECODE 1

WORKDIR /install

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
ARG ENVIRONMENT

COPY ${SERVICE}/requirements.txt ${SERVICE}/requirements.dev.txt ${SERVICE}/constraints.txt ${SERVICE}/
COPY make.py ./
RUN python3 make.py shell install --environment=${ENVIRONMENT} --app=${SERVICE}

COPY ${SERVICE} /code/${SERVICE}
COPY utils /code/utils
# consoler needs to be added for listener. Rewrite dockerfile si it would not be added anywhere else
COPY consoler /code/consoler 
COPY pytest.ini conftest.py make.py ./