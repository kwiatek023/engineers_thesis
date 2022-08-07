package simulation

type MessageQueue struct {
	queue []*Pack
}

func NewMessageQueue() *MessageQueue {
	return &MessageQueue{queue: make([]*Pack, 0)}
}

func (m *MessageQueue) Len() int {
	return len(m.queue)
}

func (m *MessageQueue) Enqueue(msgPack *Pack) {
	m.queue = append(m.queue, msgPack)
}

func (m *MessageQueue) Dequeue() *Pack {
	msg := m.queue[0]
	m.queue = m.queue[1:]
	return msg
}
