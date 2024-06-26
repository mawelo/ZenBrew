/**
 * @file package.go
 * @description
 * @author
 * @copyright
 */

package pkg

import (
	zb_types "OpenCortex/ZenBrew/types"
	"OpenCortex/ZenBrew/utils"
	"encoding/json"
	"fmt"
	"io"
	log "log/slog"
	"net/http"
	"os"
	"os/exec"
	"path"
)

type Package struct {
	zb_types.Package
}

type PackageLink struct {
	zb_types.PackageLink
}

type InstalledPackage struct {
	zb_types.InstalledPackage
}

func FromInstalled(installed_package InstalledPackage) Package {
	var pkg Package
	pkg.Name = installed_package.Name
	pkg.Latest = ""
	pkg.Format = installed_package.Format
	pkg.Maintainer = installed_package.Maintainer
	pkg.Versions = append(pkg.Versions, installed_package.Version)
	return pkg
}

func DownloadPackageMetadata(package_link PackageLink) Package {
	json_url := package_link.URL
	log.Info(fmt.Sprintf("Downloading package metadata from: %s", json_url))
	//hash_url := package_link.URL + "package.sha256"

	json_bytes := utils.DownloadFile(json_url)
	//hash_bytes := utils.DownloadFile(hash_url)

	//if !utils.CheckHash(json_bytes, hash_bytes) {
	//	log.Error("Hashes do not match.")
	//	panic("Hashes do not match.")
	//}

	var pkg Package
	err := json.Unmarshal(json_bytes, &pkg)
	if err != nil {
		log.Error("Failed to unmarshal JSON:", err)
		panic("Failed to unmarshal JSON")
	}

	return pkg
}

func (pkg Package) Download(version string) int {
	if version == "" || version == "latest" {
		version = pkg.Latest
	}
	var version_int int
	for i, v := range pkg.Versions {
		if v.Version == version {
			version_int = i
			break
		}
	}
	package_url := pkg.Versions[version_int].URL
	package_path := path.Join(utils.Preferences.RootDir, "ZenBrew", pkg.Name)

	// Download the package
	resp, err := http.Get(package_url)
	if err != nil {
		log.Error("Failed to download package:", err)
		panic("Failed to download package")
	}
	defer resp.Body.Close()

	// Create the directory if it doesn't exist
	err = os.MkdirAll(package_path, os.ModePerm)
	if err != nil {
		log.Error("Failed to create package directory:", err)
		panic("Failed to create package directory")
	}

	// Create the tar.gz file
	file_path := path.Join(package_path, "package.tar.gz")
	file, err := os.Create(file_path)
	if err != nil {
		log.Error("Failed to create package file:", err)
		panic("Failed to create package file")
	}
	defer file.Close()

	// Copy the response body to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Error("Failed to save package file:", err)
		panic("Failed to save package file")
	}

	// Extract the tar.gz file
	err = utils.ExtractTar(file_path, package_path)
	if err != nil {
		log.Error("Failed to extract package:", err)
		panic("Failed to extract package")
	}

	return version_int
}

func (pkg Package) Install() {
	package_path := path.Join(utils.Preferences.RootDir, "ZenBrew", pkg.Name)

	// Run the install file as a subprocess
	log.Info(fmt.Sprintf("Running install file for package: %s", pkg.Name))
	cmd := exec.Command(fmt.Sprintf("%s/install", package_path))
	cmd_err := cmd.Run()
	if cmd_err != nil {
		log.Error(fmt.Sprintf("Failed to run install file: %s", cmd_err))
		panic("Failed to run install file")
	}
}

func (pkg Package) Uninstall() {
	package_path := path.Join(utils.Preferences.RootDir, "ZenBrew", pkg.Name)

	// Run the install file as a subprocess
	log.Info(fmt.Sprintf("Running uninstall file for package: %s", pkg.Name))
	cmd := exec.Command(fmt.Sprintf("%s/uninstall", package_path))
	err := cmd.Run()
	if err != nil {
		log.Error("Failed to run install file:", err)
		panic("Failed to run install file")
	}
	os.RemoveAll(package_path)
}

func (pkg Package) Update() {
	package_path := path.Join(utils.Preferences.RootDir, "ZenBrew", pkg.Name)

	// Run the install file as a subprocess
	log.Info(fmt.Sprintf("Running update file for package: %s", pkg.Name))
	cmd := exec.Command(fmt.Sprintf("%s/update", package_path))
	err := cmd.Run()
	if err != nil {
		log.Error("Failed to run install file:", err)
		panic("Failed to run install file")
	}
}
