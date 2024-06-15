package copyfile

import "github.com/rrgmc/debefix"

func New(options ...Option) *CopyFile {
	ret := &CopyFile{}
	for _, opt := range options {
		opt(ret)
	}
	return ret
}

func NewOptions(options ...Option) (*CopyFile, []debefix.LoadOption, []debefix.ResolveOption) {
	c := New(options...)
	return c,
		[]debefix.LoadOption{debefix.WithLoadValueParser(c)},
		[]debefix.ResolveOption{debefix.WithRowResolvedCallback(c)}
}

type Callback func(ctx debefix.ValueResolveContext, fieldname string, fileData FileData) error

type Option func(*CopyFile)

func WithCallback(callback Callback) Option {
	return func(c *CopyFile) {
		c.callback = callback
	}
}
