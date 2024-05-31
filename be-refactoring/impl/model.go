package impl

// 解析请求体中的JSON数据
type REGISTERDATA struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PhoneNumber    string `json:"phoneNumber"`
	AvatarFilename string `json:"avatarFilename"`
	Remarks        string `json:"remarks"`
	adminTableName string
}
type LOGINDATA struct {
	Name     string `json:"loginName"`
	Password string `json:"password"`
}
type USERLOGINDATA struct {
	Name      string `json:"loginName"`
	Password  string `json:"password"`
	AdminName string `json:"adminName"`
}
type ADMINNAME struct {
	Name string `json:"loginName"`
}
type DELETEINFO struct {
	UserName  string `json:"userName"`
	AdminName string `json:"adminName"`
}
type ADDFORM struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PhoneNumber    string `json:"phoneNumber"`
	AdminName      string `json:"adminName"`
	AvatarFilename string `json:"avatarFilename"`
	Remarks        string `json:"remarks"`
}
type CHANGEINFO struct {
	UserName  string `json:"userName"`
	AdminName string `json:"adminName"`
}
type USER struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PhoneNumber    string `json:"phoneNumber"`
	AvatarFilename string `json:"avatarFilename"`
	Remarks        string `json:"remarks"`
	LoginPrivilege string `json:"loginPrivilege"`
}
type SHOWUSER struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
	Remarks     string `json:"remarks"`
}
type CHANGEDATA struct {
	ID               string
	Name             string `json:"name"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	PhoneNumber      string `json:"phoneNumber"`
	AdminName        string `json:"adminName"`
	OriginalUserName string `json:"originalUserName"`
	AvatarFilename   string `json:"userAvatarName"`
	Remarks          string `json:"remarks"`
	LoginPrivilege   string `json:"loginPrivilege"`
}
type CHANGEADMINDATA struct {
	Name              string `json:"name"`
	Email             string `json:"email"`
	Password          string `json:"password"`
	PhoneNumber       string `json:"phoneNumber"`
	Remarks           string `json:"remarks"`
	OriginalAdminName string `json:"originalAdminName"`
	AvatarFilename    string `json:"adminAvatarName"`
}
type USERAVATARFILENAME struct {
	Name string `json:"avatarFilename"`
}
