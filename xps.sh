#!/bin/bash

disk='/dev/nvme0n1'
efi="${disk}p1"
system="${disk}p2"
vol='vg0'
swapSize='16G'
root="/dev/${vol}/root"
swap="/dev/${vol}/swap"
bootDir="/mnt/boot"

vgremove "${vol}"

umount -R /mnt

sed -e 's/\s*\([\+0-9a-zA-Z]*\).*/\1/' << EOF | fdisk "${disk}"
	g # clear the in memory partition table
	n # new boot partition
	1 # partition number 1
	  # default - start at beginning of disk 
	+512M # boot partition size
	t # type of the partition
	1 # 1 is the type for EFI
	n # new partition
	2 # partion number 2
	  # default, start immediately after preceding partition
	  # default, extend partition to end of disk
	p # print the in-memory partition table
	w # write the partition table
	q # and we're done
EOF
partprobe "${disk}"

pvcreate "${system}"
vgcreate "${vol}" "${system}"

lvcreate -L "${swapSize}" "${vol}" -n swap
lvcreate -l 100%FREE "${vol}" -n root

mkfs.ext4 "${root}"
mount "${root}" /mnt

mkfs.fat -F32 "${efi}"
mkdir "${bootDir}"
mount "${efi}" "${bootDir}"

mkswap "${swap}"
swapon "${swap}"
