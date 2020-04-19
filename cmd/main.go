package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/owen-gxz/douyin_download"
)

func main() {
	fmt.Println("ctrl+C退出")
	for {
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Printf("请输入你的分享网址（https://v.douyin.com/3jj62D/）:")
		input, _, err := inputReader.ReadLine()
		if err != nil {
			fmt.Println("There were errors reading, exiting program.")
			return
		}
		info := douyin_download.GetDouyinInfo(string(input))
		fmt.Println("视频标题为:", info.GetTitle())
		fmt.Println("视频原始地址为:", info.GetOriginalVideoUrl())
		fmt.Println("视频动图地址为:", info.GetDynamicCoverUrl())
		fmt.Println("视频静图地址为:", info.GetOriginCoverUrl())
	}
}
