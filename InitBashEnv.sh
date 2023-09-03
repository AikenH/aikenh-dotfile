#!/bin/bash

# if u using ubuntu, u make face error what cause by error shell(dash), change it to bash here
grep -q dash /bin/sh
if [[ $? -eq 0 ]];then
  echo "carry out function to change default shell"
  sudo dpkg-reconfigure dash 
  if [[ $? -ne 0 ]];then
    echo "manual change soft link to bash"
    sudo ln -snf /bin/bash /bin/sh
  fi
fi

echo "done"
