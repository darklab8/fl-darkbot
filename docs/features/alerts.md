# Alerts

![](/fl-darkbot/index_assets/alert_example.png)

- `. alert` or `. alert --help` helps to find out all avilable alerts. Continue using `--help` for discovered sub commands, to find out all options.
![](/fl-darkbot/index_assets/alerts_commands.png)

## Base alerts

### Base health below threshold

- `. alert base_health_is_lower_than set 90` - the most convinient command to receive alerts about. It pings you if any base health below 90% in the current configuration.
- `. alert base_health_is_lower_than status` to check current configuration
- `. alert base_health_is_lower_than unset` to disable alert

### Base money below threshold

- `. alert base_money_is_lower_than set 5000` - Set threshold of base money, below which you will receive alert
- `. alert base_money_is_lower_than status` to check current configuration
- `. alert base_money_is_lower_than unset` to disable alert

### Base cargo below threshold

- `. alert base_cargo_space_left_is_lower_than set 5000` - Set threshold of base cargo space left, below which you will receive alert
- `. alert base_cargo_space_left_is_lower_than status` to check current configuration
- `. alert base_cargo_space_left_is_lower_than unset` to disable alert

### Base under attack

- `. alert base_is_under_attack enable` - turns on alert if spotting the base in forum thread for base attack declarations.
enable, disable, status sub commands are similar.

### Base Health decreasal

![](/fl-darkbot/index_assets/alerts_base1_commands.png)

- you can turn on `. alert base_health_is_decreasing enable` command for making alert if base loses its health.

## Ping message

- By default Discord server owner is pinged ( `<@DiscordServer.Owner.ID>` )
- `. alert ping_message set @here` - You can change to `@here` or `@specific_role` by setting this
- `. alert ping_message status` - you can use sub command `status` to check configuration like 
- `. alert ping_message unset` - you can unset to default ping message
