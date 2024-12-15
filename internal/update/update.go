package update

import (
	"encoding/json"
	"fmt"
	"gitlab.com/locke-codes/container-cli/internal/archiver"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

const gitlabAPIURL = "https://gitlab.com/api/v4/projects/47137983/releases"

type Release struct {
	Name            string    `json:"name"`
	TagName         string    `json:"tag_name"`
	Description     string    `json:"description"`
	CreatedAt       time.Time `json:"created_at"`
	ReleasedAt      time.Time `json:"released_at"`
	UpcomingRelease bool      `json:"upcoming_release"`
	Assets          struct {
		Count   int `json:"count"`
		Sources []struct {
			Format string `json:"format"`
			Url    string `json:"url"`
		} `json:"sources"`
		Links []struct {
			Id             int    `json:"id"`
			Name           string `json:"name"`
			Url            string `json:"url"`
			DirectAssetUrl string `json:"direct_asset_url"`
			LinkType       string `json:"link_type"`
		} `json:"links"`
	} `json:"assets"`
}

func sortReleasesByReleaseDateDesc(releases []Release) {
	// Sort the releases in descending order based on the ReleasedAt field
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].ReleasedAt.After(releases[j].ReleasedAt)
	})
}

// getLatestRelease fetches the latest release information from GitLab
func getLatestRelease() (Release, error) {
	log.Println("Fetching latest release from GitLab")
	resp, err := http.Get(gitlabAPIURL)
	if err != nil {
		return Release{}, fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Release{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return Release{}, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(releases) == 0 {
		return Release{}, fmt.Errorf("no releases found")
	}
	//order the releases by ReleasedAt desc
	sortReleasesByReleaseDateDesc(releases)

	return releases[0], nil
}

// downloadFile downloads a file from the given URL to the specified path
func downloadFile(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// mapArch converts runtime.GOARCH values to desired arch formats (e.g., "amd64" -> "x86_64").
func mapArch(arch string) string {
	switch arch {
	case "amd64":
		return "x86_64"
	case "arm64":
		return "arm64"
	// Add more mappings as needed
	default:
		return arch
	}
}

func getReleaseDownloadLink(release Release) (string, string, error) {
	// Get the OS and architecture of the executing machine
	runtime_os := runtime.GOOS
	arch := mapArch(runtime.GOARCH) // Convert the arch to match the desired format

	// Build the expected substring to match in the name of the asset
	title := cases.Title(language.AmericanEnglish)
	searchKey := fmt.Sprintf("%s_%s", title.String(runtime_os), arch)

	// Iterate through the links to find the one that matches the OS and Arch
	for _, link := range release.Assets.Links {
		if strings.Contains(link.Name, searchKey) {
			// Return the direct asset URL if a match is found
			return link.DirectAssetUrl, link.Name, nil
		}
	}

	return "", "", fmt.Errorf("no matching asset found for OS: %s, Arch: %s", runtime_os, arch)
}

func DownloadBinary(destinationPath string) (string, string, error) {
	latestRelease, err := getLatestRelease()
	if err != nil {
		fmt.Println("Error fetching latest release:", err)
		return "", "", err
	}

	fmt.Println("Latest GitLab Release Tag:", latestRelease.TagName)

	link, name, err := getReleaseDownloadLink(latestRelease)
	if err != nil {
		fmt.Println("Error:", err)
		return "", "", err
	}
	// Construct the download URL and file path
	destPath := filepath.Join(destinationPath, name)

	fmt.Println("Downloading binary from:", link)

	if err := downloadFile(link, destPath); err != nil {
		fmt.Println("Error downloading binary:", err)
		return "", "", err
	}

	fmt.Println("Downloaded binary to:", destPath)
	return destPath, latestRelease.TagName, nil
}

func Update() {
	_, _, err := DownloadBinary(".")
	if err != nil {
		log.Fatal(err)
	}
}
func Install() error {
	binaryPath, tagName, err := DownloadBinary("/tmp")
	if err != nil {
		return err
	}
	err = InstallBinary(binaryPath, "container-cli", tagName)
	if err != nil {
		return err
	}
	return nil
}

// InstallBinary extracts an archive and installs the binary into a versioned folder with a symlink.
// - Symlink `~/.local/bin/ccli` always points to the latest version.
// - Additionally, `/usr/local/bin/ccli` is symlinked to `~/.local/bin/ccli` for global access.
// Parameters:
// - `source` is the path to the archive file (e.g., `.tar.gz` or `.zip`).
// - `binaryName` is the name of the binary file in the archive (e.g., "ccli").
// - `version` version of the binary
func InstallBinary(source, binaryName string, version string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %v", err)
	}

	// Define paths
	baseDir := filepath.Join(homeDir, ".local", "bin")
	versionDir := filepath.Join(baseDir, "container-cli", version)
	localSymlinkPath := filepath.Join(baseDir, "ccli")
	globalSymlinkPath := filepath.Join("/usr/local/bin", "ccli")

	// Step 1: Extract the archive
	fmt.Printf("Extracting %s...\n", source)
	if err := archiver.ExtractArchive(source, versionDir); err != nil {
		return fmt.Errorf("failed to extract archive: %v", err)
	}

	// Step 2: Locate the binary file
	fmt.Println("Locating the binary...")
	binaryPath, err := findBinary(versionDir, binaryName)
	if err != nil {
		return fmt.Errorf("failed to locate binary %s: %v", binaryName, err)
	}

	// Step 3: Move the binary to the versioned folder
	fmt.Println("Installing the binary...")
	finalBinaryPath := filepath.Join(versionDir, binaryName)
	if err := os.Rename(binaryPath, finalBinaryPath); err != nil {
		return fmt.Errorf("failed to move binary to versioned directory: %v", err)
	}

	// Make the binary executable
	if err := os.Chmod(finalBinaryPath, 0755); err != nil {
		return fmt.Errorf("failed to make binary executable: %v", err)
	}

	// Step 4: Create/update the symlink in ~/.local/bin
	fmt.Println("Updating local symlink...")
	if err := updateSymlink(finalBinaryPath, localSymlinkPath); err != nil {
		return fmt.Errorf("failed to update local symlink: %v", err)
	}

	// Step 5: Create/update the global symlink in /usr/local/bin
	fmt.Println("Updating global symlink...")
	// For now just output the command for the symlink. If the user already has
	// ~/.local/bin in path then it should already work
	fmt.Println("You must either ensure that ~/.local/bin is in your path or run the following command:")
	fmt.Printf("sudo ln -s %s %s\n", localSymlinkPath, globalSymlinkPath)
	//if err := updateSymlink(localSymlinkPath, globalSymlinkPath); err != nil {
	//	return fmt.Errorf("failed to update global symlink: %v", err)
	//}

	fmt.Println("Installation successful!")
	return nil
}

// findBinary searches for the binary file in the extracted directory
func findBinary(directory, binaryName string) (string, error) {
	var binaryPath string
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Match the binary name
		if info.Mode().IsRegular() && info.Name() == binaryName {
			binaryPath = path
			return filepath.SkipDir // Stop searching once found
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if binaryPath == "" {
		return "", fmt.Errorf("binary %s not found in extracted files", binaryName)
	}
	return binaryPath, nil
}

// updateSymlink updates the symlink to point to the latest target.
// - `target` is the file for the symlink to point to.
// - `symlinkPath` is the path where the symlink should be created.
func updateSymlink(target, symlinkPath string) error {
	// Remove the symlink if it already exists
	if _, err := os.Lstat(symlinkPath); err == nil {
		if err := os.Remove(symlinkPath); err != nil {
			return fmt.Errorf("failed to remove existing symlink: %v", err)
		}
	}

	// Create the new symlink
	if err := os.Symlink(target, symlinkPath); err != nil {
		return fmt.Errorf("failed to create symlink: %v", err)
	}

	// Verify the symlink
	resolvedPath, err := os.Readlink(symlinkPath)
	if err != nil {
		return fmt.Errorf("failed to verify symlink: %v", err)
	}
	if resolvedPath != target {
		return fmt.Errorf("symlink was not set correctly: expected %s, got %s", target, resolvedPath)
	}

	return nil
}
