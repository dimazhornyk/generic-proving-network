package logic

import "context"

type InitialSyncer struct {
}

func NewInitialSyncer() *InitialSyncer {
	return &InitialSyncer{}
}

func (is *InitialSyncer) Sync(ctx context.Context) error {
	// TODO: implement
	// sync storage, sync nodes map
	return nil
}
