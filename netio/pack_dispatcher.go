package netio

type PackDispatcher interface {
	HandleMsg(data []byte)
}

type DefaultPackDispatcher struct {
}
