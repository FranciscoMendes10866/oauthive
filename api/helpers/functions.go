package helpers

import "context"

func GetSessionID(ctx context.Context) int {
	return ctx.Value(CtxSessionID).(int)
}
