package router

import (
	authApi "github.com/Infinite-Sum-Games/Am.EventKit/handlers/api/auth"
	eventApi "github.com/Infinite-Sum-Games/Am.EventKit/handlers/api/event"
	orgApi "github.com/Infinite-Sum-Games/Am.EventKit/handlers/api/organizers"
	mw "github.com/Infinite-Sum-Games/Am.EventKit/middleware"
	"github.com/gin-gonic/gin"
)

func AdminRouter(r *gin.RouterGroup) {
	adminRouter := r.Group("/admin")

	// Authentication routes
	// adminRouter.POST("/auth/login", authApi.LoginAdmin)
	auth := adminRouter.Group("/auth", mw.Auth, mw.RequireRoles("admin"))
	{
		auth.GET("/session", authApi.FetchAdminSession)
		auth.GET("/logout", authApi.Logout)
	}

	// Event routes
	event := adminRouter.Group("/event", mw.Auth, mw.RequireRoles("admin"))
	{
		// General-Event management
		event.GET("", eventApi.GetAllAdminEvents)
		event.GET("/:eventId", eventApi.GetAdminEventsById)
		event.GET("/new", eventApi.NewEvent)
		event.POST("/:eventId", eventApi.AddEventDetails)
		event.POST("/size/:eventId", eventApi.AddEventDimension)
		event.POST("/toggle/:eventId", eventApi.AddEventToggles)
		event.POST("/publish/:eventId", eventApi.PublishEvent)
		event.DELETE("/publish/:eventId", eventApi.UnpublishEvent)
		event.POST("/completed/:eventId", eventApi.MarkEventAsCompleted)
		event.DELETE("/completed/:eventId", eventApi.UnmarkEventAsCompleted)
		// Event-Poster management
		event.POST("/poster/:eventId", eventApi.AddEventPoster)
		event.DELETE("/poster/:eventId", eventApi.DeleteEventPoster)
		// Event-Organizer management
		event.POST("/org", eventApi.ConnectEventAndOrganizer)
		event.DELETE("/org", eventApi.DisconnectEventAndOrganizer)
		// Event-Tag management
		event.POST("/tag", eventApi.ConnectEventAndTags)
		event.DELETE("/tag", eventApi.DisconnectEventAndTags)
		// Event-Dignitary management
		// event.POST("/people", eventApi.ConnectEventAndPeople)
		// event.DELETE("/people", eventApi.DisconnectEventAndPeople)
		// Event-Schedule management
		event.POST("/schedule/:eventId")
		event.PUT("/schedule/:scheduleId")
		event.DELETE("/schedule/:scheduleId")
	}

	// Booking routes
	// book := adminRouter.Group("/booking", mw.Auth, mw.RequireRoles("admin"))
	// {
	// 	book.GET("/transactions", bookingApi.FetchAdminTransactions)
	// }

	// Dispute routes
	// dispute := adminRouter.Group("/dispute", mw.Auth, mw.RequireRoles("admin"))
	// {
	// 	dispute.GET("", disputeApi.GetAllDisputes)
	// 	dispute.POST("/:txnId", disputeApi.CreateDispute)
	// 	dispute.PUT("/:disputeId", disputeApi.UpdateDispute)
	// 	dispute.POST("/close-as-true/:disputeId", disputeApi.CloseAsTrueDispute)
	// 	dispute.POST("/close-as-false/:disputeId", disputeApi.CloseAsFalseDispute)
	// }

	// People routes
	// people := adminRouter.Group("/people", mw.Auth, mw.RequireRoles("admin"))
	// {
	// 	people.GET("", peopleApi.FetchAllPeople, mw.Auth, mw.RequireRoles("admin"))
	// 	people.POST("", peopleApi.AddNewPerson)
	// 	people.PUT("/:personId", peopleApi.UpdatePersonDetails)
	// 	people.DELETE("/:personId", peopleApi.UpdatePersonDetails)
	// }

	// Organizer routes
	org := adminRouter.Group("/org", mw.Auth, mw.RequireRoles("admin"))
	{
		org.GET("", orgApi.GetAllOrganizers)
		org.POST("", orgApi.CreateOrganizer)
		org.PUT("/:orgId", orgApi.EditOrganizer)
		org.PUT("/password/:orgId", orgApi.EditOrganizer)
		org.DELETE("/password/:orgId", orgApi.DeleteOrganizer)
	}

	// Tag routes
	// tag := adminRouter.Group("/tag", mw.Auth, mw.RequireRoles("admin"))
	// {
	// 	tag.GET("/", tagApi.FetchEventTags)
	// 	tag.POST("/", tagApi.CreateEventTag)
	// 	tag.PUT("/:tagId", tagApi.EditEventTag)
	// 	tag.DELETE("/:tagId", tagApi.DeleteEventTag)
	// }

	// TODO: Analytics routes
	// analytics := adminRouter.Group("/analytics", mw.Auth, mw.RequireRoles("admin"))
	{
		// analytics.GET("/quick", analyticsApi.GetQuickDashboard)
		// analytics.GET("/revenue", analyticsApi.GetRevenueAnalytics)
		// analytics.GET("/registrations", analyticsApi.GetEventRegistrationAnalytics)
		// analytics.GET("/transactions", analyticsApi.GetTransactionAnalytics)
	}
}
