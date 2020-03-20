#!/bin/bash

USAGE='''Commands:
	setupGit
	reset
	partition
	format
	mounting
	update
	install
	enterSys
	internalInstall
	sysSetup
	bootOrder
	grubSetup
	prepareReboot'''

disk='/dev/nvme0n1'
efi="${disk}p1"
boot="${disk}p2"
system="${disk}p3"
vol='vg0'
swapSize='16G'
root="/dev/${vol}/root"
swap="/dev/${vol}/swap"
bootDir="/mnt/boot"
efiDir="/mnt/boot/efi"
hostname="ua"
username="deku"
city='Detroit'

HOSTS=$(cat <<-END
127.0.0.1	localhost
::1		localhost
127.0.1.1	${hostname}.localdomain  ${hostname}
END
)

function setupGit() {
git config --global user.email "mrgarelli@gmail.com"
git config --global user.name "Matthew Garelli"
}

function reset() {
umount -R /mnt
swapoff "${swap}"
vgremove "${vol}"

sed -e 's/\s*\([\+0-9a-zA-Z]*\).*/\1/' << EOF | fdisk "${disk}"
	g # clear the in memory partition table
	w # write the partition table
	q # and we're done
EOF
}

function partition() {
sed -e 's/\s*\([\+0-9a-zA-Z]*\).*/\1/' << EOF | fdisk "${disk}"
	g # clear the in memory partition table
	n # new boot partition
	1 # partition number 1
	  # default - start at beginning of disk 
	+250M # boot partition size
	t # type of the partition
	1 # 1 is the type for EFI
	n # new partition
	2 # partion number 2
	  # default, start immediately after preceding partition
	+512M # default, extend partition to end of disk
	n # new partition
	3 # partion number 3
	  # default, start immediately after preceding partition
	  # default, extend partition to end of disk
	p # print the in-memory partition table
	w # write the partition table
	q # and we're done
EOF
partprobe "${disk}"
}

function format() {
pvcreate "${system}"
vgcreate "${vol}" "${system}"

lvcreate -L "${swapSize}" "${vol}" -n swap
lvcreate -l 100%FREE "${vol}" -n root

mkfs.fat -F32 "${efi}"
mkfs.ext4 "${boot}"
mkfs.ext4 "${root}"

mkswap "${swap}"
}

function mounting() {
mount "${root}" /mnt
mkdir "${bootDir}"
mount "${boot}" "${bootDir}"
mkdir "${efiDir}"
mount "${efi}" "${efiDir}"
swapon "${swap}"
}

case "${1}" in
	"setupGit")
		setupGit
		;;
	"reset")
		reset
		;;
	"partition")
		partition
		;;
	"format")
		format
		;;
	"mounting")
		mounting
		;;
	"update")
		pacman -Syy
		pacman -Sy archlinux-keyring
		;;
	"install")
		pacstrap /mnt base base-devel linux linux-firmware efibootmgr grub networkmanager
		;;
	"enterSys")
		genfstab -U /mnt >> /mnt/etc/fstab
		arch-chroot /mnt # chroot into the system
		;;
	"internalInstall")
		pacman -Sy vim git dhcpcd dhclient networkmanager man-db man-pages sudo openssh netctl dialog python3 python-pip xonsh i3-gaps xorg-xinit xorg-server picom lxappearance pcmanfm code unclutter konsole firefox
		;;
	"sysSetup")
		ln -sf /usr/share/zoneinfo/America/${city} /etc/localtime
		hwclock --systohc

		echo "${hostname}" > /etc/hostname

		echo "LANG=en_US>UTF-8" >> /etc/local.conf
		locale-gen

		echo "${HOSTS}" >> /etc/hosts

		systemctl enable NetworkManager
		passwd

		useradd -m -g users -Gwheel "${username}"
		passwd "${username}"
		;;
	"bootOrder")
		efibootmgr -v
		echo "To change the boot order use:"
		echo "efibootmgr -o 0002,0001,0003"
		;;
	"grubSetup")
		grub-install
		grub-mkconfig -o /boot/grub/grub.cfg
		;;
	"prepareReboot")
		umount -R /mnt
		swapoff -a
		;;
	*)
		echo "${USAGE}"
		;;
esac
