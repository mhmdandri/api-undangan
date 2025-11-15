package routes

import (
	"api-undangan/controller"

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
		api.POST("/reservations/confirm", controller.ConfirmReservation)
	}
}