package main

import (
	"log"
	wechat "toktik/kitex_gen/douyin/wechat/wechatservice"
)

func main() {
	svr := wechat.NewServer(new(WechatServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
