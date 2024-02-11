package rest

type TaskCreateRequest struct {
	Title string `json:"title"`
}

type TaskCreateResponse struct {
	ID string `json:"id"`
}
