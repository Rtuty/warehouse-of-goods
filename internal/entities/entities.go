package entities

type (
	Stock struct {
		ID        string `json:"stock_id"`
		Name      string `json:"stock_name"`
		Available bool   `json:"stock_available"`
	}

	Good struct {
		Code  int64   `json:"good_code"`
		Name  string  `json:"good_name"`
		Size  float64 `json:"good_size"`
		Value int64   `json:"good_value"`
	}
)
