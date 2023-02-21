package biz

import "toktik/kitex_gen/douyin/user"

var openAIImage = "https://bkimg.cdn.bcebos.com/pic/8b13632762d0f703918f0d436fac463d269758ee6faf?x-bce-process=image/watermark,image_d2F0ZXIvYmFpa2U4MA==,g_7,xp_5,yp_5"
var openAIIntro = "ChatGPT is a language model created by OpenAI with the ability to understand and generate human-like text responses to a wide range of topics and questions."
var zero uint32 = 0

var ChatGPTUser = &user.User{
	Id:              0,
	Name:            "ChatGPT",
	FollowCount:     0,
	FollowerCount:   100000000,
	IsFollow:        true,
	Avatar:          &openAIImage,
	BackgroundImage: &openAIImage,
	Signature:       &openAIIntro,
	TotalFavorited:  &zero,
	WorkCount:       &zero,
	FavoriteCount:   &zero,
}
