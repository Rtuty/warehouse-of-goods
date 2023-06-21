package goods

// WarehouseService - сервис склада
type WarehouseService struct{}

type Response struct {
	Result string
}

// ReserveProducts - резервирование товаров на складе для доставки
func (w *WarehouseService) ReserveProducts(args *[]string, reply *Response) error {
	*reply = Response{Result: ";kdadas;kjfasdk;jfk;ljsafdlkjsahdfljksahdashjdgas;dga;hd;jd;sdg"}

	return nil
}

// ReleaseProducts - освобождение резерва товаров
func (w *WarehouseService) ReleaseProducts(args *[]string, reply *bool) error {
	*reply = true
	return nil
}

// GetAvailableProductsCount - получение количества оставшихся товаров на складе
func (w *WarehouseService) GetAvailableProductsCount(args *int, reply *int) error {
	*reply = 0
	return nil
}
