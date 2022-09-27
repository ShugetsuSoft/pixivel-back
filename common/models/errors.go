package models

import "errors"

var (
	ErrorRetrivingFinishedTask     = errors.New("Error In Retryving Finished Task.")
	ErrorIndexExist                = errors.New("Error Index Already Existed")
	ErrorItemBanned                = errors.New("Error Item Banned")
	ErrorNoResult                  = errors.New("Error No Result")
	ErrorChannelClosed             = errors.New("channel closed")
	ErrorTimeOut                   = errors.New("Time Out")
	ErrorRecommendationNotPrepared = errors.New("Error Recommendation Not Prepared")
)
