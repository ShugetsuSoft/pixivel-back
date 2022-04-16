package incrementor

import (
	"context"
	"time"

	"github.com/ShugetsuSoft/pixivel-back/common/database/operations"
	"github.com/ShugetsuSoft/pixivel-back/modules/responser/task"
)

func CrawlRank(taskgen *task.TaskGenerator, ope *operations.DatabaseOperations) error {
	ctx := context.Background()

	contents := map[string][]string{
		"all":    {"daily", "weekly", "monthly", "rookie", "original", "male", "female"},
		"illust": {"daily", "weekly", "monthly", "rookie"},
		"manga":  {"daily", "weekly", "monthly", "rookie"},
		"ugoira": {"daily", "weekly"},
	}
	today := time.Now().AddDate(0, 0, -2).Format("20060102")
	for content := range contents {
		for index := range contents[content] {
			succ, err := ope.InsertRank(ctx, contents[content][index], today, content)
			if err != nil {
				return err
			}
			if !succ {
				continue
			}
			err = taskgen.RankInitTask(ctx, contents[content][index], today, content)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
