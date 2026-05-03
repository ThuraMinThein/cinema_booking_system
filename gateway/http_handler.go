package main

import "github.com/gin-gonic/gin"

type handler struct {
}

func NewHandler() *handler {
	return &handler{}
}

func (h *handler) registerRoutes(r *gin.Engine) {

	user := r.Group("/users")
	{
		user.POST("/", h.HandleCreateUser)
	}

	booking := r.Group("/bookings")
	{
		booking.POST("/", h.HandleCreateBooking)
	}

}

func (h *handler) HandleCreateUser(c *gin.Context) {

}

func (h *handler) HandleCreateBooking(c *gin.Context) {

}
