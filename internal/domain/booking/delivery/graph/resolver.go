// internal/domain/booking/delivery/graph/resolver.go

package graph

import (
	"booking/internal/domain/booking/usecase"
)

// Resolver служит зависимостью для резолверов
type Resolver struct {
	Usecase *usecase.Usecase
}
