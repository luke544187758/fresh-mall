package forms

type BannerForm struct {
	Image string `json:"image" binding:"url"`
	Url   string `json:"url" binding:"url"`
	Index int32  `json:"index" binding:"required"`
}
