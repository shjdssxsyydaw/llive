package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

type Match struct {
	GameID     string `json:"GameId"`
	GameName   string `json:"GameName"`
	BMatchName string `json:"bMatchName"`
	MatchDate  string `json:"MatchDate"`
	Video3     string `json:"Video3"`
}

func findAllSubsMatches(re string, str string) (matches [][]string) {
	regex := regexp.MustCompile(re)
	strings := regex.FindAllString(str, -1)
	for _, substr := range strings {
		matches = append(matches, regex.FindStringSubmatch(substr))
	}
	return
}

func getBody(url string) (body []byte, err error) {
	http.DefaultClient.Timeout = time.Second
	var res *http.Response
	for res == nil {
		res, err = http.Get(url)
		if err != nil {
			continue
		}
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}

func main() {
	body, err := getBody("http://lpl.qq.com/web201612/data/LOL_MATCH2_LIVE_BMATCH_LIST.js")
	if err != nil {
		panic(err)
	}
	results := findAllSubsMatches(`({"GameId.+?})`, string(body))
	if len(results) == 0 {
		println("找不到正在直播的赛事")
		os.Exit(0)
	}

	for _, result := range results {
		match := new(Match)
		err := json.Unmarshal([]byte(result[1]), match)
		if err != nil {
			panic(err)
		}
		name := match.GameName + " " + match.BMatchName + " " + match.MatchDate
		if match.Video3 == "" {
			println(name + " 找不到直播流")
			continue
		}
		body, err := getBody("http://livematches.qt.qq.com/get_video_url_v2?module=" + match.Video3 + "&videotype=flv&use_https=0")
		urls := findAllSubsMatches(`"urllist":"(.+?)"`, string(body))
		if len(urls) == 0 {
			println(name + " 找不到直播流")
			continue
		}
		println(name + " " + strings.Split(urls[0][1], ";")[0])
	}
}
