package forms

type CategoryForm struct {
	Name           string `json:"name" binding:"required,min=3,max=20"`
	ParentCategory int64  `json:"parent_category,string"`
	Level          int32  `json:"level" binding:"required,oneof=1 2 3"`
	IsTab          *bool  `json:"is_tab" binding:"required"`
}

type UpdateCategoryForm struct {
	Name  string `json:"name" binding:"required,min=3,max=20"`
	IsTab *bool  `json:"is_tab"`
}
