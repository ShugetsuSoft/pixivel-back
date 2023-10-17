package incrementor

import (
	"context"
	"time"

	"github.com/ShugetsuSoft/pixivel-back/common/database/operations"
	"github.com/ShugetsuSoft/pixivel-back/modules/responser/task"
)

func CrawlRank(taskgen *task.TaskGenerator, ope *operations.DatabaseOperations, date time.Time) error {
	ctx := context.Background()

	contents := map[string][]string{
		"all":    {"daily", "weekly", "monthly", "rookie", "original", "male", "female"},
		"illust": {"daily", "weekly", "monthly", "rookie"},
		"manga":  {"daily", "weekly", "monthly", "rookie"},
		"ugoira": {"daily", "weekly"},
	}
	today := date.Format("20060102")
	for content := range contents {
		for index := range contents[content] {
			exist, err := ope.IsRankExist(ctx, contents[content][index], today, content)
			if err != nil {
				return err
			}
			if exist {
				continue
			}
			err = taskgen.RankInitTask(ctx, contents[content][index], today, content)
			if err != nil {
				return err
			}
			time.Sleep(time.Second * 30)
		}
	}
	return nil
}
