package entities

type Hotel struct {
	ID      int
	Name    string
	Address string
	Rooms   []Room
}

type Room struct {
	ID      int
	HotelID int
	Number  string
	Type    string
	Price   float64
}

type Booking struct {
	ID         int
	RoomID     int
	CustomerID int
	CheckIn    string
	CheckOut   string
	Status     string
}

type Customer struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
	Phone     string
}
