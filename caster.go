package caster

type Caster struct {
	URL   string
	Root  string
	Feeds []*Feed
}

func MakeCaster(URL string, Root string) (*Caster, error) {
	c := &Caster{URL: URL, Root: Root}

	return c, nil
}

func (c *Caster) ScanFeeds() error {

	return nil
}
