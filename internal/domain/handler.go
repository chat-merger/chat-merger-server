package domain

import "context"

type Handler interface {
	Serve(ctx context.Context) error
}
