#!/bin/bash

USAGE='''Commands:
	internalInstall
	npmInstall
	sysSetup <hostname> <city>
		defaults: ichiraku, Denver
	createUser <username>
	grubSetup
	bootOrder
	setupGit
  checkoutGit
  configureGit'''

case "${1}" in
	"internalInstall")
    pacman -Sy $(cat ./rsrc/internal_packages.txt | tr '\n' ' ' )
		# pacman -Sy neovim git zsh zsh-completions zsh-syntax-highlighting dhcpcd dhclient man-db man-pages sudo openssh netctl tree dialog python3 python-pip i3-gaps i3status feh dmenu xorg-xinit xorg-server picom lxappearance unclutter alacritty pulseaudio pulseaudio-bluetooth pulseaudio-alsa alsa-utils bluez bluez-utils go gopls nodejs npm lxappearance xsel ripgrep lazygit neofetch exa zoxide entr xfce4-power-manager firefox bat

		;;
	"npmInstall")
		sudo npm i -g typescript-language-server typescript pyright
		;;
	"sysSetup")
		hostname="${2}"
		city="${3}"
		userSetup "${hostname:=ichiraku}" "${city:=Denver}"
		hostSetup "${hostname:=ichiraku}"

		systemctl enable bluetooth.service
		systemctl enable NetworkManager.service
		systemctl --user enable pulseaudio
		amixer sset Master unmute
		passwd
		;;
	"createUser")
		username="${2}"
		useradd -m -g users -G wheel "${username}"
		passwd "${username}"
		;;
	"grubSetup")
		grub-install --target=x86_64-efi --efi-directory=/boot/efi
		grub-mkconfig -o /boot/grub/grub.cfg
		;;
	"bootOrder")
		efibootmgr -v
		echo "To change the boot order use:"
		echo "efibootmgr -o 0002,0001,0003"
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
	*)
		echo "${USAGE}"
		;;
esac

