package can_test

import (
	"errors"
	"testing"

	"github.com/ryanfaerman/can"
)

type Account struct {
	ID        int
	Anonymous bool
}

type Post struct {
	OwnerID int
}

func (p *Post) Verbs() []string {
	return []string{"create", "read", "update", "delete"}
}

func (p *Post) Can(account *Account, action string) error {
	switch action {
	case "create":
		if account.Anonymous {
			return errors.New("anonymous accounts cannot create")
		}

	case "read":
		return nil

	case "update":
	case "delete":
		if account.ID != p.OwnerID {
			return errors.New("accounts can only update their own posts")
		}
	}

	return nil
}

var OverallPolicy = can.Policy[*Account]{
	"view-metrics": func(a *Account) error {
		if a.Anonymous {
			return errors.New("anonymous accounts cannot view metrics")
		}
		return nil
	},
}

func TestCan(t *testing.T) {
	// Define a policy for a resource.
	bob := &Account{ID: 1}
	sally := &Account{ID: 2}
	anonymous := &Account{Anonymous: true}

	examples := map[string]struct {
		actor   *Account
		verb    string
		subject *Post
		allowed bool
	}{
		"bob can create a new post": {
			actor:   bob,
			verb:    "create",
			subject: &Post{},
			allowed: true,
		},
		"sally can read a post": {
			actor:   sally,
			verb:    "read",
			subject: &Post{OwnerID: bob.ID},
			allowed: true,
		},
		"anon cannot create a post": {
			actor:   anonymous,
			verb:    "create",
			subject: &Post{OwnerID: bob.ID},
			allowed: false,
		},
		"no one can 'favorite' a post": {
			actor:   anonymous,
			verb:    "favorite",
			subject: &Post{OwnerID: bob.ID},
			allowed: false,
		},
		"bob can view metrics": {
			actor:   bob,
			verb:    "view-metrics",
			allowed: true,
		},
	}

	for desc, example := range examples {
		desc := desc
		example := example
		t.Run(desc, func(t *testing.T) {
			t.Parallel()
			var err error
			if example.subject != nil {
				err = can.Can[*Account](example.actor, example.verb, example.subject)
			} else {
				err = can.Can[*Account](example.actor, example.verb)
			}
			if example.allowed && err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !example.allowed && err == nil {
				t.Errorf("expected error, got nothing")
			}
		})
	}
}
