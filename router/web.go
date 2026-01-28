package router

import (
	accommodationApi "github.com/Infinite-Sum-Games/Am.EventKit/api/accomodation"
	authApi "github.com/Infinite-Sum-Games/Am.EventKit/api/auth"
	bookingApi "github.com/Infinite-Sum-Games/Am.EventKit/api/booking"
	eventApi "github.com/Infinite-Sum-Games/Am.EventKit/api/event"
	profileApi "github.com/Infinite-Sum-Games/Am.EventKit/api/profile"
	tagApi "github.com/Infinite-Sum-Games/Am.EventKit/api/tag"
	mw "github.com/Infinite-Sum-Games/Am.EventKit/middleware"
	"github.com/gin-gonic/gin"
)

func WebRouter(r *gin.RouterGroup) {
	// Public routes (no authentication required)
	public := r.Group("/")
	{
		// Event routes - public access
		public.GET("/events", eventApi.FetchAllEvents)
		public.GET("/events/:eventId", eventApi.FetchEventById)
		public.GET("/tags", tagApi.FetchEventTags)

		// Authentication routes - no auth required
		public.POST("/auth/check-email", authApi.CheckEmailExist)
		public.GET("/auth/register/csrf", authApi.RegisterUserAccountCsrf)
		public.POST("/auth/register", authApi.RegisterUserAccount)
		public.GET("/auth/login/csrf", authApi.LoginUserCsrf)
		public.POST("/auth/login", authApi.LoginUser)
		public.GET("/auth/forgot/csrf", authApi.ForgotUserPasswordCsrf)
		public.POST("/auth/forgot", authApi.ForgotUserPassword)
	}

	// Temporary authentication routes
	tempAuth := r.Group("/", mw.TempAuth)
	{
		tempAuth.GET("/auth/verify/csrf", authApi.VerifyUserOtpCsrf)
		tempAuth.POST("/auth/verify", authApi.VerifyUserOtp)
		tempAuth.POST("/auth/resend-otp", authApi.ResendUserOtp)
		tempAuth.GET("/auth/confirm/csrf", authApi.ConfirmPasswordChangeCsrf)
		tempAuth.POST("/auth/confirm", authApi.ConfirmPasswordChange)
		tempAuth.POST("/auth/resend-password-otp", authApi.ResendPasswordChangeOtp)
	}

	// User authenticated routes
	user := r.Group("/", mw.Auth, mw.RequireRoles("student"))
	{
		// Event routes with authentication
		user.GET("/events/auth", eventApi.FetchAllEventsWithAuth)
		user.GET("/events/auth/:eventId", eventApi.FetchEventByEventIdWithAuth)

		// Event star/favorite functionality
		user.POST("/events/:eventId/star", eventApi.StarEvent)
		user.POST("/events/:eventId/unstar", eventApi.UnstarEvent)

		// Booking routes
		user.GET("/booking/events/:eventId/csrf", bookingApi.BookEventCsrf)
		user.POST("/booking/events/:eventId", bookingApi.BookEvent)
		user.POST("/booking/verify", bookingApi.VerifyTransaction)

		// Profile routes
		user.GET("/profile", profileApi.FetchUserProfile)
		user.GET("/profile/edit/csrf", profileApi.EditUserProfileCsrf)
		user.PUT("/profile/edit", profileApi.EditUserProfile)
		user.GET("/profile/transactions", profileApi.GetAllUserTransactions)
		user.GET("/profile/tickets", profileApi.GetTickets)

		// Session management
		user.GET("/auth/session/user", authApi.FetchUserSession)

		// Accommodation form routes
		user.GET("/accommodation/exists", accommodationApi.AccomodationExists)
		user.GET("/accommodation/form/csrf", accommodationApi.AccomodationFormCsrf)
		user.POST("/accommodation/form/submit", accommodationApi.AccomodationFormSubmission)

		// Logout
		user.POST("/auth/logout", authApi.Logout)
	}
}
