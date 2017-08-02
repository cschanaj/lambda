package Rules

type Ruleset struct {
	Name          string         `xml:"name,attr"`
	Platform      string         `xml:"platform,attr"`
	Default_off   string         `xml:"default_off,attr"`

	Targets       []Target       `xml:"target"`
	Exclusions    []Exclusion    `xml:"exclusion"`
	Rules         []Rule         `xml:"rule"`
	TestUrls      []TestUrl      `xml:"test"`
	SecureCookies []SecureCookie `xml:"securecookie"`
}

type Target struct {
	Host string `xml:"host,attr"`
}

type Exclusion struct {
	Pattern string `xml:"pattern,attr"`
}

type Rule struct {
	From string `xml:"from,attr"`
	To   string `xml:"to,attr"`
}

type TestUrl struct {
	Url  string `xml:"url,attr"`
}

type SecureCookie struct {
	Host string `xml:"host,attr"`
	Name string `xml:"name,attr"`
}
