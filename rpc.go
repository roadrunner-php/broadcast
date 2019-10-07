package broadcast

type rpcService struct {
	s *Service
}

// Broadcast messages.
func (r *rpcService) Broadcast(msg []*Message, ok *bool) error {
	if broker := r.s.Broker(); broker != nil {
		*ok = true
		return broker.Broadcast(msg...)
	}

	return nil
}

// Broadcast messages in async mode.
func (r *rpcService) BroadcastAsync(msg []*Message, ok *bool) error {
	if broker := r.s.Broker(); broker != nil {
		*ok = true
		go broker.Broadcast(msg...)
	}

	return nil
}
