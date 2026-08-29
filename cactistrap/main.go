package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

func verifySHA256(path string, expected string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	sum := sha256.Sum256(data)
	actual := hex.EncodeToString(sum[:])

	return actual == expected
}

func downloadFile(url string, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	file, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func main() {
	if os.Geteuid() != 0 {
		fmt.Println("cactistrap must be run as root")
		return
	}
	if len(os.Args) < 2 {
		fmt.Println("usage: cactistrap <mountpoint>")
		return
	}

	target := os.Args[1]

	info, err := os.Stat(target)

	if err != nil {
		fmt.Println("target does not exists:", target)
		return
	}

	if !info.IsDir() {
		fmt.Println("target is not a directory:", target)
		return
	}

	cmd := exec.Command("mountpoint", "-q", target)

	if err := cmd.Run(); err != nil {
		fmt.Println("target is not mounted:", target)
		return
	}

	url := "https://github.com/cacti7322/cactios/releases/download/v0.1/cactios-rootfs.tar.xz"
	dest := "/tmp/cactios-rootfs.tar.xz"

	fmt.Println("Downloading rootfs...")

	if err := downloadFile(url, dest); err != nil {
		fmt.Println("download failed:", err)
		return
	}

	fmt.Println("Download complete")

	expectedSHA := "sha256:38aa12c7625bcf9c983ab0eead48ed1ea27da5a152fb3e24fa94f1de0d820b40" 

	fmt.Println("Verifying rootfs...")

	if !verifySHA256(dest, expectedSHA) {
		fmt.Println("rootfs checksum failed")
		return
	}
	fmt.Println("Extracting rootfs...")

	extractCmd := exec.Command("tar", "-xJf", dest, "-C", target)

	output, err := extractCmd.CombinedOutput()
	if err != nil {
		fmt.Println("failed to extract rootfs:", err)
		fmt.Println(string(output))
		return
	}

	if err := os.Remove(dest); err != nil {
		fmt.Println("warning: could not remove temporary archive:", err)
	}

	fmt.Println("Base system installed successfully")
}
