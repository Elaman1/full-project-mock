package service

import "context"

func JoinContexts(ctx1, ctx2 context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx1)
	go func() {
		select {
		case <-ctx1.Done():
		case <-ctx2.Done():
		}
		cancel()
	}()
	return ctx
}
