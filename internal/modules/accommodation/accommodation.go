package accommodation

import "context"

type Accommodation struct {
	ctx context.Context
}

func New(ctx context.Context) IAccomodation {
	return &Accommodation{
		ctx: ctx,
	}
}
