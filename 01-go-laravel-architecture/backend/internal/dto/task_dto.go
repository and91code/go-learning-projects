package dto

type CreateTaskRequest struct {
	Title    string `json:"title" binding:"required,min=3,max=100"`
	Priority string `json:"priority" binding:"required,oneof=low medium high"`
}

type TaskResponse struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Priority string `json:"priority"`
	Status   string `json:"status"`
}
