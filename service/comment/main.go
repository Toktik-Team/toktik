package main

import (
	"log"
	comment "toktik/kitex_gen/douyin/comment/commentservice"
)

func main() {
	svr := comment.NewServer(new(CommentServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
