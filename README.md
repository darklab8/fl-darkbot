# Name

Darkbot 3.0

**[Invite Link](https://discord.com/api/oauth2/authorize?client_id=838460303581904949&permissions=8&scope=bothttps:/)**

# Description

This project is a discord bot for open source game community [Freelancer Discovery](https://discoverygc.com/)

![](docs/_images/general.png)

User connects darkbot to some discord channel, and sets settings which space bases, player tags or space systems to track.
Darkbot repeatedly updates information to discord channel

The project has 5 package parts parts:

- scrappy scraps third party REST APIs and web forum for data. Stores in in storage (memory) to other modules
- configurator is interacting with SQL database to store user settings
- listener is Discord API connected application, that accepts user commands from channel and sends to consoler, and gives from it answer
- consoler is a CLI interface to accept input from listener and render its CLI answer back
- viewer is application that in a loop gets user settings and connected channels, and rerenders to them view based on available data in scrappy and configurator settings

[First version of project had wiki](https://darklab8.github.io/darklab_darkbot/)

# Project status:

Having first public release and still in active development.
Future plans:
- Adding player tracking
- Adding forum posts tracking
- Other stuff which was present in first darkbot version
- Hmm. may be implementing my old dijkstra stuff and other helping scripts as part of its functionality. (note for myself)

# Tech stack

- golang
- discordgo
- cobra-cli
- gorm (sqlite3)
- docker
- terraform (hetzner / helm)
- microk8s (helm)

# Dev commands

install taskfile.dev for dev commands

go run . --help

task test # to test
task --list-all # to list other commands 

# Goals to follow:

- having it clean architectured as possible
- testing coverage should be no less than 80%, ideally 90%+
- following different OOP principles, like having minimum exposed interfaces to rest of a code

# Architecture of the project

![architecture](architecture/architecture.drawio.svg)