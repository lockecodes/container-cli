package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const gitlabAPIURL = "https://gitlab.com/api/v4/projects/gitlab-org%2Fgitlab-foss/releases"

type Release struct {
	TagName string `json:"tag_name"`
}

// getLatestRelease fetches the latest release information from GitLab
func getLatestRelease() (string, error) {
	resp, err := http.Get(gitlabAPIURL)
	if err != nil {
		return "", fmt.Errorf("failed to fetch releases: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var releases []Release
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(releases) == 0 {
		return "", fmt.Errorf("no releases found")
	}

	return releases[0].TagName, nil
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

func Update() {
	latestTag, err := getLatestRelease()
	if err != nil {
		fmt.Println("Error fetching latest release:", err)
		return
	}

	fmt.Println("Latest GitLab Release Tag:", latestTag)

	// Construct the download URL and file path
	binaryURL := fmt.Sprintf("https://gitlab.com/gitlab-org/gitlab-foss/-/releases/%s/downloads/gitlab-%s.tar.gz", latestTag, latestTag)
	destPath := filepath.Join(".", fmt.Sprintf("gitlab-%s.tar.gz", latestTag))

	fmt.Println("Downloading binary from:", binaryURL)

	if err := downloadFile(binaryURL, destPath); err != nil {
		fmt.Println("Error downloading binary:", err)
		return
	}

	fmt.Println("Downloaded binary to:", destPath)
}
