package types

type ApiResponse[T any] struct {
	Code         int           `json:"code"`
	Content      T             `json:"content"`
	Paging       *PageMetaData `json:"paging,omitempty"`
	Message      any           `json:"message,omitempty"`
	ShouldNotify bool          `json:"should_notify"`
	Success      bool          `json:"success"`
}

type PageResponse[T any] struct {
	Content      []T          `json:"content,omitempty"`
	PageMetaData PageMetaData `json:"paging,omitempty"`
}

type PageMetaData struct {
	Page       int `json:"page"`
	Size       int `json:"size"`
	TotalItems int `json:"total_items"`
	TotalPages int `json:"total_pages"`
}
