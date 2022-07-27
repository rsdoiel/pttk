//
// rss.go is a library for generating RSS 2 documents and is part of the
// ptdk project.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package rss

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
)

type CustomAttrs []xml.Attr

type RSS2 struct {
	XMLName xml.Name `xml:"rss" json:"-"`
	Version string   `xml:"version,attr" json:"version"`

	// Required
	Title       string `xml:"channel>title" json:"title"`
	Link        string `xml:"channel>link" json:"link"`
	Description string `xml:"channel>description" json:"description"`

	// Optional
	Language       string `xml:"channel>language,omitempty" json:"language,omitempty"`
	Copyright      string `xml:"channel>copyright,omitempty" json:"copyright,omitempty"`
	ManagingEditor string `xml:"channel>managingEditor,omitempty" json:"managingEditor,omitempty"`
	WebMaster      string `xml:"channel>webMaster,omitempty" json:"webMaster,omitempty"`
	PubDate        string `xml:"channel>pubDate,omitempty" json:"pubDate,omitempty"`
	LastBuildDate  string `xml:"channel>lastBuildDate,omitempty" json:"lastBuildDate,omitempty"`
	Category       string `xml:"channel>category,omitempty" json:"category,omitempty"`
	Generator      string `xml:"channel>generator,omitempty" json:"generator,omitempty"`
	Docs           string `xml:"channel>docs,omitempty" json:"docs,omitempty"`
	Cloud          string `xml:"channel>cloud,omitempty" json:"cloud,omitempty"`
	TTL            string `xml:"channel>ttl,omitempty" json:"ttl,omitempty"`
	Image          string `xml:"channel>image,omitempty" json:"image,omitempty"`
	Rating         string `xml:"channel>rating,omitempty" json:"rating,omitempty"`
	SkipHours      string `xml:"channel>skipHours,omitempty" json:"skipHours,omitempty"`
	SkipDays       string `xml:"channel>skipDays,omitempty" json:"skipDays,omitempty"`
	Items          []Item `xml:"channel>item,omitempty" json:"item,omitempty"`
}

type Item struct {
	XMLName xml.Name `xml:"item,omitempty" json:"-"`

	// Optional, source was Dave Winer and scripting.com
	Title string `xml:"title" json:"title,omitempty"`

	// Required
	Link string `xml:"link" json:"link"`

	// Optional
	Author      string      `xml:"author,omitempty" json:"author,omitempty"`
	Description string      `xml:"description,omitempty" json:"description,omitempty"`
	Category    string      `xml:"category,omitempty" json:"category,omitempty"`
	Content     string      `xml:"encoded,omitempty" json:"encoded,omitempty"`
	PubDate     string      `xml:"pubDate,omitempty" json:"pubDate,omitempty"`
	Comments    string      `xml:"comments,omitempty" json:"comments,omitempty"`
	Enclosure   string      `xml:"enclosure,omitempty" json:"enclosure,omitempty"`
	GUID        string      `xml:"guid,omitempty" json:"guid,omitempty"`
	Source      string      `xml:"source,omitempty" json:"source,omitempty"`
	OtherAttr   CustomAttrs `xml:",any,attr" json:"other_attrs,omitempty"`
}

type CData struct {
	Value string `xml:",cdata,omitempty" json:"value,omitempty"`
}

func (cdata *CData) Set(src string) {
	cdata.Value = src
}

func (cdata *CData) String() string {
	return cdata.Value
}

func (cdata *CData) ToJSON() string {
	return cdata.Value
}

// MarshalJSON() marshals the custom attributes that might
// be included in an RSS feed.
func (cattr CustomAttrs) MarshalJSON() ([]byte, error) {
	m := map[string]string{}
	for _, attr := range cattr {
		k := attr.Name.Local
		v := attr.Value
		if k != "" {
			m[k] = v
		}
	}
	return json.Marshal(m)
}

// Parse return an RSS2 document as a RSS2 structure.
func Parse(buf []byte) (*RSS2, error) {
	data := new(RSS2)
	err := xml.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RSS2) channel(dataPath string) (map[string]interface{}, error) {
	results := make(map[string]interface{})
	switch {
	case strings.Compare(dataPath, ".channel") == 0:
		// package and return all the channel fields
		results[".title"] = r.Title
		results[".link"] = r.Link
		results[".description"] = r.Description
		if r.PubDate != "" {
			results[".pubDate"] = r.PubDate
		}
	case strings.HasSuffix(dataPath, ".title"):
		results[".title"] = r.Title
	case strings.HasSuffix(dataPath, ".link"):
		results[".link"] = r.Link
	case strings.HasSuffix(dataPath, ".description"):
		results[".description"] = r.Description
	case strings.HasSuffix(dataPath, ".pubDate"):
		results[".pubDate"] = r.PubDate
	default:
		return nil, fmt.Errorf("Unknown data path %s", dataPath)
	}
	return results, nil
}

type rangeExpression struct {
	first int
	last  int
}

func getRange(listLength int, exp string) *rangeExpression {
	rexp := new(rangeExpression)
	rexp.first = 0
	rexp.last = listLength - 1

	if strings.Contains(exp, "-") == true {
		nums := strings.SplitN(exp, "-", 2)
		i, err := strconv.Atoi(nums[0])
		if err == nil {
			rexp.first = i
		}
		i, err = strconv.Atoi(nums[1])
		if err == nil {
			rexp.last = i
		}
	} else {
		i, err := strconv.Atoi(exp)
		if err == nil {
			rexp.first = i
			rexp.last = i
		}
	}
	return rexp
}

func (rexp *rangeExpression) inRange(val int) bool {
	if val >= rexp.first && val <= rexp.last {
		return true
	}
	return false
}

func (r *RSS2) items(dataPath string) (map[string]interface{}, error) {
	rexp := new(rangeExpression)
	rexp.first = 0
	rexp.last = len(r.Items) - 1

	// Get the range expression so we know when to add it to results.
	s := strings.Index(dataPath, "[")
	e := strings.Index(dataPath, "]")
	if s >= 0 && e >= 0 {
		rexp = getRange(len(r.Items), dataPath[s:e])
	}

	results := make(map[string]interface{})
	switch {
	case strings.HasSuffix(dataPath, ".title") == true:
		vals := []string{}
		for i, item := range r.Items {
			if rexp.inRange(i) == true {
				vals = append(vals, item.Title)
			}
		}
		results["title"] = vals
	case strings.HasSuffix(dataPath, ".link") == true:
		vals := []string{}
		for i, item := range r.Items {
			if rexp.inRange(i) == true {
				vals = append(vals, item.Link)
			}
		}
		results["link"] = vals
	case strings.HasSuffix(dataPath, ".description") == true:
		vals := []string{}
		for i, item := range r.Items {
			if rexp.inRange(i) == true {
				vals = append(vals, item.Description)
			}
		}
		results["description"] = vals
	case strings.HasSuffix(dataPath, ".content") == true:
		vals := []string{}
		for i, item := range r.Items {
			if rexp.inRange(i) == true {
				vals = append(vals, item.Content)
			}
		}
		results["content"] = vals
	case strings.HasSuffix(dataPath, ".pubDate") == true:
		vals := []string{}
		for i, item := range r.Items {
			if rexp.inRange(i) == true {
				vals = append(vals, item.PubDate)
			}
		}
		results["pubDate"] = vals
	case strings.HasSuffix(dataPath, ".comments") == true:
		vals := []string{}
		for i, item := range r.Items {
			if rexp.inRange(i) == true {
				vals = append(vals, item.Comments)
			}
		}
		results["comments"] = vals
	}
	return results, nil
}

// Filter given an RSS2 document return all the entries matching so we
// can apply return each of the data paths requested.
// e.g. .version, .channel.title, .channel.link, .item[].link,
// .item[].guid, .item[].title, .item[].description
func (r *RSS2) Filter(dataPaths []string) (map[string]interface{}, error) {
	var (
		err  error
		data map[string]interface{}
	)
	result := make(map[string]interface{})
	for _, dataPath := range dataPaths {
		switch {
		case strings.Compare(dataPath, ".version") == 0:
			result["version"] = r.Version
		case strings.HasPrefix(dataPath, ".channel"):
			data, err = r.channel(dataPath)
			// Merge data into results keyed' by path
			for _, val := range data {
				result[dataPath] = val
			}
		case strings.HasPrefix(dataPath, ".item[]"):
			data, err = r.items(dataPath)
			// Merge data into results keyed' by path
			for _, val := range data {
				result[dataPath] = val
			}
		default:
			return nil, fmt.Errorf("path %q not found", dataPath)
		}
	}
	if result == nil {
		return nil, fmt.Errorf("No data paths found")
	}
	return result, err
}
