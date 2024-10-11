#!/bin/bash

disk="${1}"
bootSize="1" # usually 1
rootSize="200" # usually 20
swapSize="8" # usually rest of disk

homeSize="" # default is to end of disk

sed -e 's/\s*\([\+0-9a-zA-Z]*\).*/\1/' << EOF | fdisk ${disk}
	g # clear the in memory partition table
	n # new boot partition
	1 # partition number 1
	  # default - start at beginning of disk 
	+${bootSize}G # boot partition size
	t # type of the partition
	1 # 1 is the type for EFI
	n # new root partition
	2 # partition number 2
	  # default, start immediately after preceding partition
	+${rootSize}G # root size
	t # type of the partition
	2 # partition number 2
	20 # linux filesystem
	n # new swap partition
	3 # partition number 3
	  # default, start immediately after preceding partition
	+${swapSize}G # swap size
	t # type of the partition
	3 # partition number 3
	19 # linux swap
	n # new partition
	4 # partion number 4
	  # default, start immediately after preceding partition
	${homeSize} # default, extend partition to end of disk
	p # print the in-memory partition table
	w # write the partition table
	q # and we're done
EOF
