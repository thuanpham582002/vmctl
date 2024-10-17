#!/bin/zsh
set -o pipefail

COMMAND_TEMPLATE=$(cat /Users/noroom113/IdeaProjects/Init-My-Mac/zsh_oh_my_zsh/auto_complete/vmctl/command.yaml);
CONFIG_FILE=""
CONFIG_FOLDER=""

# Define color codes for output
RED="\033[0;31m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
RESET="\033[0m"

function startWith() {
  case $2 in
    $1*) true;;
    *) false;;
  esac
}

function get_config() {
  local _path=$HOME/.vmctl/config.yaml
  if [[ ! -f "$_path" ]]; then
    mkdir "$HOME/.vmctl" 2&> /dev/null
    touch "$_path"
    yq e '
    .current-context="null" |
    .contexts = []
    ' $_path -i
  fi

  local context=$1
  case $context in
  all)
    yq e '.contexts' $_path
    ;;
  *)
    yq e '.current-context' $_path
    ;;
  esac
}

function set_config() {
  local _path=$HOME/.vmctl/config.yaml
  if [[ ! -f "$_path" ]]; then
    mkdir "$HOME/.vmctl" 2&> /dev/null
    touch "$_path"
    yq e '
    .current-context="null" |
    .contexts = []
    ' $_path -i
  fi

  config=$1
  if [[ ! -f "$config" ]]; then
    echo "${RED}Error: Config file not found: $config${RESET}" >&2
    exit 1
  fi

  yq e '.current-context = "'$config'"' $_path -i
  # add to contexts
  yq e '.contexts += ["'$config'"]' $_path -i
}

function clear_config() {
  local _path=$HOME/.vmctl/config.yaml
  if [[ ! -f "$_path" ]]; then
    mkdir "$HOME/.vmctl" 2&> /dev/null
    touch "$_path"
    yq e '
    .current-context="null" |
    .contexts = []
    ' $_path -i
  fi

  yq e '
  .current-context="" |
  .contexts = []
  ' $_path -i
}


# Function to check if a command is available
function check_command() {
  local cmd="$1"
  if ! command -v "$cmd" &> /dev/null; then
    echo "${RED}Error: $cmd is not installed.${RESET}" >&2
    exit 1
  fi
}

show_help_command(){
  local _path="$1"
  IFS='. ' read -r -A path_arr <<< "$_path"
  echo "$COMMAND_TEMPLATE" | yq e "$_path._description"
  echo "COMMANDS:"
  while read -r line; do
    if startWith "_" "$line"; then
      continue
    fi
    echo -n "   "
    echo -n "$line: "
    echo "$COMMAND_TEMPLATE" | yq e "$_path.$line._description"
  done <<< "$(echo "$COMMAND_TEMPLATE" | yq e "$_path | keys | .[]")"
}

show_help_flag(){
  help_str=()
  local _path="$1"

  echo "FLAGS:"

  while true; do
    flags_size=$(echo "$COMMAND_TEMPLATE" | yq e "$_path._flags.list | length")
    if [[ $flags_size -gt 0 ]]; then
      for ((i = 0; i < $flags_size; i++)); do
        flag=$(echo "$COMMAND_TEMPLATE" | yq e "$_path._flags.list[$i].name")
        flag_desc=$(echo "$COMMAND_TEMPLATE" | yq e "$_path._flags.list[$i]._description")
        help_str+=("$flag: $flag_desc")
      done
    fi
    local inherit="$(echo "$COMMAND_TEMPLATE" | yq e "$_path._flags.inherit")"
    if [[ $inherit == "true" ]]; then
      _path=$(echo "$path" | sed 's/\.[^.]*$//')
    else
      break
    fi
  done
  # Sort the help string
  IFS=$'\n' sorted=($(sort <<<"${help_str[*]}"))
  unset IFS
  for line in "${sorted[@]}"; do
    echo "   $line"
  done
}

# Function to display help
function show_help() {
  local _path="$1"
  IFS='. ' read -r -A path_arr <<< "$_path"
  local _from_flag="$2"
  if [[ $_from_flag != true ]]; then
    _path=$(echo "$_path" | sed 's/\.[^.]*$//')
  fi
  show_help_command "$_path"
  show_help_flag "$_path"
  local usage=$(echo "$COMMAND_TEMPLATE" | yq e "$_path._usage")
  if [[ "$usage" != "null" ]]; then
    echo -n "Usage: $usage"
  fi
}

function build_path_yaml(){
  local args=("$@")
  local _path=""
  for arg in "${args[@]}"; do
    _path="$_path.$arg"
  done
  echo "$_path"
}

function have_path_yaml(){
  if [[ $(echo "${COMMAND_TEMPLATE}" | yq e "$1") != "null" ]]; then
      return 0
  else
      return 1
  fi
}

# Function to get groups from YAML
function get_groups() {
  yq e 'keys | .[]' "$CONFIG_FILE"
}

# Function to get VM names for a group
function get_vm_names() {
  local group="$1"
  yq e ".${group} | keys | .[]" "$CONFIG_FILE"
}

# Function to check if a file exists relative to the config file
function is_file_exist() {
  local file="$1"
  local from="$2"
  if [[ "$from" == "file" ]]; then
    [[ -f "$file" ]]
  else
    [[ -f "$CONFIG_FOLDER/$file" ]]
  fi
}

# Function to check if a file exists from the executable _path
function is_file_exist_from_executable_path() {
  local file="$1"
  [[ -f "$file" ]]
}

function execute_lima_shell(){
  local group="$1"
  local name="$2"
  local script_path=$3
  local script=''
  script=$(yq e '.'$group'.'$name'.'$script_path'.command' "$CONFIG_FILE" ) || script="null"
  if [[ $script == "null" ]]; then
    script=$script_path
    is_file_exist "$script" "file" && script=$(cat "$script")
  else
    is_file_exist "$script" "config" && script=$(cat "$CONFIG_FOLDER/$script")
    local su=$(yq e '.'$group'.'$name'.'$script_path'.su' "$CONFIG_FILE" || false)

    local _index=0
    while true ; do
      local envi=""
      envi=$(yq e ".$group.$name.$script_path.envs[$_index]" "$CONFIG_FILE")
      if [[ $? -ne 0 || "$envi" == "null" ]]; then
        break
      fi
      envi_name=$(yq e ".$group.$name.$script_path.envs[$_index].name" "$CONFIG_FILE")
      envi_type=$(yq e ".$group.$name.$script_path.envs[$_index].type" "$CONFIG_FILE")

      case $envi_type in
      null|config)
        envi_value=$(yq e ".$group.$name.$script_path.envs[$_index].value" "$CONFIG_FILE")
        envi_value=$(yq e ".$group.$name.$envi_value" "$CONFIG_FILE")
        ;;
      file)
        envi_value=$(yq e ".$group.$name.$script_path.envs[$_index].value" "$CONFIG_FILE")
        is_file_exist "$envi_value" "file" && envi_value=$(cat "$CONFIG_FOLDER/$envi_value")
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
function execute_lima_command() {
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
      done < <(yq e ".$group.$name.init_script | keys | .[]" "$CONFIG_FILE")
      for init_script in "${init_scripts[@]}"; do
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
function process_group() {
  local action="$1"
  local group="$2"
  local vm_names=()

  while IFS= read -r line; do
    vm_names+=("$line")
  done <<< "$(get_vm_names "$group")"

  for vm_name in "${vm_names[@]}"; do
    process_vm "$action" "$group" "$vm_name"
    echo
  done
}

# Function to process a single VM
function process_vm() {
  local action="$1"
  local group="$2"
  local name="$3"
  local template
  template=$(yq e ".$group.$name.template" "$CONFIG_FILE")
  execute_lima_command "$action" "$group" "$name" "$CONFIG_FOLDER/$template"
}

# Function to process all groups and VMs
function process_all() {
  local action="$1"
  local groups=()
  while IFS= read -r line; do
    groups+=("$line")
  done <<< "$(get_groups)"
  for group in "${groups[@]}"; do
    process_group "$action" "$group"
  done
}

function translate_execute(){
  local config_str="$1"
  # shellcheck disable=SC2001
  local _path=$(echo "$2" | sed 's/\._execute$//')
  local IFS=' ' # Tách các phần tử bằng khoảng trắng
  read -r -A config_arr <<< "$config_str" # Chia chuỗi thành mảng

  # Kiểm tra và thay thế các phần tử từ vị trí thứ 2 trở đi
  for ((i = 2; i <= ${#config_arr[@]}; i++)); do
      if [[ ! "${config_arr[i]}" =~ ^\" ]]; then
          # Lấy giá trị từ YAML và thay thế
          local yaml_value
          if [[ ${config_arr[i]} == ".val" ]]; then
            yaml_value=$(echo "$COMMAND_TEMPLATE" | yq e "${_path}.value")
          elif [[ ${config_arr[i]} == ".path" ]]; then
            yaml_value=$_path
          else
            yaml_value=$(echo "$COMMAND_TEMPLATE" | yq e ".${config_arr[i]}")
          fi
          config_arr[i]="$yaml_value" # Thay thế giá trị
      fi
  done

  # Nối lại mảng thành chuỗi
  local updated_config_str
  updated_config_str=$(IFS=' '; echo "${config_arr[*]}")
  echo "Executing: $GREEN $updated_config_str $RESET"
  eval "$updated_config_str"
}

function parse_args() {
  local args=("vmctl")
  typeset -A flags
  local flags=()
  local prev_flag=""
  local temp_args=()

  for arg in "$@"; do
    temp_args=("${args[@]}")
    temp_args+=("$arg")
    if have_path_yaml "$(build_path_yaml "${temp_args[@]}")"; then
      args+=("$arg")
      prev_flag=""
    elif [[ $arg == -* ]]; then
      flags[$arg]=""
    elif [[ -n "$prev_flag" ]]; then
      flags[${prev_flag}]="$arg"
      prev_flag=""
    else
      prev_flag=""
      local found_flag=false
      if have_path_yaml "$(build_path_yaml "${args[@]}")"; then
        while read -r line; do
          if startWith "_" "$line"; then
            continue
          fi
          temp_args=("${args[@]}")
          temp_args+=("$line")
          temp_args+=("value")
          if have_path_yaml "$(build_path_yaml "${temp_args[@]}")"; then
            args+=("${line}")
            COMMAND_TEMPLATE=$(echo "${COMMAND_TEMPLATE}" | yq e "$(build_path_yaml "${temp_args[@]}")=\"$arg\"")
            found_flag=true
            break
          fi
        done <<< "$(echo "${COMMAND_TEMPLATE}" | yq e "$(build_path_yaml "${args[@]}") | keys | .[]")"
      fi

      if [[ $found_flag == false ]]; then
        echo "${RED}Error: Invalid argument: $arg${RESET}" >&2
        exit 1
      fi
    fi
  done

  args+=("_execute")
  have_path_yaml "$(build_path_yaml "${args[@]}")" || echo "${RED}Error: Invalid argument: $arg${RESET}" >&2

  translate_execute "$(echo "$COMMAND_TEMPLATE" | yq e "$(build_path_yaml "${args[@]}")")" "$(build_path_yaml "${args[@]}")"
}

function init() {
  CONFIG_FILE=$(get_config "current")
  CONFIG_FOLDER=$(dirname $CONFIG_FILE)
  if [[ -z "$CONFIG_FILE" ]]; then
    CONFIG_FOLDER="null"
  fi
}

# Main execution
check_command "yq"
check_command "limactl"
init
parse_args "$@"