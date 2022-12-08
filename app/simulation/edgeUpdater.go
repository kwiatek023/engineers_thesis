package simulation

// IEdgeUpdater - interface used for reliability models
type IEdgeUpdater interface {
	// RunEdgeUpdating - runs edge updating task accordingly to chosen reliability model
	RunEdgeUpdating(update chan bool)
}
