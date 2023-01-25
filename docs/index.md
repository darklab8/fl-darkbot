# Description

- This project is a discord bot **Darkbot 3.0** for open source game community [Freelancer Discovery](https://discoverygc.com/)
- It underwent a major refactorization and now reimplemented in golang with clean architecture for code scalability to add new features
- and for general stability, because all errors are now handled way better with Golang approach that handles most of them at compile time.

User connects darkbot to some discord channel, and sets settings which space bases, player tags or space systems to track.
Darkbot repeatedly updates information to discord channel

# Features

- Adding player bases for tracking in discord channel

# Project status

- in active development until at least all planned features are added from the list "Comin soon" and "For later"
- May have some issues which will be resolved as soon as possible
- May went through database wipes during deployments of new releases. Announcement about new releases will be made into Discord channel where it joined and here.
# Future plans: Comming Soon

- Adding players for tracking based on tag in nickname or star systems / regions to track
- Adding configurable alert triggers to base and player tracking

# Future plans: For later

- Adding forum post for tracking

# Future plans: For more later

- Adding resource tracking at player bases
- Adding some useful help scripts to calculate analytics

# How to get started

- invite both to server [by link](https://discord.com/api/oauth2/authorize?client_id=838460303581904949&permissions=8&scope=bot)
- add to some channel by writing `. connect`
- get help which commands are available by `. --help` or requesting help on sub commands `. base --help`
- add base tag for tracking `. base add Research Station`
- confirm it was added `. base list`
- in around 20 seconds you should see rendered and constantly updated view at this channel

![](index_assets/base_render.png)

- remove tag by `. base remove Research Station` or by `. base clear`

# Acknowledgements

- Freelancer Discovery API provided by [Alex](https://github.com/dsyalex) as a way to deliver info from Flhook to the bot
- Pobbot originally made by dr.lameos sparkled this project


# Contacts

- [join Darklab discord server](https://discord.gg/zFzSs82y3W)
- [write to Discovery forum account](https://discoverygc.com/forums/member.php?action=profile&uid=42166)
- [or write to email dark.dreamflyer@gmail.com]