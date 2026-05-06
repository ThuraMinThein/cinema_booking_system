package main

import (
	pb "github.com/ThuraMinThein/common/api"
	"github.com/gin-gonic/gin"
)

type handler struct {
	userClient    pb.UserServiceClient
	bookingClient pb.BookingServiceClient
}

func NewHandler(userClient pb.UserServiceClient, bookingClient pb.BookingServiceClient) *handler {
	return &handler{
		userClient:    userClient,
		bookingClient: bookingClient,
	}
}

func (h *handler) registerRoutes(r *gin.Engine) {

	user := r.Group("/users")
	{
		user.POST("", h.HandleCreateUser)
	}

	booking := r.Group("/bookings")
	{
		booking.POST("", h.HandleCreateBooking)
	}

}

func (h *handler) HandleCreateUser(c *gin.Context) {
	var req *pb.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response, err := h.userClient.CreateUser(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
}

func (h *handler) HandleCreateBooking(c *gin.Context) {

}
