package api

import (
	mw "github.com/Thanus-Kumaar/anokha-2025-backend/middleware"
	"github.com/gin-gonic/gin"
)

func PeopleRoutes(r *gin.RouterGroup) {
	r.GET("/", FetchAllPeople)
	r.GET("/department/:dept", FetchPeopleByDepartment)
	r.GET("/event/:event", FetchPeopleByEvent)
	r.GET("/day/:day", mw.Auth, FetchPeopleByDay)

	r.POST("/", AddNewPerson)

	r.PUT("/:id", UpdatePersonDetails)

	r.DELETE("/:id", DeletePerson)
}
