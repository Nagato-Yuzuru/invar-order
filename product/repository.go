package product

type SKURepository interface {
	FindSKUByCoordinate(coordinate SKUCoordinate) (SKU, error)
}
