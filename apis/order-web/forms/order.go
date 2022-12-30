package forms

type CreateOrderForm struct {
	Address string `json:"address" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Mobile  string `json:"mobile" binding:"required,mobile"`
	Remark  string `json:"remark" binding:"required"`
}
