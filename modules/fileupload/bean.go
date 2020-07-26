package fileupload

type UploadFileRequest struct {
	FileName string `json:"file_name,omitempty"`
	FileSize uint   `json:"file_size,omitempty"`
	Md5      string `json:"md5,omitempty"`
	Tag      string `json:"tag,omitempty"`
}

type UploadFileResponse struct {
	UploadUrl   string `json:"upload_url,omitempty"`
	DownloadUrl string `json:"download_url,omitempty"`
	Exp         string `json:"expired_at,omitempty"`
}
