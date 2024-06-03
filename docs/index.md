# Darkbot - Introduction

![logo](./index_assets/fulllogo.png)

- It implements Discord bot to track player bases, players themselves and forum posting with notifications to Discord.
- The logic of it: User connects darkbot to some discord channel, and sets settings which space bases, player tags or space systems to track. Darkbot repeatedly updates information to discord channel

# Features

- Adding player bases for tracking in discord channel
- Adding players for tracking based on tag in nickname or star systems / regions to track
- Adding configurable alert triggers to base, player of forum tracking
- Adding forum post for tracking

# Important links

- [Documentation for the bot, how to use it](https://darklab8.github.io/fl-darkbot/)
- [Discovery Forum Thread and anouncements](https://discoverygc.com/forums/showthread.php?tid=188040)
- [Github of the project](https://github.com/darklab8/fl-darkbot)
- [See other Dark Tools for Freelancer here](https://darklab8.github.io/blog/pet_projects.html#DiscoveryFreelancercommunity)

# How to get started

- invite both to server [by link](https://discord.com/api/oauth2/authorize?client_id=838460303581904949&permissions=207952&scope=bot)
- You must be Server owner or having `bot_controller` role in order to command the bot.
- add to some channel by writing `. connect` (if u wish to disconnect bot from channel, write `.disconnect`)
- get help which commands are available by `. --help` or requesting help on sub commands `. base --help`
- add base tag for tracking `. base tags add Research Station`
- confirm it was added `. base tags list`
- in around 20 seconds you should see rendered and constantly updated view at this channel

![](index_assets/base_render.png)

- remove tag by `. base remove Research Station` or by `. base tags clear`

Continue with [documentation there](https://darklab8.github.io/fl-darkbot/)

# Permissions for running

- Bot expects having sufficient permissions for its working, which is seeing the channel it works in, writing, editing, deleting msgs and reading msg history. [The default link](https://discord.com/api/oauth2/authorize?client_id=838460303581904949&permissions=207952&scope=bot) configured for permissions that should be enough in theory:
    - Read Messages
    - Send Messages
    - Manage Messages
    - Read Message History
    - Mention everyone here and all roles
    - Manage channels (Probably not needed permission)
    - Add reactions (Just added for future if ever will be needed)

- If default inviting link is not good enough, [use invite link as admin](https://discord.com/api/oauth2/authorize?client_id=838460303581904949&permissions=8&scope=bot)

# Development specific stuff

See [Github Readme.md](<https://github.com/darklab8/fl-darkbot/blob/master/README.md>) for details

# Contacts

- [join Darklab discord server](https://discord.gg/zFzSs82y3W)
- [write to Discovery forum account](https://discoverygc.com/forums/member.php?action=profile&uid=42166)
- or write to email `dark.dreamflyer@gmail.com`
- [Optionally open Github Issue at repository](https://github.com/darklab8/fl-darkbot)
