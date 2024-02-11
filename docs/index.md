# Description

- This project is a discord bot **Darkbot3** for open source game community [Freelancer Discovery](https://discoverygc.com/)
- It underwent a major refactorization and now reimplemented in golang with clean architecture for code scalability to add new features
- general stability added, because all errors are now handled way better with Golang approach that handles most of them at compile time.
- project saves ready for deployment docker images [at docker hub](https://hub.docker.com/repository/docker/darkwind8/darkbot/general), see [settings here](https://github.com/darklab8/fl-darkbot/blob/master/infra/kubernetes/charts/darkbot/templates/main.yml)
- darkbot has its own server now only for itself

User connects darkbot to some discord channel, and sets settings which space bases, player tags or space systems to track.
Darkbot repeatedly updates information to discord channel

# How to get started

- invite both to server [by link](https://discord.com/api/oauth2/authorize?client_id=838460303581904949&permissions=8&scope=bot)
- You must be Server owner or having `bot_controller` role in order to command the bot.
- add to some channel by writing `. connect`
- get help which commands are available by `. --help` or requesting help on sub commands `. base --help`
- add base tag for tracking `. base tags add Research Station`
- confirm it was added `. base tags list`
- in around 20 seconds you should see rendered and constantly updated view at this channel

![](index_assets/base_render.png)

- remove tag by `. base remove Research Station` or by `. base tags clear`

See other documentation in [Darkbot forum posts](https://discoverygc.com/forums/showthread.php?tid=188040)
- it has documented new appeared features and commands

# Features

- Adding player bases for tracking in discord channel
- Adding players for tracking based on tag in nickname or star systems / regions to track
- Adding configurable alert triggers to base, player of forum tracking
- Adding forum post for tracking

# Project status

- Finished its core development of features
    - Still processes minor feature requests to implement
- May have some issues which will be resolved as soon as possible
- Announcement about new releases will be made into Discord channel where it joined and here.

# Possible future plans

- Adding resource tracking at player bases (blocked)
  - Requires API access, which is lacked for implementation
- Adding some useful help scripts to calculate analytics

# Acknowledgements

- Freelancer Discovery API provided by [Alex](https://github.com/dsyalex) as a way to deliver info from Flhook to the bot
- [Pobbot](https://github.com/dr-lameos/Pobbot) originally made by dr.lameos sparkled this project
- Forum tracking is inspired by Biqqles project [forumlancer](https://github.com/biqqles/forumlancer)

# Contacts

- [join Darklab discord server](https://discord.gg/zFzSs82y3W)
- [write to Discovery forum account](https://discoverygc.com/forums/member.php?action=profile&uid=42166)
- [or write to email dark.dreamflyer@gmail.com]
- [Github repository of the project](https://github.com/darklab8/fl-darkbot)