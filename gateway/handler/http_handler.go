package handler

import (
	pb "github.com/ThuraMinThein/common/api"
	"github.com/ThuraMinThein/gateway/middlewares"
	"github.com/ThuraMinThein/gateway/pkg/helper"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type handler struct {
	userClient    pb.UserServiceClient
	bookingClient pb.BookingServiceClient
	seatsClient   pb.SeatsServiceClient
}

func NewHandler(userClient pb.UserServiceClient, bookingClient pb.BookingServiceClient, seatsClient pb.SeatsServiceClient) *handler {
	return &handler{
		userClient:    userClient,
		bookingClient: bookingClient,
		seatsClient:   seatsClient,
	}
}

func (h *handler) RegisterRoutes(r *gin.Engine) {

	user := r.Group("/users")
	{
		user.POST("", h.HandleCreateUser)
		user.POST("/login", h.HandleLoginUser)
	}

	booking := r.Group("/bookings")
	booking.Use(middlewares.AuthMiddleware())
	{
		booking.POST("", h.HandleCreateBooking)
	}

	seat := r.Group("/seats")
	{
		seat.POST("/set", h.HandleSetSeats)
		seat.GET("", h.HandleGetSeats)
		seat.DELETE("", h.HandleDeleteSeat)
	}

}

// user related handlers
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
	token, err := helper.GetAccessToken(response.UserId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, gin.H{"user": response, "token": token})
}

func (h *handler) HandleLoginUser(c *gin.Context) {
	var req *pb.LoginUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response, err := h.userClient.LoginUser(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	token, err := helper.GetAccessToken(response.UserId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"user": response, "token": token})
}

// booking related handlers
func (h *handler) HandleCreateBooking(c *gin.Context) {

	logrus.Info("Create booking handler called")

}

func (h *handler) HandleGetBooking(c *gin.Context) {

}

// seats related handlers
func (h *handler) HandleSetSeats(c *gin.Context) {

	var req *pb.SetSeatsRequest
	response, err := h.seatsClient.SetSeats(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, response)

}

func (h *handler) HandleGetSeats(c *gin.Context) {
	var req *pb.GetSeatsRequest
	response, err := h.seatsClient.GetSeats(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

func (h *handler) HandleDeleteSeat(c *gin.Context) {

}
