package netio

type PackDispatcher interface {
	HandleMsg(data []byte)
}

type DefaultPackDispatcher struct {
}

func NewDefaultPackDispatcher() *DefaultPackDispatcher {
	return &DefaultPackDispatcher{}
}

func (this *DefaultPackDispatcher) HandleMsg(data []byte) {
	nLen := len(data)
	if nLen <= 0 {
		return
	}

}
