#!/bin/zsh

# vmctl: A script to manage virtual machines using Lima
# Usage: vmctl [command] <file_configuration> [all|<group>|<group> <name>] [flags]

# Define color codes for output
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
RESET="\033[0m"

# Define available commands and their allowed flagstypeset -A COMMANDS
typeset -A COMMANDS

COMMANDS=(
  [create]="0 2"
  [start]="0"
  [delete]="0"
  [stop]="0"
  [execute]="0 1 2"
  [config]="0"
  [help]="0"
)

# Define command descriptions
typeset -A COMMANDS_DESC
COMMANDS_DESC=(
  [start]="Start the VM with the specified name or group."
  [delete]="Delete an entry or group from the YAML file."
  [stop]="Stop the VM with the specified name or group."
  [execute]="Execute a script in the VM (script defined in the YAML file)."
  [help]="Display this help message."
)

# Define flags and their associated arguments
typeset -A FLAGS
FLAGS=(
  [-h]=0
  [--help]=0
  [-s]=1
  [--script]=1
  [-r]=2
  [--root]=2
)

# Initialize flag arguments
typeset -A FLAGS_ARGS
FLAGS_ARGS=(
  [0]=false
  [1]=""
  [2]=false
)

# Define flag descriptions
typeset -A FLAGS_DESC
FLAGS_DESC=(
  [0]="-h --help: Display this help message."
  [1]="-s --script: Execute a script in the VM (script defined in the YAML file)."
  [2]="-c --command: Execute a command in the VM."
)

# Global variables
YAML_FILE=""
YAML_FOLDER=""
COMMAND=""
NAME=""
GROUP=""

# Function to check if a command is available
check_command() {
  local cmd="$1"
  if ! command -v "$cmd" &> /dev/null; then
    echo "${RED}Error: $cmd is not installed.${RESET}" >&2
    exit 1
  fi
}

# Function to display help
show_help() {
  local target="$1"
  if [[ -z "$target" ]]; then
    for command in "${(k)COMMANDS[@]}"; do
      show_help "$command"
    done
  else
    echo "Usage: vmctl [$target] <file_configuration> [all|<group>|<group> <name>] [flags]"
    read -rA available_flags <<< "${COMMANDS[$target]}"
    for flag in $available_flags; do
      echo "${FLAGS_DESC[$flag]}"
    done
    echo
  fi
}

# Function to parse command line arguments
parse_args() {
  if [[ $# -lt 2 ]]; then
    echo "${RED}Error: Insufficient number of arguments.${RESET}" >&2
    show_help
    exit 1
  fi

  COMMAND="$1"
  YAML_FILE="$2"
  YAML_FOLDER=$(dirname "$YAML_FILE")
  shift 2

  case "$1" in
    all)
      ;;
    *)
      if [[ "$1" == */* ]]; then
        GROUP=$(echo "$1" | cut -d'/' -f1)
        NAME=$(echo "$1" | cut -d'/' -f2)
      else
        GROUP="$1"
      fi
      ;;
  esac
  shift

  while [[ "$#" -gt 0 ]]; do
    case "$1" in
      -h|--help)
        show_help "$COMMAND"
        exit 0
        ;;
      -s|--script)
        if [[ "${COMMANDS[$COMMAND]}" != *"${FLAGS[$1]}"* ]]; then
          echo "${RED}Error: Invalid flag: $1${RESET}" >&2
          show_help "$COMMAND"
          exit 1
        fi
        if [[ $# -gt 1 ]]; then
          shift
        else
          echo "${RED}Error: Missing argument for flag: $1${RESET}" >&2
          show_help "$COMMAND"
          exit 1
        fi
        FLAGS_ARGS[1]="$1"
        ;;
      -r|--root)
        if [[ "${COMMANDS[$COMMAND]}" != *"${FLAGS[$1]}"* ]]; then
          echo "${RED}Error: Invalid flag: $1${RESET}" >&2
          show_help "$COMMAND"
          exit 1
        fi
        FLAGS_ARGS[2]=true
        ;;
      *)
        echo "${RED}Error: Invalid argument: $1${RESET}" >&2
        show_help "$COMMAND"
        exit 1
        ;;
    esac
    shift
  done
}

# Function to validate YAML file
validate_yaml() {
  if [[ ! -f "$YAML_FILE" ]]; then
    echo "${RED}Error: YAML file not found: $YAML_FILE${RESET}" >&2
    exit 1
  fi
}

# Function to get groups from YAML
get_groups() {
  yq e 'keys | .[]' "$YAML_FILE"
}

# Function to get VM names for a group
get_vm_names() {
  local group="$1"
  yq e ".${group} | keys | .[]" "$YAML_FILE"
}

# Function to check if a file exists relative to the config file
is_file_exist() {
  local file="$1"
  local from="$2"
  if [[ "$from" == "file" ]]; then
    [[ -f "$file" ]]
  else
    [[ -f "$YAML_FOLDER/$file" ]]
  fi
}

# Function to check if a file exists from the executable path
is_file_exist_from_executable_path() {
  local file="$1"
  [[ -f "$file" ]]
}

execute_lima_shell(){
  local group="$1"
  local name="$2"
  local script_path=$3
  local script=''
  script=$(yq e '.'$group'.'$name'.'$script_path'.command' "$YAML_FILE" ) || script="null"
  if [[ $script == "null" ]]; then
    script=$script_path
    is_file_exist "$script" "file" && script=$(cat "$script")
  else
    is_file_exist "$script" "config" && script=$(cat "$YAML_FOLDER/$script")
    local su=$(yq e '.'$group'.'$name'.'$script_path'.su' "$YAML_FILE" || false)

    local _index=0
    while true ; do
      local envi=""
      envi=$(yq e ".$group.$name.$script_path.envs[$_index]" "$YAML_FILE")
      if [[ $? -ne 0 || "$envi" == "null" ]]; then
        break
      fi
      envi_name=$(yq e ".$group.$name.$script_path.envs[$_index].name" "$YAML_FILE")
      envi_type=$(yq e ".$group.$name.$script_path.envs[$_index].type" "$YAML_FILE")

      case $envi_type in
      null|config)
        envi_value=$(yq e ".$group.$name.$script_path.envs[$_index].value" "$YAML_FILE")
        envi_value=$(yq e ".$group.$name.$envi_value" "$YAML_FILE")
        ;;
      file)
        envi_value=$(yq e ".$group.$name.$script_path.envs[$_index].value" "$YAML_FILE")
        is_file_exist "$envi_value" "file" && envi_value=$(cat "$YAML_FOLDER/$envi_value")
        ;;
      *)
        echo "${RED}Error: Invalid envi type: $envi_type${RESET}" >&2
        exit 1
        ;;
      esac

      script="export $envi_name=$envi_value; $script"
      _index=$((_index+1))
    done
  fi

  if [[ -z $script ]]; then
    echo "${RED}Error: Script not found: $script_path${RESET}" >&2
    exit 1
  fi

  echo "${BLUE}$group.$name Executing init script: limactl shell $name bash -c '$script'${RESET}"
  if [[ ${FLAGS_ARGS[2]} == true || $su == true ]]; then
    limactl shell $name sudo bash -c "$script"
  else
    limactl shell $name bash -c "$script"
  fi
}

# Function to execute lima command
execute_lima_command() {
  local action="$1"
  local group="$2"
  local name="$3"
  local template="$4"

  echo "${BLUE}Processing VM with action $action in group: $group with name: $name and template: $template${RESET}"
  case "$action" in
    create)
      execute_lima_command delete "$group" "$name" "$template"
      limactl start --name="$name" "$template" --tty=false || exit 1
      local init_scripts=()
      while read -r line; do
        init_scripts+=("$line")
      done < <(yq e ".$group.$name.init_script | keys | .[]" "$YAML_FILE")
      for init_script in $init_scripts; do
          execute_lima_shell "$group" "$name" "init_script.$init_script"
      done
      ;;
    start)
      limactl start --name="$name" --tty=false
      ;;
    stop)
      limactl stop "$name"
      ;;
    delete)
      limactl stop "$name"
      limactl delete "$name"
      ;;
    execute)
      if [[ -n "${FLAGS_ARGS[1]}" ]]; then
        execute_lima_shell "$group" "$name" "${FLAGS_ARGS[1]}"
        return
      fi
      echo "${RED}Error: No script or command provided.${RESET}" >&2
      exit 1
      ;;
  esac
  echo "${GREEN}Processing VM with action $action in group $group with name $name done${RESET}"
}

# Function to process a single group
process_group() {
  local action="$1"
  local group="$2"
  local vm_names=()

  while IFS= read -r line; do
    vm_names+=("$line")
  done <<< "$(get_vm_names "$group")"

  for vm_name in $vm_names; do
    local template
    template=$(yq e ".$group.$vm_name.template" "$YAML_FILE")
    execute_lima_command "$action" "$group" "$vm_name" "$YAML_FOLDER/$template"
    echo
  done
}

# Function to process a single VM
process_vm() {
  local action="$1"
  local group="$2"
  local name="$3"
  local template
  template=$(yq e ".$group.$name.template" "$YAML_FILE")
  execute_lima_command "$action" "$group" "$name" "$YAML_FOLDER/$template"
}

# Function to process all groups and VMs
process_all() {
  local action="$1"
  local groups=()
  while IFS= read -r line; do
    groups+=("$line")
  done <<< "$(get_groups)"
  for group in $groups; do
    process_group "$action" "$group"
  done
}

# Main function to execute commands
execute_command() {
  local action

  case "$COMMAND" in
    create|start|delete|stop|execute)
      action="$COMMAND"
      ;;
    *)
      echo "${RED}Error: Invalid command: $COMMAND${RESET}" >&2
      show_help
      exit 1
      ;;
  esac

  if [[ -z "$GROUP" ]]; then
    echo "${BLUE}Processing all VMs.${RESET}"
    process_all "$action"
  elif [[ -z "$NAME" ]]; then
    echo "${BLUE}Processing group: $GROUP${RESET}"
    process_group "$action" "$GROUP"
  else
    echo "${BLUE}Processing VM: $NAME${RESET}"
    process_vm "$action" "$GROUP" "$NAME"
  fi
}

# Main execution
check_command "yq"
check_command "limactl"
parse_args "$@"
validate_yaml
execute_command