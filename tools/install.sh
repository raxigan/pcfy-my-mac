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
  jq --arg PROFILE_NAME "$PROFILE_NAME" '(.profiles[] | select(.name == $PROFILE_NAME) | .complex_modifications.rules) += $rules[].rules' \
    $KARABINER_CONFIG --slurpfile rules ../karabiner-elements/"$1" --indent 4 >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG
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
  jq --arg PROFILE_NAME "$PROFILE_NAME" '(.profiles[] | select(.name == $PROFILE_NAME) | .complex_modifications.rules) += $rules[].rules' \
    $KARABINER_CONFIG --slurpfile rules ../karabiner-elements/main-rules.json --indent 4 >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG

  #echo "install Albert rule"
  echo "Switch from Spotlight to Alfred? [Y/n]"

  read -r opt

  case ${opt:u} in
  N*)
    echo "Installing spotlight rules"
    apply_rules spotlight-rules.json
    ;;
  Y*)
    echo "Installing alfred rules"
    apply_rules alfred-rules.json
    ;;
  *)
    echo "Invalid choice. Shell change skipped."
    return
    ;;
  esac

  echo "${RESET}Is your ${RESET}${BOLD}external${RESET} keyboard mac or PC? [Mac/PC]"

  read -r opt

  case ${opt:u} in
  MAC*)
    echo "Preparing for Mac keyboard..."
    jq --arg PROFILE_NAME "$PROFILE_NAME" '.profiles |= map(if .name == $PROFILE_NAME then walk(if type == "object" and .conditions then del(.conditions[] | select(.identifiers[]?.is_built_in_keyboard)) else . end) else . end)' $KARABINER_CONFIG --indent 4 \
      >INPUT.tmp && mv INPUT.tmp $KARABINER_CONFIG
    ;;
  PC*) echo "Preparing for PC keyboard...";;
  *)
    echo "Invalid choice. Shell change skipped."
    return
    ;;
  esac

  echo "Do you use Terminal/iTerm/Warp? [Terminal/iTerm/Warp/none]"

  read -r opt

  case ${opt:u} in
  TERMINAL*)
    echo "Installing Terminal rules"
    apply_rules terminal-rules.json
    ;;
  ITERM*)
    echo "Installing iTerm rules"
    apply_rules iterm-rules.json
    ;;
  WARP*)
    echo "Installing Warp rules"
    apply_rules iterm-rules.json
    ;;
  *)
    echo "Invalid choice. Shell change skipped."
    return
    ;;
  esac

  IJ_VER=$(< ~/Library/Application' 'Support/JetBrains/Toolbox/scripts/idea grep IntelliJ | cut -d/ -f11)
  ij_config_dir=""

  echo "Installing IntelliJ (ver. ${IJ_VER}) keymap ..."
  IJ_CONFIGS=~/Library/Application' 'Support/JetBrains


  for entry in ~/Library/Caches/Jetbrains/*; do

      version=$(find "$entry" -type d -name "*IntelliJ*" -exec grep "app.build.number=" {}/.appinfo \; | sed 's/.*app.build.number=\([^&]*\).*/\1/')

      if [ "$version" = "$IJ_VER" ]; then
          ij_config_dir=$(echo "$entry" | cut -d/ -f7)
          break
      fi
  done

  # if ij_config_dir empty then exit

  KEYMAPS_DIR=${IJ_CONFIGS}/${ij_config_dir}/keymaps

  curl --silent -o "${KEYMAPS_DIR}/test.xml" https://raw.githubusercontent.com/raxigan/macos-pc-mode/main/keymaps/intellij-idea.xml

  echo "Restart IntelliJ. Then choose XWin IntelliJ IDEA in Preferences > Keymaps > Xwin"

}

main "$@"
