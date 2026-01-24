package bootstrap

import (
	"net/http"
	"time"

	mw "github.com/Infinite-Sum-Games/Am.EventKit/internal/middleware"
	"github.com/Infinite-Sum-Games/Am.EventKit/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {

	config := cors.Config{
		AllowOrigins:              []string{cmd.Env.ClientDomain},
		AllowWildcard:             true,
		AllowMethods:              []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders:              []string{"X-Csrf-Token", "Origin", "Content-Type"},
		AllowCredentials:          true,
		OptionsResponseStatusCode: 204,
		MaxAge:                    12 * time.Hour,
	}

	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(cors.New(config)) // Setup CORS() first before other middlewares
	// r.Use(mw.MaintainanceMiddleware)
	r.Use(pkg.Log.LogMiddleware)
	r.Use(common.TagRequestWithId)
	r.Use(mw.RecoveryPanics)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Server is live ◪_◪",
		})
		pkg.Log.SuccessCtx(c)
	})

	v1 := r.Group("/api/v1")
	{
		authRouter := v1.Group("/auth")
		attendanceRouter := v1.Group("/attendance")
		userRouter := v1.Group("/user")
		eventRouter := v1.Group("/events")
		peopleRouter := v1.Group("/people")
		tagRouter := v1.Group("/tags")
		organizerRouter := v1.Group("/organizers")
		analyticsRouter := v1.Group("/analytics")
		disputeRouter := v1.Group("/disputes")
		accomodationRouter := v1.Group("/accommodation")

		apiAuth.StudentAuthRoutes(authRouter)
		apiAuth.OrganizerAuthRoutes(authRouter)
		apiAuth.AdminAuthRoutes(authRouter)
		apiProfile.ProfileRoutes(userRouter)
		apiEvent.EventRoutes(eventRouter)
		apiTag.TagRoutes(tagRouter)
		apiAttend.AttendanceRoutes(attendanceRouter)
		apiPeople.PeopleRoutes(peopleRouter)
		apiOrganizers.OrganizerRoutes(organizerRouter)
		apiOrganizers.OrganizerDashboardRoutes(organizerRouter)
		apiBooking.BookingRoutes(eventRouter)
		apiAnalytics.AnalyticsRoutes(analyticsRouter)
		apiDispute.DisputeRoutes(disputeRouter)

		apiAccomodation.AccomodationFormRoutes(accomodationRouter)
		apiAccomodation.AccomodationAuthRoutes(accomodationRouter)
		apiAccomodation.AccomodationPanelRoutes(accomodationRouter)
		apiAccomodation.FinanceRoutes(accomodationRouter)
		apiAccomodation.GateRoutes(accomodationRouter)
		apiAccomodation.SecurityRoutes(accomodationRouter)
	}

	v2 := r.Group("/api/v2")
	{
		router.WebRouter(v2)
		router.AdminRouter(v2)
		router.OrganizerWebRouter(v2)
		router.OrganizerAppRouter(v2)
		router.LogisticsWebRouter(v2)
		router.LogisticsAppRouter(v2)
	}

	return r
}
