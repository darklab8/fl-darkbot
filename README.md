# Name

Darkbot 2.0

# Description

This project is a discord bot for open source game community [Freelancer Discovery](https://discoverygc.com/)

![](docs/_images/general.png)

User connects darkbot to some discord channel, and sets settings which space bases, player tags or space systems to track.
Darkbot repeatedly updates information to discord channel

The project has 4 microservice parts:
- scrappy scraps third party REST APIs and web forum for data. Stores in database. Outputs at REST API or gRPC endpoints.
- configurator is a REST API or gRPC application. That stores all possible user settings
- listener is Discord API connected application, that accepts user commands from channel and sends to configurator
- viewer is main celery-beat application that in a loop gets user settings and connected channels, and rerenders to them view

[First version of project had wiki](https://darklab8.github.io/darklab_darkbot/)

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

```
Warning: docker-compose needs to be installed
Warning: python3.10 is supposed to be installed
```

- python3 make.py {service_name} {action_name}
- python3 make.py --help # to discovery other available services
- python3 make.py {service_name} # to discovery available actions to server

### examples:

- python3 make.py scrappy shell # Running dev env
- python3 make.py scrappy test # Running tests


