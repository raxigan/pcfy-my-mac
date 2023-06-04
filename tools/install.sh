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

    echo "Restart IDE_FULL_NAME. Then choose XWin $IDE_FULL_NAME in Preferences > Keymaps > Xwin"
}

main() {

  setup_color
  install_brew
  install_jq
  install_karabiner

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

  #echo "install Albert rule"
  #  echo "Switch from Spotlight to Alfred? [Y/n]"

  echo -e "Your app launcher:\n"

  echo "(1) Spotlight"
  echo "(2) Launchpad"
  echo "(3) Alfred"

  echo -e "\n(R) Restart  (Q) Quit"
  echo -e "\nChoice [1|2|3|r|q]:"

  read -r opt

  case ${opt:u} in
  1*) apply_rules spotlight-rules.json "Installing spotlight rules" ;;
  2*) apply_rules launchpad-rules.json "Installing launchpad rules" ;;
  3*) apply_rules alfred-rules.json "Installing launchpad rules" ;;
  *)
    echo "Invalid choice. Shell change skipped."
    return
    ;;
  esac

  echo -e "${RESET}Your ${RESET}${BOLD}external${RESET} keyboard type\n"

  echo "(1) PC"
  echo -e "(2) Mac"
  echo -e "\n(R) Restart  (Q) Quit"
  echo -e "\nChoice [1|2|r|q]:"

  read -r opt

  case ${opt:u} in
  1*) prepare_for_mac_keyboard "Preparing for Mac keyboard...";;
  2*) echo "Preparing for PC keyboard..." ;;
  *) echo "Invalid choice. Shell change skipped."; return;;
  esac

  echo -e "What is your terminal of choice:\n"

  echo "(1) Apple Terminal"
  echo "(2) iTerm"
  echo "(3) Warp"

  echo -e "\n(R) Restart  (Q) Quit"
  echo -e "\nChoice [1|2|3|r|q]:"

  read -r opt

  case ${opt:u} in
  1*) apply_rules terminal-rules.json "Installing Terminal rules" ;;
  2*) apply_rules iterm-rules.json "Installing iTerm rules" ;;
  3*) apply_rules warp-rules.json "Installing Warp rules" ;;
  *) echo "Invalid choice. Shell change skipped."; return;;
  esac

  install_ide_keymap "IntelliJ" "IntelliJ IDEA" "idea"
  install_ide_keymap "PyCharm" "PyCharm CE" "pycharm"

}

main "$@"
