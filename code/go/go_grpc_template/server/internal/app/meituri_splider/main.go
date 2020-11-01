package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const FilePath = "/Users/coffee/Pictures/tujigu/"
const webPath = "https://www.tujigu.com"
const recordName = "/Users/coffee/Pictures/tujigu/tujigu.xlsx"

var blockKeywords = []string{}
var whiteKeywords = []string{}
var newDoc = []string{}
var allDoc = map[string]string{
	"https://www.tujigu.com/a/28235/": "匿名",
	"https://www.tujigu.com/a/33875/": "萝莉",
	"https://www.tujigu.com/a/34209/": "馨馨",
	"https://www.tujigu.com/a/34661/": "沐娅",
	"https://www.tujigu.com/a/34702/": "匿名",
	"https://www.tujigu.com/a/34704/": "匿名",
	"https://www.tujigu.com/a/34831/": "dodo",
	"https://www.tujigu.com/t/1178/":  "玉兔miki",
	"https://www.tujigu.com/t/1185/":  "妲己_Toxic",
	"https://www.tujigu.com/t/1194/":  "周妍希",
	"https://www.tujigu.com/t/1289/":  "黄楽然",
	"https://www.tujigu.com/t/161/":   "糯美子Mini",
	"https://www.tujigu.com/t/190/":   "易阳",
	"https://www.tujigu.com/t/2189/":  "朴善慧",
	"https://www.tujigu.com/t/2196/":  "许允美",
	"https://www.tujigu.com/t/2234/":  "张楚珊",
	"https://www.tujigu.com/t/242/":   "沈梦瑶",
	"https://www.tujigu.com/t/2693/":  "天使萌",
	"https://www.tujigu.com/t/2783/":  "大安妮",
	"https://www.tujigu.com/t/292/":   "夏小秋",
	"https://www.tujigu.com/t/293/":   "王雨纯",
	"https://www.tujigu.com/t/295/":   "朱可儿",
	"https://www.tujigu.com/t/296/":   "许诺",
	"https://www.tujigu.com/t/298/":   "夏美酱",
	"https://www.tujigu.com/t/3156/":  "周于希",
	"https://www.tujigu.com/t/3160/":  "little贝壳",
	"https://www.tujigu.com/t/3171/":  "小狐狸Sica",
	"https://www.tujigu.com/t/3261/":  "杨暖",
	"https://www.tujigu.com/t/3307/":  "宁宁",
	"https://www.tujigu.com/t/4218/":  "任莹樱",
	"https://www.tujigu.com/t/446/":   "陈潇",
	"https://www.tujigu.com/t/4476/":  "芸斐",
	"https://www.tujigu.com/t/4530/":  "姜璐",
	"https://www.tujigu.com/t/4568/":  "Lavinia肉肉",
	"https://www.tujigu.com/t/459/":   "杨晨晨",
	"https://www.tujigu.com/t/4640/":  "徐微微",
	"https://www.tujigu.com/t/465/":   "菲儿",
	"https://www.tujigu.com/t/4708/":  "米米",
	"https://www.tujigu.com/t/5034/":  "姜仁卿",
	"https://www.tujigu.com/t/5109/":  "桜桃喵",
	"https://www.tujigu.com/t/5110/":  "疯猫ss",
	"https://www.tujigu.com/t/5466/":  "克拉女神芊芊",
	"https://www.tujigu.com/t/5496/":  "星之迟迟",
	"https://www.tujigu.com/t/5497/":  "面饼仙儿",
	"https://www.tujigu.com/t/5499/":  "雯妹不讲道理",
	"https://www.tujigu.com/t/5500/":  "一小央泽",
	"https://www.tujigu.com/t/5501/":  "鬼畜瑶",
	"https://www.tujigu.com/t/5504/":  "魔物喵",
	"https://www.tujigu.com/t/5505/":  "你的负卿",
	"https://www.tujigu.com/t/5511/":  "雪琪",
	"https://www.tujigu.com/t/5513/":  "白银81",
	"https://www.tujigu.com/t/5514/":  "米线线sama",
	"https://www.tujigu.com/t/5515/":  "黑川",
	"https://www.tujigu.com/t/5529/":  "李浅浅",
	"https://www.tujigu.com/t/5533/":  "袁圆",
	"https://www.tujigu.com/t/5559/":  "爱丽丝",
	"https://www.tujigu.com/t/5609/":  "小牛奶",
	"https://www.tujigu.com/t/5613/":  "林文文",
	"https://www.tujigu.com/t/5652/":  "郭子蜜",
	"https://www.tujigu.com/t/5672/":  "奈汐酱",
	"https://www.tujigu.com/t/5674/":  "UU酱",
	"https://www.tujigu.com/t/5697/":  "秋楚楚",
	"https://www.tujigu.com/t/5712/":  "古川kagura",
	"https://www.tujigu.com/t/5714/":  "戚顾儿",
	"https://www.tujigu.com/t/5716/":  "南桃Momoko",
	"https://www.tujigu.com/t/5720/":  "从从从从鸾",
	"https://www.tujigu.com/t/5721/":  "是依酱呀",
	"https://www.tujigu.com/t/5725/":  "清纯妹子西瓜",
	"https://www.tujigu.com/t/5728/":  "雪晴Astra",
	"https://www.tujigu.com/t/5729/":  "-白烨-",
	"https://www.tujigu.com/t/5736/":  "绮太郎",
	"https://www.tujigu.com/t/5740/":  "Nyako喵子",
	"https://www.tujigu.com/t/5743/":  "蠢沫沫",
	"https://www.tujigu.com/t/5744/":  "弥音音ww",
	"https://www.tujigu.com/t/5745/":  "镜酱",
	"https://www.tujigu.com/t/5749/":  "樱落酱w",
	"https://www.tujigu.com/t/5755/":  "萌芽儿o0",
	"https://www.tujigu.com/t/5756/":  "十万珍吱伏特",
	"https://www.tujigu.com/t/5762/":  "眼酱大魔王w",
	"https://www.tujigu.com/t/5765/":  "蜜汁猫裘",
	"https://www.tujigu.com/t/5768/":  "是青水",
	"https://www.tujigu.com/t/5774/":  "抖娘-利世",
	"https://www.tujigu.com/t/5777/":  "南鸽",
	"https://www.tujigu.com/t/5794/":  "舒彤",
	"https://www.tujigu.com/t/5799/":  "西景",
	"https://www.tujigu.com/t/5800/":  "子吟",
	"https://www.tujigu.com/t/654/":   "玛鲁娜",
	"https://www.tujigu.com/t/657/":   "蛋糕Cake",
	"https://www.tujigu.com/t/663/":   "夏天Sienna",
	"https://www.tujigu.com/t/903/":   "娜露",
	"https://www.tujigu.com/t/919/":   "李雅",
	"https://www.tujigu.com/x/86/":    "风之领域",
	"https://www.tujigu.com/x/95/":    "喵糖映画写真",
	"https://www.tujigu.com/t/5788/":  "腿模YuHsuan",
	"https://www.tujigu.com/t/290/":   "刘娅希",
	"https://www.tujigu.com/t/5798/":  "蓓颖",
	"https://www.tujigu.com/t/150/":   "于姬",
	"https://www.tujigu.com/t/2438/":  "小热巴",
	"https://www.tujigu.com/t/5790/":  "申才恩",
	"https://www.tujigu.com/t/5472/":  "杨紫嫣",
	"https://www.tujigu.com/t/4780/":  "小尤奈",
	"https://www.tujigu.com/t/4729/":  "九月生",
	"https://www.tujigu.com/t/5553/":  "陶喜乐",
	"https://www.tujigu.com/t/5375/":  "林子欣",
	"https://www.tujigu.com/t/5177/":  "莉娜lena",
	"https://www.tujigu.com/t/455/":   "陈思琪",
	"https://www.tujigu.com/t/954/":   "芝芝Bootyw",
	"https://www.tujigu.com/t/4634/":  "美替",
	"https://www.tujigu.com/t/4242/":  "陈亦菲",
	"https://www.tujigu.com/t/5175/":  "LindaLinda",
	"https://www.tujigu.com/t/4561/":  "水花花",
	"https://www.tujigu.com/t/4441/":  "韩羽",
	"https://www.tujigu.com/t/3852/":  "Egg_尤妮丝",
	"https://www.tujigu.com/t/3243/":  "敏珺",
	"https://www.tujigu.com/t/674/":   "佟蔓",
}

func main() {

	var targetUrls []string
	for targetUrl, _ := range allDoc {
		targetUrls = append(targetUrls, targetUrl)
	}
	sort.Strings(targetUrls)

	for order, targetUrl := range targetUrls {
		fmt.Println(" 次序：", order+1, "/", len(targetUrls), " ----- ", allDoc[targetUrl], " ----- ", targetUrl)
		job(targetUrl)
	}

	log.Println("new document:", len(newDoc))
	for i, _ := range newDoc {
		log.Println(i+1, " ----- ", newDoc[i])
	}

}

func job(targetUrl string) {
	log.Println(targetUrl)
	urlObj, err := url.Parse(targetUrl)
	if err != nil {
		log.Println(err)
		return
	}
	switch {
	case strings.Contains(targetUrl, webPath+"/t"):
		mote_name, urls := FindDocument(targetUrl, "t")
		fmt.Println("出境模特：", mote_name)
		if mote_name == "" {
			return
		}
		download_path := strings.Replace(urlObj.EscapedPath()+"__"+mote_name, "/", "", -1)
		for i, val := range urls {
			fmt.Println("[", i+1, "/", len(urls), "]")
			Download(val, download_path)
		}
	case strings.Contains(targetUrl, webPath+"/x"):
		mote_name, urls := FindDocument(targetUrl, "x")
		fmt.Println("出境模特：", mote_name)
		download_path := mote_name
		if download_path == "" {
			download_path = "default"
		}
		for i, val := range urls {
			fmt.Println("[", i+1, "/", len(urls), "]")
			Download(val, download_path)
		}
	case strings.Contains(targetUrl, webPath+"/a"):
		Download(targetUrl, "default")
	}
	log.Println("done...")
}

//=====================================

func Download(url string, mote_name string) string {
	var msg string
	document := NewDocument(url, mote_name)
	if document != nil {
		document.FindAll()
		document.SaveContents()
		msg += "save content length = " + strconv.Itoa(len(document.content)) + "\n"
	}
	return msg
}

//=====================================

func FindDocument(url string, feature string) (mote string, a []string) {

	var urls = map[string]bool{url: true}
	if doc := GetDocument(url); doc != nil {
		// Find all page
		doc.Find("center").Each(func(i int, s *goquery.Selection) {
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				if strings.Contains(href, feature+"/") {
					if !strings.Contains(href, "http") {
						href = webPath + href
					}
					urls[href] = true
				}
			})
		})
		// Find mote 模特
		doc.Find("div").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			if val, exist := s.Attr("class"); !exist || val != "left" {
				return
			}
			s.Find("img").Each(func(i int, s *goquery.Selection) {
				alt, _ := s.Attr("alt")
				mote = strings.Split(strings.Replace(strings.Replace(alt, " ", "", -1), "/", "", -1), "、")[0]
			})
		})
	}

	var documents = map[string]bool{}
	for k, _ := range urls {
		doc := GetDocument(k)
		if doc == nil {
			continue
		}
		doc.Find("div").Each(func(i int, s *goquery.Selection) {
			if val, exist := s.Attr("class"); !exist || val != "hezi" {
				return
			}
			s.Find("a").Each(func(i int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				if strings.Contains(href, "a") {
					documents[href] = true
				}
			})
		})
	}

	var ret []string
	for k, _ := range documents {
		ret = append(ret, k)
	}
	return mote, ret
}

//=====================================

type document struct {
	url      string
	filepath string             //保存content的文件路径
	content  map[string]content //key=src
	pages    map[string]pages   //key=page
}

type content struct {
	src string
	alt string
}

type pages struct {
	page string
	href string
}

func NewDocument(urlStr, mote_name string) *document {
	urlObj, err := url.Parse(urlStr)
	if urlStr == "" || err != nil {
		log.Println(err)
		return nil
	}

	//判断文件夹是否存在,并创建
	filepath := path.Join(FilePath+mote_name, strings.Replace(urlObj.EscapedPath(), "a/", "", -1))
	if !PathExists(filepath) {
		if err := os.MkdirAll(filepath, os.ModePerm); err != nil {
			log.Println(err)
			return nil
		}
	} else {
		log.Println("document 已存在")
		return nil
	}

	log.Println("new document = ", urlStr, ",save path = ", filepath)
	newDoc = append(newDoc, filepath)
	return &document{
		url:      urlObj.String(),
		filepath: filepath,
		content:  map[string]content{},
		pages:    map[string]pages{},
	}
}

func (this *document) FindAll() {
	for this.url != "" {
		this.Find()

		if next := this.Next(); next == "" || next == this.url {
			this.url = ""
		} else {
			this.url = next
		}
	}
}

func (this *document) Find() {
	if doc := GetDocument(this.url); doc != nil {
		this.FindContents(doc)
		this.FindPages(doc)
	}
}

func (this *document) Next() string {
	if len(this.pages) == 0 {
		return ""
	}
	page, ok := this.pages["下一页"]
	if !ok {
		return ""
	}
	return page.href
}

func (this *document) FindContents(document *goquery.Document) {
	if this.content == nil {
		this.content = make(map[string]content)
	}

	// Find the review items
	document.Find("div").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		if val, exist := s.Attr("class"); !exist || val != "content" {
			return
		}
		s.Find("img").Each(func(i int, s *goquery.Selection) {
			src, _ := s.Attr("src")
			alt, _ := s.Attr("alt")
			class, _ := s.Attr("class")
			if class != "tupian_img" {
				return
			}
			this.content[src] = content{
				src: src,
				alt: strings.Replace(strings.Replace(alt, " ", "", -1), "/", "", -1),
			}
		})
	})
}

func (this *document) FindPages(document *goquery.Document) {
	if this.pages == nil {
		this.pages = make(map[string]pages)
	}

	// Find the review items
	document.Find("div").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		if val, exist := s.Attr("id"); !exist || val != "pages" {
			return
		}
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")
			page := s.Text()
			this.pages[page] = pages{
				page: page,
				href: href,
			}
		})
	})
}

func (this *document) SaveContents() {
	fmt.Print("content: ", len(this.content), ", download...")
	var i, percent = 0, 0
	for _, content := range this.content {
		xiezhen := strings.Split(this.filepath, "/")[len(strings.Split(FilePath, "/"))]
		content.alt = xiezhen + "__" + content.alt
		download(content, this.filepath)

		i++
		if tmp := i * 100 / len(this.content); tmp%10 == 0 && tmp != percent {
			percent = tmp
			fmt.Print(">>", percent, "% ")
		}
	}
	fmt.Println()
}

func download(content content, filepath string) {
	for i, _ := range blockKeywords {
		if strings.Contains(content.alt, blockKeywords[i]) {
			for i, _ := range whiteKeywords {
				if !strings.Contains(content.alt, whiteKeywords[i]) {
					log.Println("blockKeywords:", blockKeywords[i], "whiteKeywords", whiteKeywords[i], content.alt, content.src)
					return
				}
			}
		}
	}

	res, err := http.Get(content.src)
	if err != nil || res.StatusCode != http.StatusOK {
		// 重试
		res, err = http.Get(content.src)
		if err != nil {
			return
		}
	}
	defer res.Body.Close()
	data, _ := ioutil.ReadAll(res.Body)
	f, err := os.Create(filepath + "/" + content.alt + ".jpg")
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	f.Write(data)
}

func GetDocument(url string) (document *goquery.Document) {
	// 请求html页面
	res, err := http.Get(url)
	if err != nil {
		// 重试
		res, err = http.Get(url)
		if err != nil {
			// 错误处理
			log.Println(err)
			return
		}
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Println("status code error:", res.StatusCode, res.Status)
		return
	}

	// 加载 HTML document对象
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Println(err)
		return
	}
	return doc
}

//=====================================
// 判断文件夹是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
