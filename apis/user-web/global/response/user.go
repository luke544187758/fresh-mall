package response

type UserResponse struct {
	ID       int64  `json:"id"`
	Role     int32  `json:"role"`
	NickName string `json:"nick_name"`
	Birthday string `json:"birthday"`
	Gender   string `json:"gender"`
	Address  string `json:"address"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}
