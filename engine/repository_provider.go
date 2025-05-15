package engine

import (
	"reflect"

	"github.com/gooferOrm/goofer/repository"
	"github.com/gooferOrm/goofer/schema"
)

// RepositoryProvider defines the interface for getting repositories for entity types
type RepositoryProvider interface {
	// Repository returns a repository for the given entity type
	Repository(entity schema.Entity) interface{}

	// MustRepository returns a repository for the given entity type and panics if the entity is not registered
	MustRepository(entity schema.Entity) interface{}
}

// Repository returns a repository for the given entity type
func (c *Client) Repository(entity schema.Entity) interface{} {
	t := schema.GetEntityType(entity)
	return c.getRepositoryForType(t)
}

// getRepositoryForType returns a repository for the given reflect.Type
func (c *Client) getRepositoryForType(t reflect.Type) interface{} {
	switch t.Kind() {
	case reflect.Ptr:
		elemType := t.Elem()
		switch elemType.Kind() {
		case reflect.Struct:
			repo := repository.NewUntypedRepository(elemType, c.db, c.dialect)
			return repo
		}
	case reflect.Struct:
		// If a non-pointer struct is passed, use its pointer type
		repo := repository.NewUntypedRepository(t, c.db, c.dialect)
		return repo
	}
	return nil
}

// MustRepository returns a repository for the given entity type and panics if the entity is not registered
func (c *Client) MustRepository(entity schema.Entity) interface{} {
	t := schema.GetEntityType(entity)
	repo := c.getRepositoryForType(t)
	if repo == nil {
		panic("failed to create repository for entity: " + t.String())
	}
	return repo
}


