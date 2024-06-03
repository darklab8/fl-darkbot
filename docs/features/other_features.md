# Other features

### Msg splitting

- bot is smart enough to split content between multiple discord msgs if it can't fit into one
![msg split](https://raw.githubusercontent.com/darklab8/fl-darkbot/master/docs/index_assets/bases_split_msgs.png)

### Check commands usage

- `. base` or `. base --help` will invoke how to use the command
- `. base order_by --help` or `. base order_by`. You can navigate around which commands are available in bot interactively. By finding next sub commands and checking how to use them.

### Auto msg cleanup

- darkbot auto cleans all not related msgs in a minute of time at the channel it is using
    - which u chose with `. connect` command.
- to stop msg deletion (and bot activity), just use `. disconnect` in the same channel
