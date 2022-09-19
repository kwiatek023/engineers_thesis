package simulation

type IEdgeUpdater interface {
	RunEdgeUpdating(update chan bool)
}
