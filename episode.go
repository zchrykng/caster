package caster

import (
	"time"
	
	"github.com/vansante/go-ffprobe"
	"github.com/dhowden/tag"
)

type Episode struct {
	Title string
	File  *File
}

type File struct {
	Slug         string
	Location     string
	Name         string
	Artist       string
	ModifiedTime time.Time
	Duration     time.Duration
	Type         string
	Size         int64
}

// Even though ctx will be expired, it is good practice to call its
// cancelation function in any case. Failure to do so may keep the
// context and its parent alive longer than necessary.
defer cancel()

select {
case <-time.After(1 * time.Second):
    fmt.Println("overslept")
case <-ctx.Done():
    fmt.Println(ctx.Err())
}

func MakeEpisode(location string) (*Episode, error) {
	e := Episode{}
	
	e.File := File{}
	
	f, err := os.Open(location)
	if err != nil {
		fmt.Printf("error loading file: %v", err)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("error reading file info: %v", err)
	}
	
	d := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	
	data, err := GetProbeDataContext(ctx, location)
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		return
	}
	
	e.Slug, _ = Slugify(m.Title())
	e.Location = location
	e.Name = m.Title()
	e.Artist = m.Artist()
	e.ModifiedTime = fi.ModTime()
	e.Type = m.FileType()
	e.Size = fi.Size()
	

	



	
}