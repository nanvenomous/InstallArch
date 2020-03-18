#!/bin/bash

vgremove vg0


sed -e 's/\s*\([\+0-9a-zA-Z]*\).*/\1/' << EOF | fdisk /dev/nvme0n1
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
