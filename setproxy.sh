#!/bin/bash

# 1. ============== GET HOST IP FOR WSL2 ===================
function GetHostIp() {
	ip=$(grep nameserver /etc/resolv.conf | awk '{print $2}')
	echo "the host ip is: $ip, then we ping it to test fireware is open or not"
	ping "$ip"
}

# 2. ============== SET PROXY 4 WSL2 USING HOST=============
function SetWSL2Proxy() {
	ip=$(grep nameserver /etc/resolv.conf | awk '{print $2}')
	port=${1:-"7890"}
	export http_proxy=http://$ip:$port
	export https_proxy=http://$ip:$port
	echo "set proxy by $ip:$port"
}

# 3. ============== SET/CANCER PROXY (LINUX MAC WSL2)===========
function unsetProxy() {

	unset http_proxy
	unset https_proxy

  npm config delete proxy
  npm config delete https-proxy

  git config --global --unset http.proxy
  git config --global --unset https.proxy
	echo "cancer proxy setting"
}

function setProxy() {
  # proxy infomation
  ip=${1:-"192.168.31.201"}
  port=${2:-"7890"}
  echo "$ip:$port"
  
  # set for default,wget,curl
  export http_proxy=http://$ip:$port
  export https_proxy=http://$ip:$port

  # set for npm
  npm config set proxy="http://$ip:$port"
  npm config set https-proxy="http://$ip:$port"

  # set for git
  git config --global https.proxy "http://$ip:$port"
  git config --global http.proxy "http://$ip:$port"
}

alias gethostip=GetHostIp
alias proxyon=SetWSL2Proxy
alias proxyoff=unsetProxy
alias proxyall=setProxy

unsetProxy
setProxy "192.168.31.201" "7890"
npm config ls -l | grep proxy
