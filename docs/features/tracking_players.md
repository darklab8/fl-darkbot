# Player tracking (Deprecated)

![player view](/fl-darkbot/index_assets/player_render2.png)

### Rules

- Added Systems and Regions will show player in Neutral tab.
- Enemy players (see commands below) will be watched only in added systems or regions
- Player Friends will be watched across entire galaxy regardless of configured other stuff.

### System tracking

- `. player system add New York` - add system for tracking (You can add it as just partial name `York`, it should be enough too.)
- `. player system remove New York` - remove specific one from tracking
- `. player system list` -  check tracked systems
- `. player system clear` - clear all tracked systems

### Region tracking

- `. player region add Liberty` - to track entire liberty
add, clear, list, remove commands are similar to system tracking.

- Check [Players online at forum](https://discoverygc.com/forums/api_interface.php?action=players_online) to see valid Region names

### Enemy tracking.

Enemy added tags will be matched only for ships in tracked Systems/Regions!.

- `. player enemy add [FactionTag]` - Add factin prefixes (or suffixes)
- `. player enemy add PartialName` - or partial or any full name.

list, remove, add, clear commands are similar to other command groups

### Friends tracking

- `. player friend add My Friend` - for tracking friends acorss entire galaxy

list, remove, add clear commands are similar to other command groups

##  Player alerts

###  Enemy player alert

- `. alert player_enemy_count_above set 1` - Sets alert to ping if spotting enemies above this count in tracked regions/systems.
- `. alert player_enemy_count_above status` to check status of configuration
- `. alert player_enemy_count_above unset` to unset the alert

###  Neutral player alert

- `. player_neutral_count_above set 1` - Sets alert if finding more than X neutral players in tracked systems, regions. Convinient command if not knowing who you track, but knowing where.
- commands set, status and unset are similar to enemy alert configuration

### Friends alert

- `. alert player_friend_count_above set 0` - if your friends login, you can spot them right away across galaxy. This command is not affected by which regions/systems you track.
- commands set, status and unset are similar to enemy alert configuration.

# Event helper

![events view](/fl-darkbot/index_assets/events_table.png)

This feature was implemented on request from Event manager Barrier, for the purpose of easily tracking
that all players participating in event, arrived to the scene.

- `. player event add PlayerPreffixOrSuffix` - add player for tracking for event
- `. player event add Any Partial Name` - add player for tracking for event

- `. player event remove Any Partial Name` - remove specific one from tracking

- `. player event list` -  check tracked players
- `. player event clear` - clear all tracked players

P.S. ensure you used `. connect` first one some channel to turn one bot on.
Beware. bot automatically erases msgs for the channel for comfort of its usage.
Use bot on dedicated channel.
