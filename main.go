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

	router.POST("/page/create", controllers.CreatePage)
	router.POST("/page/update", controllers.UpdatePage)
	router.POST("/page/delete", controllers.DeletePage)

	router.GET("/site/:siteID/page/:pageID", controllers.GetElements)
	router.POST("/element/delete", controllers.DeleteElement)

	router.POST("/group/create", controllers.CreateNewGroup)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowCredentials: true,
	})

	log.Fatal(http.ListenAndServe(":8080", c.Handler(router)))

}
