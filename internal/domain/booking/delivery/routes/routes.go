package routes

import "github.com/gofiber/fiber/v2"

func SetupRouter(app *fiber.App, handlers map[string]fiber.Handler) {

	app.Post("/hotels", handlers["CreateHotel"])
	app.Get("/hotels/:id", handlers["GetHotel"])
	app.Put("/hotels/:id", handlers["UpdateHotel"])
	app.Delete("/hotels/:id", handlers["DeleteHotel"])

	app.Post("/rooms", handlers["CreateRoom"])
	app.Get("/rooms/:id", handlers["GetRoom"])
	app.Put("/rooms/:id", handlers["UpdateRoom"])
	app.Delete("/rooms/:id", handlers["DeleteRoom"])

	app.Post("/bookings", handlers["CreateBooking"])
	app.Get("/bookings/:id", handlers["GetBooking"])
	app.Put("/bookings/:id", handlers["UpdateBooking"])
	app.Delete("/bookings/:id", handlers["DeleteBooking"])

	app.Post("/customers", handlers["CreateCustomer"])
	app.Get("/customers/:id", handlers["GetCustomer"])
	app.Put("/customers/:id", handlers["UpdateCustomer"])
	app.Delete("/customers/:id", handlers["DeleteCustomer"])
}
