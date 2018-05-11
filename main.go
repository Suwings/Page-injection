package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
)

//Start...
func main() {
	//goproxy
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	proxy.OnRequest().DoFunc(
		func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			//请求事件代码
			//此处可修改请求头和信息
			return req, nil
		})

	proxy.OnResponse().DoFunc(
		func(res *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
			//响应事件代码
			//此处可修改响应头和信息
			bs, _ := ioutil.ReadAll(res.Body)
			tmpstr := string(bs)
			//修改 HTTP Body
			newbody := HandleBody(res, &tmpstr)
			res.Body = ioutil.NopCloser(bytes.NewReader([]byte(*newbody)))
			return res
		})

	fmt.Println("[Proxy] Proxy Server Listen....")
	log.Fatal(http.ListenAndServe(":8080", proxy))
}

//HTTP Body 修改
func HandleBody(res *http.Response, body *string) *string {
	host := res.Request.Host
	realDir := FindRealDir(host)
	if realDir == "" {
		fmt.Println("Not Find Script:", host)
		return body
	}
	//插入指定目录的 JavaScript 到网页的最后
	includePath := "./" + realDir + "/" + "include.js"
	fmt.Println("HandleBody Host:", host, " Include Loacl Script:", includePath)
	b, _ := ioutil.ReadFile(includePath)
	jsfilestring := string(b)
	restr := "\n\n<!-- Include Script -->\n <script>" + jsfilestring + "</script>\n <!-- Include Script End--> \n</html>"
	//自定义插入点
	str := strings.Replace(*body, "</html>", restr, -1)
	return &str
}

//通过域名寻找返回指定目录
func FindRealDir(host string) string {
	dir_list, err := ioutil.ReadDir("./")
	if err != nil {
		return ""
	}
	for _, dirinfo := range dir_list {
		dirname := dirinfo.Name()
		if dirname == host {
			return dirname
		}
		if strings.Contains(host, dirname) {
			return dirname
		}
	}
	return ""
}
