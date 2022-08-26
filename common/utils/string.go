package utils

import (
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"strings"
	"unicode"
)

func RemoveSpecialChars(in string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, in)
}

func RemoveSpecialCharsTags(tags []models.IllustTag) []models.IllustTag {
	result := make([]models.IllustTag, len(tags))
	for i := 0; i < len(tags); i++ {
		result[i].Name = RemoveSpecialChars(tags[i].Name)
		result[i].Translation = RemoveSpecialChars(tags[i].Translation)
	}
	return result
}
