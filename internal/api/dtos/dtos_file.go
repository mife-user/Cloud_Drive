package dtos

type FileDtos struct {
	FileName    string `json:"file_name"`
	Size        int64  `json:"size"`
	Path        string `json:"path"`
	Permissions string `json:"permissions"`
	Owner       string `json:"owner"`
}
