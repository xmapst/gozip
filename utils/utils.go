package utils

import (
	"github.com/shirou/gopsutil/v3/process"
	"os"
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
