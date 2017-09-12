package main

import (
	"log"
	"net/http"

	"github.com/gweinert/cms_scratch/controllers"
	"github.com/gweinert/cms_scratch/models"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {

	models.InitDB("user=Garrett dbname=cms_scratch sslmode=disable")

	router := httprouter.New()

	router.GET("/site", controllers.ShowSiteDetailFunc)
	router.GET("/site/:siteID", controllers.GetPages)
	router.POST("/site/publish", controllers.PublishSite)

	router.POST("/page/create", controllers.CreatePage)
	router.POST("/page/update", controllers.UpdatePage)
	router.POST("/page/delete", controllers.DeletePage)
	router.POST("/page/sort-order", controllers.UpdatePageSortOrder)

	router.GET("/site/:siteID/page/:pageID", controllers.GetElements)
	router.POST("/element/delete", controllers.DeleteElements)

	router.POST("/group/create", controllers.CreateNewGroup)
	router.POST("/group/delete", controllers.DeleteGroup)

	router.POST("/image/upload", controllers.UploadImage)
	router.POST("/image/delete", controllers.DeleteImage)

	router.POST("/login", controllers.Login)
	router.POST("/user/session", controllers.GetUserFromSessionID)
	// router.GET("/login/GUID", controllers.GetGUID)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	log.Fatal(http.ListenAndServe(":8080", c.Handler(router)))

}
