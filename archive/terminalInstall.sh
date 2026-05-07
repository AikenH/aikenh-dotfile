#!/bin/bash
# TODO: Add zellij install and keyshot config.
# FIXME: before Carryout >> should check before

# Proxy Setting.
proxy=http://192.168.31.201:7890
echo "proxy: $proxy"

# ============================ FUNCTION SETTING===================
source ./FunctionList.sh

# ============================MAIN PROCESS==============================
plugins_setting 
exec_cmd_status "setting plugins"

# INSTALL zsh and oh-my-zsh and plugins 
sudo apt-get install zsh
exec_cmd_status "install zsh"

sudo apt-get install curl
exec_cmd_status "install curl"

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
grep -q "hist" ~/.zshrc
if [[ $? -ne 0 ]];then
  echo "alias hist='history -i'" >> ~/.zshrc # show time of historyG
  echo "alias ltr='ls -rsthl'" >> ~/.zshrc
  echo "alias cl='clear'" >> ~/.zshrc
  echo "alias nv='nvim'" >> ~/.zshrc
  echo "alias lsd='ls -d */'" >> ~/.zshrc
fi

# install btm and htop for monitor
sudo apt-get install htop

which btm
if [[ $? -ne 0 ]];then
  cd /tmp/
  ls /tmp/bottom*.deb
  if [[ $? -ne 0 ]];then
    curl https://api.github.com/repos/ClementTsang/bottom/releases/latest --proxy $proxy| grep browser_download_url | grep amd64.deb | cut -d '"' -f 4 | wget -qi -
  else
    echo "bottom install package exit"
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
grep -q "Proxy Function Import Flag" ~/.zshrc

if [[ $? -ne 0 ]];then
  cat ./setproxy.sh >> ~/.zshrc
fi

# support bash command
echo "setopt no_nomatch" >> ~/.zshrc
