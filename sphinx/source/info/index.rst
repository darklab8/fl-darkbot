additional info
===========================


`Deployment invite link <https://discord.com/api/oauth2/authorize?client_id=838460303581904949&permissions=8&scope=bot>`_


`Tester invite link <https://discord.com/api/oauth2/authorize?client_id=839849682977161216&permissions=8&scope=bot>`_

Permission levels
#################
the bot have three categories of permission levels

global:
*********
* command like **help** or **wiki** which are visible for everyone.

restricted to user rights:
***************************
the rest of commands require possessing at least one of the next permissions:

* being a server owner
* or having access to **manage channels**
* or having role **bot_controller**


meeting one of three criterias above, allows you to use commands:

* **connect** to enable bot for this channel
* **disconnect** to disable bot for this channel
  
**warning**: default mod auto erases messages older than 40 seconds in connected channel

restricted to user rights + connected to channel
******************************************************

the last category of commands requires having user permissions
and for the bot being connected to channel

so be sure to write **.help** to get updated list of available commands
based on your user permissions and bot connecting status
