#!/bin/zsh

setup_color() {

  RED=$(printf '\033[31m')
  GREEN=$(printf '\033[32m')
  YELLOW=$(printf '\033[33m')
  BLUE=$(printf '\033[34m')
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

#KARABINER_CONFIG=~/.config/karabiner
KARABINER_CONFIG_DIR=../karabiner-config
KARABINER_CONFIG=$KARABINER_CONFIG_DIR/karabiner.json

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

main() {

  setup_color
  install_brew
  install_jq
  install_karabiner

  # shellcheck disable=SC2034
  PROFILE_NAME="PC mode"

  # do karabiner.json backup
  DATE=$(date +"%m-%d-%Y-%T")
  cp $KARABINER_CONFIG $KARABINER_CONFIG_DIR/karabiner-"$DATE".json

  # delete existing profile
  echo "delete existing profile"
  jq 'del(.profiles[] | select(.name == "$PROFILE_NAME"))' $KARABINER_CONFIG >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG

  # add profile
  echo "add profile"
  jq '.profiles += $profile' $KARABINER_CONFIG --slurpfile profile karabiner-elements-profile.json --indent 4 >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG
  #rm PROFILE-tmp.json

  # add rules to profile
  echo "add rules to profile"
  jq '(.profiles[] | select(.name == "$PROFILE_NAME") | .complex_modifications.rules) += $rules[].rules' \
    $KARABINER_CONFIG --slurpfile rules ../karabiner-elements/main-rules.json --indent 4 >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG

  #echo "install Albert rule"
  echo "Do you want to install additional launcher rules? [Spotlight/Alfred/None]"

  read -r opt

  case ${opt:u} in
  SPOTLIGHT*)
    echo "Installing spotlight rules"
    jq '(.profiles[] | select(.name == "$PROFILE_NAME") | .complex_modifications.rules) += $rules[].rules' \
      $KARABINER_CONFIG --slurpfile rules ../karabiner-elements/spotlight-rules.json --indent 4 >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG
    ;;
  RED*)
    echo "Installing alfred rules"
    jq '(.profiles[] | select(.name == "$PROFILE_NAME") | .complex_modifications.rules) += $rules[].rules' \
      $KARABINER_CONFIG --slurpfile rules ../karabiner-elements/alfred-rules.json --indent 4 >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG
    ;;
  *)
    echo "Invalid choice. Shell change skipped."
    return
    ;;
  esac

  # Install IDE keymaps
  # TODO: support IDEs not installed via toolbox

}

main "$@"
