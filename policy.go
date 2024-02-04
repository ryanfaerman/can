package can

import (
	"fmt"
	"sync"
)

// Policy is a map of actions to functions that
// return an error if the action cannot be performed for the given account.
type Policy[K any] map[string]func(K) error

// Verbs implements the Verber interface.
func (p Policy[K]) Verbs() []string {
	var verbs []string
	for verb := range p {
		verbs = append(verbs, verb)
	}
	return verbs
}

// Can implements the Canner interface.
func (p Policy[K]) Can(actor K, action string) error {
	if fn, ok := p[action]; ok {
		return fn(actor)
	}
	return nil
}

// PolicyRegistry is a map of policy names to policies. Where the name
// is the type associated withthe policy.
var (
	PolicyRegistry = make(map[string]Policy[any])
	policyMutex    sync.RWMutex
)

// AddPolicy adds a policy to the PolicyRegistry.
func AddPolicy(p Policy[any]) {
	name := fmt.Sprintf("%T", p)

	policyMutex.Lock()
	{
		PolicyRegistry[name] = p
	}
	policyMutex.Unlock()
}

// PolicyExists returns true if a policy exists in the PolicyRegistry.
func PolicyExists(k any) bool {
	name := fmt.Sprintf("%T", k)
	policyMutex.RLock()
	defer policyMutex.RUnlock()

	_, ok := PolicyRegistry[name]
	return ok
}

// RemovePolicy removes a policy from the PolicyRegistry.
func RemovePolicy(k any) {
	name := fmt.Sprintf("%T", k)
	policyMutex.Lock()
	{
		delete(PolicyRegistry, name)
	}
	policyMutex.Unlock()
}
