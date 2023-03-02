package types

type Message struct {
	Name      string `json:"name" validate:"required"`
	Version   string `json:"version" validate:"required"`
	ID        string `json:"id" validate:"required"`
	Body      string `json:"body" validate:"required"`
	CreatedAt string `json:"createdAt" validate:"required"`
}
