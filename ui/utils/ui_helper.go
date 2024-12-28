package utils

import (
	"io/fs"
	"os"
	"path/filepath"
)

func ReadFolderForVaults(name string) ([]fs.DirEntry, error) {
	files, err := os.ReadDir(name)
	if err != nil {
		return nil, err
	}

	var vaultFiles []fs.DirEntry

	for _, file := range files {
		if !file.IsDir() && (filepath.Ext(file.Name()) == ".vault") {
			vaultFiles = append(vaultFiles, file)
			file.Info()
		}
	}

	return vaultFiles, nil
}
