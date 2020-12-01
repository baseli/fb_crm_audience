package facebook

type Pagination struct {
	Pager	pager		`json:"paging"`
	HasNext func() bool
}
