package interfaces

type ConditionContextDebug interface {
}

type ConditionContext interface {
	GetAuthContext() AuthContext
	GetDebug() ConditionContextDebug
}

type Condition interface {
	Evaluate(ctx ConditionContext) (bool, error)
}
