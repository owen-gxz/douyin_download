package douyin_download

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	headers = http.Header{
		"user-agent": []string{"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1"},
	}
	dytkUrl  = "https://www.iesdouyin.com/share/video/%s/"
	inforUrl = "https://www.iesdouyin.com/web/api/v2/aweme/iteminfo/"
	dytkRe   = regexp.MustCompile(`dytk: "(.+?)" }\);`)
	vidkRe   = regexp.MustCompile(`\/share\/video\/(\d*)`)
	midkRe   = regexp.MustCompile(`mid=(.+?)&`)
	uCodeRe  = regexp.MustCompile(`u_code=(.+?)&`)
)

//uri例子： https://v.douyin.com/3jj62D/
func GetDouyinInfo(uri string) *DouyinInfo {
	uri = getLocation(uri)
	u, err := url.Parse(uri)
	if err != nil {
		return nil
	}
	locaUrl := u.String()
	if locaUrl == "" {
		return nil
	}
	vid := vidkRe.FindStringSubmatch(locaUrl)[1]
	//mid := midkRe.FindStringSubmatch(locaUrl)[1]
	//u_code := uCodeRe.FindStringSubmatch(locaUrl)[1]
	//dytk := getDytk(vid, mid, u_code)
	infoUrl := fmt.Sprintf("%s?item_ids=%s", inforUrl, vid)
	return getInfo(infoUrl)
}

func getInfo(uri string) *DouyinInfo {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil
	}
	req.Header = headers
	result, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	data,err:=ioutil.ReadAll(result.Body)
	if err != nil {
		return nil
	}
	resp := &DouyinInfo{}
	err = json.Unmarshal(data,resp)
	if err != nil {
		return nil
	}
	return resp
}

// 封面动图
func (info DouyinInfo) GetDynamicCoverUrl() string {
	if len(info.ItemList) < 1 {
		return ""
	}
	if len(info.ItemList[0].Video.DynamicCover.URLList) < 1 {
		return ""
	}
	return info.ItemList[0].Video.DynamicCover.URLList[0]
}

// 封面静图
func (info DouyinInfo) GetOriginCoverUrl() string {
	if len(info.ItemList) < 1 {
		return ""
	}
	if len(info.ItemList[0].Video.OriginCover.URLList) < 1 {
		return ""
	}
	return info.ItemList[0].Video.OriginCover.URLList[0]
}

//原始视频
func (info DouyinInfo) GetOriginalVideoUrl() string {
	if len(info.ItemList) < 1 {
		return ""
	}
	if len(info.ItemList[0].Video.PlayAddr.URLList) < 1 {
		return ""
	}
	nurl:=strings.Replace(info.ItemList[0].Video.PlayAddr.URLList[0],"playwm","play",1)
	return  nurl
}

//原始视频
func (info DouyinInfo) GetTitle() string {
	if len(info.ItemList) < 1 {
		return ""
	}

	return info.ItemList[0].Desc
}

type DouyinInfo struct {
	Extra struct {
		Logid string `json:"logid"`
		Now   int64  `json:"now"`
	} `json:"extra"`
	ItemList []struct {
		AwemeID      string      `json:"aweme_id"`
		ChaList      interface{} `json:"cha_list"`
		CommentList  interface{} `json:"comment_list"`
		Desc         string      `json:"desc"`
		Duration     int64       `json:"duration"`
		Geofencing   interface{} `json:"geofencing"`
		ImageInfos   interface{} `json:"image_infos"`
		LabelTopText interface{} `json:"label_top_text"`
		LongVideo    interface{} `json:"long_video"`
		Position     interface{} `json:"position"`
		Promotions   interface{} `json:"promotions"`
		Statistics   struct {
			AwemeID      string `json:"aweme_id"`
			CommentCount int64  `json:"comment_count"`
			DiggCount    int64  `json:"digg_count"`
		} `json:"statistics"`
		TextExtra      interface{} `json:"text_extra"`
		UniqidPosition interface{} `json:"uniqid_position"`
		Video          struct {
			BitRate interface{} `json:"bit_rate"`
			Cover   struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"cover"`
			DownloadAddr struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"download_addr"`
			Duration     int64 `json:"duration"`
			DynamicCover struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"dynamic_cover"`
			HasWatermark bool  `json:"has_watermark"`
			Height       int64 `json:"height"`
			OriginCover  struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"origin_cover"`
			PlayAddr struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"play_addr"`
			PlayAddrLowbr struct {
				URI     string   `json:"uri"`
				URLList []string `json:"url_list"`
			} `json:"play_addr_lowbr"`
			Ratio string `json:"ratio"`
			Vid   string `json:"vid"`
			Width int64  `json:"width"`
		} `json:"video"`
		VideoLabels interface{} `json:"video_labels"`
		VideoText   interface{} `json:"video_text"`
	} `json:"item_list"`
	StatusCode int64 `json:"status_code"`
}

func getDytk(vid, mid, u_code string) string {
	dytkUrl, _ := url.Parse(fmt.Sprintf(dytkUrl, vid))
	dytkVal := url.Values{}
	dytkVal.Add("region", "CN")
	dytkVal.Set("mid", mid)
	dytkVal.Set("u_code", u_code)
	dytkVal.Set("titleType", "title")
	dytkVal.Set("utm_source", "copy_link")
	dytkVal.Set("utm_campaign", "client_share")
	dytkVal.Set("utm_medium", "android")
	dytkVal.Set("app", "aweme")
	dytkUrl.RawQuery = dytkVal.Encode()
	req, err := http.NewRequest(http.MethodGet, dytkUrl.String(), nil)
	if err != nil {
		return ""
	}
	req.Header = headers
	result, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	data, err := ioutil.ReadAll(result.Body)
	if err != nil {
		return ""
	}
	fmt.Println(string(data))
	return dytkRe.FindStringSubmatch(string(data))[1]
}

func getLocation(uri string) string {
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return uri
	}
	req.Header = headers
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	result, err := client.Do(req)
	if err != nil {
		return uri
	}
	return result.Header.Get("Location")
}
