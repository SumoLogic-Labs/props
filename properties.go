package props

import (
	"sync"

	"github.com/magiconair/properties"
)

type Properties struct {
	mu    sync.RWMutex
	props *properties.Properties
}

var p *Properties

func init() {
	p = New()
}

func Replace(props *properties.Properties) { p.Replace(props) }

func (p *Properties) Replace(props *properties.Properties) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.props = props
}

func (p *Properties) Properties() *properties.Properties {
	p.mu.RLock()
	defer p.mu.RUnlock()
	res := properties.NewProperties()
	res.Prefix = p.props.Prefix
	res.Postfix = p.props.Postfix
	res.DisableExpansion = p.props.DisableExpansion
	res.WriteSeparator = p.props.WriteSeparator
	res.Merge(p.props)
	return res
}

func Get(key string) (string, bool) { return p.Get(key) }

func (p *Properties) Get(key string) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.props.Get(key)
}

func GetString(key, def string) string { return p.GetString(key, def) }

func (p *Properties) GetString(key, def string) string {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.props.GetString(key, def)
}

func GetBool(key string, def bool) bool { return p.GetBool(key, def) }

func (p *Properties) GetBool(key string, def bool) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.props.GetBool(key, def)
}

func GetInt(key string, def int) int { return p.GetInt(key, def) }

func (p *Properties) GetInt(key string, def int) int {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.props.GetInt(key, def)
}

func GetProperties() *Properties {
	return p
}

func New() *Properties {
	p := new(Properties)
	p.props = properties.NewProperties()
	return p
}
