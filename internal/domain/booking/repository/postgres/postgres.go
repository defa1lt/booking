package postgres

import (
	"booking/config"
	"booking/internal/domain/booking/entities"
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	ctx context.Context
	log *zap.Logger
	cfg *config.ConfigModel
	DB  *pgxpool.Pool
}

func NewRepository(log *zap.Logger, cfg *config.ConfigModel, ctx context.Context) *Repository {
	return &Repository{
		ctx: ctx,
		log: log,
		cfg: cfg,
	}
}

func (r *Repository) OnStart(ctx context.Context) error {
	connectionUrl := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		r.cfg.Postgres.Host,
		r.cfg.Postgres.Port,
		r.cfg.Postgres.User,
		r.cfg.Postgres.Password,
		r.cfg.Postgres.DBName,
		r.cfg.Postgres.SSLMode)

	pool, err := pgxpool.Connect(ctx, connectionUrl)
	if err != nil {
		r.log.Error("failed to connect to database", zap.Error(err))
		return err
	}
	r.DB = pool
	r.log.Info("successfully connected to database")
	return nil
}

func (r *Repository) OnStop(ctx context.Context) error {
	if r.DB != nil {
		r.DB.Close()
		r.log.Info("successfully closed database connection")
	}
	return nil
}

const (
	// Queries for Hotel
	qCreateHotel = `
		INSERT INTO hotels (name, address) 
		VALUES ($1, $2)
		RETURNING id;
	`
	qGetHotelByID = `
		SELECT id, name, address 
		FROM hotels 
		WHERE id = $1;
	`
	qUpdateHotel = `
		UPDATE hotels 
		SET name = $1, address = $2 
		WHERE id = $3;
	`
	qDeleteHotel = `
		DELETE FROM hotels 
		WHERE id = $1;
	`
	qGetAllHotels = `
		SELECT id, name, address
		FROM hotels;
	`

	// Queries for Room
	qCreateRoom = `
		INSERT INTO rooms (hotel_id, number, type, price) 
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	qGetRoomByID = `
		SELECT id, hotel_id, number, type, price 
		FROM rooms 
		WHERE id = $1;
	`
	qUpdateRoom = `
		UPDATE rooms 
		SET hotel_id = $1, number = $2, type = $3, price = $4 
		WHERE id = $5;
	`
	qDeleteRoom = `
		DELETE FROM rooms 
		WHERE id = $1;
	`
	qGetAllRooms = `
		SELECT id, hotel_id, number, type, price
		FROM rooms;
	`
	qGetRoomsByHotelID = `
		SELECT id, hotel_id, number, type, price
		FROM rooms
		WHERE hotel_id = $1;
	`

	// Queries for Booking
	qCreateBooking = `
		INSERT INTO bookings (room_id, customer_id, check_in, check_out, status) 
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	qGetBookingByID = `
		SELECT id, room_id, customer_id, check_in, check_out, status 
		FROM bookings 
		WHERE id = $1;
	`
	qUpdateBooking = `
		UPDATE bookings 
		SET room_id = $1, customer_id = $2, check_in = $3, check_out = $4, status = $5 
		WHERE id = $6;
	`
	qDeleteBooking = `
		DELETE FROM bookings 
		WHERE id = $1;
	`
	qGetAllBookings = `
		SELECT id, room_id, customer_id, check_in, check_out, status
		FROM bookings;
	`

	// Queries for Customer
	qCreateCustomer = `
		INSERT INTO customers (first_name, last_name, email, phone) 
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	qGetCustomerByID = `
		SELECT id, first_name, last_name, email, phone 
		FROM customers 
		WHERE id = $1;
	`
	qUpdateCustomer = `
		UPDATE customers 
		SET first_name = $1, last_name = $2, email = $3, phone = $4 
		WHERE id = $5;
	`
	qDeleteCustomer = `
		DELETE FROM customers 
		WHERE id = $1;
	`
	qGetAllCustomers = `
		SELECT id, first_name, last_name, email, phone
		FROM customers;
	`
)

// ----------------------
// Hotel Methods
// ----------------------

func (r *Repository) CreateHotel(ctx context.Context, hotel *entities.Hotel) (int, error) {
	var hotelID int
	err := r.DB.QueryRow(ctx, qCreateHotel, hotel.Name, hotel.Address).Scan(&hotelID)
	if err != nil {
		r.log.Error("CreateHotel: error with INSERT INTO hotels", zap.Error(err))
		return 0, err
	}
	return hotelID, nil
}

func (r *Repository) GetHotelByID(ctx context.Context, hotelID int) (*entities.Hotel, error) {
	var hotel entities.Hotel
	err := r.DB.QueryRow(ctx, qGetHotelByID, hotelID).Scan(&hotel.ID, &hotel.Name, &hotel.Address)
	if err != nil {
		r.log.Error("GetHotelByID: error with SELECT FROM hotels", zap.Error(err))
		return nil, err
	}
	return &hotel, nil
}

func (r *Repository) UpdateHotel(ctx context.Context, hotel *entities.Hotel) error {
	_, err := r.DB.Exec(ctx, qUpdateHotel, hotel.Name, hotel.Address, hotel.ID)
	if err != nil {
		r.log.Error("UpdateHotel: error with UPDATE hotels", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) DeleteHotel(ctx context.Context, hotelID int) error {
	_, err := r.DB.Exec(ctx, qDeleteHotel, hotelID)
	if err != nil {
		r.log.Error("DeleteHotel: error with DELETE FROM hotels", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) GetAllHotels(ctx context.Context) ([]*entities.Hotel, error) {
	rows, err := r.DB.Query(ctx, qGetAllHotels)
	if err != nil {
		r.log.Error("GetAllHotels: error querying hotels", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var hotels []*entities.Hotel
	for rows.Next() {
		var hotel entities.Hotel
		err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Address)
		if err != nil {
			r.log.Error("GetAllHotels: error scanning row", zap.Error(err))
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}

	return hotels, nil
}

func (r *Repository) CreateRoom(ctx context.Context, room *entities.Room) (int, error) {
	var roomID int
	err := r.DB.QueryRow(ctx, qCreateRoom, room.HotelID, room.Number, room.Type, room.Price).Scan(&roomID)
	if err != nil {
		r.log.Error("CreateRoom: error with INSERT INTO rooms", zap.Error(err))
		return 0, err
	}
	return roomID, nil
}

func (r *Repository) GetRoomByID(ctx context.Context, roomID int) (*entities.Room, error) {
	var room entities.Room
	err := r.DB.QueryRow(ctx, qGetRoomByID, roomID).Scan(&room.ID, &room.HotelID, &room.Number, &room.Type, &room.Price)
	if err != nil {
		r.log.Error("GetRoomByID: error with SELECT FROM rooms", zap.Error(err))
		return nil, err
	}
	return &room, nil
}

func (r *Repository) UpdateRoom(ctx context.Context, room *entities.Room) error {
	_, err := r.DB.Exec(ctx, qUpdateRoom, room.HotelID, room.Number, room.Type, room.Price, room.ID)
	if err != nil {
		r.log.Error("UpdateRoom: error with UPDATE rooms", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) DeleteRoom(ctx context.Context, roomID int) error {
	_, err := r.DB.Exec(ctx, qDeleteRoom, roomID)
	if err != nil {
		r.log.Error("DeleteRoom: error with DELETE FROM rooms", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) GetAllRooms(ctx context.Context) ([]*entities.Room, error) {
	rows, err := r.DB.Query(ctx, qGetAllRooms)
	if err != nil {
		r.log.Error("GetAllRooms: error querying rooms", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var rooms []*entities.Room
	for rows.Next() {
		var room entities.Room
		err := rows.Scan(&room.ID, &room.HotelID, &room.Number, &room.Type, &room.Price)
		if err != nil {
			r.log.Error("GetAllRooms: error scanning row", zap.Error(err))
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	return rooms, nil
}

func (r *Repository) GetRoomsByHotelID(ctx context.Context, hotelID int) ([]*entities.Room, error) {
	rows, err := r.DB.Query(ctx, qGetRoomsByHotelID, hotelID)
	if err != nil {
		r.log.Error("GetRoomsByHotelID: error querying rooms", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var rooms []*entities.Room
	for rows.Next() {
		var room entities.Room
		err := rows.Scan(&room.ID, &room.HotelID, &room.Number, &room.Type, &room.Price)
		if err != nil {
			r.log.Error("GetRoomsByHotelID: error scanning row", zap.Error(err))
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	return rooms, nil
}

func (r *Repository) CreateBooking(ctx context.Context, booking *entities.Booking) (int, error) {
	var bookingID int
	err := r.DB.QueryRow(ctx, qCreateBooking, booking.RoomID, booking.CustomerID, booking.CheckIn, booking.CheckOut, booking.Status).Scan(&bookingID)
	if err != nil {
		r.log.Error("CreateBooking: error with INSERT INTO bookings", zap.Error(err))
		return 0, err
	}
	return bookingID, nil
}

func (r *Repository) GetBookingByID(ctx context.Context, bookingID int) (*entities.Booking, error) {
	var booking entities.Booking
	err := r.DB.QueryRow(ctx, qGetBookingByID, bookingID).Scan(&booking.ID, &booking.RoomID, &booking.CustomerID, &booking.CheckIn, &booking.CheckOut, &booking.Status)
	if err != nil {
		r.log.Error("GetBookingByID: error with SELECT FROM bookings", zap.Error(err))
		return nil, err
	}
	return &booking, nil
}

func (r *Repository) UpdateBooking(ctx context.Context, booking *entities.Booking) error {
	_, err := r.DB.Exec(ctx, qUpdateBooking, booking.RoomID, booking.CustomerID, booking.CheckIn, booking.CheckOut, booking.Status, booking.ID)
	if err != nil {
		r.log.Error("UpdateBooking: error with UPDATE bookings", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) DeleteBooking(ctx context.Context, bookingID int) error {
	_, err := r.DB.Exec(ctx, qDeleteBooking, bookingID)
	if err != nil {
		r.log.Error("DeleteBooking: error with DELETE FROM bookings", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) GetAllBookings(ctx context.Context) ([]*entities.Booking, error) {
	rows, err := r.DB.Query(ctx, qGetAllBookings)
	if err != nil {
		r.log.Error("GetAllBookings: error querying bookings", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var bookings []*entities.Booking
	for rows.Next() {
		var booking entities.Booking
		err := rows.Scan(&booking.ID, &booking.RoomID, &booking.CustomerID, &booking.CheckIn, &booking.CheckOut, &booking.Status)
		if err != nil {
			r.log.Error("GetAllBookings: error scanning row", zap.Error(err))
			return nil, err
		}
		bookings = append(bookings, &booking)
	}

	return bookings, nil
}

func (r *Repository) CreateCustomer(ctx context.Context, customer *entities.Customer) (int, error) {
	var customerID int
	err := r.DB.QueryRow(ctx, qCreateCustomer, customer.FirstName, customer.LastName, customer.Email, customer.Phone).Scan(&customerID)
	if err != nil {
		r.log.Error("CreateCustomer: error with INSERT INTO customers", zap.Error(err))
		return 0, err
	}
	return customerID, nil
}

func (r *Repository) GetCustomerByID(ctx context.Context, customerID int) (*entities.Customer, error) {
	var customer entities.Customer
	err := r.DB.QueryRow(ctx, qGetCustomerByID, customerID).Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone)
	if err != nil {
		r.log.Error("GetCustomerByID: error with SELECT FROM customers", zap.Error(err))
		return nil, err
	}
	return &customer, nil
}

func (r *Repository) UpdateCustomer(ctx context.Context, customer *entities.Customer) error {
	_, err := r.DB.Exec(ctx, qUpdateCustomer, customer.FirstName, customer.LastName, customer.Email, customer.Phone, customer.ID)
	if err != nil {
		r.log.Error("UpdateCustomer: error with UPDATE customers", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) DeleteCustomer(ctx context.Context, customerID int) error {
	_, err := r.DB.Exec(ctx, qDeleteCustomer, customerID)
	if err != nil {
		r.log.Error("DeleteCustomer: error with DELETE FROM customers", zap.Error(err))
		return err
	}
	return nil
}

func (r *Repository) GetAllCustomers(ctx context.Context) ([]*entities.Customer, error) {
	rows, err := r.DB.Query(ctx, qGetAllCustomers)
	if err != nil {
		r.log.Error("GetAllCustomers: error querying customers", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var customers []*entities.Customer
	for rows.Next() {
		var customer entities.Customer
		err := rows.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Email, &customer.Phone)
		if err != nil {
			r.log.Error("GetAllCustomers: error scanning row", zap.Error(err))
			return nil, err
		}
		customers = append(customers, &customer)
	}

	return customers, nil
}
