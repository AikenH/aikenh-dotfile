#!/bin/bash

# 1. ============== GET HOST IP FOR WSL2 ===================
function GetHostIp(){
  ip=$(cat /etc/resolv.conf|grep nameserver|awk '{print $2}')
  echo "the host ip is: $ip, then we ping it to test fireware is open or not"
  ping $ip
}

# 2. ============== SET PROXY 4 WSL2 USING HOST=============
function SetWSL2Proxy(){
  ip=$(cat /etc/resolv.conf|grep nameserver|awk '{print $2}')
  port=${1:-"7890"}
  export http_proxy=http://$ip:$port
  export https_proxy=http://$ip:$port
  echo "set proxy by $ip:$port"
}

# 3. ============== CANCER PROXY (LINUX MAC WSL2)===========
function unsetProxy(){
  unset http_proxy
  unset https_proxy
  echo "cancer proxy setting"
}

alias gethostip=GetHostIp
alias proxyon=SetWSL2Proxy
alias proxyoff=unsetProxy


