package main

import (
	pb "github.com/ThuraMinThein/common/api"
	"github.com/gin-gonic/gin"
)

type handler struct {
	bookingClient pb.BookingServiceClient
}

func NewHandler(bookingClient pb.BookingServiceClient) *handler {
	return &handler{
		bookingClient: bookingClient,
	}
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
