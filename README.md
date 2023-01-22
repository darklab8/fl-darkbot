# Badges


| [Gitlab Pipeline Status](https://gitlab.com/darklab2/darklab_darkbot)                                               | Coverage                                                                                      |
| --------------------------------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------- |
| ![](https://gitlab.com/darklab2/darklab_darkbot/badges/master/pipeline.svg?key_text=GitlabCIPipeline&key_width=150) | ![](https://gitlab.com/darklab2/darklab_darkbot/badges/master/coverage.svg?key_text=Coverage) |

# Name

Darkbot 2.0

**[Invite Link](https://discord.com/api/oauth2/authorize?client_id=838460303581904949&permissions=8&scope=bothttps:/)**

# Description

This project is a discord bot for open source game community [Freelancer Discovery](https://discoverygc.com/)

![](docs/_images/general.png)

User connects darkbot to some discord channel, and sets settings which space bases, player tags or space systems to track.
Darkbot repeatedly updates information to discord channel

The project has 4 microservice parts:

- scrappy scraps third party REST APIs and web forum for data. Stores in database. Outputs at REST API or gRPC endpoints.
- configurator is a REST API or gRPC application. That stores all possible user settings
- listener is Discord API connected application, that accepts user commands from channel and sends to consoler, and gives from it answer
- consoler is a CLI interface to accept input from listener and render it CLI answer back
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
- python3 make.py {service_name} # to discovery available actions to service

### examples:

- python3 make.py scrappy shell # Running dev env
- python3 make.py scrappy test # Running tests

# Goals to follow:

- having it clean architectured as possible, with porto architecture main guidance to how it should be looking https://github.com/Mahmoudz/Porto
- testing coverage should be no less than 80%, ideally 90%+
- following different OOP principles, like having minimum exposed interfaces to rest of a code
- the code should be as obvious in its couplings as possible.
- main purpose of a current architecture to have it horizontally scalable(each part of architecture should be stateless) + easily extendable in features because everything relevant is gathered in one places, while having maximum isolation from Discord API as really bad fragile dependency
- closest goals to reach MVP (minimum viable product), and after that having a full sweep with static typization through mypy
- currently project also undergoes refactorization into asyncronous code
- reaching code readability mainly through well chosen function/class/variable names. If it is not enough, then resorting to docstrings, if it is not enough then resorting to comments
