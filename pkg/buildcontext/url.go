/*
Copyright 2018 Google LLC

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

package buildcontext

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/GoogleContainerTools/kaniko/pkg/constants"
	"github.com/GoogleContainerTools/kaniko/pkg/util"
)

// url unifies calls to download and unpack the build context.
type Url struct {
	context string
}

// UnpackTarFromBuildContext download and untar a file from s3
func (u *Url) UnpackTarFromBuildContext() (string, error) {
	directory := constants.BuildContextDir
	tarPath := filepath.Join(directory, constants.ContextTar)
	if err := os.MkdirAll(directory, 0750); err != nil {
		return directory, err
	}
	if err := DownloadFile(tarPath, u.context); err != nil {
		return directory, err
	}
	if err := util.UnpackCompressedTar(tarPath, directory); err != nil {
		return directory, err
	}
	// Remove the tar so it doesn't interfere with subsequent commands
	return directory, os.Remove(tarPath)
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}
