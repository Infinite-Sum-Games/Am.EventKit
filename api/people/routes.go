package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func PeopleRoutes(r *gin.RouterGroup) {
	r.GET("/", mw.Auth, FetchAllPeople)
	r.GET("/department/:dept_id", mw.Auth, FetchPeopleByDepartment)
	r.GET("/event/:event_id", mw.Auth, FetchPeopleByEvent)
	r.GET("/day/:day", mw.Auth, FetchPeopleByDay)

	r.POST("/", mw.Auth, AddNewPerson)

	r.PUT("/:id", mw.Auth, UpdatePersonDetails)

	r.DELETE("/:id", mw.Auth, DeletePerson)
}
