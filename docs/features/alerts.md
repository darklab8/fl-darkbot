# Alerts

![](/fl-darkbot/index_assets/alert_example.png)

- `. alert` or `. alert --help` helps to find out all avilable alerts. Continue using `--help` for discovered sub commands, to find out all options.
![](/fl-darkbot/index_assets/alerts_commands.png)

## Base alerts

### Base health below threshold

- `. alert base_health_is_lower_than set 90` - the most convinient command to receive alerts about. It pings you if any base health below 90% in the current configuration.
- `. alert base_health_is_lower_than status` to check current configuration
- `. alert base_health_is_lower_than unset` to disable alert

### Base under attack

- `. alert base_is_under_attack enable` - turns on alert if spotting the base in forum thread for base attack declarations.
enable, disable, status sub commands are similar.

### Base Health decreasal

![](/fl-darkbot/index_assets/alerts_base1_commands.png)

- you can turn on `. alert base_health_is_decreasing enable` command for making alert if base loses its health.

## Player alerts

### Enemy player alert

- `. alert player_enemy_count_above set 1` - Sets alert to ping if spotting enemies above this count in tracked regions/systems.
- `. alert player_enemy_count_above status` to check status of configuration
- `. alert player_enemy_count_above unset` to unset the alert

### Neutral player alert

- `. player_neutral_count_above set 1` - Sets alert if finding more than X neutral players in tracked systems, regions. Convinient command if not knowing who you track, but knowing where.
- commands set, status and unset are similar to enemy alert configuration

### Friends alert

- `. alert player_friend_count_above set 0` - if your friends login, you can spot them right away across galaxy. This command is not affected by which regions/systems you track.
- commands set, status and unset are similar to enemy alert configuration.

## Ping message

- By default Discord server owner is pinged ( `<@DiscordServer.Owner.ID>` )
- `. alert ping_message set @here` - You can change to `@here` or `@specific_role` by setting this
- `. alert ping_message status` - you can use sub command `status` to check configuration like 
- `. alert ping_message unset` - you can unset to default ping message
