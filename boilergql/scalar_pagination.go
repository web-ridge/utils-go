package boilergql

type ConnectionBackwardPagination struct {
	Last   int     `json:"last"`
	Before *string `json:"before"`
}

type ConnectionForwardPagination struct {
	First int     `json:"first"`
	After *string `json:"after"`
}

type ConnectionPagination struct {
	Forward  *ConnectionForwardPagination  `json:"forward"`
	Backward *ConnectionBackwardPagination `json:"backward"`
}
