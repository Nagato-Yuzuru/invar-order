package product

import "fmt"

type ErrorCoordinateNotFound struct {
	Coordinate SKUCoordinate
}

func (e *ErrorCoordinateNotFound) Error() string {
	return fmt.Sprintf("SKU Coordinate Not Found: %s", e.Coordinate)
}
