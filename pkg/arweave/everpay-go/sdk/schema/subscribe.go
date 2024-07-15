package schema

type FilterQuery struct {
	StartCursor   uint64
	Address       string
	TokenSymbol   string
	Action        string
	WithoutAction string
}
