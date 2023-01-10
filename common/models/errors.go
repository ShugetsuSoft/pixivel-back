package models

import "errors"

var (
	ErrorRetrivingFinishedTask = errors.New("error In Retryving Finished Task")
	ErrorIndexExist            = errors.New("error Index Already Existed")
	ErrorItemBanned            = errors.New("error Item Banned")
	ErrorNoResult              = errors.New("error No Result")
	ErrorChannelClosed         = errors.New("channel closed")
	ErrorTimeOut               = errors.New("time Out")
	ErrorArchiveMode           = errors.New("in Archive Mode")

	InternalErrorLoginNeeded = errors.New("login needed")

	PixivErrorIllustDeleted = errors.New("尚无权限浏览该作品")
	PixivErrorUserDeleted   = errors.New("抱歉，您当前所寻找的个用户已经离开了pixiv, 或者这ID不存在。")
)
