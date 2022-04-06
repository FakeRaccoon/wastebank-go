package responses

type LoginResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}
