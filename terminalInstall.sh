#!/bin/bash
#FIXME: add conditional statement. and make sure the process didnnot stop.
#FIXME: make get the ip automaticly and echo that process
#FIXME: modify the zshrc for plugins including extra.
#FIXME: backup the origin version and rewrite one, when we hit the key word
#       we will modify the content. like plugins.

# sys inso
proxy=172.30.240.1:8890

# install zsh and oh-my-zsh and plugins 
sudo apt-get install zsh
sh -c "$(curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh --proxy $proxy)"
# generate zsh dotfile here. then install plugins.
bash 
conda init zsh # init zsh
cat ~/.bashrc | grep -i nvm >> ~/.zshrc # init zsh
# zsh and plugins will contain many alias include --color=auto
echo "alias hist='history -i'" >> ~/.zshrc # show time of history
echo "alias ltr='ls -rsthl'" >> ~/.zshrc
echo "alias cl='clear'" >> ~/.zshrc
echo "alias nv='nvim'" >> ~/.zshrc
echo "alias lsd='ls -d */'" >> ~/.zshrc

# install extra plugins from github.
git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-syntax-highlighting
git clone https://github.com/zsh-users/zsh-autosuggestions ${ZSH_CUSTOM:-~/.oh-my-zsh/custom}/plugins/zsh-autosuggestions

#FIXME: register plugins of zsh | may need awk or sed command. but using python maybe best.
# like get line index and modify it using sed -i or some.
line_num=`cat test.sh | awk "/plugins=\(git\)/{print NR}"`
# write a loop to write down the plugins<Left>
plugins=('git' 'zsh-syntax-highlighting' 'zsh-autosuggestions' 'colored-man-pages' 'sage-paste' 'themes' 'tmux' 'sudo' 'z')
sed -i "${line_num}d" test.sh
line_num=$[line_num-1]
sed -i "${line_num}a \)" test.sh
sed -i "${line_num}a     plugins=\(" test.sh

line_num=$[line_num+1]
echo $line_num
for phase in ${plugins[@]}; do
  sed -i "${line_num}a $phase" test.sh
done

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
