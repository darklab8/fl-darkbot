Player Base Tracking
######################

**base commands**

Those are command to track status of player bases!
Example of output below.

check `the list <https://discoverygc.com/forums/bases.php>`_
of precised based names before using the commands

base add
***********

possible usages:

.. code-block::

    .base add MyStation

if your base have **spaces** in its name, better be using:

.. code-block::

    .base add "My Very Long Named Station"     

you can add multiple bases in one command

. code-block::

    .base add "My station #1", "My station #2", "My station #3"

also, you could add only sub string that the list of bases have

. code-block::

    .base add Depot

it will render all bases having Depot in its name

.. code-block:: JSON

    ["48.0885"]{"Copper Storage Depot"}["Freelancers"]
    ["100"]{"Wismar Shipping Depot"}["Kruger Minerals"]
    ["83.3811"]{"Bristol Depot"}["Freelancers"]
    ["100"]{"Shiojiri Storage Depot"}["Samura Industries"]
    ["4.79464e-26"]{"Malfunctioning Depot"}["No Affiliation"]
    ["54.6548"]{"Howler Depot"}["Junkers"]
    ["56.3603"]{"Aruba Depot"}["Freelancers"]

additional commands

base remove
**************

for deleting of one base names

. code-block::

    .base remove Depot

base clear
*************

to clear the whole list of bases

. code-block::

    .base clear
