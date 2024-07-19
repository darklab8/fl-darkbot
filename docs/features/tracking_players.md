# Player tracking

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
