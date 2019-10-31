package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/shurcooL/httpfs/vfsutil"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var (
	router *gin.Engine
	config Config
	users []*User
)

const (
	RootFolder = "."
)

func main() {
	log.Println("Starting asset-hoster...")
	log.Println("This program comes with ABSOLUTELY NO WARRANTY;\nThis is free software, and you are welcome to redistribute it under certain conditions.")

	// write binary files to disk
	err := vfsutil.WalkFiles(assets, "/", func(path string, fi os.FileInfo, r io.ReadSeeker, err error) error {
		if err != nil {
			log.Fatal(err)
			return nil
		}
		log.Println(path)
		if fi.IsDir() {
			err = os.Mkdir(RootFolder + "/assets" + path, 0777)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			b, err := ioutil.ReadAll(r)

			if err != nil {
				log.Fatal(err)
			}

			err = ioutil.WriteFile(RootFolder + "/assets" + path, b, 0777)
			if err != nil {
				log.Fatal(err)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatal(err.Error())
	}

	// Init web-server
	router = gin.Default()
	setupRoutes()

	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(config.Port),
		Handler: router,
	}

	// start web-server in goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down asset-hoster...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown: ", err)
	}
	log.Println("asset-hoster has stopped.")
}

func setupRoutes() {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Static("/css", RootFolder + "/assets/css")
	router.Static("/js", RootFolder + "/assets/js")
	router.Static("/images", RootFolder + "/assets/images")

	router.LoadHTMLGlob(RootFolder+"/assets/html/*")
	router.NoRoute(func(c *gin.Context) {
		c.HTML(404, "404.html", gin.H{})
	})
}