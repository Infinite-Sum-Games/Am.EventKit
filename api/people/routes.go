package api

import (
	mw "github.com/Infinite-Sum-Games/Am.EventKit/middleware"
	"github.com/gin-gonic/gin"
)

func PeopleRoutes(r *gin.RouterGroup) {
	// Digniatries
	r.GET("/", mw.Auth, mw.CheckAdmin, FetchAllPeople)
	r.POST("/", mw.Auth, mw.CheckAdmin, AddNewPerson)
	r.PUT("/:personId", mw.Auth, mw.CheckAdmin, UpdatePersonDetails)
	r.DELETE("/:personId", mw.Auth, mw.CheckAdmin, DeletePerson)

	// Student
	r.GET("/students", mw.Auth, mw.CheckAdmin, GetAllStudents)
}
