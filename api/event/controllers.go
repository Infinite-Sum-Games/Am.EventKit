package api

import (
	"net/http"

	"github.com/Thanus-Kumaar/anokha-2025-backend/pkg"
	"github.com/gin-gonic/gin"
)

func FetchEventCatalog(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Retrived event catalog successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func FetchEventById(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Retrived event successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func FetchRegisteredEvents(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Retrived registered events successfully",
	})
	pkg.Log.SuccessCtx(c)
}

func ToggleFavouriteEvent(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Toggled favourite event successfully",
	})
	pkg.Log.SuccessCtx(c)
}
