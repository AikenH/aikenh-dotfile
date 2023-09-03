#!/bin/bash
# **********************************************************FUNCTION LIST*********************************************
# 1. DESC: input is process-description.
function exec_cmd_status(){
  # this input para's process name.
  if [ $? -ne 0 ]; then
    echo "$1 fail! check this. process."
    exit 1
  else
    echo "$1 success. contiune"
  fi
}

# 2. DESC: check file exist
function check_file(){
  mode=$(($2))
  # whether file is exist or not.
  if [ -f "$1" ]; then
   echo "$1 exist. work on"
  else
    if (( $mode <= 0 )); then
      echo "$1 not exist, we will create one."
      touch "$1"
    else
      echo "$1 is necessary, pipeline failed, try again."
      exit 1 
    fi
  fi
}

# 3. DESC: check dir exist
function check_dir(){
  # whether dir is exist or not.
  if [ -d "$1" ]; then
    echo "$1 is a directory, work on"
  else
    echo "$1 not exist, create directory"
    mkdir -p "$1"
  fi
}

# 4. DESC: check file, if $? -eq 0 not exist, -eq 1 file exit.
function skip_exist(){
  if ! ls $1 > /dev/null 2>&1; then
    echo "$1 file not exist. contiune"
    return 0
  else
    echo "file exit pass this."
    return 1
  fi
}

# 5. DESC: setting plugin for zsh
function plugins_setting(){
  # find out plugins line and add plugins on it.
  # this method only suit for init situation. so we check it.
  zsh_file="$HOME/.zshrc"
  line_num=`cat $zsh_file  | awk "/plugins=\(git\)/{print NR}"`
  # write a loop to write down the plugins<Left>
  plugins=('git' 'zsh-syntax-highlighting' 'zsh-autosuggestions' 'colored-man-pages' 'safe-paste' 'themes' 'tmux' 'sudo' 'z')
  grep -A 10 "plugins" $zsh_file | grep -q "sudo\|safe-paste\|tmux"
  isinit=$?

  if [[ $isinit -eq '0' ]]; then
    echo "init had finished, skip config plugins"
  else
    echo "config plugins"
    sudo sed -i "${line_num}d" $zsh_file
    line_num=$[line_num-1]
    sudo sed -i "${line_num}a  \       \)" $zsh_file
    sudo sed -i "${line_num}a     plugins=\(" $zsh_file

    line_num=$[line_num+1]
    echo $line_num
    for phase in ${plugins[@]}; do
      sudo sed -i "${line_num}a \         $phase" $zsh_file
    done
  fi
}
# ****************************************************************************************************************
