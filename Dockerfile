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
    bash-completion \
    && rm -rf /var/lib/apt/lists/*

FROM base AS branch-shared-env-prod

FROM branch-shared-env-${ENVIRONMENT} as branch-app-scrappy

FROM branch-shared-env-${ENVIRONMENT} as branch-app-listener

FROM branch-shared-env-${ENVIRONMENT} as branch-app-discorder

FROM branch-shared-env-${ENVIRONMENT} as branch-app-configurator

FROM branch-shared-env-${ENVIRONMENT} as branch-app-viewer

FROM branch-app-${SERVICE} AS final
ARG SERVICE
ARG ENVIRONMENT

COPY ${SERVICE}/requirements.txt ${SERVICE}/requirements.dev.txt ${SERVICE}/constraints.txt ${SERVICE}/
COPY make.py ./

RUN python3 make.py shell install --environment=${ENVIRONMENT} --app=${SERVICE}

COPY ${SERVICE} /code/${SERVICE}
COPY utils /code/utils

COPY pytest.ini conftest.py make.py ./

FROM final as final-listener

COPY --from=final /code /code
COPY consoler /code/consoler
COPY configurator /code/configurator
COPY docker/.bash_profile /install/
RUN cat /install/.bash_profile >> /etc/profile

