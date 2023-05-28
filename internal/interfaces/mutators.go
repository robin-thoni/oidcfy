package interfaces

import "net/http"

type MutatorContextDebug interface {
}

type MutatorContext interface {
	GetAuthContext() AuthContext
	GetDebug() MutatorContextDebug
}

type Mutator interface {
	Mutate(rw http.ResponseWriter, ctx MutatorContext) error
}
