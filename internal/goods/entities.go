package goods

type (
	Stock struct {
		ID        string `json:"stock_id"` //TODO Возможно, необязательный параметр
		Name      string `json:"stock_name"`
		Available bool   `json:"stock_available"`
	}

	Product struct {
		Code  string `json:"product_code"`
		Name  string `json:"product_name"`
		Size  string `json:"product_size"` // TODO возможно другой тип данных
		Value string `json:"product_value"`
	}
)
