#!/bin/bash
# ****************************************INIT FUNCTION************************************
source ./FunctionList.sh
# *****************************************************************************************
# get the sus info
# set proxy for wsl2
proxy="172.30.240.1:8890"

# 0.install miniconda for python
echo "------------------install miniconda for neovim"
wget -c https://mirrors.tuna.tsinghua.edu.cn/anaconda/miniconda/Miniconda3-latest-Linux-x86_64.sh
chmod 777 Miniconda3-latest-Linux-x86_64.sh
bash Miniconda3-latest-Linux-x86_64.sh
rm Miniconda3-latest-Linux-x86_64.sh

# 1.install neovim (using unstable for diff sys)
echo "------------------install neovim"
sudo apt-get install software-propertier-common
sudo add-apt-repository ppa:neovim-ppa/unstable
sudo apt-get update
sudo apt-get install neovim

nvim_version=$(nvim -v | grep dev)
echo "neovim's version is $nvim_version"

# 2. install dependency of neovim
echo "-------------------install dependency"
sudo apt-get install python-pip python3-dev
sudo apt-get install python3
sudo pip3 install pynvim
sudo pip3 install neovim

# 3. install lazygit
echo "--------------------install lazygit"
LAZYGIT_VERSION=$(curl -s "https://api.github.com/repos/jesseduffield/lazygit/releases/latest" | grep -Po '"tag_name": "v\K[^"]*')
curl -Lo lazygit.tar.gz "https://github.com/jesseduffield/lazygit/releases/latest/download/lazygit_${LAZYGIT_VERSION}_Linux_x86_64.tar.gz" --proxy $proxy
tar xf lazygit.tar.gz lazygit
sudo install lazygit /usr/local/bin
rm lazygit.tar.gz

# 4. install ripgrep, fd-find
echo "--------------------install ripgrep fd-find"
sudo apt-get install ripgrep
sudo apt-get install fd-find

# 5. install nvm to install nodejs.(need restart)
echo "---------------------install nvm and nodejs, npm"
curl --proxy $proxy -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash
source ~/.bashrc
echo "you need restart to reboot nvm, then we use nvm to install node"

# cpy nvm 2 .zshrc
cat ~/.bashrc | grep -i nvm >>~/.zshrc # init zsh
# becus we install nvm in bash, we move config here.
echo "export NVM_DIR=\"$HOME/.nvm\"" >>~/.zshrc
echo "[ -s \"$NVM_DIR/nvm.sh\" ] && \. \"$NVM_DIR/nvm.sh\"" >>~/.zshrc                   # This loads nvm
echo "[ -s \"$NVM_DIR/bash_completion\" ] && \. \\"$NVM_DIR/bash_completion\" >>~/.zshrc # This loads nvm bash_completion" >> ~/.zshrc

# manager specific npm version
nvm list-remote
nvm -v
nvm install v18.15.0
node -v
npm -version
npm install -g neovim

# 6. install ruby and gem
sudo apt-get install ruby-full
ruby -v
gem environment
gem install neovim

# 7. start nvim
nvim
