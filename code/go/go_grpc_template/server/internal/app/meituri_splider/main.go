package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
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

func main() {

	pwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
		return
	}

	targetUrls, err := os.Open(path.Join(pwd, "code/go/go_grpc_template/server/config/meituri_splider/target_urls.txt"))
	if err != nil {
		log.Println(err)
		return
	}
	defer targetUrls.Close()

	targetUrlsReader := bufio.NewReader(targetUrls)

	for order := 0; ; order++ {
		line, c := targetUrlsReader.ReadString('\n')
		line = strings.TrimSuffix(line, "\n")
		items := strings.Split(line, "/:")
		if len(items) < 2 {
			log.Println(items)
			continue
		}

		fmt.Println(" 次序：", order+1, " ----- ", items[1], " ----- ", items[0])
		job(items[0] + "/")
		if c == io.EOF {
			break
		}
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
		return
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
		// 错误处理
		log.Println(err)
		return
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
