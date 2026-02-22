package save

import (
	"context"
	"io"

	"golang.org/x/time/rate"
)

// 限制读取速率
type readlimit struct {
	reader  io.Reader
	limiter *rate.Limiter
	ctx     context.Context
}

// newReadLimit 创建一个新的读取速率限制器
func newReadLimit(ctx context.Context, reader io.Reader, r rate.Limit, b int) *readlimit {
	return &readlimit{
		reader:  reader,
		limiter: rate.NewLimiter(r, b),
		ctx:     ctx,
	}
}

// Read 从读取速率限制器中读取数据
func (r *readlimit) Read(p []byte) (int, error) {
	var err error
	maxRead := r.limiter.Burst()
	if maxRead > len(p) {
		p = p[:maxRead]
	}
	if err = r.limiter.WaitN(r.ctx, maxRead); err != nil {
		return 0, err
	}
	return r.reader.Read(p)
}
