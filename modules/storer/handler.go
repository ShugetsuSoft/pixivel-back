package storer

import (
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"github.com/ShugetsuSoft/pixivel-back/modules/storer/source"
	"github.com/mitchellh/mapstructure"
)

func (st *Storer) handleDatabase(dataq *source.DataQueue) error {
	for {
		data, tag, priority, err := dataq.GetData()
		if err != nil {
			return err
		}
		if data == nil {
			return models.ErrorChannelClosed
		}
		switch data.Type {
		case models.CrawlIllustDetail:
			var resdata models.Illust
			if err := mapstructure.Decode(data.Response, &resdata); err != nil {
				st.tracer.FailTask(data.Group, err.Error())
				return err
			}

			err = st.ops.InsertIllust(&resdata)
			if err != nil {
				st.tracer.FailTask(data.Group, err.Error())
				return err
			}

			st.tracer.FinishTask(data.Group)
			err = dataq.Ack(tag)
			if err != nil {
				return err
			}

		case models.CrawlUserDetail:
			var resdata models.User
			if err := mapstructure.Decode(data.Response, &resdata); err != nil {
				st.tracer.FailTask(data.Group, err.Error())
				return err
			}

			err = st.ops.InsertUser(&resdata)
			if err != nil {
				st.tracer.FailTask(data.Group, err.Error())
				return err
			}

			st.tracer.FinishTask(data.Group)
			err = dataq.Ack(tag)
			if err != nil {
				return err
			}
		case models.CrawlUserIllusts:
			var resdata models.UserIllustsResponse
			if err := mapstructure.Decode(data.Response, &resdata); err != nil {
				st.tracer.FailTask(data.Group, err.Error())
				return err
			}
			illustCount := uint(len(resdata.Illusts))

			user, err := st.ops.QueryUser(resdata.UserID, true)
			if err != nil {
				st.tracer.FailTask(data.Group, err.Error())
				return err
			}

			if user != nil {
				if illustCount != user.IllustsCount {
					err = st.ops.SetIllustsCount(resdata.UserID, illustCount)
					if err != nil {
						st.tracer.FailTask(data.Group, err.Error())
						return err
					}
				}
			} else {
				task := models.CrawlTask{
					Group:      "",
					Type:       models.CrawlUserDetail,
					Params:     map[string]string{"id": utils.Itoa(resdata.UserID)},
					RetryCount: st.retrys,
				}
				err = st.task.SendTask(task, priority)
				if err != nil {
					st.tracer.FailTask(data.Group, err.Error())
					return err
				}
			}

			if user != nil && illustCount != user.IllustsCount || user == nil {
				for id := range resdata.Illusts {
					exist, err := st.ops.IsIllustExist(resdata.Illusts[id])
					if err != nil {
						st.tracer.FailTask(data.Group, err.Error())
						return err
					}
					if !exist {
						task := models.CrawlTask{
							Group:      data.Group,
							Type:       models.CrawlIllustDetail,
							Params:     map[string]string{"id": utils.Itoa(resdata.Illusts[id])},
							RetryCount: st.retrys,
						}
						err = st.tracer.NewTask(data.Group)
						if err != nil {
							st.tracer.FailTask(data.Group, err.Error())
							return err
						}
						err = st.task.SendTask(task, priority)
						if err != nil {
							st.tracer.FailTask(data.Group, err.Error())
							return err
						}
					}
				}
			}

			err = st.ops.UpdateUserIllustsTime(resdata.UserID)
			if err != nil {
				st.tracer.FailTask(data.Group, err.Error())
				return err
			}
			st.tracer.FinishTask(data.Group)
			err = dataq.Ack(tag)
			if err != nil {
				return err
			}

		case models.CrawlRankIllusts:
			var resdata models.RankIllustsResponseMessage
			if err := mapstructure.Decode(data.Response, &resdata); err != nil {
				st.tracer.FailTask(data.Group, err.Error())
				return err
			}
			rankIllusts := make([]models.RankIllust, len(resdata.Illusts))
			for i := range resdata.Illusts {
				rankIllusts[i].Rank = resdata.Illusts[i].Pos
				rankIllusts[i].ID = resdata.Illusts[i].ID
				task := models.CrawlTask{
					Group:      data.Group,
					Type:       models.CrawlIllustDetail,
					Params:     map[string]string{"id": utils.Itoa(resdata.Illusts[i].ID)},
					RetryCount: st.retrys,
				}
				err = st.tracer.NewTask(data.Group)
				if err != nil {
					st.tracer.FailTask(data.Group, err.Error())
					return err
				}
				err = st.task.SendTask(task, priority)
				if err != nil {
					st.tracer.FailTask(data.Group, err.Error())
					return err
				}
			}
			err := st.ops.AddRankIllusts(resdata.Mode, resdata.Date, resdata.Content, rankIllusts)
			if err != nil {
				st.tracer.FailTask(data.Group, err.Error())
				return err
			}
			if resdata.Next {
				task := models.CrawlTask{
					Group:      data.Group,
					Type:       models.CrawlRankIllusts,
					Params:     map[string]string{"mode": resdata.Mode, "page": utils.Itoa(resdata.Page + 1), "date": resdata.Date, "content": resdata.Content},
					RetryCount: st.retrys,
				}
				err = st.tracer.NewTask(data.Group)
				if err != nil {
					st.tracer.FailTask(data.Group, err.Error())
					return err
				}
				err = st.task.SendTask(task, priority)
				if err != nil {
					st.tracer.FailTask(data.Group, err.Error())
					return err
				}
			}

			st.tracer.FinishTask(data.Group)
			err = dataq.Ack(tag)
			if err != nil {
				return err
			}
		case models.CrawlUgoiraDetail:
			var resdata models.Ugoira
			if err := mapstructure.Decode(data.Response, &resdata); err != nil {
				st.tracer.FailTask(data.Group, err.Error())
				return err
			}

			err = st.ops.InsertUgoira(&resdata)
			if err != nil {
				st.tracer.FailTask(data.Group, err.Error())
				return err
			}

			st.tracer.FinishTask(data.Group)
			err = dataq.Ack(tag)
			if err != nil {
				return err
			}

		case models.CrawlError:
			errinfo := data.Response.(string)
			st.tracer.FailTask(data.Group, errinfo)
			err = dataq.Ack(tag)
			if err != nil {
				return err
			}
		default:
			err = dataq.Ack(tag)
			if err != nil {
				return err
			}
		}
	}
}

func (st *Storer) handleElasticSearch(dataq *source.DataQueue) error {
	for {
		data, tag, _, err := dataq.GetData()
		if err != nil {
			return err
		}
		if data == nil {
			return models.ErrorChannelClosed
		}
		switch data.Type {
		case models.CrawlIllustDetail:
			var resdata models.Illust
			if err := mapstructure.Decode(data.Response, &resdata); err != nil {
				return err
			}

			err = st.ops.InsertIllustSearch(&resdata)
			if err != nil {
				return err
			}

			err = st.ops.InsertIllustTagNearDB(&resdata)
			if err != nil {
				return err
			}

			err = dataq.Ack(tag)
			if err != nil {
				return err
			}

		case models.CrawlUserDetail:
			var resdata models.User
			if err := mapstructure.Decode(data.Response, &resdata); err != nil {
				return err
			}

			err = st.ops.InsertUserSearch(&resdata)
			if err != nil {
				return err
			}

			err = dataq.Ack(tag)
			if err != nil {
				return err
			}
		default:
			err = dataq.Ack(tag)
			if err != nil {
				return err
			}
		}
	}
}
