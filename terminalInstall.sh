#!/bin/bash
#FIXME: make get the ip automaticly and echo that process
#FIXME: modify the zshrc for plugins including extra.
#FIXME: backup the origin version and rewrite one, when we hit the key word
#       we will modify the content. like plugins.

# Proxy Setting.
proxy=172.30.240.1:8890
echo "proxy: $proxy"

# DESC: input is process-description.
function exec_cmd_status(){
  # this input para's process name.
  if [ $? -ne 0 ]; then
    echo $1" fail! check this. process."
    exit 1
  else
    echo $1" success. contiune"
  fi
}

# DESC: check file exist
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

# DESC: check dir exist
function check_dir(){
  # whether dir is exist or not.
  if [ -d "$1" ]; then
    echo "$1 is a directory, work on"
  else
    echo "$1 not exist, create directory"
    mkdir -p $1
  fi
}

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

# ====================main process=========================
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
cat ~/.bashrc | grep -i nvm >> ~/.zshrc # init zsh
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
 
cd /tmp/
curl https://api.github.com/repos/ClementTsang/bottom/releases/latest --proxy $proxy| grep browser_download_url | grep amd64.deb | cut -d '"' -f 4 | wget -qi -
sudo apt install ./bottom*.deb
cd - 

# install ranger for file manager
sudo apt-get install ranger

# install neofetch
sudo apt-get install neofetch

# export the extra function into zsh.
cat ./zsh_extra >> ~/.zshrc

# becus we install nvm in bash, we move config here.
echo "export NVM_DIR=\"$HOME/.nvm\"" >> ~/.zshrc
echo "[ -s \"$NVM_DIR/nvm.sh\" ] && \. \"$NVM_DIR/nvm.sh\"" >> ~/.zshrc # This loads nvm
echo "[ -s \"$NVM_DIR/bash_completion\" ] && \. \\"$NVM_DIR/bash_completion\" >> ~/.zshrc  # This loads nvm bash_completion" >> ~/.zshrc

# support bash command
echo "setopt no_nomatch" >> ~/.zshrc
