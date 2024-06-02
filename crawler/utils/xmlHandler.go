package utils

import (
	"encoding/xml"
	"io"
	"net/http"
)

type ParsedData struct {
	Data []Post `xml:"channel>item"`
}

type Post struct {
	Title   string `xml:"title"`
	Link    string `xml:"link"`
	PubDate string `xml:"pubDate"`
}

func GetParsedData(url string) []Post {
	res, err := http.Get(url)
	CheckErr(err)
	CheckHttpResponse(res)
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	CheckErr(err)
	var posts ParsedData
	xmlerr := xml.Unmarshal(data, &posts)
	CheckErr(xmlerr)
	return posts.Data
}

func CheckUpdatedPost(posts []Post, domainURL string, updatedDate int64) int8 {
	lastUpdatedDate := UnixTime2Time(updatedDate)
	var index int8 = 0
	pathStartIdx := len("https://") + len(domainURL)
	for index < int8(len(posts)) {
		post := posts[index]
		pubDate := Str2time(post.PubDate)
		if pubDate.Compare(lastUpdatedDate) == 1 {
			posts[index].Link = post.Link[pathStartIdx:]
			index++ // check next post when it needs to update
		} else {
			break
		}
	}
	return index - 1
}