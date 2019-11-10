package broadcast

type rpcService struct {
	s *Service
}

// Publish Messages.
func (r *rpcService) Publish(msg []*Message, ok *bool) error {
	*ok = true
	return r.s.Publish(msg...)
}

// Publish Messages in async mode.
func (r *rpcService) PublishAsync(msg []*Message, ok *bool) error {
	*ok = true
	go r.s.Publish(msg...)
	return nil
}
