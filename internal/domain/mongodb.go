package domain

type BulkUpsert struct {
	Filter interface{}
	Update interface{}
}

// PaginationMeta รายละเอียด Metadata ข้อมูลการแบ่งหน้า
type PaginationMeta struct {
	CurrentPage int   `json:"current_page" example:"1"`
	Limit       int   `json:"limit" example:"10"`
	TotalItems  int64 `json:"total_items" example:"100"`
	TotalPages  int   `json:"total_pages" example:"10"`
	HasNext     bool  `json:"has_next" example:"true"`
	HasPrev     bool  `json:"has_prev" example:"false"`
}

// PaginatedResponse โครงสร้างมาตรฐานสำหรับส่งคืนข้อมูลแบบ List ที่มี Pagination
type PaginatedResponse[T any] struct {
	Data       []T            `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}
