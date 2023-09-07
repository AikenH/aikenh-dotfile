#!/bin/bash
# ****************************************INIT FUNCTION************************************
source ./FunctionList.sh
# *****************************************************************************************
# get the sus info
# set proxy for wsl2
proxy="172.30.240.1:8890"

# 0.install miniconda for python
echo "***************install miniconda for neovim"

if [[ ! $(conda -V) ]];then 
  wget -c -O miniconda_install.sh https://mirrors.tuna.tsinghua.edu.cn/anaconda/miniconda/Miniconda3-latest-Linux-x86_64.sh
  chmod 777 miniconda_install.sh
  bash miniconda_install.sh
  exec_cmd_status "install conda process"
  rm -v miniconda_install.sh
  source ~/.zshrc
else
  echo "conda had installed, pass"
fi

# 1.install neovim (using unstable for diff sys)
echo "***************install neovim"
grep -qr neovim /etc/apt/sources.list.d/* 
if [[ $? -ne 0 ]];then
  sudo apt-get install software-propertier-common
  sudo add-apt-repository ppa:neovim-ppa/unstable
  sudo apt-get update
  sudo apt-get install neovim
else
  sudo apt-get install neovim
fi
exec_cmd_status "install neovim process"

nvim_version=$(nvim -v | grep -i nvim)
echo "nvim ver.is : $nvim_version"

# 2. install dependency of neovim
echo "***************install dependency"
pip install pynvim
exec_cmd_status "pynvim install"

pip install neovim
exec_cmd_status "pip install neovim"

# 3. install lazygit
echo "***************install lazygit"
lazygit -v
if [[ $? -ne 0 ]];then 
  LAZYGIT_VERSION=$(curl -s "https://api.github.com/repos/jesseduffield/lazygit/releases/latest" | grep -Po '"tag_name": "v\K[^"]*')
  curl -Lo lazygit.tar.gz "https://github.com/jesseduffield/lazygit/releases/latest/download/lazygit_${LAZYGIT_VERSION}_Linux_x86_64.tar.gz" --proxy $proxy
  tar xf lazygit.tar.gz lazygit
  sudo install lazygit /usr/local/bin
  exec_cmd_status "install lazygit"
  rm -v lazygit.tar.gz
  rm -v lazygit
  echo "lazygit version $LAZYGIT_VERSION had installed"
else
  echo "lazygit had installed, pass"
fi

# 4. install ripgrep, fd-find
echo "***************install ripgrep fd-find"
sudo apt-get install ripgrep
exec_cmd_status "install ripgrep"

sudo apt-get install fd-find
exec_cmd_status "install fd-find"

# 5. install nvm to install nodejs.(need restart)
echo "***************install nvm and nodejs, npm"
source ~/.nvm/nvm.sh
source ~/.profile
nvm -v 
if [[ $? -ne 0 ]]; then
  curl --proxy $proxy -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
  exec_cmd_status "install nvm"
  source ~/.bashrc
  echo "you need restart to reboot nvm, then we use nvm to install node"
  # cpy nvm 2 .zshrc
  cat ~/.bashrc | grep -i nvm >>~/.zshrc # init zsh
  # becus we install nvm in bash, we move config here.
  echo "export NVM_DIR=\"$HOME/.nvm\"" >>~/.zshrc
  echo "[ -s \"$NVM_DIR/nvm.sh\" ] && \. \"$NVM_DIR/nvm.sh\"" >>~/.zshrc                   # This loads nvm
  echo "[ -s \"$NVM_DIR/bash_completion\" ] && \. \"$NVM_DIR/bash_completion\" >>~/.zshrc # This loads nvm bash_completion" >> ~/.zshrc
else
  echo "nvm had installed"
fi

# manager specific npm version
node -v 
if [[ $? -ne 0 ]]; then 
  nvm list-remote
  nvm -v
  nvm install v18.15.0
  exec_cmd_status "install node and npm"
else
  node -v
  npm -version
fi
npm install -g neovim
exec_cmd_status "npm install neovim"

# 6. install build-essential
echo "***************install build-essential"
sudo apt-get install build-essential
exec_cmd_status "install build-essential"

sudo apt-get install gcc
exec_cmd_status "install gcc"

# 7. install ruby and gem
echo "***************install ruby and gem"
sudo apt-get install ruby-dev
exec_cmd_status "install ruby"

sudo apt-get install rubygems
exec_cmd_status "install gem for ruby"

# 8. install gem neovim
ruby -v
# gem environment
sudo gem install neovim
exec_cmd_status "gem install neovim"

# 9. start neovim 
echo "Neovim Install Success. U can Start And Try it"
