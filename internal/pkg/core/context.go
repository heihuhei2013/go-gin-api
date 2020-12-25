package core

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"

	"github.com/xinliangnote/go-gin-api/internal/pkg/errno"
	"github.com/xinliangnote/go-gin-api/internal/pkg/journal"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

type HandlerFunc func(c Context)

type Journal = journal.T

const (
	_JournalName = "_journal_"
	_LoggerName  = "_logger_"
	_BodyName    = "_body_"
	_PayloadName = "_payload_"
)

var contextPool = &sync.Pool{
	New: func() interface{} {
		return new(context)
	},
}

func newContext(ctx *gin.Context) Context {
	context := contextPool.Get().(*context)
	context.ctx = ctx
	return context
}

func releaseContext(ctx Context) {
	c := ctx.(*context)
	c.ctx = nil
	contextPool.Put(c)
}

var _ Context = (*context)(nil)

type Context interface {
	init()

	// ShouldBindQuery 反序列化querystring
	// tag: `form:"xxx"` (注：不要写成query)
	ShouldBindQuery(obj interface{}) error

	// ShouldBindPostForm 反序列化postform(querystring会被忽略)
	// tag: `form:"xxx"`
	ShouldBindPostForm(obj interface{}) error

	// ShouldBindForm 同时反序列化querystring和postform;
	// 当querystring和postform存在相同字段时，postform优先使用。
	// tag: `form:"xxx"`
	ShouldBindForm(obj interface{}) error

	// ShouldBindJSON 反序列化postjson
	// tag: `json:"xxx"`
	ShouldBindJSON(obj interface{}) error

	// ShouldBindURI 反序列化path参数(如路由路径为 /user/:name)
	// tag: `uri:"xxx"`
	ShouldBindURI(obj interface{}) error

	// Redirect 重定向
	Redirect(code int, location string)

	Journal() Journal
	setJournal(journal Journal)
	disableJournal()

	Logger() *zap.Logger
	setLogger(logger *zap.Logger)

	GetPayload() interface{}
	SetPayload(payload errno.Error)

	Header() http.Header
	GetHeader(key string) string
	SetHeader(key, value string)

	RawData() []byte

	Method() string

	Host() string

	Path() string

	URI() string
}

type context struct {
	ctx *gin.Context
}

func (c *context) init() {
	body, err := c.ctx.GetRawData()
	if err != nil {
		panic(err)
	}

	c.ctx.Set(_BodyName, body)                                   // cache body是为了journal使用
	c.ctx.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // re-construct req body
}

// ShouldBindQuery 反序列化querystring
// tag: `form:"xxx"` (注：不要写成query)
func (c *context) ShouldBindQuery(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Query)
}

// ShouldBindPostForm 反序列化postform(querystring会被忽略)
// tag: `form:"xxx"`
func (c *context) ShouldBindPostForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.FormPost)
}

// ShouldBindForm 同时反序列化querystring和postform;
// 当querystring和postform存在相同字段时，postform优先使用。
// tag: `form:"xxx"`
func (c *context) ShouldBindForm(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.Form)
}

// ShouldBindJSON 反序列化postjson
// tag: `json:"xxx"`
func (c *context) ShouldBindJSON(obj interface{}) error {
	return c.ctx.ShouldBindWith(obj, binding.JSON)
}

// ShouldBindURI 反序列化path参数(如路由路径为 /user/:name)
// tag: `uri:"xxx"`
func (c *context) ShouldBindURI(obj interface{}) error {
	return c.ctx.ShouldBindUri(obj)
}

// Redirect 重定向
func (c *context) Redirect(code int, location string) {
	c.ctx.Redirect(code, location)
}

func (c *context) Journal() Journal {
	j, ok := c.ctx.Get(_JournalName)
	if !ok || j == nil {
		return nil
	}

	return j.(Journal)
}

func (c *context) setJournal(journal Journal) {
	c.ctx.Set(_JournalName, journal)
}

func (c *context) disableJournal() {
	c.setJournal(nil)
}

func (c *context) Logger() *zap.Logger {
	logger, ok := c.ctx.Get(_LoggerName)
	if !ok {
		return nil
	}

	return logger.(*zap.Logger)
}

func (c *context) setLogger(logger *zap.Logger) {
	c.ctx.Set(_LoggerName, logger)
}

func (c *context) GetPayload() interface{} {
	payload, _ := c.ctx.Get(_PayloadName)
	return payload
}

func (c *context) SetPayload(payload errno.Error) {
	payload.WithID(c.Journal().ID())
	c.ctx.Set(_PayloadName, payload)
}

func (c *context) Header() http.Header {
	header := c.ctx.Request.Header

	clone := make(http.Header, len(header))
	for k, v := range header {
		value := make([]string, len(v))
		copy(value, v)

		clone[k] = value
	}
	return clone
}

func (c *context) GetHeader(key string) string {
	return c.ctx.GetHeader(key)
}

func (c *context) SetHeader(key, value string) {
	c.ctx.Header(key, value)
}

func (c *context) RawData() []byte {
	body, ok := c.ctx.Get(_BodyName)
	if !ok {
		return nil
	}

	return body.([]byte)
}

// Method 请求的method
func (c *context) Method() string {
	return c.ctx.Request.Method
}

// Host 请求的host
func (c *context) Host() string {
	return c.ctx.Request.Host
}

// Path 请求的路径(不附带querystring)
func (c *context) Path() string {
	return c.ctx.Request.URL.Path
}

// URI unescape后的uri
func (c *context) URI() string {
	uri, _ := url.QueryUnescape(c.ctx.Request.URL.RequestURI())
	return uri
}
