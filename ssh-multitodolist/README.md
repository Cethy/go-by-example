---
Order: 12
Title: SSHable Multi-user Todolist 
#Dependencies: ["./demo.gif"]
ImgSrc: https://images.unsplash.com/photo-1591439657848-9f4b9ce436b9?ixid=M3w2NjYzMTJ8MHwxfHJhbmRvbXx8fHx8fHx8fDE3Mjk0NTI5MDF8&ixlib=rb-4.0.3
---

# SSHable Multi-user Todolist

## Instructions

Build upon previous [Multi todo list project](./cli-multitodolist.html)
and make it reachable via ssh and multi-user

![Made with VHS](./demo.gif)

## Key Features

- multi-user ssh server
- instant data sharing
- see what the other users do (cursor & active tab)
- talk to other users (sidebar chat)
- standalone mode (with multi-users UIs disabled)
- multi-room(/files) setup
- optional redis storage
- private rooms (with invite code to share access)

## Usage

```shell
# ssh server
[PORT=23234] go run main.go server
# ssh server (redis)
# docker run -d --name redis-todo -p 6379:6379 redis
[REDIS_ADDR="localhost:6379"] [REDIS_PASSWORD=""] [PORT=23234] go run main.go server [--db="file"|"redis"]
# connect to server
ssh -p23234 -t localhost [room] [privacy<true|1>] [inviteCode]

# standalone
go run main.go standalone [room] [--db="file"|"redis"]
```

## TODO

- [ ] store/persist room data (private settings, users & chat historic)
- [ ] https landing page ("go to ssh...")
- [ ] header horizontal viewport
- [ ] ctrl-enter behavior (submit and open input again)
- [ ] refactor with lipgloss.SetDefaultRenderer() ?
- [ ] migrate to bubbletea v2
- [ ] alternate kanban view
- [ ] special access to list the rooms
- [ ] cleanup empty rooms (?)
