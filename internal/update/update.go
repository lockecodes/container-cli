package update

import (
	"encoding/json"
	"fmt"
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

func Update() {
	latestRelease, err := getLatestRelease()
	if err != nil {
		fmt.Println("Error fetching latest release:", err)
		return
	}

	fmt.Println("Latest GitLab Release Tag:", latestRelease.TagName)

	link, name, err := getReleaseDownloadLink(latestRelease)
	if err != nil {
		fmt.Println("Error:", err)
	}
	// Construct the download URL and file path
	destPath := filepath.Join(".", name)

	fmt.Println("Downloading binary from:", link)

	if err := downloadFile(link, destPath); err != nil {
		fmt.Println("Error downloading binary:", err)
		return
	}

	fmt.Println("Downloaded binary to:", destPath)
}
