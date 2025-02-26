package linuxTools

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/develar/app-builder/pkg/download"
	"github.com/develar/app-builder/pkg/util"
	"github.com/develar/go-fs-util"
)

func GetAppImageToolDir() (string, error) {
	dirName := "appimage-13.0.1"
	//noinspection SpellCheckingInspection
	result, err := download.DownloadArtifact("",
		download.GetGithubBaseUrl()+dirName+"/"+dirName+".7z",
		"ZG8U7K9Bk71cvP1VDlP+L7hO+HhRTJW6RO0kLgh5hbbJJHhPfoA/kw1hsFeq1pAyez6MxvoDyL/5/O45hX9Jaw==")
	if err != nil {
		return "", err
	}
	return result, nil
}

func GetAppImageToolBin(toolDir string) string {
	if util.GetCurrentOs() == util.MAC {
		return filepath.Join(toolDir, "darwin")
	} else {
		return filepath.Join(toolDir, "linux-"+goArchToArchSuffix())
	}
}

func GetLinuxTool(name string) (string, error) {
	toolDir, err := GetAppImageToolDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(GetAppImageToolBin(toolDir), name), nil
}

func GetMksquashfs() (string, error) {
	result := "mksquashfs"
	if !util.IsEnvTrue("USE_SYSTEM_MKSQUASHFS") {
		result = os.Getenv("MKSQUASHFS_PATH")
		if len(result) == 0 {
			var err error
			result, err = GetLinuxTool("mksquashfs")
			if err != nil {
				return "", err
			}
		}
	}

	return result, nil
}

func goArchToArchSuffix() string {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		return "x64"
	case "386":
		return "ia32"
	case "arm":
		return "arm32"
	default:
		return arch
	}
}

func ReadDirContentTo(dir string, paths []string, filter func(string) bool) ([]string, error) {
	content, err := fsutil.ReadDirContent(dir)
	if err != nil {
		return nil, err
	}

	for _, value := range content {
		if filter == nil || filter(value) {
			paths = append(paths, filepath.Join(dir, value))
		}
	}
	return paths, nil
}
