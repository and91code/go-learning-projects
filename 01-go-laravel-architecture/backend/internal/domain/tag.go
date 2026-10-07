package domain

// Tag representa a tabela 'tags' no banco de dados.
type Tag struct {
	ID     int64  `json:"id"`
	ListID int64  `json:"list_id"`
	Name   string `json:"name"`
	Color  string `json:"color"`
}
