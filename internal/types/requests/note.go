package requests

type CreateNoteRequest struct {
	Title   string `json:"title" validate:"required,min=3,max=50"`
	Content string `json:"content" validate:"required,min=3,max=500"`
}

type UpdateNoteRequest struct {
	Title   string `json:"title,omitempty" validate:"omitempty,min=3,max=50"`
	Content string `json:"content,omitempty" validate:"omitempty,min=3,max=500"`
}
