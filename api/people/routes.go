package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func PeopleRoutes(r *gin.RouterGroup) {
	r.GET("/", mw.Auth, FetchAllPeople)
	r.POST("/", mw.Auth, mw.CheckAdmin, AddNewPerson)
	r.PUT("/:id", mw.Auth, mw.CheckAdmin, UpdatePersonDetails)
	r.DELETE("/:id", mw.Auth, mw.CheckAdmin, DeletePerson)
}
