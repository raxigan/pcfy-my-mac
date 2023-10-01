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
      echo "brew is required. Exiting."
      return
      ;;
    *)
      echo "Invalid choice. Shell change skipped."
      return
      ;;
    esac

    exit 1
  else
    echo "${GREEN}brew installed${RESET}"
  fi
}

main() {
  clear
  setup_color
  install_brew

  if [ -f "main.go" ]; then
    echo "Running local..."
    go build -ldflags "-w"

    if [ $? -eq 1 ]; then
      echo "Compilation error"
      exit 1
    fi

    ./pcfy-my-mac "$@"
  else
    echo "Running from remote package..."
#     get it from releases packages instead
    curl -fsSL -o my_binary https://raw.githubusercontent.com/raxigan/macos-pc-mode/feature/test-golang/tools/macos-pc-mode && chmod +x my_binary && ./my_binary && rm my_binary
  fi
}

main "$@"