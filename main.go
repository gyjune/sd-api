package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// 频道映射表：完全对应原业务逻辑
var channelMap = map[string][2]int{
	// 济南
	"jncqxw": {171, 2}, "jncqsh": {171, 20}, "jnjrtv": {303, 1}, "jnjyzh": {85, 1},
	"jnjyys": {85, 2}, "jnlcxw": {261, 1}, "jnpyzh": {257, 1}, "jnpyxc": {257, 3},
	"jnshzh": {97, 1}, "jnshys": {97, 2}, "jnzqzh": {195, 1}, "jnzqgg": {195, 2},
	// 东营
	"dyxwzh": {537, 1}, "dygg": {537, 3}, "dygg2": {29, 90}, "dykj": {537, 7},
	"dydyqxw": {163, 5}, "dydyqkj": {163, 7}, "dygrzh": {237, 1}, "dygrkj": {237, 5},
	"dyklxw": {269, 3}, "dyljzh": {153, 1}, "dyljwh": {153, 3},
	// 青岛
	"qdcyzh": {403, 5}, "qdhdzh": {227, 1}, "qdhdsh": {227, 3}, "qdjmzh": {221, 2},
	"qdjmsh": {221, 3}, "qdjzzh": {305, 1}, "qdjzsh": {305, 3}, "qdlc": {173, 1},
	"qdls": {295, 1}, "qdlxzh": {253, 1}, "qdlxsh": {253, 3}, "qdpdxw": {45, 4},
	"qdpdsh": {45, 5},
	// 潍坊
	"wfxwzh": {635, 1}, "wfsh": {635, 5}, "wfyswy": {635, 7}, "wfkjwl": {635, 9},
	"wfgxq": {421, 14}, "wfaqzh": {137, 3}, "wfaqms": {137, 4}, "wfbhxw": {199, 1},
	"wfclzh": {1, 3}, "wfcyzh": {47, 1}, "wfcyjj": {47, 2}, "wffzxw": {285, 1},
	"wfgmzh": {71, 24}, "wfgmdj": {71, 38}, "wfhtxw": {133, 1}, "wfkwtv": {127, 17},
	"wflqxw": {205, 39}, "wfqzzh": {125, 2}, "wfqzwh": {125, 3}, "wfwc": {15, 3},
	"wfzcxw": {115, 23}, "wfzcsh": {115, 25},
	// 烟台
	"ytcd": {175, 1}, "ytfszh": {189, 4}, "ytfssh": {189, 5}, "ythyzh": {255, 1},
	"ytlkzh": {57, 1}, "ytlksh": {57, 2}, "ytlszh": {245, 4}, "ytlsys": {245, 6},
	"ytlyzh": {241, 4}, "ytlyms": {241, 7}, "ytlzzh": {239, 1}, "ytmpzh": {281, 1},
	"ytplzh": {109, 1}, "ytplzy": {109, 2}, "ytqxzh": {165, 12}, "ytqxpg": {165, 14},
	"ytzyzh": {55, 2}, "ytzyzy": {55, 4},
	// 淄博
	"zbbsxw": {17, 8}, "zbbstw": {17, 9}, "zbgqzh": {61, 1}, "zbgqys": {61, 2},
	"zbht1": {23, 15}, "zbht2": {23, 16}, "zblzxw": {151, 6}, "zblzsh": {151, 7},
	"zbyyzh": {203, 6}, "zbyysh": {203, 7}, "zbzcxw": {75, 1}, "zbzcsh": {75, 2},
	"zbzd1": {101, 1}, "zbzd2": {101, 6}, "zbzctv1": {259, 1}, "zbzctv2": {259, 3},
	// 枣庄
	"zzstzh": {243, 1}, "zzszzh": {233, 1}, "zztezxw": {185, 2}, "zztzzh": {103, 2},
	"zztzms": {103, 3}, "zzxcxw": {37, 8}, "zzyczh": {209, 1},
	// 滨州
	"bzbctv": {249, 35}, "bzbxzh": {207, 3}, "bzbxsh": {207, 4}, "bzhmzh": {211, 2},
	"bzhmys": {211, 3}, "bzwdzh": {169, 1}, "bzwdzy": {169, 21}, "bzyxxw": {217, 1},
	"bzzhzh": {277, 1}, "bzzhzy": {277, 9}, "bzzpzh": {11, 15}, "bzzpms": {11, 16},
	// 德州
	"dzxwzh": {179, 1}, "dzjjsh": {179, 2}, "dztw": {179, 9}, "dzlczh": {215, 6},
	"dzllxw": {267, 1}, "dzllcs": {267, 5}, "dzly1": {49, 3}, "dzly2": {49, 4},
	"dznjzh": {193, 1}, "dzpyzh": {19, 2}, "dzqhzh": {251, 8}, "dzqyzh": {5, 9},
	"dzqysh": {5, 7}, "dzwczh": {33, 4}, "dzwczy": {33, 6}, "dzxjzh": {223, 1},
	"dzxjgg": {223, 2}, "dzyczh": {235, 1}, "dzyczy": {235, 3},
	// 菏泽
	"hzcwzh": {131, 1}, "hzcwzy": {131, 2}, "hzcxzh": {87, 2}, "hzdmxw": {111, 2},
	"hzdt1": {27, 7}, "hzdt2": {27, 8}, "hzjczh": {141, 186}, "hzjyxw": {139, 1},
	"hzmdxw": {219, 6}, "hzmdzy": {219, 17}, "hzsxzh": {155, 2}, "hzycxw": {135, 3},
	"hzyczy": {135, 2},
	// 济宁
	"jijiazh": {273, 1}, "jijiash": {273, 3}, "jijxzh": {129, 2}, "jijxsh": {129, 4},
	"jilszh": {89, 1}, "jiqfxw": {13, 1}, "jircxw": {73, 8}, "jircys": {73, 9},
	"jissxw": {117, 5}, "jisswh": {117, 6}, "jiws1": {53, 4}, "jiws2": {53, 5},
	"jiwszh": {301, 1}, "jiytxw": {63, 5}, "jiytsh": {63, 15}, "jiyzxw": {231, 1},
	"jiyzsh": {231, 3}, "jizczh": {181, 1}, "jizcwh": {181, 4},
	// 聊城
	"lccpzh": {31, 6}, "lccpsh": {31, 8}, "lcdczh": {265, 1}, "lcdezh": {95, 22},
	"lcdezy": {95, 29}, "lcgtzh": {43, 1}, "lcgtzy": {43, 5}, "lcgxzh": {79, 1},
	"lclqzh": {65, 2}, "lclqjj": {65, 5}, "lcsxzh": {183, 1}, "lcsxsh": {183, 5},
	"lcygzh": {81, 1}, "lcygys": {81, 10},
	// 临沂
	"lyfxzh": {41, 119}, "lyfxsh": {41, 117}, "lyhdys": {191, 1}, "lyhdzh": {191, 2},
	"lyjnzh": {105, 4}, "lyjnys": {105, 5}, "lyllzh": {113, 131}, "lyllgg": {113, 133},
	"lylszh": {201, 1}, "lyls1": {167, 3}, "lyls2": {167, 4}, "lylzzh": {147, 1},
	"lylzys": {147, 17}, "lymy1": {161, 13}, "lymy2": {161, 15}, "lypyzh": {345, 4},
	"lypysh": {345, 14}, "lytc1": {83, 1}, "lytc2": {83, 2}, "lyynzh": {177, 6},
	"lyynys": {177, 7}, "lyys1": {145, 1}, "lyys2": {145, 2},
	// 日照
	"rzjx1": {159, 23}, "rzjx2": {159, 27}, "rzls": {289, 1}, "rzwlzh": {299, 10},
	"rzwlwh": {299, 12},
	// 泰安
	"tadpzh": {187, 9}, "tadpms": {187, 11}, "tady": {293, 1}, "tafczh": {51, 18},
	"tafcsh": {51, 6}, "tany1": {123, 1}, "tany2": {123, 7}, "tats": {263, 1},
	"taxtzh": {59, 2}, "taxtxc": {59, 3},
	// 威海
	"whxwzh": {157, 1}, "whdssh": {157, 3}, "whhy": {157, 12}, "whhczh": {213, 5},
	"whrczh": {77, 10}, "whrcsh": {77, 11}, "whrszh": {143, 8}, "whrssh": {143, 9},
	"whwd1": {91, 7}, "whwd2": {91, 8},
}

// 中文频道名称映射
var cnName = map[string]string{
	"jncqxw": "长清新闻", "jncqsh": "长清生活", "jnjrtv": "济铁电视台", "jnjyzh": "济阳综合",
	"jnjyys": "济阳影视", "jnlcxw": "历城新闻综合", "jnpyzh": "平阴综合", "jnpyxc": "平阴乡村振兴",
	"jnshzh": "商河综合", "jnshys": "商河影视", "jnzqzh": "章丘综合", "jnzqgg": "章丘公共",
	"dyxwzh": "东营新闻综合", "dygg": "东营公共", "dygg2": "东营公共2", "dykj": "东营科教",
	"dydyqxw": "东营区新闻综合", "dydyqkj": "东营区科教影视", "dygrzh": "广饶综合", "dygrkj": "广饶科教文艺",
	"dyklxw": "垦利新闻综合", "dyljzh": "利津综合", "dyljwh": "利津文化生活", "qdcyzh": "城阳综合",
	"qdhdzh": "黄岛综合", "qdhdsh": "黄岛生活", "qdjmzh": "即墨综合", "qdjmsh": "即墨生活服务",
	"qdjzzh": "胶州综合", "qdjzsh": "胶州生活", "qdlc": "李沧TV", "qdls": "崂山TV",
	"qdlxzh": "莱西综合", "qdlxsh": "莱西生活", "qdpdxw": "平度新闻综合", "qdpdsh": "平度生活服务",
	"wfxwzh": "潍坊新闻综合", "wfsh": "潍坊生活", "wfyswy": "潍坊影视综艺", "wfkjwl": "潍坊科教文旅",
	"wfgxq": "潍坊高新区", "wfaqzh": "安丘综合", "wfaqms": "安丘民生", "wfbhxw": "滨海新闻综合",
	"wfclzh": "昌乐综合", "wfcyzh": "昌邑综合", "wfcyjj": "昌邑经济生活", "wffzxw": "坊子新闻综合",
	"wfgmzh": "高密综合", "wfgmdj": "高密党建农科", "wfhtxw": "寒亭新闻", "wfkwtv": "奎文电视台",
	"wflqxw": "临朐新闻综合", "wfqzzh": "青州综合", "wfqzwh": "青州文化旅游", "wfwc": "潍城TV",
	"wfzcxw": "诸城新闻综合", "wfzcsh": "诸城生活娱乐", "ytcd": "长岛TV", "ytfszh": "福山综合",
	"ytfssh": "福山生活", "ythyzh": "海阳综合", "ytlkzh": "龙口综合", "ytlksh": "龙口生活",
	"ytlszh": "莱山综合", "ytlsys": "莱山影视", "ytlyzh": "莱阳综合", "ytlyms": "莱阳民生综艺",
	"ytlzzh": "莱州综合", "ytmpzh": "牟平综合", "ytplzh": "蓬莱综合", "ytplzy": "蓬莱综艺",
	"ytqxzh": "栖霞综合", "ytqxpg": "栖霞苹果", "ytzyzh": "招远综合", "ytzyzy": "招远综艺",
	"zbbsxw": "博山新闻", "zbbstw": "博山图文", "zbgqzh": "高青综合", "zbgqys": "高青影视",
	"zbht1": "桓台综合", "zbht2": "桓台影视", "zblzxw": "临淄新闻综合", "zblzsh": "临淄生活服务",
	"zbyyzh": "沂源综合", "zbyysh": "沂源生活", "zbzcxw": "淄川新闻", "zbzcsh": "淄川生活",
	"zbzd1": "张店综合", "zbzd2": "张店2", "zbzctv1": "周村新闻", "zbzctv2": "周村生活",
	"zzstzh": "山亭综合", "zzszzh": "枣庄市中综合", "zztezxw": "台儿庄新闻综合", "zztzzh": "滕州综合",
	"zztzms": "滕州民生", "zzxcxw": "薛城新闻综合", "zzyczh": "峄城综合", "bzbctv": "滨城TV",
	"bzbxzh": "博兴综合", "bzbxsh": "博兴生活", "bzhmzh": "惠民综合", "bzhmys": "惠民影视",
	"bzwdzh": "无棣综合", "bzwdzy": "无棣综艺", "bzyxxw": "阳信新闻综合", "bzzhzh": "沾化综合",
	"bzzhzy": "沾化综艺", "bzzpzh": "邹平综合", "bzzpms": "邹平民生", "dzxwzh": "德州新闻综合",
	"dzjjsh": "德州经济生活", "dztw": "德州图文", "dzlczh": "陵城综合", "dzllxw": "乐陵新闻综合",
	"dzllcs": "乐陵城市生活", "dzly1": "临邑1", "dzly2": "临邑2", "dznjzh": "宁津综合",
	"dzpyzh": "平原综合", "dzqhzh": "齐河综合", "dzqyzh": "庆云综合", "dzqysh": "庆云生活",
	"dzwczh": "武城综合", "dzwczy": "武城综艺影视", "dzxjzh": "夏津综合", "dzxjgg": "夏津公共",
	"dzyczh": "禹城综合", "dzyczy": "禹城综艺", "hzcwzh": "成武综合", "hzcwzy": "成武综艺",
	"hzcxzh": "曹县综合", "hzdmxw": "东明新闻综合", "hzdt1": "定陶新闻", "hzdt2": "定陶综艺",
	"hzjczh": "鄄城综合", "hzjyxw": "巨野新闻", "hzmdxw": "牡丹区新闻综合", "hzmdzy": "牡丹区综艺",
	"hzsxzh": "单县综合", "hzycxw": "郓城新闻", "hzyczy": "郓城综艺", "jijiazh": "嘉祥综合",
	"jijiash": "嘉祥生活", "jijxzh": "金乡综合", "jijxsh": "金乡生活", "jilszh": "梁山综合",
	"jiqfxw": "曲阜新闻综合", "jircxw": "任城新闻综合", "jircys": "任城影视娱乐", "jissxw": "泗水新闻综合",
	"jisswh": "泗水文化生活", "jiws1": "微山综合", "jiws2": "微山2套", "jiwszh": "汶上综合",
	"jiytxw": "鱼台新闻", "jiytsh": "鱼台生活", "jiyzxw": "兖州新闻", "jiyzsh": "兖州生活",
	"jizczh": "邹城综合", "jizcwh": "邹城文化生活", "lccpzh": "茌平综合", "lccpsh": "茌平生活",
	"lcdczh": "东昌综合", "lcdezh": "东阿综合", "lcdezy": "东阿综艺", "lcgtzh": "高唐综合",
	"lcgtzy": "高唐综艺", "lcgxzh": "冠县综合", "lclqzh": "临清综合", "lclqjj": "临清经济信息",
	"lcsxzh": "莘县综合", "lcsxsh": "莘县生活", "lcygzh": "阳谷综合", "lcygys": "阳谷影视",
	"lyfxzh": "费县综合", "lyfxsh": "费县生活", "lyhdys": "河东影视", "lyhdzh": "河东综合",
	"lyjnzh": "莒南综合", "lyjnys": "莒南影视", "lyllzh": "兰陵综合", "lyllgg": "兰陵公共",
	"lylszh": "兰山综合", "lyls1": "临沭综合", "lyls2": "临沭生活", "lylzzh": "罗庄综合",
	"lylzys": "罗庄影视", "lymy1": "蒙阴综合", "lymy2": "蒙阴2套", "lypyzh": "平邑综合",
	"lypysh": "平邑生活", "lytc1": "郯城综合", "lytc2": "郯城2套", "lyynzh": "沂南综合",
	"lyynys": "沂南红色影视", "lyys1": "沂水综合", "lyys2": "沂水生活", "rzjx1": "莒县综合",
	"rzjx2": "莒县2套", "rzls": "岚山TV", "rzwlzh": "五莲综合", "rzwlwh": "五莲文化旅游",
	"tadpzh": "东平综合", "tadpms": "东平民生", "tady": "岱岳TV", "tafczh": "肥城综合",
	"tafcsh": "肥城生活", "tany1": "宁阳综合", "tany2": "宁阳2套", "tats": "泰山TV",
	"taxtzh": "新泰综合", "taxtxc": "新泰乡村", "whxwzh": "威海新闻综合", "whdssh": "威海都市生活",
	"whhy": "威海海洋", "whhczh": "威海环翠综合", "whrczh": "荣成综合", "whrcsh": "荣成生活",
	"whrszh": "乳山综合", "whrssh": "乳山生活", "whwd1": "文登TV1", "whwd2": "文登TV2",
}

// 城市分组映射
var cityGroups = map[string][]string{
	"济南": {"jncqxw", "jncqsh", "jnjrtv", "jnjyzh", "jnjyys", "jnlcxw", "jnpyzh", "jnpyxc", "jnshzh", "jnshys", "jnzqzh", "jnzqgg"},
	"东营": {"dyxwzh", "dygg", "dygg2", "dykj", "dydyqxw", "dydyqkj", "dygrzh", "dygrkj", "dyklxw", "dyljzh", "dyljwh"},
	"青岛": {"qdcyzh", "qdhdzh", "qdhdsh", "qdjmzh", "qdjmsh", "qdjzzh", "qdjzsh", "qdlc", "qdls", "qdlxzh", "qdlxsh", "qdpdxw", "qdpdsh"},
	"潍坊": {"wfxwzh", "wfsh", "wfyswy", "wfkjwl", "wfgxq", "wfaqzh", "wfaqms", "wfbhxw", "wfclzh", "wfcyzh", "wfcyjj", "wffzxw", "wfgmzh", "wfgmdj", "wfhtxw", "wfkwtv", "wflqxw", "wfqzzh", "wfqzwh", "wfwc", "wfzcxw", "wfzcsh"},
	"烟台": {"ytcd", "ytfszh", "ytfssh", "ythyzh", "ytlkzh", "ytlksh", "ytlszh", "ytlsys", "ytlyzh", "ytlyms", "ytlzzh", "ytmpzh", "ytplzh", "ytplzy", "ytqxzh", "ytqxpg", "ytzyzh", "ytzyzy"},
	"淄博": {"zbbsxw", "zbbstw", "zbgqzh", "zbgqys", "zbht1", "zbht2", "zblzxw", "zblzsh", "zbyyzh", "zbyysh", "zbzcxw", "zbzcsh", "zbzd1", "zbzd2", "zbzctv1", "zbzctv2"},
	"枣庄": {"zzstzh", "zzszzh", "zztezxw", "zztzzh", "zztzms", "zzxcxw", "zzyczh"},
	"滨州": {"bzbctv", "bzbxzh", "bzbxsh", "bzhmzh", "bzhmys", "bzwdzh", "bzwdzy", "bzyxxw", "bzzhzh", "bzzhzy", "bzzpzh", "bzzpms"},
	"德州": {"dzxwzh", "dzjjsh", "dztw", "dzlczh", "dzllxw", "dzllcs", "dzly1", "dzly2", "dznjzh", "dzpyzh", "dzqhzh", "dzqyzh", "dzqysh", "dzwczh", "dzwczy", "dzxjzh", "dzxjgg", "dzyczh", "dzyczy"},
	"菏泽": {"hzcwzh", "hzcwzy", "hzcxzh", "hzdmxw", "hzdt1", "hzdt2", "hzjczh", "hzjyxw", "hzmdxw", "hzmdzy", "hzsxzh", "hzycxw", "hzyczy"},
	"济宁": {"jijiazh", "jijiash", "jijxzh", "jijxsh", "jilszh", "jiqfxw", "jircxw", "jircys", "jissxw", "jisswh", "jiws1", "jiws2", "jiwszh", "jiytxw", "jiytsh", "jiyzxw", "jiyzsh", "jizczh", "jizcwh"},
	"聊城": {"lccpzh", "lccpsh", "lcdczh", "lcdezh", "lcdezy", "lcgtzh", "lcgtzy", "lcgxzh", "lclqzh", "lclqjj", "lcsxzh", "lcsxsh", "lcygzh", "lcygys"},
	"临沂": {"lyfxzh", "lyfxsh", "lyhdys", "lyhdzh", "lyjnzh", "lyjnys", "lyllzh", "lyllgg", "lylszh", "lyls1", "lyls2", "lylzzh", "lylzys", "lymy1", "lymy2", "lypyzh", "lypysh", "lytc1", "lytc2", "lyynzh", "lyynys", "lyys1", "lyys2"},
	"日照": {"rzjx1", "rzjx2", "rzls", "rzwlzh", "rzwlwh"},
	"泰安": {"tadpzh", "tadpms", "tady", "tafczh", "tafcsh", "tany1", "tany2", "tats", "taxtzh", "taxtxc"},
	"威海": {"whxwzh", "whdssh", "whhy", "whhczh", "whrczh", "whrcsh", "whrszh", "whrssh", "whwd1", "whwd2"},
}

// 请求头配置
var headers = map[string]string{
	"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Accept":          "application/json, text/plain, */*",
	"Accept-Language": "zh-CN,zh;q=0.9,en;q=0.8",
	"Referer":         "https://app.litenews.cn/",
	"Origin":          "https://app.litenews.cn",
}

const apiBaseURL = "https://app.litenews.cn/v1/app/play/tv/live?_orgid_="

// 直播流结构体
type streamItem struct {
	ID     int    `json:"id"`
	Stream string `json:"stream"`
}
type apiResponse struct {
	Data []streamItem `json:"data"`
}

// 获取直播流地址
func getStreamURL(channelKey string) (string, error) {
	ids, ok := channelMap[channelKey]
	if !ok {
		return "", fmt.Errorf("未知的频道标识: %s", channelKey)
	}
	orgid, targetId := ids[0], ids[1]
	apiURL := fmt.Sprintf("%s%d", apiBaseURL, orgid)

	req, _ := http.NewRequest("GET", apiURL, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{Timeout: 10 * time.Second}
	time.Sleep(500 * time.Millisecond)

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API请求失败: %v", err)
	}
	defer resp.Body.Close()

	var res apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("API返回数据格式错误: %v", err)
	}

	for _, item := range res.Data {
		if item.ID == targetId && item.Stream != "" {
			return item.Stream, nil
		}
	}
	return "", fmt.Errorf("未找到对应的直播流，orgid: %d, target_id: %d", orgid, targetId)
}

// 切片转字符串
func join(slice []string, sep string) string {
	var res string
	for i, v := range slice {
		if i > 0 {
			res += sep
		}
		res += v
	}
	return res
}

func main() {
	gin.DisableConsoleColor()
	r := gin.Default()

	// 首页：TiviMate格式分组频道列表
	r.GET("/", func(c *gin.Context) {
		host := c.Request.Host
		var lines []string
		for city, keys := range cityGroups {
			lines = append(lines, fmt.Sprintf("%s频道,#genre#", city))
			for _, key := range keys {
				if name, ok := cnName[key]; ok {
					lines = append(lines, fmt.Sprintf("%s,http://%s/play?id=%s", name, host, key))
				}
			}
			lines = append(lines, "")
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(join(lines, "\n")))
	})

	// 播放接口：302重定向
	r.GET("/play", func(c *gin.Context) {
		channelKey := c.DefaultQuery("id", "jncqxw")
		streamURL, err := getStreamURL(channelKey)
		if err != nil {
			c.String(http.StatusInternalServerError, "错误: %v", err)
			return
		}
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
		c.Redirect(http.StatusFound, streamURL)
	})

	// 测试接口
	r.GET("/test/:channelKey", func(c *gin.Context) {
		channelKey := c.Param("channelKey")
		ids, ok := channelMap[channelKey]
		if !ok {
			c.String(http.StatusNotFound, "未知频道: %s", channelKey)
			return
		}
		orgid, targetId := ids[0], ids[1]
		apiURL := fmt.Sprintf("%s%d", apiBaseURL, orgid)

		req, _ := http.NewRequest("GET", apiURL, nil)
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%v", err)})
			return
		}
		defer resp.Body.Close()

		var res map[string]interface{}
		_ = json.NewDecoder(resp.Body).Decode(&res)
		c.JSON(http.StatusOK, gin.H{
			"channel":        channelKey,
			"channel_name":   cnName[channelKey],
			"orgid":          orgid,
			"target_id":      targetId,
			"api_url":        apiURL,
			"status_code":    resp.StatusCode,
			"data_structure": res,
		})
	})

	// 简单列表接口
	r.GET("/simple", func(c *gin.Context) {
		host := c.Request.Host
		var lines []string
		for key := range channelMap {
			lines = append(lines, fmt.Sprintf("%s,http://%s/play?id=%s", cnName[key], host, key))
		}
		c.Data(http.StatusOK, "text/plain; charset=utf-8", []byte(join(lines, "\n")))
	})

	// 端口配置（适配容器环境变量）
	port := os.Getenv("PORT")
	if port == "" {
		port = "9003"
	}
	r.Run(fmt.Sprintf("0.0.0.0:%s", port))
}