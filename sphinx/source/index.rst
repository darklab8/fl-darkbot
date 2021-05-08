.. darkbot documentation master file, created by
   sphinx-quickstart on Thu May  6 20:10:46 2021.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Welcome to darkbot's documentation!
===================================

.. toctree::
   :maxdepth: 3
   :caption: Contents:

   info/index.rst

Getting started
******************

* invite the bot with `deployment invite link <https://discord.com/api/oauth2/authorize?client_id=838460303581904949&permissions=4294967287&scope=bot%20applications.commands>`_
* look further instructions at the :ref:`Begining`

Example of the working bot
***************************
.. image:: images/general.png

Current features
******************

* Tracking bases health / name / affiliation
* Tracking players in chosen space regions and systems (trackable area)
* Distinguishing friends and enemies in trackable area
* Making @here alert in discord, if amount of players above thresholds (default: disabled)
* Separated alerts for unregocnized players, friends and enemies
* Auto erasing messages older than 40 seconds in connected channel
* Ability to clear all msgs manually

Plans for the future
************************

When `Alex <https://github.com/dsyalex>`_ would process `my code additions <https://github.com/DiscoveryGC/FLHook/pull/160>`_ for Flhook:

* Alerts for base attacks
* Showing which player bases have ore for sale or for buying
* Same for devices selling/buying at the player bases
* Same could be also used to track Food/Water/Oxygen at the bases

Technologies
************************

* Visual Studio Code as IDE
* `Python 3.8.5 <https://www.python.org/downloads/release/python-385/>`_
* `Discord async framework <https://discordpy.readthedocs.io/en/stable/index.html>`_
* `Pytest for unit testing, coverage <https://docs.pytest.org/en/6.2.x/>`_
* Pylint/Flake8 for linting, Yapf for style formatting
* `Sphinx for documentation <https://www.sphinx-doc.org/en/master/contents.html>`_
* `VPS with Ubuntu 20.04 LTS Server for deployment <https://releases.ubuntu.com/20.04/>`_
* `Supervisor for running in background with autorestart <http://supervisord.org/>`_
* `FLhook <https://github.com/DiscoveryGC/FLHook>`_ as source of base and player data
* Freelancer Discovery API provided by `Alex <https://github.com/dsyalex>`_ as a way to deliver info from Flhook to the bot

Developers
************************

* `dd84ai <https://github.com/dd84ai>`_ a.k.a `darkwind <https://discoverygc.com/forums/showthread.php?tid=185978>`_ (Python / Testing / Deployment)
* `Alex <https://github.com/dsyalex>`_ (FLhook and API)

Github
************************

`Can be found here <https://github.com/dd84ai/darkbot>`_

Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`
