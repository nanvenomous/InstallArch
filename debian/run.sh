#!/bin/bash

USAGE='''Commands:
	update
	aptInstall
	aptGetInstall
	setupGit
	checkoutGit
	configureGit
	zshSetup
'''

case "${1}" in
	"update")
		sudo apt-get install software-properties-common
		sudo add-apt-repository ppa:neovim-ppa/unstable
		sudo apt update
		sudo apt upgrade
		;;
	"aptInstall")
		sudo apt install build-essential zsh xauth
		;;
	"aptGetInstall")
		sudo apt-get install neovim exa neofetch xsel ripgrep zoxide software-properties-common
		sudo apt-get install ninja-build gettext libtool libtool-bin autoconf automake cmake g++ pkg-config unzip curl doxygen
		;;
	"setupGit")
		cd "${HOME}"
		git clone --bare 'https://github.com/nanvenomous/unix.git' "${HOME}/.unx"
		;;
	"checkoutGit")
		git --git-dir=${HOME}/.unx/ --work-tree=${HOME} checkout
		;;
	"configureGit")
		git config --global core.excludesfile ~/.gitignore
		git config --global --includes include.path './.keybindings_git'
		git --git-dir=${HOME}/.unx/ --work-tree=${HOME} config --local status.showUntrackedFiles no
		;;
	"zshSetup")
		git clone https://github.com/zsh-users/zsh-syntax-highlighting.git ~/.settings/zsh-syntax-highlighting
		git clone https://github.com/kutsan/zsh-system-clipboard.git ~/.settings/zsh-system-clipboard
		;;
	*)
		echo "${USAGE}"
		;;
esac

