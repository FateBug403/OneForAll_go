package oneforall

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/FateBug403/util"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

type OneForAll struct {
	Options Options
}

func NewOneForAll(options Options) (*OneForAll,error) {
	if options.TmpPath ==""{
		// 获取可执行文件的路径
		exePath, err := os.Executable()
		if err != nil {
			log.Println(err)
			return nil,err
		}
		// 获取可执行文件所在的目录
		exeDir := filepath.Dir(exePath)
		options.TmpPath = exeDir+"\\tmp"
	}

	// 检查路径是否存在,如果不存在则创建
	err := util.CreatePath(options.TmpPath)
	if err != nil {
		return nil,err
	}

	oneforall := &OneForAll{Options: options}
	return oneforall,nil
}

func (receiver OneForAll) GetSubDomains(domains []string) ([]string,error) {
	if len(domains)==0{ // 如果为空直接退出
		return nil,nil
	}
	currentTime := time.Now()
	timestamp := currentTime.Format("20060102_150405")
	InputPath := fmt.Sprintf(receiver.Options.TmpPath+"\\%s.txt",timestamp)
	OutputPath := fmt.Sprintf(receiver.Options.TmpPath+"\\%s.json",timestamp)
	err :=util.WriteLineFile(InputPath,domains)
	if err != nil {
		return nil,err
	}

	var cmd *exec.Cmd
	cmd = exec.Command("python",receiver.Options.ExePath+"\\oneforall.py","--targets",InputPath,"--brute","False","--req","False","--path",OutputPath,"--fmt","json","run")

	// 获取命令的标准输出管道
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("无法获取标准输出管道:", err)
		return nil,err
	}
	// 启动命令
	//log.Println("--------------------------------------------开始运行指纹探测程序----------------------------------------")
	if err := cmd.Start(); err != nil {
		fmt.Println("命令启动失败:", err)
		return nil,err
	}
	// 实时读取输出
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		//output(line)
	}
	// 等待命令执行完成
	if err := cmd.Wait(); err != nil {
		fmt.Println("命令执行出错:", err)
		return nil,err
	}

	// 判断是否为json文件，如何是提取json中的内容，如果不是按行读取内容
	var subdomains []string
	data, err := ioutil.ReadFile(OutputPath)
	if err != nil {
		fmt.Println("无法读取文件:", err)
		return nil,err
	}
	var subdomainsTmp []string
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		subdomainsTmp = util.ReadFile(OutputPath)
	}else {
		subdomainsTmp,err =JsonResolveSubdomain(OutputPath)
		if err != nil {
			return nil,err
		}
	}
	subdomains = append(subdomains,subdomainsTmp...)
	log.Println("成功获取域名:"+ strconv.Itoa(len(subdomains)))
	return subdomains,nil
}

func JsonResolveSubdomain(subdomainPath string) ([]string,error) {
	//从文件中获取批量获取的子域名
	//subdomain := util.ReadFile(domainPath)
	fileContent, err := ioutil.ReadFile(subdomainPath)
	if err != nil {
		fmt.Println("Error reading the file:", err)
		return nil,err
	}

	// 创建一个空切片来存储JSON数据
	var jsonArray []map[string]interface{}

	// 解析JSON数据到切片
	err = json.Unmarshal(fileContent, &jsonArray)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return nil,err
	}
	var subdomains []string
	for _,item := range jsonArray{
		subdomain,ok := item["subdomain"].(string)
		if ok{
			subdomains = append(subdomains,subdomain)
		}
	}
	return subdomains,nil
}