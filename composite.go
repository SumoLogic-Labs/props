package props

import (
	"context"

	"github.com/magiconair/properties"
)

type Composite struct {
	Sources []Poller
}

func (c Composite) Poll(ctx context.Context) (*properties.Properties, error) {
	props := properties.NewProperties()
	props.DisableExpansion = true
	for _, src := range c.Sources {
		p, err := src.Poll(ctx)
		if err != nil {
			return nil, err
		}
		props.Merge(p)
	}
	return props, nil
}

func NewCompositeSource(sources ...Poller) *Composite {
	return &Composite{
		Sources: sources,
	}
}
