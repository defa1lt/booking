package server

import (
	"booking/config"
	"booking/internal/domain/booking/delivery/routes"
	"booking/internal/domain/booking/entities"
	"booking/internal/domain/booking/usecase"
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
)

type Link struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
	Type string `json:"type"`
}

type HypermediaResponse struct {
	Data  interface{} `json:"data"`
	Links []Link      `json:"links"`
}

type Server struct {
	logger  *zap.Logger
	cfg     *config.ConfigModel
	API     *fiber.App
	Usecase *usecase.Usecase
}

func NewServer(logger *zap.Logger, cfg *config.ConfigModel, uc *usecase.Usecase) (*Server, error) {
	return &Server{
		logger:  logger,
		cfg:     cfg,
		API:     fiber.New(),
		Usecase: uc,
	}, nil
}

func (s *Server) OnStart(_ context.Context) error {
	handlers := map[string]fiber.Handler{
		"CreateHotel":    s.CreateHotel,
		"GetHotel":       s.GetHotel,
		"UpdateHotel":    s.UpdateHotel,
		"DeleteHotel":    s.DeleteHotel,
		"CreateRoom":     s.CreateRoom,
		"GetRoom":        s.GetRoom,
		"UpdateRoom":     s.UpdateRoom,
		"DeleteRoom":     s.DeleteRoom,
		"CreateBooking":  s.CreateBooking,
		"GetBooking":     s.GetBooking,
		"UpdateBooking":  s.UpdateBooking,
		"DeleteBooking":  s.DeleteBooking,
		"CreateCustomer": s.CreateCustomer,
		"GetCustomer":    s.GetCustomer,
		"UpdateCustomer": s.UpdateCustomer,
		"DeleteCustomer": s.DeleteCustomer,
	}

	routes.SetupRouter(s.API, handlers)

	go func() {
		s.logger.Debug("HTTP server started on :3000")
		if err := s.API.Listen(":3000"); err != nil {
			s.logger.Error("failed to start HTTP server: " + err.Error())
		}
	}()
	return nil
}

func (s *Server) OnStop(_ context.Context) error {
	s.logger.Debug("HTTP server stopped")
	return s.API.Shutdown()
}

func generateHotelLinks(hotelID int) []Link {
	return []Link{
		{
			Rel:  "self",
			Href: fmt.Sprintf("/hotels/%d", hotelID),
			Type: "GET",
		},
		{
			Rel:  "update",
			Href: fmt.Sprintf("/hotels/%d", hotelID),
			Type: "PUT",
		},
		{
			Rel:  "delete",
			Href: fmt.Sprintf("/hotels/%d", hotelID),
			Type: "DELETE",
		},
	}
}

func (s *Server) CreateHotel(c *fiber.Ctx) error {
	var hotel entities.Hotel
	if err := c.BodyParser(&hotel); err != nil {
		s.logger.Error("Error parsing hotel body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	hotelID, err := s.Usecase.CreateHotel(c.Context(), &hotel)
	if err != nil {
		s.logger.Error("Error creating hotel", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	response := HypermediaResponse{
		Data:  hotel,
		Links: generateHotelLinks(hotelID),
	}

	return c.JSON(response)
}

func (s *Server) GetHotel(c *fiber.Ctx) error {
	hotelID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid hotel ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid hotel ID")
	}

	hotel, err := s.Usecase.GetHotelByID(c.Context(), hotelID)
	if err != nil {
		s.logger.Error("Error getting hotel", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	response := HypermediaResponse{
		Data:  hotel,
		Links: generateHotelLinks(hotelID),
	}

	return c.JSON(response)
}

func (s *Server) UpdateHotel(c *fiber.Ctx) error {
	var hotel entities.Hotel
	hotelID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid hotel ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid hotel ID")
	}
	hotel.ID = hotelID

	if err := c.BodyParser(&hotel); err != nil {
		s.logger.Error("Error parsing hotel body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = s.Usecase.UpdateHotel(c.Context(), &hotel)
	if err != nil {
		s.logger.Error("Error updating hotel", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) DeleteHotel(c *fiber.Ctx) error {
	hotelID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid hotel ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid hotel ID")
	}

	err = s.Usecase.DeleteHotel(c.Context(), hotelID)
	if err != nil {
		s.logger.Error("Error deleting hotel", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func generateRoomLinks(roomID int) []Link {
	return []Link{
		{
			Rel:  "self",
			Href: fmt.Sprintf("/rooms/%d", roomID),
			Type: "GET",
		},
		{
			Rel:  "update",
			Href: fmt.Sprintf("/rooms/%d", roomID),
			Type: "PUT",
		},
		{
			Rel:  "delete",
			Href: fmt.Sprintf("/rooms/%d", roomID),
			Type: "DELETE",
		},
	}
}

func (s *Server) CreateRoom(c *fiber.Ctx) error {
	var room entities.Room
	if err := c.BodyParser(&room); err != nil {
		s.logger.Error("Error parsing room body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	roomID, err := s.Usecase.CreateRoom(c.Context(), &room)
	if err != nil {
		s.logger.Error("Error creating room", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	response := HypermediaResponse{
		Data:  room,
		Links: generateRoomLinks(roomID),
	}

	return c.JSON(response)
}

func (s *Server) GetRoom(c *fiber.Ctx) error {
	roomID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid room ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid room ID")
	}

	room, err := s.Usecase.GetRoomByID(c.Context(), roomID)
	if err != nil {
		s.logger.Error("Error getting room", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	response := HypermediaResponse{
		Data:  room,
		Links: generateRoomLinks(roomID),
	}

	return c.JSON(response)
}

func (s *Server) UpdateRoom(c *fiber.Ctx) error {
	var room entities.Room
	roomID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid room ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid room ID")
	}
	room.ID = roomID

	if err := c.BodyParser(&room); err != nil {
		s.logger.Error("Error parsing room body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = s.Usecase.UpdateRoom(c.Context(), &room)
	if err != nil {
		s.logger.Error("Error updating room", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) DeleteRoom(c *fiber.Ctx) error {
	roomID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid room ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid room ID")
	}

	err = s.Usecase.DeleteRoom(c.Context(), roomID)
	if err != nil {
		s.logger.Error("Error deleting room", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func generateBookingLinks(bookingID int) []Link {
	return []Link{
		{
			Rel:  "self",
			Href: fmt.Sprintf("/bookings/%d", bookingID),
			Type: "GET",
		},
		{
			Rel:  "update",
			Href: fmt.Sprintf("/bookings/%d", bookingID),
			Type: "PUT",
		},
		{
			Rel:  "delete",
			Href: fmt.Sprintf("/bookings/%d", bookingID),
			Type: "DELETE",
		},
	}
}

func (s *Server) CreateBooking(c *fiber.Ctx) error {
	var booking entities.Booking
	if err := c.BodyParser(&booking); err != nil {
		s.logger.Error("Error parsing booking body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	bookingID, err := s.Usecase.CreateBooking(c.Context(), &booking)
	if err != nil {
		s.logger.Error("Error creating booking", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	response := HypermediaResponse{
		Data:  booking,
		Links: generateBookingLinks(bookingID),
	}

	return c.JSON(response)
}

func (s *Server) GetBooking(c *fiber.Ctx) error {
	bookingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid booking ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid booking ID")
	}

	booking, err := s.Usecase.GetBookingByID(c.Context(), bookingID)
	if err != nil {
		s.logger.Error("Error getting booking", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	response := HypermediaResponse{
		Data:  booking,
		Links: generateBookingLinks(bookingID),
	}

	return c.JSON(response)
}

func (s *Server) UpdateBooking(c *fiber.Ctx) error {
	var booking entities.Booking
	bookingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid booking ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid booking ID")
	}
	booking.ID = bookingID

	if err := c.BodyParser(&booking); err != nil {
		s.logger.Error("Error parsing booking body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = s.Usecase.UpdateBooking(c.Context(), &booking)
	if err != nil {
		s.logger.Error("Error updating booking", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) DeleteBooking(c *fiber.Ctx) error {
	bookingID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid booking ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid booking ID")
	}

	err = s.Usecase.DeleteBooking(c.Context(), bookingID)
	if err != nil {
		s.logger.Error("Error deleting booking", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func generateCustomerLinks(customerID int) []Link {
	return []Link{
		{
			Rel:  "self",
			Href: fmt.Sprintf("/customers/%d", customerID),
			Type: "GET",
		},
		{
			Rel:  "update",
			Href: fmt.Sprintf("/customers/%d", customerID),
			Type: "PUT",
		},
		{
			Rel:  "delete",
			Href: fmt.Sprintf("/customers/%d", customerID),
			Type: "DELETE",
		},
	}
}

func (s *Server) CreateCustomer(c *fiber.Ctx) error {
	var customer entities.Customer
	if err := c.BodyParser(&customer); err != nil {
		s.logger.Error("Error parsing customer body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	customerID, err := s.Usecase.CreateCustomer(c.Context(), &customer)
	if err != nil {
		s.logger.Error("Error creating customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	response := HypermediaResponse{
		Data:  customer,
		Links: generateCustomerLinks(customerID),
	}

	return c.JSON(response)
}

func (s *Server) GetCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid customer ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid customer ID")
	}

	customer, err := s.Usecase.GetCustomerByID(c.Context(), customerID)
	if err != nil {
		s.logger.Error("Error getting customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	response := HypermediaResponse{
		Data:  customer,
		Links: generateCustomerLinks(customerID),
	}

	return c.JSON(response)
}

func (s *Server) UpdateCustomer(c *fiber.Ctx) error {
	var customer entities.Customer
	customerID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid customer ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid customer ID")
	}
	customer.ID = customerID

	if err := c.BodyParser(&customer); err != nil {
		s.logger.Error("Error parsing customer body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = s.Usecase.UpdateCustomer(c.Context(), &customer)
	if err != nil {
		s.logger.Error("Error updating customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) DeleteCustomer(c *fiber.Ctx) error {
	customerID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		s.logger.Error("Invalid customer ID", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).SendString("Invalid customer ID")
	}

	err = s.Usecase.DeleteCustomer(c.Context(), customerID)
	if err != nil {
		s.logger.Error("Error deleting customer", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
