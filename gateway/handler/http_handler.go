package handler

import (
	"net/http"

	pb "github.com/ThuraMinThein/common/api"
	"github.com/ThuraMinThein/gateway/middlewares"
	"github.com/ThuraMinThein/gateway/pkg/helper"
	"github.com/gin-gonic/gin"
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
		booking.POST("/hold", h.HandleHoldBooking)
		booking.GET("", h.HandleGetBooking)
		// booking.GET("", h.HandleIsSeatAvailable)
		booking.DELETE("", h.HandleCancelBooking)
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

	var req *pb.CreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetString("user_id")

	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	req.UserId = userId
	response, err := h.bookingClient.Create(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, response)

}

func (h *handler) HandleHoldBooking(c *gin.Context) {
	var req *pb.HoldBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetString("user_id")

	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	req.UserId = userId
	response, err := h.bookingClient.HoldBooking(c, req)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

func (h *handler) HandleGetBooking(c *gin.Context) {
	var req *pb.FindAllRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetString("user_id")

	if userId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	req.UserId = userId
	response, err := h.bookingClient.FindAll(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
}

func (h *handler) HandleIsSeatAvailable(c *gin.Context) {
	var req *pb.IsSeatAvailableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	response, err := h.bookingClient.IsSeatAvailable(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
}

func (h *handler) HandleCancelBooking(c *gin.Context) {
	var req *pb.CancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	req.UserId = c.GetString("user_id")
	response, err := h.bookingClient.Cancel(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, response)
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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	response, err := h.seatsClient.GetSeats(c, req)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, response)
}

func (h *handler) HandleDeleteSeat(c *gin.Context) {

}
