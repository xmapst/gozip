package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/xmapst/gozip/ratelimit"
	"github.com/xmapst/gozip/symwalk"
	"github.com/xmapst/gozip/utils"
	"io"
	"os"
	"strings"
	"time"
)

var (
	paths []string
	limit int64 = 2000 * 1024 // 限制输出 2000KB/s
)

func init() {
	if len(os.Args) < 2 {
		_, _ = fmt.Fprintf(os.Stderr, "参数不全, 请检查后重试!")
		os.Exit(0)
	}
	paths = os.Args[1:]
}

func main() {
	// 检测当前进程数
	procTotal := utils.CountProcessIdByName(os.Args[0])
	if procTotal > 3 {
		_, _ = fmt.Fprintf(os.Stderr, "当前并发下载连接已达上限, 请稍后再试.")
		os.Exit(0)
	}
	go closeEvent()
	zw := zip.NewWriter(ratelimit.Writer(os.Stdout, ratelimit.New(limit)))
	defer zw.Close()

	for _, p := range paths {
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
		var zwPath = utils.ToLinuxPath(filePath)
		// 去除路径以根开始, 解决7z及windows资源管理器打开为空问题
		zipFile, err := zw.Create(strings.TrimPrefix(zwPath, "/"))
		if err != nil {
			return err
		}
		fsFile, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer fsFile.Close()
		_, err = io.Copy(zipFile, fsFile)
		if err != nil {
			return err
		}
		return nil
	})
}

func closeEvent() {
	input := bufio.NewReader(os.Stdin)
	for {
		bs, err := input.ReadByte()
		if err != nil {
			if err == io.EOF {
				return
			}
			continue
		}
		switch bs {
		case 4:
			os.Exit(0)
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
