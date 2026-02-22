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
func newReadLimit(ctx context.Context, reader io.Reader, rating int, max int) *readlimit {
	return &readlimit{
		reader:  reader,
		limiter: rate.NewLimiter(rate.Limit(rating), max),
		ctx:     ctx,
	}
}

// Read 从读取速率限制器中读取数据
func (r *readlimit) Read(p []byte) (int, error) {
	var err error
	// 计算最大读取量，不超过缓冲区大小和令牌桶突发大小
	maxRead := r.limiter.Burst()
	if maxRead > len(p) {
		maxRead = len(p)
	}
	// 等待足够的令牌
	if err = r.limiter.WaitN(r.ctx, maxRead); err != nil {
		return 0, err
	}
	return r.reader.Read(p[:maxRead])
}
