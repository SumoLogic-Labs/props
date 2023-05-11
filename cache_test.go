package props

import (
	"context"
	"fmt"
	"testing"

	"github.com/magiconair/properties"
)

type CtxKey string

var Key = new(CtxKey)

type SimpleSource struct{}

func (s *SimpleSource) Poll(ctx context.Context) (*properties.Properties, error) {
	return properties.MustLoadString(fmt.Sprintf("key=%v", ctx.Value(Key))), nil
}

func TestAsyncPoller_Concurrent(t *testing.T) {
	t.Run("Concurrent access of props should be safe", func(t *testing.T) {
		c := &Cache{
			Store:  GetProperties(),
			Source: &SimpleSource{},
		}
		ctx := context.TODO()
		c.sync(context.WithValue(ctx, Key, "init"))
		go func() {
			for i := 0; i < 100; i++ {
				c.sync(context.WithValue(ctx, Key, i))
			}
		}()
		for i := 100; i < 200; i++ {
			v := GetString("key", "")
			c.sync(context.WithValue(ctx, Key, i))
			nv := GetString("key", "")
			if v == nv || v == "" || nv == "" {
				t.Error("updating props failed")
				return
			}
		}
	})
}
