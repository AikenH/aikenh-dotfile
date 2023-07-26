# Intro

create by aikenhong,this repo setup my personal development environment.
i want make the script more general, consider any situation in diff kinds of sys.

those are some principle of this repo:

1. make it easy to use.
2. make it clear and light
3. make sys beautify and easy to use.

## Install Scripts 

the carry out logic should be like that.

**Install_env** as the main entry, take the params and call function.

- Install_basic. (this will carry out default, install core-clis and build-essential)
- Install_dev. (this will install the dev package like python,nvm,ruby)
- Install_zsh. (this script will install zsh, oh-my-zsh, and setup plugins)
- Install_Neovim. (this script will install the dependcies by **Install_Nv_Depend** \[**update_zsh**\])

And this process should check each step is carry out correctly. or it should be stop
After we carry out the prev-step we can carry out the script again to continue.
If we want it be perfect. we should throw out the error when it stop.  

## Uninstall Scripts

if something wrong, we may want the env clear. Provide some Uninstall Scripts.
for those software which is hard to remove clear(or with configs).

- Nvim 

## Dotfiles

Keep some dotfiles here, which make it easy 2 use in any machine. 
Here is some we need to keep.

- For_SHELL: ex_func, alias, env_setup
- For_vim: light config and powerful(extra_function)
- For_nvim: NvChad and personal-themes-and-keymaps and script to replace.
- For_Ranger:
- For_Tmux: keymap setting

