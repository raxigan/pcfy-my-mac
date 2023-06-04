#!/bin/zsh

# TODO change
#BRANCH_NAME="main"
BRANCH_NAME="feature/installation_script"

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
    echo "${GREEN}brew installed"

  fi
}

install_jq() {
  if ! command_exists jq; then
    echo "${YELLOW}jq is not installed.${RESET} Do you want to install jq? [Y/n]"

    read -r opt
    case $opt in
    y* | Y* | "") brew install jq ;;
    n* | N*)
      echo "jq. returning."
      return
      ;;
    *)
      echo "Invalid choice. Shell change skipped."
      return
      ;;
    esac

    exit 1
  else
    echo "${GREEN}jq installed"
  fi
}

install_jq() {
  if ! command_exists jq; then
    echo "${YELLOW}jq is not installed.${RESET} Do you want to install jq? [Y/n]"

    read -r opt
    case $opt in
    y* | Y* | "") brew install jq ;;
    n* | N*)
      echo "jq. returning."
      return
      ;;
    *)
      echo "Invalid choice. Shell change skipped."
      return
      ;;
    esac

    exit 1
  else
    echo "${GREEN}jq installed"
  fi
}

PROFILE_NAME="PC mode"

KARABINER_CONFIG_DIR=~/.config/karabiner
KARABINER_CONFIG=$KARABINER_CONFIG_DIR/karabiner.json

apply_rules() {
  echo "$2"
  jq --arg PROFILE_NAME "$PROFILE_NAME" '(.profiles[] | select(.name == $PROFILE_NAME) | .complex_modifications.rules) += $rules[].rules' \
    $KARABINER_CONFIG --slurpfile rules ../karabiner-elements/"$1" --indent 4 >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG
}

prepare_for_mac_keyboard() {
  echo "$1"
  jq --arg PROFILE_NAME "$PROFILE_NAME" '.profiles |= map(if .name == $PROFILE_NAME then walk(if type == "object" and .conditions then del(.conditions[] | select(.identifiers[]?.is_built_in_keyboard)) else . end) else . end)' $KARABINER_CONFIG --indent 4 \
    >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG
}

install_karabiner() {
  if [ -f $KARABINER_CONFIG ]; then
    echo "${GREEN}Karabiner-Elements installed"
  else
    echo "${YELLOW}Karabiner-Elements is not installed.${RESET} Do you want to install karabiner-elements? [Y/n]"

    read -r opt
    case $opt in
    y* | Y* | "") brew install --cask karabiner-elements ;;
    n* | N*)
      echo "Karabiner-Elements required. returning."
      return
      ;;
    *)
      echo "Invalid choice. Shell change skipped."
      return
      ;;
    esac
  fi
}

install_ide_keymap() {

  IDE_NAME=$1
  IDE_FULL_NAME=$2
  IDE_LAUNCHER_SCRIPT=$3

  IDE_VERSION=$(grep <~/Library/Application' 'Support/JetBrains/Toolbox/scripts/"$IDE_LAUNCHER_SCRIPT" "$IDE_NAME" | cut -d/ -f11)
  IDE_CONFIG_DIR=""

  echo "Installing XWin plugin..."

  open -naj "$IDE_FULL_NAME" --args installPlugins com.intellij.plugins.xwinkeymap

  echo "Installing $IDE_NAME (ver. ${IDE_VERSION}) keymap ..."
  IJ_CONFIGS=~/Library/Application' 'Support/JetBrains

  for entry in ~/Library/Caches/Jetbrains/*; do

    version=$(find "$entry" -type d -name "*$IDE_NAME*" -exec grep "app.build.number=" {}/.appinfo \; | sed 's/.*app.build.number=\([^&]*\).*/\1/')

    if [ "$version" = "$IDE_VERSION" ]; then
      IDE_CONFIG_DIR=$(echo "$entry" | cut -d/ -f7)
      break
    fi
  done

  # if IDE_CONFIG_DIR empty then exit

  KEYMAPS_DIR=${IJ_CONFIGS}/${IDE_CONFIG_DIR}/keymaps
  KEYMAP_FILENAME=$(echo "${IDE_FULL_NAME:l}" | tr " " "-")

  curl --silent -o "${KEYMAPS_DIR}/${KEYMAP_FILENAME}.xml" https://raw.githubusercontent.com/raxigan/macos-pc-mode/$BRANCH_NAME/keymaps/"${KEYMAP_FILENAME}".xml

  echo "Restart $IDE_FULL_NAME. Then choose XWin $IDE_FULL_NAME in Preferences > Keymaps > Xwin"
}

quit() {
  echo "${BOLD}Quitting...${BOLD}"
  exit 0
}

main() {

  clear

  setup_color
  install_brew
  install_jq
  install_karabiner

  #  APP_LAUNCHER="spotlight"
  #  TERMINAL="warp"
  #  KEYBOARD_TYPE="pc"

  while (($#)); do

    case "$1" in
    --app-launcher)
      shift
      if (($#)); then APP_LAUNCHER="$1"; else echo "ERROR: '--app-launcher' requires an argument" exit 1 >&2; fi
      ;;
    --terminal)
      shift
      if (($#)); then TERMINAL="$1" else echo "ERROR: '--terminal' requires an argument" exit 1 >&2; fi
      ;;
    --keyboard-type)
      shift
      if (($#)); then KEYBOARD_TYPE="$1" else echo "ERROR: '--keyboard-type' requires an argument" exit 1 >&2; fi
      ;;
    *)
      echo "ERROR: Unknown argument: $1" >&2
      exit 1
      ;;
    esac
    shift
  done

  if [ -z "$APP_LAUNCHER" ]; then

    clear
    echo -e "${RESET}Your app launcher:\n"

    echo "(1) Spotlight"
    echo "(2) Launchpad"
    echo "(3) Alfred"
    echo "${YELLOW}(Q) Quit${RESET}"

    printf "\nChoice [1|2|3|q]: "
    read -r opt

    case ${opt:l} in
    1*) APP_LAUNCHER="spotlight" ;;
    2*) APP_LAUNCHER="launchpad" ;;
    3*) APP_LAUNCHER="alfred" ;;
    q*) quit;;
    *)
      echo "Invalid choice. Shell change skipped."
      return
      ;;
    esac
  fi

  if [ -z "$KEYBOARD_TYPE" ]; then

    clear
    echo -e "${RESET}Your ${RESET}${BOLD}external${RESET} keyboard type\n"

    echo "(1) PC"
    echo "(2) Mac"
    echo "${YELLOW}(Q) Quit${RESET}"

    printf "\nChoice [1|2|q]: "
    read -r opt

    case ${opt:l} in
    1*) KEYBOARD_TYPE="pc" ;;
    2*) KEYBOARD_TYPE="mac" ;;
    q*) quit;;
    *)
      echo "Invalid choice. Shell change skipped."
      return
      ;;
    esac
  fi

  if [ -z "$TERMINAL" ]; then
    clear
    echo -e "${RESET}What is your terminal of choice:\n"

    echo "(1) Apple Terminal"
    echo "(2) iTerm"
    echo "(3) Warp"
    echo "${YELLOW}(Q) Quit${RESET}"

    printf "\nChoice [1|2|3|q]: "
    read -r opt

    case ${opt:l} in
    1*) TERMINAL="default-terminal" ;;
    2*) TERMINAL="iterm" ;;
    3*) TERMINAL="warp" ;;
    q*) quit;;
    *)
      echo "Invalid choice. Shell change skipped."
      return
      ;;
    esac

    clear
  fi

  #  === HERE do all actions

  # do karabiner.json backup
  DATE=$(date +"%m-%d-%Y-%T")
  cp $KARABINER_CONFIG $KARABINER_CONFIG_DIR/karabiner-"$DATE".json

  # delete existing profile
  echo "delete existing profile"
  jq --arg PROFILE_NAME "$PROFILE_NAME" 'del(.profiles[] | select(.name == $PROFILE_NAME))' $KARABINER_CONFIG >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG

  # add profile
  echo "add profile"
  jq '.profiles += $profile' $KARABINER_CONFIG --slurpfile profile karabiner-elements-profile.json --indent 4 >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG
  #rm PROFILE-tmp.json

  # add rules to profile
  echo "add rules to profile"

  apply_rules main-rules.json
  apply_rules finder-rules.json

  case ${APP_LAUNCHER:l} in
  spotlight*) apply_rules spotlight-rules.json "Installing Spotlight rules" ;;
  launchpad*) apply_rules launchpad-rules.json "Installing Launchpad rules" ;;
  alfred*) apply_rules alfred-rules.json "Installing Alfred rules" ;;
  *)
    echo "Invalid choice. Shell change skipped."
    return
    ;;
  esac

  case ${KEYBOARD_TYPE:l} in
  pc*) echo "Preparing for PC keyboard..." ;;
  mac*) prepare_for_mac_keyboard "Preparing for Mac keyboard..." ;;
  *)
    echo "Invalid choice. Shell change skipped."
    return
    ;;
  esac

  case ${TERMINAL:l} in
  default-terminal*) apply_rules terminal-rules.json "Installing Terminal rules" ;;
  iterm*) apply_rules iterm-rules.json "Installing iTerm rules" ;;
  warp*) apply_rules warp-rules.json "Installing Warp rules" ;;
  *)
    echo "Invalid choice. Shell change skipped."
    return
    ;;
  esac

  install_ide_keymap "IntelliJ" "IntelliJ IDEA Ultimate" "idea"
  install_ide_keymap "PyCharm" "PyCharm Community Edition" "pycharm"

  echo "${GREEN}SUCCESS${GREEN}"
  exit 0
}

main "$@"
