package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func PeopleRoutes(r *gin.RouterGroup) {
	r.GET("/people", mw.Auth, FetchAllPeople)
	r.GET("/people/department/:dept", mw.Auth, FetchPeopleByDepartment)
	r.GET("/people/event/:event", mw.Auth, FetchPeopleByEvent)
	r.GET("/people/day/:day", mw.Auth, FetchPeopleByDay)

	r.POST("/people", mw.Auth, AddNewPerson)

	r.PUT("/people/:id", mw.Auth, UpdatePersonDetails)

	r.DELETE("/people/:id", mw.Auth, DeletePerson)
}
