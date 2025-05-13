package engine

import (
	"github.com/gooferOrm/goofer/repository"
	"github.com/gooferOrm/goofer/schema"
)

func Repo[T schema.Entity](c *Client) *repository.Repository[T] {
    return repository.NewRepository[T](c.db, c.dialect)
}