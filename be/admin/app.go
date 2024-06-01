package admin

var avatarFilename string

type registerData struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PhoneNumber    string `json:"phoneNumber"`
	AvatarFilename string `json:"avatarFilename"`
	Remarks        string `json:"remarks"`
	adminTableName string
}
type loginData struct {
	Name     string `json:"loginName"`
	Password string `json:"password"`
}

type adminName struct {
	Name string `json:"loginName"`
}

type showUser struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
	Remarks     string `json:"remarks"`
}

type changeAdminData struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	PhoneNumber       string `json:"phoneNumber"`
	Remarks           string `json:"remarks"`
	OriginalAdminName string `json:"originalAdminName"`
	AvatarFilename    string `json:"adminAvatarName"`
}
