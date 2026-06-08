package files

import "time"

type File struct {
	ID           string
	UserID       *string
	Bucket       string
	ObjectName   string
	OriginalName string
	ContentType  string
	SizeBytes    int64
	FileURL      *string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
