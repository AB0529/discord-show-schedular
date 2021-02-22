# Discord Anime Scheduler

A Discord bot written in Go with the purpose of notifying users about certain shows (anime) and when new episodes of them air.

---

# Table of Contents
1. [Commands](#Commands)
1. [Build from Source](#Build)
    1. [Running](#Running)

___

# Commands
- `schedule`
    - Handles creating, removing, and listing shows from user's schedule
    - Flags:
        - `add` adds a show to the schedule
        - `list` lists current users schedule
        - `remove` deletes a show from the schedule
 - `ping`
    - Generic ping pong command
___
    
# Build

### Prerequisites
```
go >= v1.15.x
```

### Running
```shell
git clone https://github.com/AB0529/discord-show-schedular
cd discord-show-schedular/src
go get
go build -o bot
```

```shell
# Make it executable
chmod +x ./bot
# Run it
./bot
```