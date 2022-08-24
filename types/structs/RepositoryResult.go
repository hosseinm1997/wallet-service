package structs

import "gorm.io/gorm"

type RepositoryResult[T any] struct {
	Data         []*T
	Error        error
	RowsAffected int64
	Model        *T
}

func (c *RepositoryResult[T]) Set(data any, response *gorm.DB) {
	switch response.RowsAffected {
	case 0:
	case 1:
		c.Model = data.(*T)
		c.Data = nil
	default:
		c.Model = nil
		c.Data = data.([]*T)
	}
	c.RowsAffected = response.RowsAffected
	c.Error = response.Error
}
