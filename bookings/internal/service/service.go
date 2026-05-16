package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ThuraMinThein/bookings/internal/model"
	"github.com/ThuraMinThein/bookings/internal/repository"
	"github.com/ThuraMinThein/bookings/pkg/memory_storage"
	"github.com/ThuraMinThein/common/api"
	"golang.org/x/sync/errgroup"
)

type service struct {
	repository         *repository.Repository
	userGrpcConnection api.UserServiceClient
	seatGrpcConnection api.SeatsServiceClient
}

func NewService(repo *repository.Repository, userGrpcConnection api.UserServiceClient, seatGrpcConnection api.SeatsServiceClient) *service {
	return &service{
		repository:         repo,
		userGrpcConnection: userGrpcConnection,
		seatGrpcConnection: seatGrpcConnection,
	}
}

func (s *service) Create(request *api.CreateRequest) error {
	redisCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var (
		bookings []*model.Booking
		mu       sync.Mutex
	)

	g, ctx := errgroup.WithContext(redisCtx)

	for _, seatID := range request.SeatIds {
		seatID := seatID

		g.Go(func() error {
			var booking *model.Booking

			key := fmt.Sprintf(
				"booking:%d:%d",
				request.MovieId,
				seatID,
			)

			if err := memory_storage.GetData(ctx, key, &booking); err != nil {
				return err
			}

			if booking == nil {
				return errors.New("hold a seat first")
			}

			if booking.UserID != request.UserId {
				return errors.New("this seat is held by someone else")
			}

			booking.Status = "Booked"
			mu.Lock()
			bookings = append(bookings, booking)
			mu.Unlock()

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return s.repository.Create(bookings)
}

func (s *service) HoldBooking(request *api.HoldBookingRequest) error {

	redisCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	gRPCCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var (
		user *api.GetUserByIdResponse
		seat *api.GetOneSeatResponse

		userErr error
		seatErr error
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		user, userErr = s.userGrpcConnection.GetUserById(
			gRPCCtx,
			&api.GetUserByIdRequest{
				UserId: request.UserId,
			},
		)
	}()

	go func() {
		defer wg.Done()
		seat, seatErr = s.seatGrpcConnection.GetOneSeat(
			gRPCCtx,
			&api.GetOneSeatRequest{
				SeatId: request.SeatId,
			},
		)
	}()

	wg.Wait()

	if userErr != nil {
		return userErr
	}

	if seatErr != nil {
		return seatErr
	}

	bookingModel := &model.Booking{
		UserID:     request.UserId,
		UserName:   user.User.Name,
		MovieID:    request.MovieId,
		SeatID:     request.SeatId,
		SeatNumber: seat.Seat.Name,
		Status:     "Hold",
		Showtime:   time.Now(),
	}

	if err := memory_storage.SaveBooking(redisCtx, request, bookingModel); err != nil {
		return err
	}

	booking, err := s.repository.FindByMovieAndSeatID(
		request.MovieId,
		request.SeatId,
	)

	rollback := func() {
		key := fmt.Sprintf(
			"booking:%d:%d",
			request.MovieId,
			request.SeatId,
		)
		_ = memory_storage.InvalidateStorage(redisCtx, key)
	}

	if err != nil {
		rollback()
		return err
	}

	if booking != nil {
		rollback()
		return errors.New("seat already booked")
	}

	return nil
}

func (s *service) FindAll(userId string, movieId int64) ([]model.Booking, error) {
	return s.repository.FindAll(userId, movieId)
}

func (s *service) FindAllBookedSeats(movieId int64) ([]model.Booking, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	dbBookings, err := s.repository.FindAllBookedSeats(movieId)
	if err != nil {
		return nil, err
	}

	redisBookings, err := memory_storage.GetMovieBookings(ctx, movieId)
	if err != nil {
		return nil, err
	}

	allBookings := append(dbBookings, redisBookings...)

	return allBookings, nil
}

func (s *service) IsSeatAvailable(movieID int64, seatID int64) (bool, string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	var booking *model.Booking

	key := fmt.Sprintf("booking:%d:%d", movieID, seatID)

	if err := memory_storage.GetData(ctx, key, &booking); err != nil {
		return false, "Redis Error", err
	}

	if booking != nil {
		return false, "Held", nil
	}

	booking, err := s.repository.FindByMovieAndSeatID(movieID, seatID)
	if err != nil {
		return false, "DB Error", err
	}

	var message string
	if booking == nil {
		message = "Available"
	} else {
		message = "Booked"
	}

	return booking == nil, message, nil
}

func (s *service) Update(booking *model.Booking) error {
	return s.repository.Update(booking)
}

func (s *service) Cancel(bookingID int64) error {
	return s.repository.Cancel(bookingID)
}
