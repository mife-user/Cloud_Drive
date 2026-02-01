package domain

type UserRepo interface {
	Register(user *User) error
	Logon(user *User) error
}

type FileRepo interface {
	UploadFile(file *File) error
	ViewFile(userID string) ([]File, error)
}
