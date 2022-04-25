package forms

type PasswordLoginForm struct {
	Mobile    string `json:"mobile" binding:"required,mobile"`
	Password  string `json:"password" binding:"required,max=20,min=3"`
	Captcha   string `json:"captcha" binding:"required,min=5,max=5"`
	CaptchaID string `json:"captcha_id" binding:"required"`
}

type RegisterForm struct {
	Mobile   string `json:"mobile" binding:"required,mobile"`
	Password string `json:"password" binding:"required,max=20,min=3"`
}
