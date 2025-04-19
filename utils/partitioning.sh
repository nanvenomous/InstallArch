#!/bin/bash

disk="${1}"
bootSize="1" # usually 1
swapSize="8" # usually rest of disk

sed -e 's/\s*\([\+0-9a-zA-Z]*\).*/\1/' << EOF | fdisk ${disk}
	g # clear the in memory partition table
	n # new boot partition
	1 # partition number 1
	  # default - start at beginning of disk 
	+${bootSize}G # boot partition size
	t # type of the partition
	1 # 1 is the type for EFI
	n # new swap partition
	2 # partition number 2
	  # default, start immediately after preceding partition
	+${swapSize}G # swap size
	t # type of the partition
	2 # partition number 2
	19 # linux swap
	n # new root partition
	3 # partition number 3
	  # default, start immediately after preceding partition
	  # default, extend partition to end of disk
	t # type of the partition
	3 # partition number 3
	20 # linux filesystem
	p # print the in-memory partition table
	w # write the partition table
	q # and we're done
EOF
