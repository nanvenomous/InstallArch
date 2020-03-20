#!/bin/bash

USAGE='''Commands:
	checkEFI
	checkInternet
	setWifi
	reset
	setClock
	partitionDisk <disk> <bootSize> <rootSize> <swapSize> <homeSize>
		defaults[GB]: 1 boot, 20 root, 12 swap, rest of filesystem home
	format <partitionIdentifier>
		example: for partition /dev/sda1, partitionIdentifier=sda
	mountInstall <partitionIdentifier>
		example: for partition /dev/sda1, partitionIdentifier=sda
	install
	enterSys
	sysSetup <hostname> <city>
		defaults: UA, Detroit
	reboot'''

HOSTS='''
127.0.0.1	localhost
::1			localhost
127.0.1.1	myhostname.localdomain  myhostname
'''

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
	wifi-menu
	echo "now run sudo netctl start <service>"
}

function reset() {
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

function mountInstall() {
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
	ln -sf /usr/share/zoneinfo/America/${CITY:=Detroit} /etc/localtime
	hwclock --systohc
	echo "LANG=en_US>UTF-8" > /etc/local.conf
	locale-gen
	echo "${HOSTNM:=UA}" > /etc/hostname
}

function appendHosts() {
	echo "${HOSTS}" >> /etc/hosts
}


case "${1}" in

	"checkEFI")
		checkEFI
		;;
	"checkInternet")
		checkInternet
		;;
	"setWifi")
		reset
		;;
	"reset")
		setWifi
		;;
	"setClock")
		timedatectl set-ntp true
		timedatectl set-timezone America/Detroit
		;;
	"partitionDisk")
		./partitioning.sh "${2}" "${3}" "${4}" "${5}" "${6}"
		;;
	"format")
		format "${2}"
		;;
	"mountInstall")
		mountInstall "${2}"
		;;
	"install")
		pacman -Syy
		pacstrap /mnt base base-devel linux linux-firmware efibootmgr vim git dhcpcd dhclient networkmanager man-db man-pages sudo openssh grub netctl dialog python3 python-pip xonsh i3-gaps xorg-xinit xorg-server picom lxappearance pcmanfm code unclutter konsole firefox
		;;
	"enterSys")
		genfstab -U /mnt >> /mnt/etc/fstab
		arch-chroot /mnt # chroot into the system
		;;
	"sysSetup")
		userSetup "${2}" "${3}"

		appendHosts

		systemctl enable NetworkManager
		passwd
		grub-install --target=x86_64-efi --efi-directory=/boot/efi
		grub-mkconfig -o /boot/grub/grub.cfg
		exit
		umount -R /mnt
		;;
	"reboot")
		reboot
		;;
	*)
		echo "${USAGE}"
		;;
esac
