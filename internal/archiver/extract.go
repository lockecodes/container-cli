package archiver

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractArchive extracts an archive file (either .tar.gz or .zip) to the target directory.
func ExtractArchive(source, target string) error {
	// Check the file extension to determine the archive type
	if strings.HasSuffix(source, ".tar.gz") {
		return extractTarGz(source, target)
	} else if strings.HasSuffix(source, ".zip") {
		return extractZip(source, target)
	} else {
		return fmt.Errorf("unsupported file type: %s", source)
	}
}

// extractTarGz extracts a .tar.gz archive to the target directory
func extractTarGz(source, target string) error {
	file, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %v", source, err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			// End of archive
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %v", err)
		}

		// Determine the path where the file will be extracted
		targetPath := filepath.Join(target, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			// Create directory
			if err := os.MkdirAll(targetPath, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", targetPath, err)
			}
		case tar.TypeReg:
			// Create regular file
			if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
				return fmt.Errorf("failed to create parent directory for file %s: %v", targetPath, err)
			}
			outFile, err := os.Create(targetPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %v", targetPath, err)
			}
			defer outFile.Close()

			if _, err := io.Copy(outFile, tarReader); err != nil {
				return fmt.Errorf("failed to write to file %s: %v", targetPath, err)
			}
		default:
			return fmt.Errorf("unsupported tar entry type: %c in file %s", header.Typeflag, source)
		}
	}
	return nil
}

// extractZip extracts a .zip archive to the target directory
func extractZip(source, target string) error {
	r, err := zip.OpenReader(source)
	if err != nil {
		return fmt.Errorf("failed to open zip file %s: %v", source, err)
	}
	defer r.Close()

	for _, file := range r.File {
		targetPath := filepath.Join(target, file.Name)

		if file.FileInfo().IsDir() {
			// Create directory
			if err := os.MkdirAll(targetPath, file.Mode()); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", targetPath, err)
			}
			continue
		}

		// Create file
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory for file %s: %v", targetPath, err)
		}
		outFile, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %v", targetPath, err)
		}
		defer outFile.Close()

		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file inside zip %s: %v", file.Name, err)
		}
		defer rc.Close()

		if _, err := io.Copy(outFile, rc); err != nil {
			return fmt.Errorf("failed to write to file %s: %v", targetPath, err)
		}
	}
	return nil
}
