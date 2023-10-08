# PC'fy my mac

<p>
<a href="https://github.com/raxigan/pcfy-my-mac/actions">
    <img src="https://github.com/raxigan/pcfy-my-mac/actions/workflows/go.yml/badge.svg?branch=feature/installation_script" alt="Build Status">
</a>
<a href="https://github.com/raxigan/pcfy-my-mac/releases">
    <img src="https://img.shields.io/github/release/raxigan/pcfy-my-mac.svg" alt="Latest Release">
</a>
</p>

All-in-one project to help you to get PC-like experience (known from Windows or Linux systems) on your macOS. Mostly for tech guys, but everyone is welcome.
What you get:
- System and Browser (Chromium-based) shortcuts
- PC keymaps for JetBrains tools
- Blazingly-fast application launching with a single Win/Opt key
- Basic window snapping with Win/Opt + arrow keys (Snap left/right, maximize)
- Use keyboard shortcuts instead of gestures with new windows switcher
- **Everything works on any keyboard layout (you can use built in mac and external PC in same time)**

<img src="docs/demo.gif" alt="demo" width="100%"/>

## Installation

### [Homebrew](https://brew.sh/)

```shell
brew install raxigan/tap/pcfy-my-mac
```

### Script
Requires [special terminal permissions](#Terminal-dev-permissions)
```shell
curl -s https://raw.githubusercontent.com/raxigan/pcfy-my-mac/main/pcfy.sh | sudo bash
```

### Go binary
```shell
go install github.com/raxigan/pcfy-my-mac@latest 
```

### From sources
```shell
go run pcfy.go
```

### Terminal dev permissions
**Necessary only if you install it by shell script**. The binary is not signed, so macOS won't let you run it without the following permissions
for your terminal. Just go to *System Settings* > *Privacy & Security* > *Developer Tools* and enable it:

![terminal_permissions.png](docs/terminal_permissions.png)

## Troubleshooting
TBD

## Missing things:
- Finder and Fleet keymaps are incomplete
- Select files using Opt+LMB instead of Ctrl like you would on PC
- Multicursor shortcut (2xCtrl in Jetbrains tools on PC) is under 2xOpt
- There's no Alt/Cmd+F4, use Win/Opt+Q instead - it's easy to add (I'm kind of used to it though)

