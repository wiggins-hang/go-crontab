package shutdown

import (
	"context"
	"sync"
	"time"

	"go-crontab/common"
	"go-crontab/log"

	"github.com/zeromicro/go-zero/core/threading"
)

type StopListenerManager struct {
	lock      sync.Mutex
	listeners []func()
}

var (
	//链接资源(mysql,rabbitmq等)
	ConnectResourceListeners = new(StopListenerManager)
)

func (this *StopListenerManager) RegisterStopListener(fn func()) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.listeners = append(this.listeners, fn)
}

func (this *StopListenerManager) NotifyStopListener(ctx context.Context) {
	this.lock.Lock()
	defer this.lock.Unlock()
	group := threading.NewRoutineGroup()
	for _, listener := range this.listeners {
		l := listener
		group.RunSafe(func() {
			shutDownServer(ctx, l)
		})
	}
	group.Wait()
}

func shutDownServer(ctx context.Context, fn func()) {
	finish := make(chan struct{}, 1)
	common.SafelyGo(func() {
		fn()
		finish <- struct{}{}
	})
	for {
		select {
		case <-ctx.Done():
			log.Info("timeout,force finish")
			return
		case <-finish:
			return
		default:
			time.Sleep(time.Second)
		}
	}
}
