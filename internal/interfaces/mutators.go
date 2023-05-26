package interfaces

type Mutator interface {
	Mutate(*AuthContext) error
}
