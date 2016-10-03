// Package main provides ...
// generate feed.xml and sitemap.xml
package main

import (
	"os"
	"text/template"
	"time"

	"github.com/eiblog/eiblog/setting"
	"github.com/eiblog/utils/logd"
	"github.com/eiblog/utils/tmpl"
)

const (
	FEED_COUNT    = 20
	TEMPLATE_GLOB = "conf/tpl/*.xml"
)

var tpls *template.Template

func init() {
	var err error
	tpls, err = template.New("").Funcs(template.FuncMap{"dateformat": tmpl.DateFormat}).ParseGlob(TEMPLATE_GLOB)
	if err != nil {
		logd.Fatal(err)
	}
	doOpensearch()
	go doFeed()
	go doSitemap()
}

func doFeed() {
	tpl := tpls.Lookup("feedTpl.xml")
	if tpl == nil {
		logd.Error("not found feedTpl.")
		return
	}
	_, _, artcs := PageList(1, FEED_COUNT)
	buildDate := time.Now()
	params := map[string]interface{}{
		"Title":     Ei.BTitle,
		"SubTitle":  Ei.SubTitle,
		"Domain":    setting.Conf.Mode.Domain,
		"BuildDate": buildDate.Format(time.RFC1123Z),
		"Artcs":     artcs,
	}

	f, err := os.OpenFile("static/feed.xml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logd.Error(err)
		return
	}
	defer f.Close()
	err = tpl.Execute(f, params)
	if err != nil {
		logd.Error(err)
		return
	}
	time.AfterFunc(time.Hour*4, doFeed)
}

func doSitemap() {
	tpl := tpls.Lookup("sitemapTpl.xml")
	if tpl == nil {
		logd.Error("not found sitemapTpl.")
		return
	}
	params := map[string]interface{}{"Artcs": Ei.Articles, "Domain": setting.Conf.Mode.Domain}
	f, err := os.OpenFile("static/sitemap.xml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logd.Error(err)
		return
	}
	defer f.Close()
	err = tpl.Execute(f, params)
	if err != nil {
		logd.Error(err)
		return
	}
	time.AfterFunc(time.Hour*24, doFeed)
}

func doOpensearch() {
	tpl := tpls.Lookup("opensearchTpl.xml")
	if tpl == nil {
		logd.Error("not found opensearchTpl.")
		return
	}
	params := map[string]string{
		"BTitle":   Ei.BTitle,
		"SubTitle": Ei.SubTitle,
	}
	f, err := os.OpenFile("static/opensearch.xml", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		logd.Error(err)
		return
	}
	defer f.Close()
	err = tpl.Execute(f, params)
	if err != nil {
		logd.Error(err)
		return
	}
}