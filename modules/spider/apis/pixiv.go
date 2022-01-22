package apis

import "fmt"

const (
	Domain = "https://www.pixiv.net"

	IllustDetail = "/ajax/illust/%s"
	UserDetail   = "/ajax/user/%s"
	UserIllusts  = "/ajax/user/%s/profile/all"
	RankIllusts  = "/ranking.php"
	UgoiraDetail = "/ajax/illust/%s/ugoira_meta"
)

func IllustDetailG(id string) string {
	return Domain + fmt.Sprintf(IllustDetail, id) + "?lang=zh&full=1"
}

func UserDetailG(id string) string {
	return Domain + fmt.Sprintf(UserDetail, id) + "?lang=zh&full=1"
}

func UserIllustsG(id string) string {
	return Domain + fmt.Sprintf(UserIllusts, id) + "?lang=zh&full=1"
}

func UgoiraDetailG(id string) string {
	return Domain + fmt.Sprintf(UgoiraDetail, id) + "?lang=zh"
}

func RankIllustsG(mode string, page string, date string, content string) string {
	// mode = 'daily', 'weekly', 'monthly', 'rookie', 'male', 'female'
	return Domain + RankIllusts + fmt.Sprintf("?mode=%s&p=%s&date=%s&format=json&lang=zh&full=1&content=%s", mode, page, date, content)
}
