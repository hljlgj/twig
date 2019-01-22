package twig

import (
	"context"
	"os"
	"os/signal"
	"time"
)

// Reloader 描述一个可以被重新加载的对象
// 在某些信号发生时候， 可以对Relaoder对象进行Reload操作，用于重新加载
type Reloader interface {
	Reload() error
}

// 信号处理函数
// 返回true 退出
// 返回false 等待处理下一个信号
type SignalFunc func(os.Signal) bool

// 正常退出，不做任何处理
func Quit() SignalFunc {
	return func(_ os.Signal) bool {
		return true
	}
}

// GracefulShutdown
func Graceful(t *Twig, timeout time.Duration) SignalFunc {
	return func(_ os.Signal) bool {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := t.Shutdown(ctx); err != nil {
			t.Logger.Println(err)
		}

		return true
	}
}

// Signal 用于监听系统信号并堵塞当前gorouting
// 参数f为信号处理函数
// 参数sig 为需要监听的系统信号，未出现在sig中的信号会被忽略
// 如果sig 为空，则监听所有信号
// 特别注意：部分操作系统的信号不可以被忽略 (SIGKILL & SIGSTOP)
func Signal(f SignalFunc, sig ...os.Signal) {
	ch := make(chan os.Signal)
	defer close(ch)

	signal.Notify(ch, sig...)

	for s := range ch {
		if f(s) {
			break
		}
	}
}
