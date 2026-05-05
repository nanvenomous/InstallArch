#!/bin/bash
# requires: xorriso, squashfs-tools
set -euo pipefail

ISO_IN="iso/arch.iso"
ISO_OUT="iso/arch-custom.iso"
WORK="work"
SFS="${WORK}/airootfs.sfs"

echo "==> Extracting squashfs from ISO..."
mkdir -p "${WORK}"
xorriso -osirrox on -indev "${ISO_IN}" \
  -extract /arch/x86_64/airootfs.sfs "${SFS}" \
  2>"${WORK}/xorriso.log" || { cat "${WORK}/xorriso.log"; exit 1; }

echo "==> Unsquashing..."
rm -rf "${WORK}/airootfs"
unsquashfs -d "${WORK}/airootfs" "${SFS}"

echo "==> Injecting files into /root/..."
cp user_configuration.json "${WORK}/airootfs/root/"
cp rsrc/setup_dotfiles.sh "${WORK}/airootfs/root/"
chmod +x "${WORK}/airootfs/root/setup_dotfiles.sh"

echo "==> Repacking squashfs (this takes a few minutes)..."
rm "${SFS}"
mksquashfs "${WORK}/airootfs" "${SFS}" -noappend -comp xz -processors "$(nproc)"

echo "==> Updating checksum..."
(cd "${WORK}" && sha512sum airootfs.sfs > airootfs.sfs.sha512)

echo "==> Baking into ISO (preserving boot structure)..."
rm -f "${ISO_OUT}"
xorriso \
  -indev "${ISO_IN}" \
  -outdev "${ISO_OUT}" \
  -map "${SFS}" /arch/x86_64/airootfs.sfs \
  -map "${WORK}/airootfs.sfs.sha512" /arch/x86_64/airootfs.sfs.sha512 \
  -end 2>"${WORK}/xorriso.log" || { cat "${WORK}/xorriso.log"; exit 1; }

echo "==> Done: ${ISO_OUT}"
