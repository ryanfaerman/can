package can

import (
	"errors"
	"fmt"
	"slices"
)

// Canner allows a resource to define if an account can perform a given action.
type Canner[K any] interface {
	Can(actor K, action string) error
}

// Verber is used to allow a resource to define the actions it supports.
type Verber interface {
	Verbs() []string
}

var (
	ErrNotAuthorized = errors.New("not authorized")
	ErrInvalidAction = errors.New("invalid action")
)

func Can[K any](actor K, action string, resources ...any) error {
	if len(resources) == 0 {
		name := fmt.Sprintf("%T", actor)
		if p, ok := PolicyRegistry[name]; ok {
			resources = append(resources, p)
		}
	}

	for _, resource := range resources {
		if r, ok := resource.(Verber); ok {
			if !slices.Contains(r.Verbs(), action) {
				return ErrInvalidAction
			}
		}

		if r, ok := resource.(Canner[K]); ok {
			if err := r.Can(actor, action); err != nil {
				err = errors.Join(ErrNotAuthorized, err)
				return err
			}
		}
	}
	return nil
}

func Not[K any](actor K, action string, resources ...any) bool {
	return Can(actor, action, resources...) != nil
}
