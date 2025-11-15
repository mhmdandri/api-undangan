package routes

import (
	"api-undangan/controller"
	"api-undangan/middleware"

	"github.com/gin-gonic/gin"
)
func RegisterRoutes(r *gin.Engine){
	api := r.Group("/api")
	{
		api.GET("/comments", controller.GetComments)
		api.POST("/comments", controller.PostComment)

		api.GET("/reservations", controller.GetReservations)
		api.POST("/reservations", controller.CreateReservation)
		api.GET("/reservations/:code", controller.FindReservationByCode)
		api.POST("/login", controller.Login)
		
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/reservations/confirm", controller.ConfirmReservation)
	}
}