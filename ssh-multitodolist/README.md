---
Order: 12
Title: SSHable Multi Todo List 
#Dependencies: ["./demo.gif"]
ImgSrc: https://images.unsplash.com/photo-1591439657848-9f4b9ce436b9?ixid=M3w2NjYzMTJ8MHwxfHJhbmRvbXx8fHx8fHx8fDE3Mjk0NTI5MDF8&ixlib=rb-4.0.3
---

# SSHable multi todo list

## Instructions

Build upon previous [Multi todo list project](./cli-multitodolist.html)
and make it reachable via ssh

![Made with VHS](./demo.gif)

## Key Features

- multi-user ssh server
- instant data sharing
- see what the other users do
- standalone mode

## Usage

```shell
# ssh server
go run main.go server [--port=23234]

# standalone
go run main.go standalone
```

## TODO

- [ ] sidebar chat ([example](https://github.com/charmbracelet/wish/blob/main/examples/multichat/main.go))
- [ ] authentication
- [ ] new todo by url
- [ ] https landing page (go to ssh...)
- [ ] redis repository
- [ ] source filename flag
- [ ] header viewport
- [ ] shift-enter behavior (submit and open input again)
- [ ] refactor with lipgloss.SetDefaultRenderer() ?
- [ ] migrate to bubbletea v2
- [ ] better help grouping
