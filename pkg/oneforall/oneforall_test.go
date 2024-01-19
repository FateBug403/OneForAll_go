package oneforall

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestOneForAll(t *testing.T) {
	ofaRun,err:=NewOneForAll(Options{
		ExePath: "E:\\Pentest\\FateBugScan\\FateBug\\pkg\\tools\\domain\\OneForAll",
	})
	if err != nil {
		log.Println(err)
		return
	}
	domains,err:=ofaRun.GetSubDomains([]string{"wesvr.cn"})
	if err != nil {
		log.Println(err)
	}
	log.Println(domains)
}

func TestFunc(t *testing.T) {
	// 获取可执行文件的路径
	exePath, err := os.Executable()
	if err != nil {
		log.Println(err)
		return
	}
	// 获取可执行文件所在的目录
	exeDir := filepath.Dir(exePath)

	// 输出当前项目路径
	fmt.Println("当前项目路径:", exeDir)
}