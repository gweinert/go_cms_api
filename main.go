package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gweinert/cms_scratch/controllers"
	"github.com/gweinert/cms_scratch/models"
	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "user=Garrett dbname=cms_scratch sslmode=disable"
	}
	models.InitDB(dbURL)

	router := httprouter.New()

	router.GET("/site", controllers.ShowSiteDetailFunc)
	router.GET("/site/:siteID", controllers.GetPages)
	router.POST("/site/publish", controllers.BasicAuth(controllers.PublishSite))

	router.POST("/page/create", controllers.BasicAuth(controllers.CreatePage))
	router.POST("/page/update", controllers.BasicAuth(controllers.UpdatePage))
	router.POST("/page/delete", controllers.BasicAuth(controllers.DeletePage))
	router.POST("/page/sort-order", controllers.BasicAuth(controllers.UpdatePageSortOrder))

	router.GET("/site/:siteID/page/:pageID", controllers.GetElements)
	router.POST("/element/delete", controllers.BasicAuth(controllers.DeleteElements))

	router.POST("/group/create", controllers.BasicAuth(controllers.CreateNewGroup))
	router.POST("/group/delete", controllers.BasicAuth(controllers.DeleteGroup))

	router.POST("/image/upload", controllers.BasicAuth(controllers.UploadImage))
	router.POST("/image/delete", controllers.BasicAuth(controllers.DeleteImage))

	router.POST("/login", controllers.Login)
	router.POST("/user/session", controllers.GetUserFromSessionID)

	router.POST("/contact", controllers.SendContactMail)
	router.POST("/contact/price-quote", controllers.SendBookingMail)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "http://tiny.garrettdev.xyz", "http://thetinytattooshop.com/", "http://www.thetinytattooshop.com/"},
		AllowCredentials: true,
	})

	log.Fatal(http.ListenAndServe(":"+port, c.Handler(router)))

}
