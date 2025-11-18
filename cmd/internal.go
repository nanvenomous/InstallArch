/*
Copyright © 2025 nanvenomous mrgarelli@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// Note: hostname, city, and username variables are declared in external.go

// setupGitCmd clones the bare unix repo
var setupGitCmd = &cobra.Command{
	Use:   "setup-git",
	Short: "Clone bare unix dotfiles repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		home := os.Getenv("HOME")
		c := exec.Command("git", "clone", "--bare", "https://github.com/nanvenomous/unix.git", filepath.Join(home, ".unx"))
		c.Dir = home
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

// checkoutGitCmd checks out the dotfiles
var checkoutGitCmd = &cobra.Command{
	Use:   "checkout-git",
	Short: "Checkout dotfiles from bare repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		home := os.Getenv("HOME")
		c := exec.Command("git", "--git-dir="+filepath.Join(home, ".unx")+"/", "--work-tree="+home, "checkout")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

// configureGitCmd configures git settings
var configureGitCmd = &cobra.Command{
	Use:   "configure-git",
	Short: "Configure git settings for dotfiles",
	RunE: func(cmd *cobra.Command, args []string) error {
		home := os.Getenv("HOME")
		gitDir := filepath.Join(home, ".unx") + "/"

		commands := [][]string{
			{"git", "config", "--global", "core.excludesfile", "~/.gitignore"},
			{"git", "config", "--global", "--includes", "include.path", "./.keybindings_git"},
			{"git", "--git-dir=" + gitDir, "--work-tree=" + home, "config", "--local", "status.showUntrackedFiles", "no"},
			{"git", "--git-dir=" + gitDir, "--work-tree=" + home, "config", "--get", "remote.origin.fetch"},
			{"git", "--git-dir=" + gitDir, "--work-tree=" + home, "config", "remote.origin.fetch", "+refs/heads/*:refs/remotes/origin/*"},
		}

		for _, cmdArgs := range commands {
			c := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				// Continue on error for --get command
				fmt.Printf("Warning: %v\n", err)
			}
		}

		return nil
	},
}

// installInternalCmd installs packages from internal_packages.txt
var installInternalCmd = &cobra.Command{
	Use:   "install-internal",
	Short: "Install packages from rsrc/internal_packages.txt",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile("./rsrc/internal_packages.txt")
		if err != nil {
			return fmt.Errorf("failed to read internal_packages.txt: %w", err)
		}

		packages := strings.Fields(string(data))
		if len(packages) == 0 {
			return fmt.Errorf("no packages found in internal_packages.txt")
		}

		cmdArgs := append([]string{"-Sy"}, packages...)
		c := exec.Command("pacman", cmdArgs...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin

		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to install packages: %w", err)
		}

		fmt.Println("Internal packages installed successfully")
		return nil
	},
}

// goInstallCmd installs Go packages
var goInstallCmd = &cobra.Command{
	Use:   "go-install",
	Short: "Install Go development tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		packages := []string{
			"github.com/go-task/task/v3/cmd/task@latest",
			"github.com/a-h/templ/cmd/templ@latest",
			"github.com/nanvenomous/e@latest",
			"github.com/nanvenomous/where-to@latest",
			"github.com/moson-mo/pacseek@latest",
			"github.com/ChausseBenjamin/termpicker@latest",
		}

		for _, pkg := range packages {
			fmt.Printf("Installing %s...\n", pkg)
			c := exec.Command("go", "install", pkg)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				return fmt.Errorf("failed to install %s: %w", pkg, err)
			}
		}

		fmt.Println("Go packages installed successfully")
		return nil
	},
}

// npmInstallCmd installs npm packages
var npmInstallCmd = &cobra.Command{
	Use:   "npm-install",
	Short: "Install npm language servers and tools",
	RunE: func(cmd *cobra.Command, args []string) error {
		packages := []string{
			"typescript-language-server",
			"typescript",
			"pyright",
			"vscode-langservers-extracted",
			"@tailwindcss/language-server",
		}

		cmdArgs := append([]string{"i", "-g"}, packages...)
		c := exec.Command("sudo", append([]string{"npm"}, cmdArgs...)...)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin

		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to install npm packages: %w", err)
		}

		fmt.Println("npm packages installed successfully")
		return nil
	},
}

// gitInstallCmd installs packages from git/AUR
var gitInstallCmd = &cobra.Command{
	Use:   "git-install",
	Short: "Install zsh-system-clipboard, yay, and bun",
	RunE: func(cmd *cobra.Command, args []string) error {
		home := os.Getenv("HOME")

		// Clone zsh-system-clipboard
		scriptsDir := filepath.Join(home, ".scripts")
		os.MkdirAll(scriptsDir, 0755)
		c := exec.Command("git", "clone", "https://github.com/kutsan/zsh-system-clipboard.git")
		c.Dir = scriptsDir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to clone zsh-system-clipboard: %w", err)
		}

		// Create projects directory and clone yay
		projectsDir := filepath.Join(home, "projects")
		os.MkdirAll(projectsDir, 0755)

		yayDir := filepath.Join(projectsDir, "yay")
		c = exec.Command("git", "clone", "https://aur.archlinux.org/yay.git")
		c.Dir = projectsDir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to clone yay: %w", err)
		}

		// Build and install yay
		c = exec.Command("makepkg", "-si")
		c.Dir = yayDir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to build yay: %w", err)
		}

		// Install bun
		c = exec.Command("bash", "-c", "curl -fsSL https://bun.sh/install | bash")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to install bun: %w", err)
		}

		fmt.Println("Git-based packages installed successfully")
		return nil
	},
}

// yayInstallCmd installs AUR packages using yay
var yayInstallCmd = &cobra.Command{
	Use:   "yay-install",
	Short: "Install AUR packages (font-symbola, enpass-bin)",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := exec.Command("yay", "-S", "font-symbola", "enpass-bin")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Stdin = os.Stdin
		return c.Run()
	},
}

// amdGPUCmd installs AMD GPU packages
var amdGPUCmd = &cobra.Command{
	Use:   "amd-gpu",
	Short: "Install AMD GPU drivers and ROCm packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		packages := [][]string{
			{"sudo", "pacman", "-S", "ollama-rocm"},
			{"sudo", "pacman", "-S", "xf86-video-amdgpu"},
			{"sudo", "pacman", "-S", "rocm-opencl-sdk", "rocm-hip-sdk"},
		}

		for _, cmdArgs := range packages {
			c := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			if err := c.Run(); err != nil {
				return fmt.Errorf("failed to install AMD packages: %w", err)
			}
		}

		fmt.Println("AMD GPU packages installed successfully")
		return nil
	},
}

// setClockCmd configures time synchronization
var setClockCmd = &cobra.Command{
	Use:   "set-clock",
	Short: "Configure time synchronization (default: America/Chicago)",
	RunE: func(cmd *cobra.Command, args []string) error {
		commands := [][]string{
			{"timedatectl", "set-ntp", "true"},
			{"timedatectl", "set-timezone", "America/Chicago"},
			{"sudo", "systemctl", "start", "systemd-timesyncd.service"},
			{"sudo", "systemctl", "enable", "systemd-timesyncd.service"},
		}

		for _, cmdArgs := range commands {
			c := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				return fmt.Errorf("failed to configure clock: %w", err)
			}
		}

		// Copy timezone dispatcher script
		if err := exec.Command("cp", "rsrc/09-timezone", "/etc/NetworkManager/dispatcher.d/").Run(); err != nil {
			fmt.Printf("Warning: failed to copy timezone script: %v\n", err)
		}

		fmt.Println("Clock configured successfully")
		return nil
	},
}

// sysSetupCmd performs system setup
var sysSetupCmd = &cobra.Command{
	Use:   "sys-setup",
	Short: "System setup: timezone, locale, hostname, services (default: fairytail, Chicago)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if hostname == "" {
			hostname = "fairytail"
		}
		if city == "" {
			city = "Chicago"
		}

		// Set timezone
		if err := os.Symlink("/usr/share/zoneinfo/America/"+city, "/etc/localtime"); err != nil {
			// File may already exist, remove and retry
			os.Remove("/etc/localtime")
			if err := os.Symlink("/usr/share/zoneinfo/America/"+city, "/etc/localtime"); err != nil {
				return fmt.Errorf("failed to set timezone: %w", err)
			}
		}

		// Set hardware clock
		if err := exec.Command("hwclock", "--systohc").Run(); err != nil {
			return fmt.Errorf("failed to set hardware clock: %w", err)
		}

		// Set locale
		if err := os.WriteFile("/etc/locale.conf", []byte("LANG=en_US.UTF-8\n"), 0644); err != nil {
			return fmt.Errorf("failed to write locale.conf: %w", err)
		}

		// Generate locale
		if err := exec.Command("locale-gen").Run(); err != nil {
			return fmt.Errorf("failed to generate locale: %w", err)
		}

		// Set hostname
		if err := os.WriteFile("/etc/hostname", []byte(hostname+"\n"), 0644); err != nil {
			return fmt.Errorf("failed to write hostname: %w", err)
		}

		// Set hosts file
		hostsContent := fmt.Sprintf(`127.0.0.1	localhost
::1		localhost
127.0.1.1	%s.localdomain  %s
`, hostname, hostname)
		if err := os.WriteFile("/etc/hosts", []byte(hostsContent), 0644); err != nil {
			return fmt.Errorf("failed to write hosts file: %w", err)
		}

		// Enable services
		services := []string{
			"bluetooth.service",
			"NetworkManager.service",
			"systemd-timesyncd.service",
			"sshd.service",
		}

		for _, service := range services {
			if err := exec.Command("systemctl", "enable", service).Run(); err != nil {
				fmt.Printf("Warning: failed to enable %s: %v\n", service, err)
			}
		}

		// Enable user pulseaudio
		if err := exec.Command("systemctl", "--user", "enable", "pulseaudio").Run(); err != nil {
			fmt.Printf("Warning: failed to enable pulseaudio: %v\n", err)
		}

		// Unmute audio
		if err := exec.Command("amixer", "sset", "Master", "unmute").Run(); err != nil {
			fmt.Printf("Warning: failed to unmute audio: %v\n", err)
		}

		// Set default browser
		if err := exec.Command("xdg-settings", "set", "default-web-browser", "org.qutebrowser.qutebrowser.desktop").Run(); err != nil {
			fmt.Printf("Warning: failed to set default browser: %v\n", err)
		}

		// Set root password
		fmt.Println("\nSet root password:")
		c := exec.Command("passwd")
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to set root password: %w", err)
		}

		fmt.Println("System setup completed successfully")
		return nil
	},
}

// createUserCmd creates a new user
var createUserCmd = &cobra.Command{
	Use:   "create-user",
	Short: "Create a new user account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if username == "" {
			return fmt.Errorf("username flag is required")
		}

		// Create user
		if err := exec.Command("useradd", "-m", "-g", "users", "-G", "wheel", username).Run(); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// Add to wheel group (redundant but matches original)
		if err := exec.Command("usermod", "-a", "-G", "wheel", username).Run(); err != nil {
			return fmt.Errorf("failed to add user to wheel: %w", err)
		}

		// Set password
		fmt.Printf("\nSet password for %s:\n", username)
		c := exec.Command("passwd", username)
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to set password: %w", err)
		}

		// Change shell to zsh (note: original has hardcoded 'gin')
		if err := exec.Command("chsh", "-s", "/bin/zsh", username).Run(); err != nil {
			return fmt.Errorf("failed to change shell: %w", err)
		}

		fmt.Printf("User %s created successfully\n", username)
		return nil
	},
}

// configUserCmd opens sudoers for editing
var configUserCmd = &cobra.Command{
	Use:   "config-user",
	Short: "Edit sudoers file to configure wheel group",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := exec.Command("nvim", "-c", "/wheel", "/etc/sudoers")
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		return c.Run()
	},
}

// grubSetupCmd installs and configures GRUB
var grubSetupCmd = &cobra.Command{
	Use:   "grub-setup",
	Short: "Install and configure GRUB bootloader",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := exec.Command("grub-install", "--target=x86_64-efi", "--efi-directory=/boot/efi").Run(); err != nil {
			return fmt.Errorf("failed to install grub: %w", err)
		}

		if err := exec.Command("grub-mkconfig", "-o", "/boot/grub/grub.cfg").Run(); err != nil {
			return fmt.Errorf("failed to generate grub config: %w", err)
		}

		fmt.Println("GRUB installed successfully")
		return nil
	},
}

// bootOrderCmd shows boot order information
var bootOrderCmd = &cobra.Command{
	Use:   "boot-order",
	Short: "Display EFI boot order",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := exec.Command("efibootmgr", "-v")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to show boot order: %w", err)
		}

		fmt.Println("\nTo change the boot order use:")
		fmt.Println("efibootmgr -o 0002,0001,0003")
		return nil
	},
}

// sshKeyCmd generates SSH key
var sshKeyCmd = &cobra.Command{
	Use:   "ssh-key",
	Short: "Generate SSH key",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := exec.Command("ssh-keygen")
		c.Stdin = os.Stdin
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to generate SSH key: %w", err)
		}

		fmt.Println("Now add the key to your GitHub account")
		return nil
	},
}

// fromSourceCmd builds packages from source
var fromSourceCmd = &cobra.Command{
	Use:   "from-source",
	Short: "Build where-to and e from source",
	RunE: func(cmd *cobra.Command, args []string) error {
		home := os.Getenv("HOME")
		projectsDir := filepath.Join(home, "projects")

		// Build where-to
		whereToDir := filepath.Join(projectsDir, "where-to")
		c := exec.Command("git", "clone", "https://github.com/nanvenomous/where-to.git")
		c.Dir = projectsDir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to clone where-to: %w", err)
		}

		makeCommands := [][]string{
			{"make"},
			{"sudo", "make", "install"},
			{"sudo", "make", "zsh-completions"},
		}

		for _, cmdArgs := range makeCommands {
			c = exec.Command(cmdArgs[0], cmdArgs[1:]...)
			c.Dir = whereToDir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			if err := c.Run(); err != nil {
				return fmt.Errorf("failed to build where-to: %w", err)
			}
		}

		// Build e
		eDir := filepath.Join(projectsDir, "e")
		c = exec.Command("git", "clone", "git@github.com:nanvenomous/e.git")
		c.Dir = projectsDir
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to clone e: %w", err)
		}

		for _, cmdArgs := range [][]string{{"make"}, {"sudo", "make", "install"}} {
			c = exec.Command(cmdArgs[0], cmdArgs[1:]...)
			c.Dir = eDir
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Stdin = os.Stdin
			if err := c.Run(); err != nil {
				return fmt.Errorf("failed to build e: %w", err)
			}
		}

		fmt.Println("Source packages built successfully")
		return nil
	},
}

// runAllInternalCmd runs all internal commands in sequence
var runAllInternalCmd = &cobra.Command{
	Use:   "run-all-internal",
	Short: "Run all internal setup commands in order (stops on first error)",
	Long: `Runs all internal setup commands in the proper sequence inside the new system.
This should be run after entering the new system with 'enter-sys'.

Steps:
1. install-internal (install packages)
2. sys-setup (system configuration)
3. create-user (create user account)
4. config-user (configure sudoers)
5. grub-setup (install bootloader)
6. setup-git (clone dotfiles)
7. checkout-git (checkout dotfiles)
8. configure-git (configure git)
9. go-install (install Go tools)
10. npm-install (install npm tools)
11. git-install (install from git/AUR)
12. set-clock (configure time)
13. ssh-key (generate SSH key)

Optional commands you can run manually:
- yay-install (AUR packages)
- amd-gpu (AMD GPU drivers)
- from-source (build from source)
- boot-order (show boot order)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if hostname == "" {
			hostname = "fairytail"
		}
		if city == "" {
			city = "Chicago"
		}
		if username == "" {
			return fmt.Errorf("--username flag is required for run-all-internal")
		}

		steps := []struct {
			name string
			fn   func() error
		}{
			{"install-internal", func() error { return installInternalCmd.RunE(cmd, []string{}) }},
			{"sys-setup", func() error { return sysSetupCmd.RunE(cmd, []string{}) }},
			{"create-user", func() error { return createUserCmd.RunE(cmd, []string{}) }},
			{"config-user", func() error { return configUserCmd.RunE(cmd, []string{}) }},
			{"grub-setup", func() error { return grubSetupCmd.RunE(cmd, []string{}) }},
			{"setup-git", func() error { return setupGitCmd.RunE(cmd, []string{}) }},
			{"checkout-git", func() error { return checkoutGitCmd.RunE(cmd, []string{}) }},
			{"configure-git", func() error { return configureGitCmd.RunE(cmd, []string{}) }},
			{"go-install", func() error { return goInstallCmd.RunE(cmd, []string{}) }},
			{"npm-install", func() error { return npmInstallCmd.RunE(cmd, []string{}) }},
			{"git-install", func() error { return gitInstallCmd.RunE(cmd, []string{}) }},
			{"set-clock", func() error { return setClockCmd.RunE(cmd, []string{}) }},
			{"ssh-key", func() error { return sshKeyCmd.RunE(cmd, []string{}) }},
		}

		for i, step := range steps {
			fmt.Printf("\n[%d/%d] Running %s...\n", i+1, len(steps), step.name)
			if err := step.fn(); err != nil {
				return fmt.Errorf("step '%s' failed: %w", step.name, err)
			}
			fmt.Printf("✓ %s completed\n", step.name)
		}

		fmt.Println("\n✓ All internal setup completed successfully!")
		fmt.Println("\nOptional steps you may want to run:")
		fmt.Println("  - InstallArch yay-install (AUR packages)")
		fmt.Println("  - InstallArch amd-gpu (AMD GPU drivers)")
		fmt.Println("  - InstallArch from-source (build custom packages)")
		fmt.Println("  - InstallArch boot-order (check boot order)")
		return nil
	},
}

func init() {
	// Add flags
	sysSetupCmd.Flags().StringVarP(&hostname, "hostname", "n", "fairytail", "System hostname")
	sysSetupCmd.Flags().StringVarP(&city, "city", "c", "Chicago", "City for timezone")

	createUserCmd.Flags().StringVarP(&username, "username", "u", "", "Username to create")
	createUserCmd.MarkFlagRequired("username")

	runAllInternalCmd.Flags().StringVarP(&hostname, "hostname", "n", "fairytail", "System hostname")
	runAllInternalCmd.Flags().StringVarP(&city, "city", "c", "Chicago", "City for timezone")
	runAllInternalCmd.Flags().StringVarP(&username, "username", "u", "", "Username to create")
	runAllInternalCmd.MarkFlagRequired("username")

	// Add commands to root
	rootCmd.AddCommand(setupGitCmd)
	rootCmd.AddCommand(checkoutGitCmd)
	rootCmd.AddCommand(configureGitCmd)
	rootCmd.AddCommand(installInternalCmd)
	rootCmd.AddCommand(goInstallCmd)
	rootCmd.AddCommand(npmInstallCmd)
	rootCmd.AddCommand(gitInstallCmd)
	rootCmd.AddCommand(yayInstallCmd)
	rootCmd.AddCommand(amdGPUCmd)
	rootCmd.AddCommand(setClockCmd)
	rootCmd.AddCommand(sysSetupCmd)
	rootCmd.AddCommand(createUserCmd)
	rootCmd.AddCommand(configUserCmd)
	rootCmd.AddCommand(grubSetupCmd)
	rootCmd.AddCommand(bootOrderCmd)
	rootCmd.AddCommand(sshKeyCmd)
	rootCmd.AddCommand(fromSourceCmd)
	rootCmd.AddCommand(runAllInternalCmd)
}
