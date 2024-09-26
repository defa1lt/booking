package usecase

import (
	"booking/internal/domain/booking/entities"
	"booking/internal/domain/booking/repository/postgres"
	"context"
	"go.uber.org/zap"
)

type Usecase struct {
	log  *zap.Logger
	Repo *postgres.Repository
}

func NewUsecase(logger *zap.Logger, Repo *postgres.Repository) (*Usecase, error) {
	return &Usecase{
		log:  logger,
		Repo: Repo,
	}, nil
}

func (u *Usecase) CreateHotel(ctx context.Context, hotel *entities.Hotel) (int, error) {
	hotelID, err := u.Repo.CreateHotel(ctx, hotel)
	if err != nil {
		u.log.Error("fail to create hotel", zap.Error(err))
		return 0, err
	}
	return hotelID, nil
}

func (u *Usecase) GetHotelByID(ctx context.Context, hotelID int) (*entities.Hotel, error) {
	hotel, err := u.Repo.GetHotelByID(ctx, hotelID)
	if err != nil {
		u.log.Error("fail to get hotel by ID", zap.Error(err))
		return nil, err
	}
	return hotel, nil
}

func (u *Usecase) UpdateHotel(ctx context.Context, hotel *entities.Hotel) error {
	err := u.Repo.UpdateHotel(ctx, hotel)
	if err != nil {
		u.log.Error("fail to update hotel", zap.Error(err))
		return err
	}
	return nil
}

func (u *Usecase) DeleteHotel(ctx context.Context, hotelID int) error {
	err := u.Repo.DeleteHotel(ctx, hotelID)
	if err != nil {
		u.log.Error("fail to delete hotel", zap.Error(err))
		return err
	}
	return nil
}

// Методы для работы с моделью Room

func (u *Usecase) CreateRoom(ctx context.Context, room *entities.Room) (int, error) {
	roomID, err := u.Repo.CreateRoom(ctx, room)
	if err != nil {
		u.log.Error("fail to create room", zap.Error(err))
		return 0, err
	}
	return roomID, nil
}

func (u *Usecase) GetRoomByID(ctx context.Context, roomID int) (*entities.Room, error) {
	room, err := u.Repo.GetRoomByID(ctx, roomID)
	if err != nil {
		u.log.Error("fail to get room by ID", zap.Error(err))
		return nil, err
	}
	return room, nil
}

func (u *Usecase) UpdateRoom(ctx context.Context, room *entities.Room) error {
	err := u.Repo.UpdateRoom(ctx, room)
	if err != nil {
		u.log.Error("fail to update room", zap.Error(err))
		return err
	}
	return nil
}

func (u *Usecase) DeleteRoom(ctx context.Context, roomID int) error {
	err := u.Repo.DeleteRoom(ctx, roomID)
	if err != nil {
		u.log.Error("fail to delete room", zap.Error(err))
		return err
	}
	return nil
}

// Методы для работы с моделью Booking

func (u *Usecase) CreateBooking(ctx context.Context, booking *entities.Booking) (int, error) {
	bookingID, err := u.Repo.CreateBooking(ctx, booking)
	if err != nil {
		u.log.Error("fail to create booking", zap.Error(err))
		return 0, err
	}
	return bookingID, nil
}

func (u *Usecase) GetBookingByID(ctx context.Context, bookingID int) (*entities.Booking, error) {
	booking, err := u.Repo.GetBookingByID(ctx, bookingID)
	if err != nil {
		u.log.Error("fail to get booking by ID", zap.Error(err))
		return nil, err
	}
	return booking, nil
}

func (u *Usecase) UpdateBooking(ctx context.Context, booking *entities.Booking) error {
	err := u.Repo.UpdateBooking(ctx, booking)
	if err != nil {
		u.log.Error("fail to update booking", zap.Error(err))
		return err
	}
	return nil
}

func (u *Usecase) DeleteBooking(ctx context.Context, bookingID int) error {
	err := u.Repo.DeleteBooking(ctx, bookingID)
	if err != nil {
		u.log.Error("fail to delete booking", zap.Error(err))
		return err
	}
	return nil
}

// Методы для работы с моделью Customer

func (u *Usecase) CreateCustomer(ctx context.Context, customer *entities.Customer) (int, error) {
	customerID, err := u.Repo.CreateCustomer(ctx, customer)
	if err != nil {
		u.log.Error("fail to create customer", zap.Error(err))
		return 0, err
	}
	return customerID, nil
}

func (u *Usecase) GetCustomerByID(ctx context.Context, customerID int) (*entities.Customer, error) {
	customer, err := u.Repo.GetCustomerByID(ctx, customerID)
	if err != nil {
		u.log.Error("fail to get customer by ID", zap.Error(err))
		return nil, err
	}
	return customer, nil
}

func (u *Usecase) UpdateCustomer(ctx context.Context, customer *entities.Customer) error {
	err := u.Repo.UpdateCustomer(ctx, customer)
	if err != nil {
		u.log.Error("fail to update customer", zap.Error(err))
		return err
	}
	return nil
}

func (u *Usecase) DeleteCustomer(ctx context.Context, customerID int) error {
	err := u.Repo.DeleteCustomer(ctx, customerID)
	if err != nil {
		u.log.Error("fail to delete customer", zap.Error(err))
		return err
	}
	return nil
}
