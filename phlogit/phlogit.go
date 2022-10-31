// phlogit.go is a sub-package pttk. A packages for managing static content
// phlogs and documentation via Pandoc.
//
// @Author R. S. Doiel, <rsdoiel@gmail.com>
//
// copyright 2022 R. S. Doiel
// All rights reserved.
//
// License under the 3-Clause BSD License
// See https://opensource.org/licenses/BSD-3-Clause
package phlogit

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	// My packages
	"github.com/rsdoiel/pttk/frontmatter"
	stn "github.com/rsdoiel/stngo"
	"github.com/rsdoiel/stngo/report"

	// 3rd Party Packages
	"gopkg.in/yaml.v3"
)

// by start '%' at the being of the line in the start of a text file.
type MetadataBlock struct {
	Title   string   `json:"title"`
	Authors []string `json:"authors"`
	Date    string   `json:"date"`
}

func (block *MetadataBlock) String() string {
	return fmt.Sprintf("%% %s\n%% %s\n%% %s", block.Title, strings.Join(block.Authors, "; "), block.Date)
}

func (block *MetadataBlock) Marshal() ([]byte, error) {
	return json.Marshal(block)
}

// MetadataBlock holds the Pandoc style Metadata block delimited
func (block *MetadataBlock) Unmarshal(src []byte) error {
	lines := bytes.Split(src, []byte("\n"))
	fieldCnt := 0
	key := ""
	block.Title = ""
	block.Authors = []string{}
	block.Date = ""
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if bytes.HasPrefix(line, []byte("% ")) {
			fieldCnt += 1
			switch fieldCnt {
			case 1:
				key = "title"
			case 2:
				key = "authors"
			case 3:
				key = "date"
			default:
				key = ""
			}
			line = bytes.TrimPrefix(line, []byte("% "))
		}
		if (len(key) > 0) && (fieldCnt <= 3) {
			switch key {
			case "title":
				if len(block.Title) > 0 {
					block.Title = fmt.Sprintf("%s\n%s", block.Title, bytes.TrimSpace(line))
				} else {
					block.Title = fmt.Sprintf("%s", bytes.TrimSpace(line))
				}
			case "authors":
				if bytes.Contains(line, []byte(";")) {
					parts := bytes.Split(line, []byte(";"))
					for _, part := range parts {
						block.Authors = append(block.Authors, fmt.Sprintf("%s", bytes.TrimSpace(part)))
					}
				} else {
					block.Authors = append(block.Authors, fmt.Sprintf("%s", bytes.TrimSpace(line)))
				}
			case "date":
				block.Date = fmt.Sprintf("%s", bytes.TrimSpace(line))
				key = ""
			}
		}
	}
	if fieldCnt != 3 {
		return fmt.Errorf("Missing or ill formed metablock, expecting title, author(s), date")
	}
	return nil
}

const (
	DateFmt = "2006-01-02"
)

type CreatorObj struct {
	ORCID string `json:"orcid,omitempty" yaml:"orcid,omitempty"`
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
}

type PostObj struct {
	Slug        string       `json:"slug" yaml:"slug"`
	Document    string       `json:"document" yaml:"document"`
	Title       string       `json:"title,omitempty" yaml:"title,omitempty"`
	SubTitle    string       `json:"subtitle,omitempty" yaml:"subtitle,omitempty"`
	Author      string       `json:"author,omitempty" yaml:"author,omitempty"`
	Byline      string       `json:"byline,omitempty" yaml:"byline,omitempty"`
	Series      string       `json:"series,omitempty" yaml"series,omitempty"`
	Number      string       `json:"number,omitempty" yaml:"number,omitempty"`
	Subject     string       `json:"subject,omitempty" yaml:"subject,omitempty"`
	Keywords    []string     `json:"keywords,omitempty" yaml:"keywords,omitempty"`
	Abstract    string       `json:"abstract,omitempty" yaml:"abstract,omitempty"`
	Description string       `json:"description,omitempty" yaml:"description,omitempty"`
	Category    string       `json:"category,omitempty" yaml:"catagory,omitempty"`
	Lang        string       `json:"lang,omitempty" yaml:"lang,omitempty"`
	Direction   string       `json:"direction,omitempty" yaml:"direction,omitempty"`
	Draft       bool         `json:"draft,omitempty" yaml:"draft,omitempty"`
	Creators    []CreatorObj `json:"creators,omitempty" yaml:"creators,omitempty"`
	Created     string       `json:"date,omitempty" yaml:"date,omitempty"`
	Updated     string       `json:"updated,omitempty" yaml:"updated,omitempty"`
}

type DayObj struct {
	Day   string     `json:"day" yaml:"day"`
	Posts []*PostObj `json:"posts" yaml:"post"`
}

type MonthObj struct {
	Month string    `json:"month" yaml:"month"`
	Days  []*DayObj `json:"days" yaml:"days"`
}

type YearObj struct {
	Year   string      `json:"year" yaml:"year"`
	Months []*MonthObj `json:"months" yaml:"months"`
}

type PhlogMeta struct {
	Name        string     `json:"name,omitempty" yaml:"name,omitempty"`
	Quip        string     `json:"quip,omitempty" yaml:"quip,omitempty"`
	Description string     `json:"description,omitempty" yaml:"description,omitempty"`
	BaseURL     string     `json:"url,omitempty" yaml:"url,omitempty"`
	Copyright   string     `json:"copyright,omitempty" yaml:"copyright,omitempty"`
	License     string     `json:"license,omitempty" yaml:"license,omitempty"`
	Language    string     `json:"language,omitempty" yaml:"language,omitempty"`
	Started     string     `json:"started,omitempty" yaml:"started,omitempty"`
	Ended       string     `json:"ended,omitempty" yaml:"ended,omitempty"`
	Updated     string     `json:"updated,omitempty" yaml:"updated,omitempty"`
	IndexTmpl   string     `json:"index_tmpl,omitempty" yaml:"index_tmpl,omitempty"`
	PostTmpl    string     `json:"post_tmpl,omitempty" yaml:"post_tmpl,omitempty"`
	Years       []*YearObj `json:"years" yaml:"years"`
}

//
// Support funcs
//

func calcYMD(dateString string) ([]string, error) {
	ymd := strings.SplitN(dateString, "-", 3)
	return ymd, nil
}

func calcPath(prefix string, ymd []string) (string, error) {
	if len(ymd) != 3 {
		return "", fmt.Errorf("Invalid Year, Month and Date")
	}
	dPath := path.Join(prefix, ymd[0], ymd[1], ymd[2])
	return dPath, nil
}

func unpackCreators(objects []interface{}) []CreatorObj {
	creators := []CreatorObj{}
	for _, obj := range objects {
		switch obj.(type) {
		case string:
			creator := CreatorObj{}
			creator.Name = obj.(string)
			creators = append(creators, creator)
		case map[string]string:
			m := obj.(map[string]string)
			creator := CreatorObj{}
			if name, ok := m["name"]; ok {
				creator.Name = name
			}
			if orcid, ok := m["orcid"]; ok {
				creator.ORCID = orcid
			}
			creators = append(creators, creator)
		}
	}
	return creators
}

//
// Exported funcs
//

// postIndex looks through the day's posts and find
// the position that matches the slug or returns -1
func (d DayObj) postIndex(slug string) int {
	postIndex := -1
	for i, obj := range d.Posts {
		if obj.Slug == slug {
			postIndex = i
			break
		}
	}
	return postIndex
}

// dayIndex looks through the months days and returns
// index position found or -1 if not found.
func (m MonthObj) dayIndex(day string) int {
	dayIndex := -1
	for i, obj := range m.Days {
		if obj.Day == day {
			dayIndex = i
			break
		}
	}
	return dayIndex
}

// monthIndex looks through a years' months and returns
// the index position found or -1 if not found.
func (y YearObj) monthIndex(month string) int {
	monthIndex := -1
	for i, obj := range y.Months {
		if obj.Month == month {
			monthIndex = i
			break
		}
	}
	return monthIndex
}

// yearIndex
func (meta PhlogMeta) yearIndex(year string) int {
	yearIndex := -1
	for i, obj := range meta.Years {
		if obj.Year == year {
			yearIndex = i
			break
		}
	}
	return yearIndex
}

// updatePosts will create a new post if necessary and insert in to the
// post list.
func (dy *DayObj) updatePosts(ymd []string, targetName string) error {

	// Read in front matter from targetName
	src, err := os.ReadFile(targetName)
	if err != nil {
		return fmt.Errorf("Failed to read post %q, %s", targetName, err)
	}
	obj := map[string]interface{}{}
	src, _ = frontmatter.ReadAll(bytes.NewBuffer(src))
	if len(src) > 0 {
		if err := json.Unmarshal(src, &obj); err != nil {
			return fmt.Errorf("Failed to unmarshal frontmatter %q, %s", targetName, err)
		}
	}
	// Create a new PostObj
	today := time.Now().Format(DateFmt)
	created := strings.Join(ymd, "-")
	post := new(PostObj)
	post.Document = targetName
	post.Updated = today
	post.Created = created
	post.Slug = strings.TrimSuffix(path.Base(targetName), filepath.Ext(targetName))
	if title, ok := obj["title"]; ok {
		post.Title = title.(string)
	}
	if subtitle, ok := obj["subtitle"]; ok {
		post.SubTitle = subtitle.(string)
	}
	if byline, ok := obj["byline"]; ok {
		post.Byline = byline.(string)
	}
	if author, ok := obj["author"]; ok {
		post.Author = author.(string)
	}
	if series, ok := obj["series"]; ok {
		post.Series = series.(string)
	}
	if number, ok := obj["number"]; ok {
		switch number.(type) {
		case int:
			post.Number = fmt.Sprintf("%d", number.(int))
		case string:
			post.Number = number.(string)
		}
	}
	if subject, ok := obj["subject"]; ok {
		post.Subject = subject.(string)
	}
	if objects, ok := obj["keywords"]; ok {
		switch objects.(type) {
		case []string:
			keywords := objects.([]string)
			for _, kw := range keywords {
				post.Keywords = append(post.Keywords, kw)
			}
		}
	}
	if abstract, ok := obj["abstract"]; ok {
		post.Abstract = abstract.(string)
	}
	if description, ok := obj["description"]; ok {
		post.Description = description.(string)
	}
	if category, ok := obj["category"]; ok {
		post.Category = category.(string)
	}
	if lang, ok := obj["lang"]; ok {
		post.Lang = lang.(string)
	}
	if direction, ok := obj["direction"]; ok {
		post.Direction = direction.(string)
	}
	if draft, ok := obj["draft"]; ok {
		post.Draft = draft.(bool)
	} else {
		post.Draft = false
	}
	if creators, ok := obj["creators"]; ok {
		post.Creators = unpackCreators(creators.([]interface{}))
	}
	if dt, ok := obj["date"]; ok {
		switch dt.(type) {
		case string:
			post.Created = dt.(string)
		case time.Time:
			t := dt.(time.Time)
			post.Created = t.Format(DateFmt)
		}
	}
	if updated, ok := obj["updated"]; ok {
		switch updated.(type) {
		case string:
			post.Updated = updated.(string)
		case time.Time:
			t := updated.(time.Time)
			post.Updated = t.Format(DateFmt)
		}
	}

	i := dy.postIndex(post.Slug)
	if i < 0 {
		// Add a post
		post.Created = today
		posts := dy.Posts[0:]
		dy.Posts = append([]*PostObj{post}, posts...)
	} else {
		// Update a post
		dy.Posts[i] = post
	}
	return nil
}

// updateDays will create a new day and insert in order
// before passing the post data to UpdatePost()
func (mn *MonthObj) updateDays(ymd []string, targetName string) error {
	dy := new(DayObj)
	dy.Day = ymd[2]
	i := mn.dayIndex(dy.Day)
	if i < 0 {
		if len(mn.Days) == 0 {
			mn.Days = append(mn.Days, dy)
			i = 0
		} else {
			for j, obj := range mn.Days {
				if dy.Day > obj.Day {
					i = j
					if i == 0 {
						mn.Days = append([]*DayObj{dy}, mn.Days...)
					} else {
						days := mn.Days[0 : j-1]
						days = append(days, dy)
						mn.Days = append(days, mn.Days[j:]...)
					}
					break
				}
			}
			if i < 0 {
				mn.Days = append(mn.Days, dy)
				i = len(mn.Days) - 1
			}
		}
	}
	return mn.Days[i].updatePosts(ymd, targetName)
}

// updateMonths will create/update month
// before passing the post data to UpdateDays()
func (yr *YearObj) updateMonths(ymd []string, targetName string) error {
	mn := new(MonthObj)
	mn.Month = ymd[1]
	i := yr.monthIndex(mn.Month)
	if i < 0 {
		if len(yr.Months) == 0 {
			yr.Months = append(yr.Months, mn)
			i = 0
		} else {
			for j, obj := range yr.Months {
				if mn.Month > obj.Month {
					i = j
					if i == 0 {
						yr.Months = append([]*MonthObj{mn}, yr.Months...)
					} else {
						months := yr.Months[0 : j-1]
						months = append(months, mn)
						yr.Months = append(months, yr.Months[j:]...)
					}
					break
				}
			}
			if i < 0 {
				yr.Months = append(yr.Months, mn)
				i = len(yr.Months) - 1
			}
		}
	}
	return yr.Months[i].updateDays(ymd, targetName)
}

// updateYears will create/update year in `meta.Years`
// before passing the post data to UpdateMonths()
func (meta *PhlogMeta) updateYears(ymd []string, targetName string) error {
	yr := new(YearObj)
	yr.Year = ymd[0]
	i := meta.yearIndex(yr.Year)
	if i < 0 {
		if len(meta.Years) == 0 {
			meta.Years = append(meta.Years, yr)
			i = 0
		} else {
			for j, obj := range meta.Years {
				if yr.Year > obj.Year {
					i = j
					if i == 0 {
						meta.Years = append([]*YearObj{yr}, meta.Years...)
					} else {
						years := meta.Years[0 : j-1]
						years = append(years, yr)
						meta.Years = append(years, meta.Years[j:]...)
					}
					break
				}
			}
			if i < 0 {
				// We need to append the year, it's earlier than
				// known years.
				meta.Years = append(meta.Years, yr)
				i = len(meta.Years) - 1
			}
		}
	}
	return meta.Years[i].updateMonths(ymd, targetName)
}

// PhlogAsset copies a asset file to the directory as a phlog post
// calculated from prefix and dataString. It is used by PhlogIt to
// copy the Markdown document to the right path.
func (meta *PhlogMeta) PhlogAsset(prefix string, fName string, dateString string) error {
	var (
		targetName string
	)
	// Check to see if dateStr is empty, if so default to today.
	if dateString == "" {
		dateString = time.Now().Format("2015-05-07")
	}
	// Check to see if path.join(prefix, datePath(dateStr)) exists
	// and create it if needed.
	ymd, err := calcYMD(dateString)
	if err != nil {
		return err
	}
	dPath, err := calcPath(prefix, ymd)
	if err != nil {
		return err
	}
	// copy fName to target path.
	var (
		in, out *os.File
	)
	in, err = os.Open(fName)
	if err != nil {
		return err
	} else {
		os.MkdirAll(dPath, 0775)
		targetName = path.Join(dPath, path.Base(fName))
		out, err = os.Create(targetName)
		if err != nil {
			return fmt.Errorf("Creating %q, %s", targetName, err)
		}
		if _, err := io.Copy(out, in); err != nil {
			return err
		}
		in.Close()
		out.Close()
	}
	return nil
}

// PhlogSTN is a tool for posting and updating a phlog directory
// structure based on the contents of an [stn](https://rsdoiel.github.io/stngo) (Simple Timesheet Notation). Entries are mapped to phlog
// posts to populate the phlog.
func (meta *PhlogMeta) PhlogSTN(prefix string, fName string, author string) error {
	in, err := os.Open(fName)
	if err != nil {
		return err
	}
	defer in.Close()
	scanner := bufio.NewScanner(in)

	entry := new(stn.Entry)
	aggregation := new(report.EntryAggregation)
	activeDate := time.Now().Format(DateFmt)

	lineNo := 0

	for scanner.Scan() {
		line := scanner.Text()
		lineNo++
		if stn.IsDateLine(line) == true {
			activeDate = stn.ParseDateLine(line)
		} else if stn.IsEntry(line) {
			entry, err = stn.ParseEntry(activeDate, line)
			if err != nil {
				return fmt.Errorf("line %5d: can't parse entry %q\n", lineNo, line)
			}
			aggregation.Aggregate(entry)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	dayFmt := `Monday, January 2, 2006`
	timeFmt := `3:04 PM`
	//hourMinFmt := `15:03`
	postFmt := `Monday, January 2, 2006 15:03 MST`
	fileExtFmt := `2006-01-02_150304`
	ext := path.Ext(fName)
	// postFName (based on STN provided filename) but written
	// to the appropriate date path but without an extension.
	// The end date and time of entry will complete the name.
	postFName := strings.TrimSuffix(path.Base(fName), ext)
	tot := len(aggregation.Entries)
	for i, entry := range aggregation.Entries {
		// entry metadata
		series := ""
		keywords := []string{}
		if len(entry.Annotations) > 0 {
			series = entry.Annotations[0]
		}
		if len(entry.Annotations) > 1 {
			if strings.Contains(entry.Annotations[1], ",") {
				for _, word := range strings.Split(entry.Annotations[1], ",") {
					keywords = append(keywords, word)
				}
			} else {
				keywords = append(keywords, entry.Annotations[1])
			}
		}
		title := fmt.Sprintf("%s", entry.End.Format(timeFmt))
		if series != "" {
			title = fmt.Sprintf("%s, %s", entry.End.Format(timeFmt), series)
		}
		if len(keywords) > 0 {
			title = fmt.Sprintf("%s: %s", title, strings.Join(keywords, ", "))
		}
		issueNo := tot - i
		pubDate := entry.End.Format(DateFmt)
		fileExtDate := entry.End.Format(fileExtFmt)
		// Make sure we have a path to write the file
		ymd, err := calcYMD(pubDate)
		if err != nil {
			return fmt.Errorf("entry %5d: %q %s\n", i, pubDate, err)
		}
		dPath, err := calcPath(prefix, ymd)
		if err != nil {
			return fmt.Errorf("entry %5d: %q %s\n", i, pubDate, err)
		}
		os.MkdirAll(dPath, 0775)
		targetName := path.Join(dPath, fmt.Sprintf("%s-%s.md", postFName, fileExtDate))
		out, err := os.Create(targetName)
		if err != nil {
			return fmt.Errorf("entry %5d: %q %s\n", i, targetName, err)
		}
		// Write front matter in YAML to the file
		fmt.Fprintf(out, "---\n")
		fmt.Fprintf(out, "title: %q\n", title)
		if author != "" {
			fmt.Fprintf(out, "author: %q\n", author)
		}
		fmt.Fprintf(out, "pubDate: %s\n", pubDate)
		if series != "" {
			fmt.Fprintf(out, "series: %q\n", series)
		}
		fmt.Fprintf(out, "no: %d\n", issueNo)
		if len(keywords) > 0 {
			fmt.Fprintf(out, "keywords:\n")
			for _, word := range keywords {
				fmt.Fprintf(out, "  - %q\n", strings.TrimSpace(word))
			}
		}
		fmt.Fprintf(out, "---\n\n")
		// Write the body content to the file
		fmt.Fprintf(out, "# %s\n\n", title)
		if author != "" {
			// Add a by line
			fmt.Fprintf(out, "By %s, %s\n\n", author, entry.End.Format(postFmt))
		} else {
			fmt.Fprintf(out, "Post: %s, %s\n\n", entry.End.Format(dayFmt), entry.End.Format(timeFmt))
		}
		l := len(entry.Annotations)
		switch {
		case l > 2:
			fmt.Fprintf(out, "%s\n\n", strings.Join(entry.Annotations[2:], "\n"))
		case l > 1:
			fmt.Fprintf(out, "%s\n\n", strings.Join(entry.Annotations[1:], "\n"))
		default:
			fmt.Fprintf(out, "%s\n\n", strings.Join(entry.Annotations, "\n"))
		}
		out.Close()
		// Refresh the phlog betatadata structure as needed.
		meta.Updated = time.Now().Format(DateFmt)
		if err := meta.updateYears(ymd, targetName); err != nil {
			return fmt.Errorf("%q %s, %s\n", postFName, entry.Start.Format(DateFmt), err)
		}
	}
	return nil
}

// PhlogIt is a tool for posting and updating a phlog directory
// structure on local disc.  It includes maintaining additional
// metadata resources to make it easy to script phlogsites and
// podcast sites.
// @param prefix - a prefix path for the phlog (relative to working dir)
// @param fName - the name of the file to publish to generated phlog path
// @param dateString - A date string used to calculate the phlog path, e.g.
//
//	YYYY-MM-DD maps to YYYY/MM/DD.
//
// @returns an error type
func (meta *PhlogMeta) PhlogIt(prefix string, fName string, dateString string) error {
	var (
		targetName string
	)
	// Check to see if dateStr is empty, if so default to today.
	if dateString == "" {
		dateString = time.Now().Format(DateFmt)
	}
	// Check to see if path.join(prefix, datePath(dateStr)) exists
	// and create it if needed.
	ymd, err := calcYMD(dateString)
	if err != nil {
		return err
	}
	dPath, err := calcPath(prefix, ymd)
	if err != nil {
		return err
	}
	// copy fName to target path.
	var (
		in, out *os.File
	)
	in, err = os.Open(fName)
	if err != nil {
		return err
	} else {
		os.MkdirAll(dPath, 0775)
		targetName = path.Join(dPath, path.Base(fName))
		out, err = os.Create(targetName)
		if err != nil {
			return fmt.Errorf("Creating %q, %s", targetName, err)
		}
		if _, err := io.Copy(out, in); err != nil {
			return err
		}
		in.Close()
		out.Close()
	}
	// NOTE: Updated is always today.
	meta.Updated = time.Now().Format(DateFmt)
	return meta.updateYears(ymd, targetName)
}

// Save writes a JSON phlog meta document
func (meta *PhlogMeta) Save(fName string) error {
	var (
		src []byte
		err error
	)
	ext := path.Ext(fName)
	switch ext {
	case ".json":
		src, err = json.MarshalIndent(meta, "", "    ")
		if err != nil {
			return fmt.Errorf("Marshaling %q, %s", fName, err)
		}
	case ".yml":
		src, err = yaml.Marshal(meta)
		if err != nil {
			return fmt.Errorf("Marshaling %q, %s", fName, err)
		}
	case ".yaml":
		src, err = yaml.Marshal(meta)
		if err != nil {
			return fmt.Errorf("Marshaling %q, %s", fName, err)
		}
	default:
		return fmt.Errorf("%q is an unsupported output format", ext)
	}
	err = os.WriteFile(fName, src, 0666)
	if err != nil {
		return fmt.Errorf("Writing %q, %s", fName, err)
	}
	return nil
}

// Reads a JSON phlog meta document and popualtes a phlog meta structure
func LoadPhlogMeta(fName string, meta *PhlogMeta) error {
	src, err := os.ReadFile(fName)
	if err != nil {
		return fmt.Errorf("Reading %q, %s", fName, err)
	}
	if len(src) > 0 {
		ext := path.Ext(fName)
		switch ext {
		case ".json":
			if err := json.Unmarshal(src, meta); err != nil {
				return fmt.Errorf("Unmarshing %q, %s", fName, err)
			}
		case ".yml":
			if err := yaml.Unmarshal(src, meta); err != nil {
				return fmt.Errorf("Unmarshing %q, %s", fName, err)
			}
		case ".yaml":
			if err := yaml.Unmarshal(src, meta); err != nil {
				return fmt.Errorf("Unmarshing %q, %s", fName, err)
			}
		default:
			return fmt.Errorf("%q is an unsupported input format", ext)
		}
	}
	return nil
}

// hasExt checks if ext is in list of target ext.
func hasExt(ext string, targetExts []string) bool {
	for _, e := range targetExts {
		if strings.Compare(ext, e) == 0 {
			return true
		}
	}
	return false
}

// monthName
func monthName(month string) string {
	switch month {
	case "01":
		return "January"
	case "02":
		return "February"
	case "03":
		return "March"
	case "04":
		return "April"
	case "05":
		return "May"
	case "06":
		return "June"
	case "07":
		return "July"
	case "08":
		return "August"
	case "09":
		return "September"
	case "10":
		return "October"
	case "11":
		return "November"
	case "12":
		return "December"
	}
	return month
}

// RefreshFromPath crawls the dircetory tree and rebuilds
// the `phlog.json` file based on what is found. It takes a
// File extension to target (e.g. .md for Markdown) and
// analyzes the path for YYYY/MM/DD and transforms the
// information found into an entry in `phlog.json`.
func (meta *PhlogMeta) RefreshFromPath(prefix string, year string) error {
	var (
		ymd []string
	)
	absPrefix, err := filepath.Abs(prefix)
	if err != nil {
		absPrefix = "./" + prefix
	}
	targetExts := []string{
		".md",
		".rst",
		".textile",
		".jira",
		".txt",
	}
	months := map[string]int{
		"01": 31, "02": 29, "03": 31, "04": 30,
		"05": 31, "06": 30, "07": 31, "08": 31,
		"09": 30, "10": 31, "11": 30, "12": 31,
	}
	gophermapName := path.Join(prefix, year, "gophermap")
	gophermap := []string{fmt.Sprintf(`
 Pages for %s
 ==============
`, year),
	}
	ymd = append(ymd, year, "", "")
	for month, cnt := range months {
		ymd[1] = month
		entries := []string{}
		for day := 1; day <= cnt; day++ {
			ymd[2] = fmt.Sprintf("%02d", day)
			// CalcPath and find files.
			folder := path.Join(absPrefix, ymd[0], ymd[1], ymd[2])
			// Scan the fold for files ending in ext,
			files, err := os.ReadDir(folder)
			if err == nil {
				// for each file with matching extension run updateYear(ymd, targetName)
				for _, file := range files {
					targetName := path.Join(absPrefix, ymd[0], ymd[1], ymd[2], file.Name())
					ext := filepath.Ext(targetName)
					if hasExt(ext, targetExts) {
						// FIXME: Should pull the title from the file.
						fName := path.Base(targetName)
						entries = append(entries, fmt.Sprintf("0%s\t%s/%s/%s\n", fName, ymd[1], ymd[2], fName))
						if err := meta.updateYears(ymd, targetName); err != nil {
							return err
						}
					}
				}
			}
		}
		if len(entries) > 0 {
			monthEntry := fmt.Sprintf("1%s\t%s\n", monthName(month), month)
			gophermap = append(gophermap, monthEntry)
			monthMapName := path.Join(prefix, year, month, "gophermap")
			monthMap := bytes.NewBuffer([]byte{})
			monthMap.Write([]byte(fmt.Sprintf("%s\r\n", monthName(month))))
			for _, entry := range entries {
				parts := strings.Split(entry, "\t")
				if len(parts) == 2 {
					monthEntry := fmt.Sprintf("%s\t%s\r\n", parts[0], strings.TrimPrefix(parts[1], month+"/"))
					monthMap.Write([]byte(monthEntry))
					gophermap = append(gophermap, entry)
				}
			}
			err := os.WriteFile(monthMapName, monthMap.Bytes(), 0664)
			if err != nil {
				return err
			}

		}
	}
	src := []byte(strings.ReplaceAll(strings.Join(gophermap, "\n"), "\n", "\r\n"))
	return os.WriteFile(gophermapName, src, 0664)
}
