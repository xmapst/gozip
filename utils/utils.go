package utils

import (
	"github.com/shirou/gopsutil/v3/process"
	"os"
	"runtime"
	"strings"
)

func FileOrPathExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func CountProcessIdByName(name string) int {
	var count int
	p, _ := process.Processes()
	for _, v := range p {
		n, err := v.Name()
		if err != nil {
			continue
		}
		if n == name {
			count += 1
		}
	}
	return count
}

// ToLinuxPath windows 路径转Linux格式路径
func ToLinuxPath(path string) string {
	if runtime.GOOS == "windows" {
		if len(strings.Split(path, "C:")) != 0 {
			path = strings.ReplaceAll(path, "C:", "")
		}
		if len(strings.Split(path, "c:")) != 0 {
			path = strings.ReplaceAll(path, "c:", "")
		}
		path = strings.ReplaceAll(path, `\\`, "/")
		path = strings.ReplaceAll(path, `\`, "/")
	}
	return path
}
