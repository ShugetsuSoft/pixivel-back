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
)
