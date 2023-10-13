#!/bin/zsh

setup_color() {
  RED=$(printf '\033[31m')
  GREEN=$(printf '\033[32m')
  YELLOW=$(printf '\033[33m')
  BOLD=$(printf '\033[1m')
  RESET=$(printf '\033[0m')
}

command_exists() {
  command -v "$@" >/dev/null 2>&1
}

install_brew() {

  if ! command_exists brew; then
    echo "${YELLOW}brew is not installed.${RESET} Do you want to install brew? [Y/n]"

    read -r opt
    case ${opt:u} in
    Y* | "") /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)" ;;
    N*)
      echo "brew is required. Quitting..."
      return
      ;;
    *)
      echo "Invalid choice. Quitting..."
      return
      ;;
    esac

    exit 1
  fi
}

get_latest_release() {
  curl --silent "https://api.github.com/repos/raxigan/pcfy-my-mac/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

run_latest_release() {
  latest=$(get_latest_release)
  arch=$(uname -m)

  if [[ "$arch" == "x86_64" ]]; then
    arch="amd64"
  fi

  url="https://github.com/raxigan/pcfy-my-mac/releases/download/${latest}/pcfy-my-mac-${latest}-darwin-${arch}.tar.gz"

  echo "Installing ${latest} release..."

  curl -L $url | tar xz
  chmod +x pcfy-my-mac
  clear
  ./pcfy-my-mac
}

main() {

  setup_color
  install_brew

  if [ -f "pcfy.go" ]; then
    echo "Found pcfy.go file, running from sources..."
    go build -ldflags '-w'

    if [ $? -eq 1 ]; then
      echo "Compilation error"
      exit 1
    fi

    ./pcfy-my-mac "$@"
  else
    run_latest_release
  fi
}

main "$@"
