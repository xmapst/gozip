package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"zip/ratelimit"
	"zip/symwalk"
	"zip/utils"
)

var (
	paths []string
	limit int64 = 2000 * 1024 // 限制输出 2000KB/s
)

func init() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stdout, "参数不全, 请检查后重试!")
		os.Exit(0)
	}
	paths = os.Args[1:]
}

func main() {
	// 检测当前进程数
	procTotal := utils.CountProcessIdByName(os.Args[0])
	if procTotal > 3 {
		_, _ = fmt.Fprintf(os.Stdout, "当前并发下载连接已达上限, 请稍后再试.")
		os.Exit(0)
	}
	_, _ = fmt.Fprintf(os.Stdout, "0")

	// debug code
	//var fw *os.File
	//fw, err = os.Create("xxx.zip")
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//zw := zip.NewWriter(ratelimit.Writer(fw, ratelimit.New(limit)))

	zw := zip.NewWriter(ratelimit.Writer(os.Stdout, ratelimit.New(limit)))
	defer zw.Close()

	for _, p := range paths {
		if runtime.GOOS == "windows" {
			if len(strings.Split(p, "C:")) != 0 {
				p = strings.ReplaceAll(p, "C:", "")
			}
			if len(strings.Split(p, "c:")) != 0 {
				p = strings.ReplaceAll(p, "c:", "")
			}
		}
		if !utils.FileOrPathExist(p) {
			continue
		}
		makeZip(p, zw)
	}
}

func makeZip(inFilepath string, zw *zip.Writer) error {
	return symwalk.Walk(inFilepath, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		// 目录拉平
		//relPath := strings.TrimPrefix(filePath, filepath.Dir(inFilepath))
		zipFile, err := zw.Create(filePath)
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
}
