package base

// Page 分页基础结构
type Page struct {
	PageNum  int   `json:"page_num"`  // 当前页码
	PageSize int   `json:"page_size"` // 每页条数
	Total    int64 `json:"total"`     // 总条数
	Pages    int   `json:"pages"`     // 总页数
}

// NewPage 创建分页对象，设置默认值
func NewPage() *Page {
	return &Page{
		PageNum:  1,
		PageSize: 10,
	}
}

// Offset 获取 MySQL 分页的起始行
func (p *Page) Offset() int {
	return (p.PageNum - 1) * p.PageSize
}

// Limit 获取每页记录数
func (p *Page) Limit() int {
	return p.PageSize
}
