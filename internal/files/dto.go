package files

type FileResponse struct {
	ID           string  `json:"id"`
	UserID       *string `json:"user_id,omitempty"`
	Bucket       string  `json:"bucket"`
	ObjectName   string  `json:"object_name"`
	OriginalName string  `json:"original_name"`
	ContentType  string  `json:"content_type"`
	SizeBytes    int64   `json:"size_bytes"`
	FileURL      *string `json:"file_url,omitempty"`
}
