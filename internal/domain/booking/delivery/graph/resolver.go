// internal/domain/booking/delivery/graph/resolver.go

package graph

import (
	"booking/internal/domain/booking/usecase"
)

type Resolver struct {
	Usecase *usecase.Usecase
}
