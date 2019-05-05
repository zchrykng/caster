package caster

type Feed struct {
	URL      string
	Root     string
	Title    string
	Episodes []*Episode
}

func MakeFeed(URL string, Root string, Title string) (*Feed, error) {
	f := &Feed{URL: URL, Root: Root, Title: Title}

	return f, nil
}

func (f *Feed) ScanEpisodes() error {

	return nil
}
