package response

type HeaderPathRPS struct {
	HeaderPath string `json:"header_path"`
}

func ToDTHeaderPath(headerPath string) *HeaderPathRPS {
	return &HeaderPathRPS{
		HeaderPath: headerPath,
	}
}
