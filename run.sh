#!/bin/bash

USAGE='''Commands:
	terminalSetup
	checkEFI
	checkInternet
	setWifi
	reset <disk>
	setClock
		Default: Denver
	partitionDisk <disk> <bootSize> <rootSize> <swapSize> <homeSize>
		defaults[GB]: 1 boot, 20 root, 12 swap, rest of filesystem home
	format <partitionIdentifier>
		example: for partition /dev/sda1, partitionIdentifier=sda
		example: for partition /dev/nvme0n1p1 partitionIdentifier=nvme0n1p
	mounting <partitionIdentifier>
		example: for partition /dev/sda1, partitionIdentifier=sda
		example: for partition /dev/nvme0n1p1 partitionIdentifier=nvme0n1p
	update
	install
	tab
	enterSys
	internalInstall
	installWithNPM
	sysSetup <hostname> <city>
		defaults: ichiraku, Denver
	createUser <username>
	grubSetup
	bootOrder
	setupGit
  checkoutGit
  configureGit
	prepareReboot'''


function checkEFI() {
if [[ $(ls '/sys/firmware/efi/efivars' &>/dev/null; echo $?) -eq 0 ]]; then
	echo 'This is an EFI system.'
else
	echo 'This is not an EFI system. This guide does not apply :('
fi
}

function checkInternet() {
ping -q -c 1 google.com &>/dev/null
if [[ $? -ne 0 ]]; then
	echo 'There is no internet.'
else
	echo 'There is internet! Skip setWifi command.'
fi
}

function setWifi() {
iwctl
}

function reset() {
disk="${1}"
echo "disk to reformat: ${disk}"
umount -R /mnt
swapoff -a

sed -e 's/\s*\([\+0-9a-zA-Z]*\).*/\1/' << EOF | fdisk "${disk}"
	g # clear the in memory partition table
	w # write the partition table
	q # and we're done
EOF
}

function format() {
part="${1}"

# formatting
mkfs.fat -F32 /dev/${part}1
mkfs.ext4 /dev/${part}2
mkswap /dev/${part}3
mkfs.ext4 /dev/${part}4
}

function mounting() {
part="${1}"

# mounting
mount /dev/${part}2 /mnt
mkdir -p /mnt/boot
mkdir -p /mnt/boot/efi
mount /dev/${part}1 /mnt/boot/efi
mkdir -p /mnt/home
mount /dev/${part}4 /mnt/home
swapon /dev/${part}3
}

function userSetup() {
HOSTNM="${1}"
CITY="${2}"
ln -sf /usr/share/zoneinfo/America/${CITY} /etc/localtime
hwclock --systohc
echo "LANG=en_US>UTF-8" > /etc/local.conf
locale-gen
echo "${HOSTNM}" > /etc/hostname
}

function hostSetup() {
hostname="${1}"
HOSTS=$(cat <<-END
127.0.0.1	localhost
::1		localhost
127.0.1.1	${hostname}.localdomain  ${hostname}
END
)
echo "${HOSTS}" >> /etc/hosts
}

function customGit() {
dir="${1}"
git --git-dir=$HOME/${dir}/ --work-tree=$HOME "${@:2}"
}

case "${1}" in

	"terminalSetup")
		echo "set -o vi"
		echo "alias c='clear'"
		;;
	"checkEFI")
		checkEFI
		;;
	"checkInternet")
		checkInternet
		;;
	"setWifi")
		setWifi
		;;
	"reset")
		reset "${@:2}"
		;;
	"setClock")
		timedatectl set-ntp true
		timedatectl set-timezone America/Denver

    sudo systemctl start systemd-timesyncd.service
    sudo systemctl enable systemd-timesyncd.service
    cp rsrc/09-timezone /etc/NetworkManager/dispatcher.d/
		;;
	"partitionDisk")
		./partitioning.sh "${@:2}"
		;;
	"format")
		format "${@:2}"
		;;
	"mounting")
		mounting "${@:2}"
		;;
	"update")
		pacman -Syy
		pacman -Sy archlinux-keyring
		;;
	"install")
		pacstrap /mnt base base-devel linux linux-headers linux-docs linux-firmware efibootmgr grub networkmanager bash-completion zsh zsh-completions
		;;
	"tab")
		genfstab -U /mnt >> /mnt/etc/fstab
		;;
	"enterSys")
		arch-chroot /mnt # chroot into the system
		;;
	"internalInstall")
		pacman -Sy neovim git zsh zsh-completions zsh-syntax-highlighting dhcpcd dhclient man-db man-pages sudo openssh netctl tree dialog python3 python-pip i3-gaps i3status feh dmenu xorg-xinit xorg-server picom lxappearance unclutter alacritty pulseaudio pulseaudio-bluetooth pulseaudio-alsa alsa-utils bluez bluez-utils go gopls nodejs npm lxappearance xsel ripgrep lazygit neofetch exa zoxide entr xfce4-power-manager firefox bat
		;;
	"installWithNPM")
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
	"prepareReboot")
		umount -R /mnt
		swapoff -a
		;;
	*)
		echo "${USAGE}"
		;;
esac

