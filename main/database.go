package main

var routeCache map[string]*IFile

type IDatabase struct {
	Users []*IUser
	Files []*IFile
}

type IUser struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type IFile struct {
	FileLocation string `json:"file_location"`
	DownloadOnLoad bool `json:"download_on_load"`
	WebRoute string `json:"web_route"`
	DateCreated int64 `json:"date_created"`
}