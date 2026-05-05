# InstallArch

Builds a custom Arch Linux live USB with an `archinstall` configuration and dotfiles setup script baked in.

## Dependencies

```bash
task deps
```

Installs: `xorriso`, `squashfs-tools`, `curl`, `gnupg`, `archinstall`

## Usage

### Clone

```bash
git clone https://github.com/nanvenomous/InstallArch.git
cd InstallArch
```

### Full build

```bash
task
```

Downloads the latest Arch ISO, verifies its GPG signature, and bakes `user_configuration.json` and `setup_dotfiles.sh` into the squashfs. Produces `iso/arch-custom.iso`.

### Write to USB

```bash
task usb -- /dev/sdX
```

Shows the target device via `lsblk`, waits 5 seconds, then writes the ISO with `dd`. Replace `/dev/sdX` with your USB device.

### Edit archinstall config

```bash
task config
```

Opens the `archinstall` TUI with the current `user_configuration.json` loaded. Use **Save configuration** before quitting — the task copies the result back automatically.

## On boot

After booting the live USB, `/root/` contains:

- `user_configuration.json` — passed to archinstall with `archinstall --config /root/user_configuration.json`
- `setup_dotfiles.sh` — clones the bare dotfiles repo and checks out
