package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type UploadReq struct {
	FileLocation string `form:"file_location" binding:"required"`
	DownloadOnLoad bool `form:"download_on_load" binding:"required"`
	WebRoute string `form:"web_route" binding:"required"`
	LocationFromRoot bool `form:"location_from_root" binding:"required"`
	File *multipart.FileHeader `form:"file"`
}

func UploadFile(c *gin.Context) {

	var form UploadReq

	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "500.html", gin.H{"error": err.Error()})
		return
	}

	if !form.LocationFromRoot {
		form.FileLocation = UploadsFolder + "/" + form.FileLocation
	}

	file := IFile{
		FileLocation:   form.FileLocation,
		DownloadOnLoad: form.DownloadOnLoad,
		WebRoute:       form.WebRoute,
		DateCreated:    time.Now().Unix(),
	}

	if form.File == nil {
		err := c.SaveUploadedFile(form.File, form.FileLocation)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
			_, _ = fmt.Fprintf(os.Stderr, err.Error())
			return
		}
	}

	StoreFile(&file)

	c.Redirect(302, "/manage")

}

func ServeFile(c *gin.Context) {
	if routeCache[c.Param("route")] == nil {
		c.HTML(http.StatusNotFound, "404.html", gin.H{})
		return
	}

	file := *routeCache[c.Param("route")]

	// TODO file is not found
	f, err := os.Open(file.FileLocation)
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	contentType, err := GetFileContentType(f)
	if err != nil {
		log.Println(err.Error())
		c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
		return
	}

	//c.Header("Content-Transfer-Encoding", "binary")
	if file.DownloadOnLoad {
		c.Header("Content-Disposition", "attachment")
	} else {
		c.Header("Content-Disposition", "inline")
	}
	c.Header("Content-Type", contentType)
	c.File(file.FileLocation)
}
