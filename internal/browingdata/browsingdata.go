package browingdata

import (
	"path"

	"github.com/wenkechen/hack-browser-data/internal/browingdata/bookmark"
	"github.com/wenkechen/hack-browser-data/internal/browingdata/cookie"
	"github.com/wenkechen/hack-browser-data/internal/browingdata/creditcard"
	"github.com/wenkechen/hack-browser-data/internal/browingdata/download"
	"github.com/wenkechen/hack-browser-data/internal/browingdata/extension"
	"github.com/wenkechen/hack-browser-data/internal/browingdata/history"
	"github.com/wenkechen/hack-browser-data/internal/browingdata/localstorage"
	"github.com/wenkechen/hack-browser-data/internal/browingdata/password"
	"github.com/wenkechen/hack-browser-data/internal/item"
	"github.com/wenkechen/hack-browser-data/internal/log"
	"github.com/wenkechen/hack-browser-data/internal/utils/fileutil"
)

type Data struct {
	sources map[item.Item]Source
}

type Source interface {
	Parse(masterKey []byte) error

	Name() string

	Length() int
}

func New(sources []item.Item) *Data {
	bd := &Data{
		sources: make(map[item.Item]Source),
	}
	bd.addSource(sources)
	return bd
}

func (d *Data) Recovery(masterKey []byte) error {
	for _, source := range d.sources {
		if err := source.Parse(masterKey); err != nil {
			log.Errorf("parse %s error %s", source.Name(), err.Error())
		}
	}
	return nil
}

func (d *Data) Output(dir, browserName, flag string) {
	output := NewOutPutter(flag)

	for _, source := range d.sources {
		if source.Length() == 0 {
			// if the length of the export data is 0, then it is not necessary to output
			continue
		}
		filename := fileutil.Filename(browserName, source.Name(), output.Ext())

		f, err := output.CreateFile(dir, filename)
		if err != nil {
			log.Errorf("create file error %s", err)
		}
		if err := output.Write(source, f); err != nil {
			log.Errorf("%s write to file %s error %s", source.Name(), filename, err.Error())
		}
		log.Noticef("output to file %s success", path.Join(dir, filename))
	}
}

func (d *Data) addSource(Sources []item.Item) {
	for _, source := range Sources {
		switch source {
		case item.ChromiumPassword:
			d.sources[source] = &password.ChromiumPassword{}
		case item.ChromiumCookie:
			d.sources[source] = &cookie.ChromiumCookie{}
		case item.ChromiumBookmark:
			d.sources[source] = &bookmark.ChromiumBookmark{}
		case item.ChromiumHistory:
			d.sources[source] = &history.ChromiumHistory{}
		case item.ChromiumDownload:
			d.sources[source] = &download.ChromiumDownload{}
		case item.ChromiumCreditCard:
			d.sources[source] = &creditcard.ChromiumCreditCard{}
		case item.ChromiumLocalStorage:
			d.sources[source] = &localstorage.ChromiumLocalStorage{}
		case item.ChromiumExtension:
			d.sources[source] = &extension.ChromiumExtension{}
		case item.YandexPassword:
			d.sources[source] = &password.YandexPassword{}
		case item.YandexCreditCard:
			d.sources[source] = &creditcard.YandexCreditCard{}
		case item.FirefoxPassword:
			d.sources[source] = &password.FirefoxPassword{}
		case item.FirefoxCookie:
			d.sources[source] = &cookie.FirefoxCookie{}
		case item.FirefoxBookmark:
			d.sources[source] = &bookmark.FirefoxBookmark{}
		case item.FirefoxHistory:
			d.sources[source] = &history.FirefoxHistory{}
		case item.FirefoxDownload:
			d.sources[source] = &download.FirefoxDownload{}
		case item.FirefoxLocalStorage:
			d.sources[source] = &localstorage.FirefoxLocalStorage{}
		case item.FirefoxExtension:
			d.sources[source] = &extension.FirefoxExtension{}
		}
	}
}
