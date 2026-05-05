#!/bin/bash
# requires: xorriso, squashfs-tools
set -euo pipefail

ISO_IN="iso/arch.iso"
ISO_OUT="iso/arch-custom.iso"
WORK="work"
SQUASHFS="${WORK}/extracted/arch/x86_64/airootfs.sfs"

echo "==> Extracting ISO filesystem..."
rm -rf "${WORK}/extracted"
mkdir -p "${WORK}"
xorriso -osirrox on -indev "${ISO_IN}" -extract / "${WORK}/extracted" \
  2>"${WORK}/xorriso.log" || { cat "${WORK}/xorriso.log"; exit 1; }
chmod -R u+w "${WORK}/extracted/"

echo "==> Extracting boot artifacts from original ISO..."
# Syslinux hybrid MBR boot code (bytes 0-431, before the partition table)
dd if="${ISO_IN}" bs=1 count=432 of="${WORK}/isohdpfx.bin" 2>/dev/null
# EFI system partition (appended after the ISO filesystem, sector range is version-specific)
EFI_RANGE=$(xorriso -indev "${ISO_IN}" -report_system_area cmd 2>/dev/null \
  | grep "append_partition 2 0xef" | grep -o '[0-9]*d-[0-9]*d' | head -1)
EFI_START=$(echo "${EFI_RANGE}" | cut -d- -f1 | tr -d 'd')
EFI_END=$(echo "${EFI_RANGE}" | cut -d- -f2 | tr -d 'd')
EFI_COUNT=$((EFI_END - EFI_START + 1))
dd if="${ISO_IN}" bs=512 skip="${EFI_START}" count="${EFI_COUNT}" of="${WORK}/efiboot.img" 2>/dev/null

echo "==> Unsquashing..."
rm -rf "${WORK}/airootfs"
unsquashfs -d "${WORK}/airootfs" "${SQUASHFS}"

echo "==> Injecting files into /root/..."
cp user_configuration.json "${WORK}/airootfs/root/"
cp rsrc/setup_dotfiles.sh "${WORK}/airootfs/root/"
chmod +x "${WORK}/airootfs/root/setup_dotfiles.sh"

echo "==> Repacking squashfs (this takes a few minutes)..."
rm "${SQUASHFS}"
mksquashfs "${WORK}/airootfs" "${SQUASHFS}" -noappend -comp xz -processors "$(nproc)"

echo "==> Updating checksum..."
(cd "${WORK}/extracted" && sha512sum arch/x86_64/airootfs.sfs > arch/x86_64/airootfs.sfs.sha512)

echo "==> Repacking ISO with hybrid boot structure..."
rm -f "${ISO_OUT}"
xorriso -as mkisofs \
  -iso-level 3 \
  -full-iso9660-filenames \
  -joliet \
  -joliet-long \
  -rational-rock \
  -volid "ARCH_$(date +%Y%m)" \
  -isohybrid-mbr "${WORK}/isohdpfx.bin" \
  --protective-msdos-label \
  -partition_offset 16 \
  --mbr-force-bootable \
  -append_partition 2 0xef "${WORK}/efiboot.img" \
  -appended_part_as_gpt \
  -iso_mbr_part_type 0x00 \
  -c boot/syslinux/boot.cat \
  -b boot/syslinux/isolinux.bin \
  -no-emul-boot \
  -boot-load-size 4 \
  -boot-info-table \
  -eltorito-alt-boot \
  -e '--interval:appended_partition_2:all::' \
  -no-emul-boot \
  -output "${ISO_OUT}" \
  "${WORK}/extracted/" \
  2>"${WORK}/xorriso-repack.log" || { cat "${WORK}/xorriso-repack.log"; exit 1; }

echo "==> Done: ${ISO_OUT}"
