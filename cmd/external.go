/*
Copyright © 2025 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var (
	disk     string
	bootSize string
	swapSize string
	username string
	hostname string
	city     string
)

// getRAMBasedSwapSize reads RAM from /proc/meminfo and calculates swap size
// Logic: if RAM < 2GB, use 2*RAM; if RAM <= 8GB, use RAM; else use RAM/2
// Returns swap size in GB (rounded up)
func getRAMBasedSwapSize() (string, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return "", fmt.Errorf("failed to open /proc/meminfo: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			// Parse MemTotal line (format: "MemTotal:       16384000 kB")
			fields := strings.Fields(line)
			if len(fields) < 2 {
				return "", fmt.Errorf("unexpected MemTotal format")
			}

			memKB, err := strconv.ParseFloat(fields[1], 64)
			if err != nil {
				return "", fmt.Errorf("failed to parse memory value: %w", err)
			}

			// Convert KB to GB
			memGB := memKB / 1024 / 1024

			// Calculate swap size based on RAM
			var swapGB float64
			if memGB < 2 {
				swapGB = memGB * 2
			} else if memGB <= 8 {
				swapGB = memGB
			} else {
				swapGB = memGB / 2
			}

			// Round up to nearest integer
			swapGBInt := int(math.Ceil(swapGB))
			return strconv.Itoa(swapGBInt), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading /proc/meminfo: %w", err)
	}

	return "", fmt.Errorf("MemTotal not found in /proc/meminfo")
}

func setAutoSwapOrDefault() {
	defaultSwapSize := "32"
	autoSwap, err := getRAMBasedSwapSize()
	if err != nil {
		fmt.Printf(" Falling back to default %sGB swap size: %v\n", defaultSwapSize, err)
		swapSize = defaultSwapSize
	} else {
		swapSize = autoSwap
		fmt.Printf("󰄬 Auto-calculated swap size based on system RAM: %sGB\n", swapSize)
	}
}

// checkInputsCmd represents the checkEFI command
var checkInputsCmd = &cobra.Command{
	Use:   "check-inputs",
	Short: "Check default inputs before run-all",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat("/sys/firmware/efi/efivars"); err != nil {
			fmt.Println(" This is not an EFI system. This guide does not apply :(")
		} else {
			fmt.Println("󰄬 This is an EFI system.")
		}

		setAutoSwapOrDefault()

		c := exec.Command("ping", "-q", "-c", "1", "google.com")
		if err := c.Run(); err != nil {
			fmt.Println(" There is no internet")
		}
		fmt.Println("󰄬 There is internet! Skip setWifi command.")

		return nil
	},
}

// checkInternetCmd represents the checkInternet command
var checkInternetCmd = &cobra.Command{
	Use:   "check-internet",
	Short: "Check if there is internet connectivity",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := exec.Command("ping", "-q", "-c", "1", "google.com")
		if err := c.Run(); err != nil {
			return fmt.Errorf("There is no internet")
		}
		fmt.Println("There is internet! Skip setWifi command.")
		return nil
	},
}

// setWifiCmd represents the setWifi command
var setWifiCmd = &cobra.Command{
	Use:   "set-wifi",
	Short: "Launch iwctl to configure WiFi",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := exec.Command("iwctl")
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset disk partition table",
	RunE: func(cmd *cobra.Command, args []string) error {
		if disk == "" {
			return fmt.Errorf("disk flag is required")
		}

		fmt.Printf("disk to reformat: %s\n", disk)

		// Unmount and turn off swap
		exec.Command("umount", "-R", "/mnt").Run()
		exec.Command("swapoff", "-a").Run()

		// Reset partition table
		fdiskScript := `g
w
q
`
		c := exec.Command("fdisk", disk)
		c.Stdin = strings.NewReader(fdiskScript)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

// partitionDiskCmd represents the partitionDisk command
var partitionDiskCmd = &cobra.Command{
	Use:   "partition-disk",
	Short: "Partition the disk (defaults: 1GB boot, auto-calculated swap based on RAM, rest for root)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if disk == "" {
			return fmt.Errorf("disk flag is required")
		}

		if bootSize == "" {
			bootSize = "1"
		}
		if swapSize == "" {
			setAutoSwapOrDefault()
		}

		fdiskScript := fmt.Sprintf(`g
n
1

+%sG
t
1
n
2

+%sG
t
2
19
n
3


t
3
20
p
w
q
`, bootSize, swapSize)

		c := exec.Command("fdisk", disk)
		c.Stdin = strings.NewReader(fdiskScript)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

// formatCmd represents the format command
var formatCmd = &cobra.Command{
	Use:   "format [partition-identifier]",
	Short: "Format the partitions (e.g., 'sda' for /dev/sda1, 'nvme0n1p' for /dev/nvme0n1p1)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		part := args[0]

		// Format boot partition (FAT32)
		if err := exec.Command("mkfs.fat", "-F32", "/dev/"+part+"1").Run(); err != nil {
			return fmt.Errorf("failed to format boot partition: %w", err)
		}

		// Format swap partition
		if err := exec.Command("mkswap", "/dev/"+part+"2").Run(); err != nil {
			return fmt.Errorf("failed to format swap partition: %w", err)
		}

		// Format root partition (ext4)
		if err := exec.Command("mkfs.ext4", "/dev/"+part+"3").Run(); err != nil {
			return fmt.Errorf("failed to format root partition: %w", err)
		}

		fmt.Println("All partitions formatted successfully")
		return nil
	},
}

// mountingCmd represents the mounting command
var mountingCmd = &cobra.Command{
	Use:   "mounting [partition-identifier]",
	Short: "Mount the partitions (e.g., 'sda' for /dev/sda1, 'nvme0n1p' for /dev/nvme0n1p1)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		part := args[0]

		// Mount root partition
		if err := exec.Command("mount", "/dev/"+part+"3", "/mnt").Run(); err != nil {
			return fmt.Errorf("failed to mount root partition: %w", err)
		}

		// Create and mount boot/efi
		os.MkdirAll("/mnt/boot/efi", 0755)
		if err := exec.Command("mount", "/dev/"+part+"1", "/mnt/boot/efi").Run(); err != nil {
			return fmt.Errorf("failed to mount boot partition: %w", err)
		}

		// Create home directory
		os.MkdirAll("/mnt/home", 0755)

		// Enable swap
		if err := exec.Command("swapon", "/dev/"+part+"2").Run(); err != nil {
			return fmt.Errorf("failed to enable swap: %w", err)
		}

		fmt.Println("All partitions mounted successfully")
		return nil
	},
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update pacman and keyring",
	RunE: func(cmd *cobra.Command, args []string) error {
		pacmanCmd := exec.Command("pacman", "-Syy")
		pacmanCmd.Stderr = os.Stderr
		pacmanCmd.Stdout = os.Stdout
		if err := pacmanCmd.Run(); err != nil {
			return fmt.Errorf("failed to sync package databases: %w", err)
		}

		pacmanCmd = exec.Command("pacman", "-Sy", "archlinux-keyring")
		pacmanCmd.Stderr = os.Stderr
		pacmanCmd.Stdout = os.Stdout
		if err := pacmanCmd.Run(); err != nil {
			return fmt.Errorf("failed to update keyring: %w", err)
		}

		fmt.Println("System updated successfully")
		return nil
	},
}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install packages from rsrc/external_packages.txt to /mnt",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Read external packages file
		data, err := os.ReadFile("./rsrc/external_packages.txt")
		if err != nil {
			return fmt.Errorf("failed to read external_packages.txt: %w", err)
		}

		// Parse packages
		packages := strings.Fields(string(data))
		if len(packages) == 0 {
			return fmt.Errorf("no packages found in external_packages.txt")
		}

		// Build pacstrap command
		cmdArgs := append([]string{"/mnt"}, packages...)
		c := exec.Command("pacstrap", cmdArgs...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin

		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to run pacstrap: %w", err)
		}

		fmt.Println("Packages installed successfully")
		return nil
	},
}

// copyBinaryCmd copies the InstallArch binary and rsrc directory to the new system
var copyBinaryCmd = &cobra.Command{
	Use:   "copy-binary",
	Short: "Copy InstallArch binary and resources to /mnt/root",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the current executable path
		exePath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("failed to get executable path: %w", err)
		}

		// Copy binary to /mnt/root
		if err := exec.Command("cp", exePath, "/mnt/root/InstallArch").Run(); err != nil {
			return fmt.Errorf("failed to copy binary: %w", err)
		}

		// Make it executable
		if err := os.Chmod("/mnt/root/InstallArch", 0755); err != nil {
			return fmt.Errorf("failed to make binary executable: %w", err)
		}

		// Copy rsrc directory
		if err := exec.Command("cp", "-r", "./rsrc", "/mnt/root/").Run(); err != nil {
			return fmt.Errorf("failed to copy rsrc directory: %w", err)
		}

		fmt.Println("Binary and resources copied to /mnt/root successfully")
		return nil
	},
}

// tabCmd represents the tab command
var tabCmd = &cobra.Command{
	Use:   "tab",
	Short: "Generate fstab file",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Open fstab file for appending
		f, err := os.OpenFile("/mnt/etc/fstab", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open fstab: %w", err)
		}
		defer f.Close()

		// Run genfstab
		c := exec.Command("genfstab", "-U", "/mnt")
		c.Stdout = f
		c.Stderr = os.Stderr

		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to generate fstab: %w", err)
		}

		fmt.Println("fstab generated successfully")
		return nil
	},
}

// enterSysCmd represents the enterSys command
var enterSysCmd = &cobra.Command{
	Use:   "enter-sys",
	Short: "Chroot into the new system",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := exec.Command("arch-chroot", "/mnt")
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

// prepareRebootCmd represents the prepareReboot command
var prepareRebootCmd = &cobra.Command{
	Use:   "prepare-reboot",
	Short: "Unmount partitions and disable swap before reboot",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := exec.Command("umount", "-R", "/mnt").Run(); err != nil {
			fmt.Printf("Warning: failed to unmount: %v\n", err)
		}

		if err := exec.Command("swapoff", "-a").Run(); err != nil {
			fmt.Printf("Warning: failed to disable swap: %v\n", err)
		}

		fmt.Println("System prepared for reboot")
		return nil
	},
}

// runAllCmd runs all commands in sequence
var runAllCmd = &cobra.Command{
	Use:   "run-all",
	Short: "Run complete installation from start to finish (stops on first error)",
	Long: `Runs the complete Arch Linux installation process:

EXTERNAL STEPS (on install medium):
1. check-internet
2. reset (requires --disk flag)
3. partition-disk (requires --disk flag, optional --boot-size and --swap-size)
4. format
5. mounting
6. update
7. install
8. copy-binary (copy InstallArch to new system)
9. tab

INTERNAL STEPS (inside new system via chroot):
10. install-internal
11. sys-setup (requires --hostname and --city)
12. create-user (requires --username)
13. config-user
14. grub-setup
15. setup-git
16. checkout-git
17. configure-git
18. go-install
19. npm-install
20. git-install
21. set-clock
22. ssh-key

FINAL STEP:
23. prepare-reboot

Required flags: --disk, --username
Optional flags: --hostname (default: fairytail), --city (default: Chicago), --boot-size (default: 1), --swap-size (default: 32)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if disk == "" {
			return fmt.Errorf("--disk flag is required for run-all")
		}
		if username == "" {
			return fmt.Errorf("--username flag is required for run-all")
		}
		if hostname == "" {
			hostname = "fairytail"
		}
		if city == "" {
			city = "Chicago"
		}

		// Determine partition identifier from disk
		// e.g., /dev/sda -> sda, /dev/nvme0n1 -> nvme0n1p
		partID := strings.TrimPrefix(disk, "/dev/")
		if strings.Contains(partID, "nvme") {
			partID += "p"
		}

		fmt.Println("=== STARTING EXTERNAL INSTALLATION (Install Medium) ===")

		externalSteps := []struct {
			name string
			fn   func() error
		}{
			{"check-internet", func() error { return checkInternetCmd.RunE(cmd, []string{}) }},
			{"reset", func() error { return resetCmd.RunE(cmd, []string{}) }},
			{"partition-disk", func() error { return partitionDiskCmd.RunE(cmd, []string{}) }},
			{"format", func() error { return formatCmd.RunE(cmd, []string{partID}) }},
			{"mounting", func() error { return mountingCmd.RunE(cmd, []string{partID}) }},
			{"update", func() error { return updateCmd.RunE(cmd, []string{}) }},
			{"install", func() error { return installCmd.RunE(cmd, []string{}) }},
			{"copy-binary", func() error { return copyBinaryCmd.RunE(cmd, []string{}) }},
			{"tab", func() error { return tabCmd.RunE(cmd, []string{}) }},
		}

		for i, step := range externalSteps {
			fmt.Printf("\n[EXTERNAL %d/%d] Running %s...\n", i+1, len(externalSteps), step.name)
			if err := step.fn(); err != nil {
				return fmt.Errorf("external step '%s' failed: %w", step.name, err)
			}
			fmt.Printf("✓ %s completed\n", step.name)
		}

		fmt.Println("\n=== STARTING INTERNAL INSTALLATION (Inside New System) ===")
		fmt.Println("Executing commands inside chroot environment...")

		// Build the internal command to run inside chroot
		internalCmd := fmt.Sprintf("cd /root && ./InstallArch run-all-internal --username=%s --hostname=%s --city=%s",
			username, hostname, city)

		c := exec.Command("arch-chroot", "/mnt", "bash", "-c", internalCmd)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin

		if err := c.Run(); err != nil {
			return fmt.Errorf("internal installation failed: %w\nYou can manually enter the system with 'InstallArch enter-sys' and run:\n  cd /root && ./InstallArch run-all-internal --username=%s --hostname=%s --city=%s", err, username, hostname, city)
		}

		fmt.Println("\n=== FINALIZING INSTALLATION ===")

		// Prepare for reboot
		fmt.Println("\n[FINAL] Running prepare-reboot...")
		if err := prepareRebootCmd.RunE(cmd, []string{}); err != nil {
			return fmt.Errorf("prepare-reboot failed: %w", err)
		}
		fmt.Println("✓ prepare-reboot completed")

		fmt.Println("\n" + strings.Repeat("=", 60))
		fmt.Println("✓ INSTALLATION COMPLETE!")
		fmt.Println(strings.Repeat("=", 60))
		fmt.Println("\nYour Arch Linux system is ready to boot!")
		fmt.Println("\nNext steps:")
		fmt.Println("  1. Remove the installation medium")
		fmt.Println("  2. Reboot with: reboot")
		fmt.Println("\nOptional post-install steps (run after reboot):")
		fmt.Println("  - InstallArch yay-install (AUR packages)")
		fmt.Println("  - InstallArch amd-gpu (if you have AMD GPU)")
		fmt.Println("  - InstallArch from-source (build custom packages)")
		fmt.Println("  - InstallArch boot-order (adjust boot order)")

		return nil
	},
}

func init() {
	// Add flags
	resetCmd.Flags().StringVarP(&disk, "disk", "d", "", "Disk to reset (e.g., /dev/sda)")
	resetCmd.MarkFlagRequired("disk")

	partitionDiskCmd.Flags().StringVarP(&disk, "disk", "d", "", "Disk to partition (e.g., /dev/sda)")
	partitionDiskCmd.Flags().StringVarP(&bootSize, "boot-size", "b", "1", "Boot partition size in GB")
	partitionDiskCmd.Flags().StringVarP(&swapSize, "swap-size", "s", "", "Swap partition size in GB")
	partitionDiskCmd.MarkFlagRequired("disk")

	runAllCmd.Flags().StringVarP(&disk, "disk", "d", "", "Disk to install to (e.g., /dev/sda)")
	runAllCmd.Flags().StringVarP(&bootSize, "boot-size", "b", "1", "Boot partition size in GB")
	runAllCmd.Flags().StringVarP(&swapSize, "swap-size", "s", "", "Swap partition size in GB")
	runAllCmd.Flags().StringVarP(&username, "username", "u", "", "Username to create")
	runAllCmd.Flags().StringVarP(&hostname, "hostname", "n", "fairytail", "System hostname")
	runAllCmd.Flags().StringVarP(&city, "city", "c", "Chicago", "City for timezone")
	runAllCmd.MarkFlagRequired("disk")
	runAllCmd.MarkFlagRequired("username")

	// Add commands to root
	rootCmd.AddCommand(checkInputsCmd)
	rootCmd.AddCommand(checkInternetCmd)
	rootCmd.AddCommand(setWifiCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(partitionDiskCmd)
	rootCmd.AddCommand(formatCmd)
	rootCmd.AddCommand(mountingCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(copyBinaryCmd)
	rootCmd.AddCommand(tabCmd)
	rootCmd.AddCommand(enterSysCmd)
	rootCmd.AddCommand(prepareRebootCmd)
	rootCmd.AddCommand(runAllCmd)
}
