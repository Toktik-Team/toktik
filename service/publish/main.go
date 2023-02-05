package main

import (
	"log"
	publish "toktik/kitex_gen/douyin/publish/publishservice"
	"toktik/service/publish/storage"
)

func main() {
	storage.Init()

	svr := publish.NewServer(new(PublishServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
