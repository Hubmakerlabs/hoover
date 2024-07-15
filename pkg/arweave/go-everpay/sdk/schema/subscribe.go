package schema

type FilterQuery struct {
	StartCursor   int64
	Address       string
	TokenTag      string
	Action        string
	WithoutAction string
}
