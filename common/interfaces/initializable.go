package interfaces

import (
	"context"
)

type Initializable interface {
	Init(context.Context) error
}
