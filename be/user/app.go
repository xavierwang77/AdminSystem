package user

type userLoginData struct {
	Name      string `json:"loginName"`
	Password  string `json:"password"`
	AdminName string `json:"adminName"`
}

type addForm struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PhoneNumber    string `json:"phoneNumber"`
	AdminName      string `json:"adminName"`
	AvatarFilename string `json:"avatarFilename"`
	Remarks        string `json:"remarks"`
}
type changeInfo struct {
	UserName  string `json:"userName"`
	AdminName string `json:"adminName"`
}
type user struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	PhoneNumber    string `json:"phoneNumber"`
	AvatarFilename string `json:"avatarFilename"`
	Remarks        string `json:"remarks"`
	LoginPrivilege string `json:"loginPrivilege"`
}

type userChangeData struct {
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

type userAvatarFilename struct {
	Name string `json:"avatarFilename"`
}

type adminName struct {
	Name string `json:"loginName"`
}

type deleteInfo struct {
	UserName  string `json:"userName"`
	AdminName string `json:"adminName"`
}

type showUser struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
	Remarks     string `json:"remarks"`
}
