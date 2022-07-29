# Project status:
Old code is declared as old, and will be refacorized fully to new version with considering all gathered experience in this application
And having this time more reliable and scalable archirecture. Scalable in a code growth and workload performance.

# Architecture of the project
![](architecture/architecture.drawio.svg)

# Tech stack
- FastAPI
- Celery
- Kafka
- Click
- Helm
- Docker-compose for dev env
- Black (auto format on save)

# Commands

Warning: docker-compose needs to be installed
Warning: python3.10 is supposed to be installed

- python3 make.py {service_name} {action_name}
- python3 make.py --help # to discovery other available services and commands

### examples:

- python3 make.py scrappy shell # Running dev env
- python3 make.py scrappy test # Running tests


