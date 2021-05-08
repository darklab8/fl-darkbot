Begining
#####################

At first you have only four commands available

**wiki** and **help** available for everyone

wiki
******

gives link to the current documentation

.. code-block::

    .wiki

help
*****

| shows available command for your current user permissions
| it is changed also based on bot being **connect** ed to channel or not

possible usage:

list all root level commands:

.. code-block::

    .help

list commands from sub category beloning to base for example:

.. code-block::

    .help base

connect
*********

| initializes bot work for the current channel
| **attention!**: this command and all further commands require having user permissions

* being a server owner
* or having permission to manage channels
* or having role **bot_controller**

**warning**: default setting enables auto erasing of all messages
older than 40 seconds in the connected channel!

.. code-block::

    .connect

disconnect
*************

disconnects the bot from channel. It could be also used for mass setting
nullification. Disconnect the bot and wait 10 seconds in order for all your
current channel settings being gone.

.. code-block::

    .disconnect
