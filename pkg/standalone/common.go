/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package standalone

import (
	"os"
	"os/exec"
	path_filepath "path/filepath"
	"runtime"
)

const (
	defaultDaprDirName       = ".dapr"
	defaultDaprDirNoHomeName = "dapr"
	defaultDaprBinDirName    = "bin"
	defaultComponentsDirName = "components"
	defaultConfigFileName    = "config.yaml"
)

func DefaultDaprDirPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// in some environments where os.UserHomeDir() can fail fall back to /usr/local/dapr
		return path_filepath.Join("usr", "local", defaultDaprDirNoHomeName)
	}
	return path_filepath.Join(homeDir, defaultDaprDirName)
}

func daprBinPath(daprDir string) string {
	return path_filepath.Join(daprDir, defaultDaprBinDirName)
}

func binaryFilePathWithDir(binaryDir string, binaryFilePrefix string) string {
	binaryPath := path_filepath.Join(binaryDir, binaryFilePrefix)
	if runtime.GOOS == daprWindowsOS {
		binaryPath += ".exe"
	}
	return binaryPath
}

func lookupBinaryFilePath(binaryFilePrefix string) string {
	binaryFile := binaryFilePrefix
	if runtime.GOOS == daprWindowsOS {
		binaryFile += ".exe"
	}

	// lookup the binary in $PATH first; if we cannot find it in $PATH, fall
	// back to the default dapr install directory in case the user has not added
	// that to $PATH
	path, err := exec.LookPath(binaryFile)
	if err != nil {
		return path_filepath.Join(daprBinPath(DefaultDaprDirPath()), binaryFile)
	}
	return path
}

func DaprComponentsPath(daprDir string) string {
	return path_filepath.Join(daprDir, defaultComponentsDirName)
}

func DaprConfigPath(daprDir string) string {
	return path_filepath.Join(daprDir, defaultConfigFileName)
}
