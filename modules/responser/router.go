package responser

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ShugetsuSoft/pixivel-back/common/convert"
	"github.com/ShugetsuSoft/pixivel-back/common/database/drivers"
	"github.com/ShugetsuSoft/pixivel-back/common/utils/telemetry"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/ShugetsuSoft/pixivel-back/common/database/operations"
	"github.com/ShugetsuSoft/pixivel-back/common/database/tasktracer"
	"github.com/ShugetsuSoft/pixivel-back/common/models"
	"github.com/ShugetsuSoft/pixivel-back/common/utils"
	"github.com/ShugetsuSoft/pixivel-back/modules/responser/reader"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Router struct {
	reader *reader.Reader
	cache  *drivers.Cache
	debug  bool
}

func success(data interface{}) *Response {
	return &Response{Error: false, Message: "", Data: data}
}

func NewRouter(dbops *operations.DatabaseOperations, mq models.MessageQueue, taskchaname string, retrys uint, tracer *tasktracer.TaskTracer, redis *drivers.RedisPool, debugflag bool) *Router {
	return &Router{
		reader: reader.NewReader(dbops, mq, taskchaname, retrys, tracer),
		cache:  drivers.NewCache(redis),
		debug:  debugflag,
	}
}

func fail(err string) *Response {
	return &Response{Error: true, Message: err, Data: nil}
}

func (r *Router) Fail(err error) *Response {
	erra := ""
	if r.debug {
		erra = fmt.Sprintf("%s", err)
	} else {
		erra = "服务器酱出错啦！"
	}
	logerr := fmt.Sprintf("%s", err)
	if strings.Contains(logerr, "context canceled") {
		return &Response{Error: true, Message: erra, Data: nil}
	}
	telemetry.Log(telemetry.Label{"pos": "ResponseError"}, logerr)
	return &Response{Error: true, Message: erra, Data: nil}
}

func (r *Router) GetIllustHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "illust"}).Inc()
	id := utils.Atoi(c.Param("id"))
	if id == 0 {
		return
	}

	forcefetch := false
	if i, e := strconv.ParseBool(c.Query("forcefetch")); i && e == nil {
		forcefetch = true
	}

	if !forcefetch {
		cached, err := r.cache.Get("illust", utils.Itoa(id))
		if err != nil {
			telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
		}
		if cached != nil {
			c.JSON(200, success(cached))
			return
		}
	}

	illust, err := r.reader.IllustResponse(ctx, id, forcefetch)

	if err != nil {
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "illust"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(illust))
	err = r.cache.Set("illust", illust, 60*60*12, utils.Itoa(id))
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
	}
}

func (r *Router) GetUgoiraHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "ugoira"}).Inc()
	id := utils.Atoi(c.Param("id"))
	if id == 0 {
		return
	}

	forcefetch := false
	if i, e := strconv.ParseBool(c.Query("forcefetch")); i && e == nil {
		forcefetch = true
	}

	if !forcefetch {
		cached, err := r.cache.Get("ugoira", utils.Itoa(id))
		if err != nil {
			telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
		}
		if cached != nil {
			c.JSON(200, success(cached))
			return
		}
	}

	ugoira, err := r.reader.UgoiraResponse(ctx, id, forcefetch)

	if err != nil {
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "ugoira"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(ugoira))
	err = r.cache.Set("ugoira", ugoira, 60*60*12, utils.Itoa(id))
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
	}
}

func (r *Router) GetUserDetailHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "user"}).Inc()
	id := utils.Atoi(c.Param("id"))
	if id == 0 {
		return
	}

	cached, err := r.cache.Get("user", utils.Itoa(id))
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
	}
	if cached != nil {
		c.JSON(200, success(cached))
		return
	}

	user, err := r.reader.UserDetailResponse(ctx, id)

	if err != nil {
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "user"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(user))
	err = r.cache.Set("user", user, 60*60*12, utils.Itoa(id))
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
	}
}

func (r *Router) GetUserIllustsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "user-illust"}).Inc()
	id := utils.Atoi(c.Param("id"))
	if id == 0 {
		return
	}

	page := utils.Atoi(c.Query("page"))
	if page < 0 {
		page = 0
	}

	limit := utils.Atoi(c.Query("limit"))
	if limit > 40 || limit < 1 {
		limit = 30
	}

	cached, err := r.cache.Get("user-illust", utils.Itoa(id), utils.Itoa(page), utils.Itoa(limit))
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
	}
	if cached != nil {
		c.JSON(200, success(cached))
		return
	}

	illusts, err := r.reader.UserIllustsResponse(ctx, id, int64(page), int64(limit))

	if err != nil {
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "user-illust"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(illusts))
	err = r.cache.Set("user-illust", illusts, 60*60*12, utils.Itoa(id), utils.Itoa(page), utils.Itoa(limit))
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
	}
}

func (r *Router) SearchIllustHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "search-illust"}).Inc()
	keyword := c.Param("keyword")
	if keyword == "" {
		return
	}

	page := utils.Atoi(c.Query("page"))
	if page < 0 {
		page = 0
	}

	limit := utils.Atoi(c.Query("limit"))
	if limit > 40 || limit < 1 {
		limit = 30
	}

	sortpop := false
	if i, e := strconv.ParseBool(c.Query("sortpop")); i && e == nil {
		sortpop = true
	}

	sortdate := false
	if i, e := strconv.ParseBool(c.Query("sortdate")); i && e == nil {
		sortdate = true
	}

	illusts, err := r.reader.SearchIllustsResponse(ctx, keyword, int(page), int(limit), sortpop, sortdate)

	if err != nil {
		if err == models.ErrorNoResult {
			c.JSON(200, r.Fail(err))
			return
		}
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "search-illust"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(illusts))
}

func (r *Router) SearchIllustSuggestHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "search-illust-suggest"}).Inc()
	keyword := c.Param("keyword")
	if keyword == "" {
		return
	}

	suggests, err := r.reader.SearchIllustsSuggestResponse(ctx, keyword)

	if err != nil {
		if err == models.ErrorNoResult {
			c.JSON(200, r.Fail(err))
			return
		}
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "search-illust-suggest"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(suggests))
}

func (r *Router) SearchUserHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "search-user"}).Inc()
	keyword := c.Param("keyword")
	if keyword == "" {
		return
	}

	page := utils.Atoi(c.Query("page"))
	if page < 0 {
		page = 0
	}

	limit := utils.Atoi(c.Query("limit"))
	if limit > 40 || limit < 1 {
		limit = 30
	}

	users, err := r.reader.SearchUsersResponse(ctx, keyword, int(page), int(limit))

	if err != nil {
		if err == models.ErrorNoResult {
			c.JSON(200, r.Fail(err))
			return
		}
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "search-user"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(users))
}

func (r *Router) SearchUserSuggestHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "search-user-suggest"}).Inc()
	keyword := c.Param("keyword")
	if keyword == "" {
		return
	}

	suggests, err := r.reader.SearchUsersSuggestResponse(ctx, keyword)

	if err != nil {
		if err == models.ErrorNoResult {
			c.JSON(200, r.Fail(err))
			return
		}
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "search-user-suggest"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(suggests))
}

func (r *Router) SearchTagSuggestHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "search-tag-suggest"}).Inc()
	keyword := c.Param("keyword")
	if keyword == "" {
		return
	}

	suggests, err := r.reader.SearchTagsSuggestResponse(ctx, keyword)

	if err != nil {
		if err == models.ErrorNoResult {
			c.JSON(200, r.Fail(err))
			return
		}
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "search-tag-suggest"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(suggests))
}

func (r *Router) SearchIllustByTagHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "search-illust-by-tag"}).Inc()
	keywords := c.Param("keyword")
	if keywords == "" {
		return
	}
	twotags := strings.Split(keywords, "|")
	var musttags []string
	var shouldtags []string
	if len(twotags) > 0 {
		if twotags[0] != "" {
			musttags = strings.Split(twotags[0], ",")
		}
	}
	if len(twotags) > 1 {
		if twotags[1] != "" {
			shouldtags = strings.Split(twotags[1], ",")
		}
	}

	if len(musttags) == 0 && len(shouldtags) == 0 {
		c.JSON(400, fail(fmt.Sprintf("%s", models.ErrorNoResult)))
		return
	}

	page := utils.Atoi(c.Query("page"))
	if page < 0 {
		page = 0
	}

	limit := utils.Atoi(c.Query("limit"))
	if limit > 40 || limit < 1 {
		limit = 30
	}

	sortpop := false
	if i, e := strconv.ParseBool(c.Query("sortpop")); i && e == nil {
		sortpop = true
	}

	sortdate := false
	if i, e := strconv.ParseBool(c.Query("sortdate")); i && e == nil {
		sortdate = true
	}

	perfectmatch := true
	if i, e := strconv.ParseBool(c.Query("perfectmatch")); !i && e == nil {
		perfectmatch = false
	}

	illusts, err := r.reader.SearchIllustsByTagsResponse(ctx, musttags, shouldtags, perfectmatch, int(page), int(limit), sortpop, sortdate)

	if err != nil {
		if err == models.ErrorNoResult {
			c.JSON(200, r.Fail(err))
			return
		}
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "search-illust-by-tag"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(illusts))
}

func (r *Router) GetIllustsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "illusts"}).Inc()
	keywords := c.Param("ids")
	if keywords == "" {
		return
	}
	illustsstr := strings.Split(keywords, ",")
	if len(illustsstr) > 100 {
		c.JSON(400, fail("Query is too Large."))
		return
	}
	illustsids := make([]uint64, len(illustsstr))
	for i, id := range illustsstr {
		ida := utils.Atoi(id)
		if ida == 0 {
			return
		}
		illustsids[i] = ida
	}

	illusts, err := r.reader.IllustsResponse(ctx, illustsids)

	if err != nil {
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "illusts"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(illusts))
}

func (r *Router) RecommendIllustsByIllustIdHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "illust-recommend"}).Inc()
	const maxpage = 5
	id := utils.Atoi(c.Param("id"))
	if id == 0 {
		return
	}

	page := utils.Atoi(c.Query("page"))
	if page < 0 {
		page = 0
	}

	if page >= maxpage {
		c.JSON(400, fail("No More Info"))
		return
	}

	cached, err := r.cache.Get("illust-recommend", utils.Itoa(id), utils.Itoa(page))
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
	}
	if cached != nil {
		c.JSON(200, success(cached))
		return
	}

	const limit = 30

	illusts, err := r.reader.RecommendIllustsByIllustId(ctx, id, limit*maxpage)

	if err != nil {
		if err == models.ErrorNoResult {
			c.JSON(200, r.Fail(err))
			return
		}
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "illust-recommend"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	for i := 0; i < maxpage; i++ {
		pagenow := illusts[limit*i : limit*(i+1)]
		if len(illusts) < limit*(i+2) {
			err = r.cache.Set("illust-recommend", convert.Illusts2IllustsResponse(pagenow, false), 60*60*2, utils.Itoa(id), utils.Itoa(i))
			break
		}
		err = r.cache.Set("illust-recommend", convert.Illusts2IllustsResponse(pagenow, i < maxpage-1), 60*60*2, utils.Itoa(id), utils.Itoa(i))
		if err != nil {
			telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
		}
	}

	c.JSON(200, success(convert.Illusts2IllustsResponse(illusts[limit*page:limit*(page+1)], page < maxpage-1)))
}

func (r *Router) GetRankHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "rank"}).Inc()
	mode := c.Query("mode")
	modes := map[string]bool{"daily": true, "weekly": true, "monthly": true, "rookie": true, "original": true, "male": true, "female": true}
	if _, ok := modes[mode]; !ok {
		c.JSON(400, fail("Illegal mode"))
		return
	}

	content := c.Query("content")
	contents := map[string]bool{"all": true, "illust": true, "manga": true, "ugoira": true}
	if _, ok := contents[content]; !ok {
		c.JSON(400, fail("Illegal content"))
		return
	}

	date := c.Query("date")
	dateI, err := time.Parse("20060102", date)
	if err != nil {
		c.JSON(400, fail(fmt.Sprintf("Time parse Error. %s", err)))
		return
	}

	if dateI.After(time.Now().AddDate(0, 0, -1)) {
		c.JSON(400, fail("Rank Info DNE. 你是不是傻啊"))
		return
	}

	page := utils.Atoi(c.Query("page"))
	if page < 0 {
		page = 0
	}

	if page > 9 {
		c.JSON(400, fail("No More Info"))
		return
	}

	limit := 50

	illusts, err := r.reader.RankIllustsResponse(ctx, mode, date, int(page), content, limit)

	if err != nil {
		if err == models.ErrorNoResult {
			c.JSON(200, r.Fail(err))
			return
		}
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "rank"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	if page == 0 && illusts.HasNext == false {
		c.JSON(400, fail("No More Info"))
		return
	}

	c.JSON(200, success(illusts))
}

func (r *Router) GetSampleIllustsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "sample-illusts"}).Inc()
	page := utils.Atoi(c.Query("p"))
	if page > 20 || page < 0 {
		page = 0
	}

	cached, err := r.cache.Get("illust-sample", utils.Itoa(page))
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
	}
	if cached != nil {
		c.JSON(200, success(cached))
		return
	}

	illusts, err := r.reader.SampleIllustsResponse(ctx, 15000, 50)

	if err != nil {
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "sample-illusts"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	err = r.cache.Set("illust-sample", illusts, 60*60*6, utils.Itoa(page))
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
	}

	c.JSON(200, success(illusts))
}

func (r *Router) GetSampleIllustHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "sample-illust"}).Inc()

	rawtags := c.Query("tags")
	var tags []string
	if rawtags != "" {
		tags = strings.Split(rawtags, ",")
	}

	quality := utils.Atoi(c.Query("quality"))
	if quality < 0 {
		quality = 0
	}

	illust, err := r.reader.SampleIllustResponse(ctx, tags, int(quality))

	if err != nil {
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "sample-illust"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	c.JSON(200, success(illust))
}

func (r *Router) GetSampleUsersHandler(c *gin.Context) {
	ctx := c.Request.Context()

	telemetry.RequestsCount.With(prometheus.Labels{"handler": "sample-users"}).Inc()
	page := utils.Atoi(c.Query("p"))
	if page > 20 || page < 0 {
		page = 0
	}

	cached, err := r.cache.Get("user-sample", utils.Itoa(page))
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
	}
	if cached != nil {
		c.JSON(200, success(cached))
		return
	}

	users, err := r.reader.SampleUsersResponse(ctx, 50)

	if err != nil {
		telemetry.RequestsErrorCount.With(prometheus.Labels{"handler": "sample-users"}).Inc()
		c.JSON(500, r.Fail(err))
		return
	}

	err = r.cache.Set("user-sample", users, 60*60*6, utils.Itoa(page))
	if err != nil {
		telemetry.Log(telemetry.Label{"pos": "cache"}, err.Error())
	}

	c.JSON(200, success(users))
}

func (r *Router) mount(rout *gin.Engine) {
	r1 := rout.Group("/v2")
	// pixiv

	pixiv := r1.Group("/pixiv")
	{
		pixiv.GET("/illust/:id", r.GetIllustHandler)
		pixiv.GET("/user/:id", r.GetUserDetailHandler)
		pixiv.GET("/user/:id/illusts", r.GetUserIllustsHandler)
		pixiv.GET("/illust/search/:keyword", r.SearchIllustHandler)
		pixiv.GET("/illust/search/:keyword/suggest", r.SearchIllustSuggestHandler)
		pixiv.GET("/user/search/:keyword", r.SearchUserHandler)
		pixiv.GET("/user/search/:keyword/suggest", r.SearchUserSuggestHandler)
		pixiv.GET("/tag/search/:keyword", r.SearchIllustByTagHandler)
		pixiv.GET("/tag/search/:keyword/suggest", r.SearchTagSuggestHandler)
		pixiv.GET("/illust/:id/recommend", r.RecommendIllustsByIllustIdHandler)
		//pixiv.GET("/illusts/:ids", r.GetIllustsHandler)
		pixiv.GET("/rank/", r.GetRankHandler)
		pixiv.GET("/illusts/sample", r.GetSampleIllustsHandler)
		pixiv.GET("/user/sample", r.GetSampleUsersHandler)
		pixiv.GET("/ugoira/:id", r.GetUgoiraHandler)
		pixiv.GET("/illust/sample/single", r.GetSampleIllustHandler)
	}
}
