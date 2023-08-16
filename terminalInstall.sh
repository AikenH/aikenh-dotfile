#!/bin/bash
#FIXME: make get the ip automaticly and echo that process

# Proxy Setting.
proxy=172.30.240.1:8890
echo "proxy: $proxy"

# ============================ FUNCTION SETTING===================
# 1. DESC: input is process-description.
function exec_cmd_status(){
  # this input para's process name.
  if [ $? -ne 0 ]; then
    echo $1" fail! check this. process."
    exit 1
  else
    echo $1" success. contiune"
  fi
}

# 2. DESC: check file exist
function check_file(){
  mode=$(($2))
  # whether file is exist or not.
  if [ -f "$1" ]; then
   echo "$1 exist. work on"
  else
    if (( $mode <= 0 ))
      echo "$1 not exist, we will create one."
      touch $1
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
    mkdir -p $1
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
  zsh_file='~/.zshrc'
  line_num=`cat $zsh_file  | awk "/plugins=\(git\)/{print NR}"`
  # write a loop to write down the plugins<Left>
  plugins=('git' 'zsh-syntax-highlighting' 'zsh-autosuggestions' 'colored-man-pages' 'safe-paste' 'themes' 'tmux' 'sudo' 'z')
  grep -A 10 "plugins" $zsh_file | grep -q "sudo\|safe-paste\|themes\|tmux"
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

# ============================MAIN PROCESS==============================
# INSTALL zsh and oh-my-zsh and plugins 
sudo apt-get install zsh
exec_cmd_status "install zsh"

sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh --proxy $proxy)"
exec_cmd_status "install & download oh-my-zsh"

# install extra plugins from github.
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting
exec_cmd_status "download zsh-syntax-highlighting to omz"

git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions
exec_cmd_status "download zsh-autosuggestions to omz"

# FIXME: =======================this part should move to conda install =====================
# generate zsh dotfile here. then install plugins.
bash 
conda init zsh # init zsh
exec_cmd_status "conda init zsh "
# ==========================================================================================

# zsh and plugins will contain many alias include --color=auto
check_file ~/.zshrc 1
echo "alias hist='history -i'" >> ~/.zshrc # show time of history
echo "alias ltr='ls -rsthl'" >> ~/.zshrc
echo "alias cl='clear'" >> ~/.zshrc
echo "alias nv='nvim'" >> ~/.zshrc
echo "alias lsd='ls -d */'" >> ~/.zshrc

plugins_setting()
exec_cmd_status "setting plugins"

# install btm and htop for monitor
sudo apt-get install htop

command -v btm
if [[ $? -eq 0 ]];then
  cd /tmp/
  skip_exist /tmp/bottom*.deb
  if [[ $? -eq 0 ]];then
    curl https://api.github.com/repos/ClementTsang/bottom/releases/latest --proxy $proxy| grep browser_download_url | grep amd64.deb | cut -d '"' -f 4 | wget -qi -
  else
    echo "bottom exit"
  fi
  sudo apt install ./bottom*.deb
  cd -
else
  echo "bottom has exist, not need to install"
fi

# install ranger for file manager
sudo apt-get install ranger

# install neofetch
sudo apt-get install neofetch

# export the extra function into zsh.
cat ./setproxy.sh >> ~/.zshrc

# support bash command
echo "setopt no_nomatch" >> ~/.zshrc
