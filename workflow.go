package temporaltestcancel

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func sleepUntil(ctx workflow.Context, wakeUpTime time.Time, triggerTimerChannel workflow.ReceiveChannel) (err error) {
	timerCtx, timerCancel := workflow.WithCancel(ctx)
	duration := wakeUpTime.Sub(workflow.Now(timerCtx))
	timer := workflow.NewTimer(timerCtx, duration)

	workflow.NewSelector(timerCtx).
		AddFuture(timer, func(f workflow.Future) {
			_ = f.Get(timerCtx, nil)
		}).
		AddReceive(triggerTimerChannel, func(c workflow.ReceiveChannel, more bool) {
			timerCancel()
			c.Receive(timerCtx, nil)
		}).
		Select(timerCtx)
	return ctx.Err()
}

func Workflow(ctx workflow.Context) error {
	return sleepUntil(ctx, workflow.Now(ctx).Add(128*time.Hour), workflow.GetSignalChannel(ctx, "FireTimer"))
}
