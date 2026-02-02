package accommodation

import "context"

type Accommodation struct {
	ctx        context.Context
	userModule usermodule.IUser
}

func New(ctx context.Context) IAccommodation {
	return &User{
		ctx:        ctx,
		userModule: usermodule.New(ctx),
	}
}
