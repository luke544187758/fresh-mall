package forms

type BrandForm struct {
	Name string `json:"name" binding:"required,min=3,max=10"`
	Logo string `json:"logo" binding:"url"`
}

type CategoryBrandForm struct {
	CategoryId int64 `json:"category_id,string" binding:"required"`
	BrandId    int64 `json:"brand_id,string" binding:"required"`
}
