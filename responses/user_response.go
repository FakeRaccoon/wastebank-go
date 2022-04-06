package responses

type UserResponse struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}
