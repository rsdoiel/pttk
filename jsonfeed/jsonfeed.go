package jsonfeed

type Feed struct {
	Version     string    `json:"version,required"`
	Title       string    `json:"title,required"`
	HomePageURL string    `json:"home_page_url,omitempty"`
	FeedURL     string    `json:"feed_url,omitempty"`
	Description string    `json:"description,omitempty"`
	UserComment string    `json:"user_comment,omitempty"`
	NextURL     string    `json:"next_url,omitempty"`
	Icon        string    `json:"icon,omitempty"`
	Favicon     string    `json:"favicon,omitempty"`
	Authors     []*Author `json:"authors,omitempty"`
	Language    string    `json:"language,omitempty"`
	Expired     string    `json:"expired,omitempty"`
	Hubs        string    `json:"hubs,omitempty"`
}

type Author struct {
	Name   string `json:"name,omitempty"`
	URL    string `json:"url,omitempty"`
	Avatar string `json:"avatar,omitempty"`
}

type Items struct {
	ID            string        `json:"id,required"`
	URL           string        `json:"url,omitempty"`
	ExternalURL   string        `json:"external_url,omitempty"`
	Title         string        `json:"title,omitempty"`
	ContentHTML   string        `json:"content_html,omitempty"`
	ContentText   string        `json:"content_text,omitempty"`
	Image         string        `json:"image,omitempty"`
	BannerImage   string        `json:"banner_image,omitempty"`
	DatePublished string        `json:"date_published,omitempty"`
	DateModified  string        `json:"date_modified,omitempty"`
	Authors       []*Author     `json:"authors,omitemtpy"`
	Tags          []string      `json:"tags,omitempty"`
	Language      string        `json:"language,omitempty"`
	Attachments   []*Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	URL               string `json:"url,required"`
	MimeType          string `json:"mime_type,required"`
	Title             string `json:"title,omitempty"`
	SizeInBytes       int    `json:"size_in_bytes,omitempty"`
	DurationInSeconds int    `json:"duration_in_seconds,omitempty"`
}
