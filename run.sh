#!/bin/bash

USAGE='''Commands:
	checkEFI
	checkInternet
	setWifi
	setClock
	partitionDisk <disk> <bootSize> <rootSize> <swapSize> <homeSize>
		defaults[GB]: 1 boot, 20 root, 12 swap, rest of filesystem home'''

function checkEFI() {
	if [[ $(ls '/sys/firmware/efi/efivars' &>/dev/null; echo $?) -eq 0 ]]; then
		echo 'This is an EFI system.'
	else
		echo 'This is not an EFI system.'
	fi
}

function checkInternet() {
	ping -q -c 1 google.com &>/dev/null
	if [[ $? -ne 0 ]]; then
		echo 'There is no internet.'
	else
		echo 'There is internet!'
	fi
}

function setWifi() {
	wifi-menu
	echo "now run sudo netctl start <service>"
}

case "${1}" in

	"checkEFI")
		checkEFI
		;;
	"checkInternet")
		checkInternet
		;;
	"setWifi")
		setWifi
		;;
	"setClock")
		timedatectl set-ntp true
		;;
	"partitionDisk")
		./partitioning.sh "${2}" "${3}" "${4}" "${5}" "${6}"
		;;
	*)
		echo "${USAGE}"
		;;
esac
