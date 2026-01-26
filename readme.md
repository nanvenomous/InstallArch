# Arch Linux Installation Instructions

to get the cli on a new system run

```bash
curl -L https://github.com/nanvenomous/InstallArch/releases/latest/download/InstallArch > InstallArch
chmod +x InstallArch
```

# Extra

### Audio Visualizers

- [projectM](https://wiki.archlinux.org/title/ProjectM)
- [cava](https://github.com/karlstav/cava)

### XPS 13

- [xps 13 9310](<https://wiki.archlinux.org/title/Dell_XPS_13_(9310)>)

### framework

- [Framework_Laptop_13](https://wiki.archlinux.org/title/Framework_Laptop_13)
  - [framework-system](https://archlinux.org/packages/extra/x86_64/framework-system/)
  - [docs](https://github.com/FrameworkComputer/linux-docs/tree/main/framework12)

screen rotation

```bash
yay -S rot8
```

Create a config file at `~/.config/rot8/rot8.toml`:

```ini
[laptop]
display = "eDP-1"
touchscreen = "eDP-1"
threshold = 0.5
```

then `rot8 &`

### Sway/Wayland

#### Core

sudo pacman -S sway swaylock swayidle swaybg

#### Wayland equivalents for common tools

sudo pacman -S waybar # Status bar (alternative to i3status/i3bar)
sudo pacman -S grim slurp # Screenshots (alternative to maim/scrot)
sudo pacman -S wl-clipboard # Clipboard (alternative to xclip/xsel)
sudo pacman -S wofi # App launcher (alternative to dmenu/rofi)

#### Optional but recommended

sudo pacman -S xorg-xwayland # Run X11 apps on Wayland
sudo pacman -S polkit-kde-agent # For authentication dialogs
Copy your i3 config
mkdir -p ~/.config/sway
cp ~/.config/i3/config ~/.config/sway/config
Edit the Sway config
Open ~/.config/sway/config and make these minimal changes:

### for screenshare

```bash
sudo pacman -S xdg-desktop-portal xdg-desktop-portal-wlr # screen share
systemctl --user start xdg-desktop-portal.service
systemctl --user enable xdg-desktop-portal.service
systemctl --user start xdg-desktop-portal-wlr.service
systemctl --user enable xdg-desktop-portal-wlr.service
```
