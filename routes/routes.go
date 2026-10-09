package routes

import (
	"api-undangan/controller"
	"api-undangan/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/comments", controller.GetComments)
		api.POST("/comments", controller.PostComment)

		api.GET("/reservations", controller.GetReservations)
		api.GET("/reservations/preview-email", controller.PreviewReservationEmail)
		api.GET("/reservations/check-in/:code", controller.CheckInLookup)
		api.POST("/reservations/check-in", controller.GuestCheckIn)
		api.POST("/reservations/check-in/:code", controller.GuestCheckInByCode)
		api.POST("/reservations", controller.CreateReservation)
		api.GET("/reservations/:code", controller.FindReservationByCode)
		api.POST("/login", controller.Login)

		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/reservations/confirm", controller.ConfirmReservation)
	}
}
