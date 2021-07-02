package page

type Order = string

const (
	Asc  Order = "ASC"
	Desc Order = "DESC"
)

var maxSize = 100

func SetMaxSize(max int) {
	maxSize = max
}

type Pageable interface {
	SetNumber(int) Pageable                     // set page number
	SetSize(int) Pageable                       // set page size
	SetOrder(key string, order string) Pageable // set order keys
	SetNoPaging() Pageable                      // set no paing
	Offset() int
	Limit() int
	Order() []map[string]string
	IsPaging() bool
}

type page struct {
	page     int
	size     int
	order    []map[string]string
	noPaging bool
}

func NewPage(number, size int) Pageable {
	return page{page: number, size: size}
}

func (p page) SetNumber(page int) Pageable {
	p.page = page
	return p
}

func (p page) SetSize(size int) Pageable {
	p.size = size
	return p
}

func (p page) SetOrder(key string, order string) Pageable {
	p.order = append(p.order, map[string]string{key: order})
	return p
}

func (p page) Offset() int {
	if p.page <= 1 {
		return 0
	}

	if p.page > 1 && p.size > 0 {
		return (p.page - 1) * p.size
	}

	return 0
}

func (p page) Limit() int {
	if p.size <= 0 {
		return maxSize
	}

	return p.size
}

func (p page) Order() []map[string]string {
	return p.order
}

func (p page) SetNoPaging() Pageable {
	p.noPaging = true
	return p
}

func (p page) IsPaging() bool {
	return !p.noPaging
}
