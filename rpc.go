package broadcast

type rpcService struct {
	svc *Service
}

// Publish Messages.
func (r *rpcService) Publish(msg []*Message, ok *bool) error {
	*ok = true
	return r.svc.Publish(msg...)
}

// Publish Messages in async mode.
func (r *rpcService) PublishAsync(msg []*Message, ok *bool) error {
	*ok = true
	go r.svc.Publish(msg...)
	return nil
}
