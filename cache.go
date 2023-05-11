package props

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/magiconair/properties"
)

type Poller interface {
	Poll(context.Context) (*properties.Properties, error)
}

type Cache struct {
	Store           *Properties
	Source          Poller
	RefreshInterval time.Duration
	ExpireAfter     time.Duration
}

func (c *Cache) Poll(ctx context.Context) (*properties.Properties, error) {
	return c.Store.Properties(), nil
}

func (c *Cache) Start(stopCh <-chan struct{}) error {
	if _, err := c.sync(context.TODO()); err != nil {
		return fmt.Errorf("unable to poll: %w", err)
	}
	go c.syncLoop(stopCh)
	return nil
}

func (c *Cache) sync(ctx context.Context) (*properties.Properties, error) {
	p, err := c.Source.Poll(ctx)
	if err != nil {
		return p, fmt.Errorf("unable to poll: %w", err)
	}
	c.Store.Replace(p)
	return p, nil
}

func (c *Cache) syncLoop(stopCh <-chan struct{}) {
	t := time.NewTicker(c.RefreshInterval)
	expire := time.AfterFunc(c.ExpireAfter, func() {
		log.Printf("clearing cache due to expiry (last successful sync more than %v ago)\n", c.ExpireAfter)
		c.clear()
	})
	defer func() {
		t.Stop()
		expire.Stop()
	}()
	for {
		select {
		case <-t.C:
			if _, err := c.sync(context.TODO()); err != nil {
				log.Println("unable to sync cache:", err)
			} else {
				expire.Reset(c.ExpireAfter)
			}
		case <-stopCh:
			log.Println("stopping")
			return
		}
	}
}

func (c *Cache) clear() {
	c.Store.Replace(properties.NewProperties())
}

func NewAsyncPollerSource(p *Properties, src Poller, period time.Duration) *Cache {
	return &Cache{
		Store:           p,
		Source:          src,
		RefreshInterval: period,
	}
}
