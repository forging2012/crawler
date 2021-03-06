package console

import (
	"errors"
	"github.com/golang-collections/collections/set"
	"github.com/golang-collections/collections/stack"
	"github.com/urfave/cli"
	"log"
	"net/url"
)

// Set of visited links
var VisitedLinks = set.New()

// Start analyse URL
var StartUrl *url.URL

// Links for analyse
var LinksStack = stack.New()

// Result links
var ResultLinks = set.New()

// Start parse site
func Parse(c *cli.Context) error {
	prepareUrl, prepareErr := parseStartUrl(c.String("url"))
	if prepareErr != nil {
		return prepareErr
	}
	StartUrl = prepareUrl

	if "" == c.String("filename") {
		return errors.New("Filename is unset")
	}

	LinksStack.Push(Link{
		Link:   c.String("url"),
		Source: "",
	})

	for LinksStack.Len() > 0 {
		analyse(LinksStack.Pop().(Link))
	}

	saveResult(c.String("filename"))

	return nil
}
func parseStartUrl(inputUrl string) (*url.URL, error) {
	if "" == inputUrl {
		return nil, errors.New("URL isn't set.")
	}

	parseUrl, parseErr := url.Parse(inputUrl)
	if parseErr != nil {
		return nil, errors.New("Incorrect URL.")
	}

	return parseUrl, nil
}

func analyse(link Link) {
	if VisitedLinks.Has(link.Link) {
		return
	}
	VisitedLinks.Insert(link.Link)

	parseUrl, parseErr := url.Parse(link.Link)
	if parseErr != nil {
		log.Println("Error parse url")
		return
	}

	if parseUrl.Host == StartUrl.Host {
		RunCrawler(parseUrl.String(), link.Source)
	} else {
		ResultLinks.Insert(Page{
			Link:   parseUrl.String(),
			Source: link.Source,
			Type:   "ext",
			Status: 0,
		})
	}

	return
}

func saveResult(filename string) {
	SaveVisitedLinks(VisitedLinks, filename+".txt")
	SaveResultLinks(ResultLinks, filename+".csv")
}
