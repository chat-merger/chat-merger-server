package domain

import "context"

type Controller interface {
	Run(ctx context.Context) error
}
