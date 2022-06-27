# Arch Linux Installation Instructions

customizable Step-by-step instructions to install arch linux

* command completion
```
. ./completions/ext_ins
. ./completions/int_ins
```

* commands to run from live usb before entering arch system
```
./ext_ins

Commands:
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
	prepareReboot
```

* commands to run on newly installed system
```
./int_ins

Commands:
	install
	npmInstall
	sysSetup <hostname> <city>
		defaults: ichiraku, Denver
	createUser <username>
	grubSetup
	bootOrder
	setupGit
  checkoutGit
  configureGit
```

Note: run the commands in sequence
