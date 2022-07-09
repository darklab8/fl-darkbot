Player Tracking
#######################

in order to start tracking players in area
you need to set systems or regions for tracking

* for a precise list of names `look here <https://discoverygc.com/forums/api_interface.php?action=players_online>`_

Examples of commands:

system add
**************

same rules as to base commands, apply to all **add**, **remove** commands
regarding multiple adding

.. code-block::

    .system add "New York"

multiple additions:

.. code-block::

    .system add "New York", "New London"

system remove
***************

. code-block::

    .system remove "New York", "New London"

system clear
***************

clears all added systems

. code-block::

    .system clear

region add
**************

together with system adding, you could also select additional systems with whole region

.. code-block::

    .region add "Liberty Space"

.. code-block::

    .region add "Tau Border Worlds"

region remove
**************

same as in base commands

. code-block::

    .region remove "Liberty Space"

region clear
**************

same as in base commands

. code-block::

    .region clear