package models

import "errors"

var (
	ErrorRetrivingFinishedTask = errors.New("Error In Retryving Finished Task.")
	ErrorFailToCompleteTask    = errors.New("Error Fail To Complete Task.")
	ErrorIndexExist            = errors.New("Error Index Already Existed")
	ErrorItemBanned            = errors.New("Error Item Banned")
	ErrorNoResult              = errors.New("Error No Result")
	ErrorChannelClosed         = errors.New("channel closed")
	ErrorTimeOut               = errors.New("Time Out")
)
