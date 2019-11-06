package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

type UploadReq struct {
	FileLocation string `form:"file_location" json:"file_location" binding:"required"`
	DownloadOnLoad string `form:"download_on_load" json:"download_on_load"`
	WebRoute string `form:"web_route" binding:"required" json:"web_route"`
	LocationFromRoot string `form:"location_from_root" json:"location_from_route"`
}

func UploadFile(c *gin.Context) {

	var form UploadReq

	if err := c.ShouldBind(&form); err != nil {
		c.HTML(http.StatusBadRequest, "500.html", gin.H{"error": err.Error()})
		log.Println("400: " + err.Error())
		return
	}

	f, err := c.FormFile("file")
	if err != nil {
		c.HTML(http.StatusBadRequest, "500.html", gin.H{"error": err.Error()})
		log.Println("400 file: " + err.Error())
		return
	}

	if form.LocationFromRoot != "on" {
		form.FileLocation = UploadsFolder + "/" + form.FileLocation
	}

	file := IFile{
		FileLocation:   form.FileLocation,
		DownloadOnLoad: form.DownloadOnLoad == "on",
		WebRoute:       form.WebRoute,
		DateCreated:    time.Now().Unix(),
	}

	if f != nil {
		err := c.SaveUploadedFile(f, form.FileLocation)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "500.html", gin.H{})
			_, _ = fmt.Fprintf(os.Stderr, err.Error() + "\n")
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
