package history

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

// level constants
const (
	RedLvl Level = iota
	YellowLvl
	ErrorLvl
	InfoLvl
)

// errors ...
var (
	ErrYellowAlert = errors.New("warning: error or invalid state occured")
	ErrRedAlert    = errors.New("very bad: earth damaging error occured, check now")
)

var (
	hl = struct {
		ml  sync.Mutex
		hls Handler
	}{}
)

// SetDefaultHandlers sets default Handlers to be included
// in all instances of Sources to be used to process provided
// BugLogs.
func SetDefaultHandlers(hs ...Handler) bool {
	if len(hs) == 0 {
		return false
	}

	hl.ml.Lock()
	defer hl.ml.Unlock()
	hl.hls = Handlers(hs)
	return true
}

//**************************************************************
// Ctx Interface
//**************************************************************

var _ = &BugLog{}

// Attrs defines a map type.
type Attrs map[string]interface{}

// Ctx exposes an interface which provides a means to collate giving
// data for a giving data.
type Ctx interface {
	Red(string, ...interface{}) Ctx
	Info(string, ...interface{}) Ctx
	Yellow(string, ...interface{}) Ctx
	Error(error, string, ...interface{}) Ctx

	WithTags(...string) Ctx
	WithFields(Attrs, ...string) Ctx
	WithCompute([]string, ...Compute) Ctx
	WithHandler([]string, ...Handler) Ctx
	WithTitle(string, ...interface{}) Ctx
	With(string, interface{}, ...string) Ctx
}

// WithHandlers  returns a new Source with giving Handlers as receivers of
// all provided instances of B.
func WithHandlers(hs ...Handler) Ctx {
	return &BugLog{bugs: Handlers(hs), Fields: make(map[string]Field), Signature: randName(20), From: time.Now()}
}

// WithTags returns a new Source with giving tags as default tags
// added to all B instances submitted through Source's Ctx instances.
func WithTags(hs ...string) Ctx {
	return &BugLog{Tags: hs, Fields: make(map[string]Field), Signature: randName(20), From: time.Now()}
}

// With returns a new Ctx adding the giving key-value.
func With(k string, v interface{}, tags ...string) Ctx {
	return (&BugLog{Signature: randName(20), Tags: tags, Fields: make(map[string]Field), From: time.Now()}).addKV(k, v)
}

// WithFields returns a new Source with giving fields as default tags
// added to all B instances submitted through Source's Ctx instances.
func WithFields(attr Attrs, tags ...string) Ctx {
	return (&BugLog{Signature: randName(20), Tags: tags, Fields: make(map[string]Field), From: time.Now()}).addKVS(attr)
}

// WithTitle returns a new Ctx setting the BugLog.Title.
func WithTitle(title string, v ...interface{}) Ctx {
	if len(v) != 0 {
		title = fmt.Sprintf(title, v...)
	}

	return &BugLog{Signature: randName(20), Fields: make(map[string]Field), Title: title, From: time.Now()}
}

//**************************************************************
// Handler Interface
//**************************************************************

// Recv defines a function type which receives a pointer of B.
type Recv func(BugLog) error

// FilterRecv defines a function type which returns true/false for
// a given BugLog.
type FilterRecv func(BugLog) bool

// Handler exposes a single method to deliver giving B value to
type Handler interface {
	Recv(BugLog) error
}

// HandlerFunc returns a new Handler using provided function for
// calls to the Handler.Recv method.
func HandlerFunc(recv Recv) Handler {
	return fnHandler{rc: recv}
}

// FilterByLevel returns a Handler which will only allow BugLogs
// with Status.Level equal or above provided level.
func FilterByLevel(recv Recv, level Level) Handler {
	return FilterFunc(recv, func(b BugLog) bool {
		if b.Status.Level < level {
			return false
		}
		return true
	})
}

// FilterFunc returns a new Handler which will filter all BugLog received
// by the recv function.
func FilterFunc(recv Recv, filter FilterRecv) Handler {
	return HandlerFunc(func(log BugLog) error {
		if filter(log) {
			return recv(log)
		}
		return nil
	})
}

//Handlers defines a slice of Handlers as a type.
type Handlers []Handler

// Recv calls individual Handler.Recv in slice with BugLog instance.
func (h Handlers) Recv(b BugLog) error {
	for _, hl := range h {
		if hl == nil {
			continue
		}
		if err := hl.Recv(b); err != nil {
			return err
		}
	}
	return nil
}

type fnHandler struct {
	rc Recv
}

// Recv implements the Handler and calls the underline Recv
// function provided with provided B pointer.
func (fn fnHandler) Recv(b BugLog) error {
	return fn.rc(b)
}

//**************************************************************
// Level
//**************************************************************

// Level defines a int type which represent the a giving level of entry for a giving entry.
type Level int

// GetLevel returns Level value for the giving string.
// It returns -1 if it does not know the level string.
func GetLevel(lvl string) Level {
	switch strings.ToLower(lvl) {
	case "red":
		return RedLvl
	case "yellow":
		return YellowLvl
	case "error":
		return ErrorLvl
	case "info":
		return InfoLvl
	}

	return -1
}

// String returns the string version of the Level.
func (l Level) String() string {
	switch l {
	case RedLvl:
		return "red"
	case YellowLvl:
		return "yellow"
	case ErrorLvl:
		return "error"
	case InfoLvl:
		return "info"
	}

	return "UNKNOWN"
}

//**************************************************************
// Location
//**************************************************************

// Location defines the location which an history occured in.
type Location struct {
	Function string `json:"function"`
	Line     int    `json:"line"`
	File     string `json:"file"`
}

// CallGraph embodies a graph representing the areas where a method
// call occured.
type CallGraph struct {
	In Location
	By Location
}

//**************************************************************
// Field
//**************************************************************

// Field represents a giving key-value pair with location details.
type Field struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

//**************************************************************
// Status struct.
//**************************************************************

// Progress embodies giving messages where giving progress with
// associated message and level for status.
type Status struct {
	Message string      `json:"message"`
	Level   Level       `json:"level"`
	Err     interface{} `json:"err"`
	Graph   CallGraph   `json:"graph"`
}

//**************************************************************
// Compute struct.
//**************************************************************

// Metric is the result of a computation.
type Metric struct {
	Title string      `json:"title"`
	Value interface{} `json:"title"`
	Meta  interface{} `json:"meta"`
}

// Compute defines an interface which exposes a method to get
// the title of giving computation with computed value.
type Compute interface {
	Compute() Metric
}

//**************************************************************
// BugLog struct.
//**************************************************************

// BugLog represent a giving record of data at a giving period of time.
type BugLog struct {
	Title     string           `json:"title"`
	Signature string           `json:"signature"`
	From      time.Time        `json:"from"`
	Tags      []string         `json:"tags"`
	Metrics   []Metric         `json:"metrics"`
	Status    Status           `json:"status"`
	Fields    map[string]Field `json:"fields"`

	bugs     Handler
	computes []Compute
}

// WithHandler adds giving handler into returned Ctx.
func (b *BugLog) WithHandler(tags []string, c ...Handler) Ctx {
	if len(c) == 0 {
		return b
	}

	var cr Handlers

	if len(c) == 1 {
		cr = Handlers{b.bugs, c[0]}
	} else {
		cr = Handlers{Handlers(c), b.bugs}
	}

	br := b.branch().addTags(tags...)
	br.bugs = cr
	return br
}

// WithCompute adds giving computation into giving bug.
func (b *BugLog) WithCompute(tags []string, c ...Compute) Ctx {
	if len(c) == 0 {
		return b
	}

	br := b.branch().addTags(tags...)
	br.computes = append(b.computes, c...)
	return br
}

// WithFields returns a new instance of BugLog from source but unique
// and sets giving fields.
func (b *BugLog) WithFields(kv Attrs, tags ...string) Ctx {
	if len(kv) == 0 {
		return b
	}

	return b.branch().addKVS(kv).addTags(tags...)
}

func (b *BugLog) addKVS(kv Attrs) *BugLog {
	for k, v := range kv {
		if before, ok := b.Fields[k]; ok {
			before.Value = v
			b.Fields[k] = before
		} else {
			b.Fields[k] = Field{
				Key:   k,
				Value: v,
			}
		}
	}
	return b
}

// With returns a new instance of BugLog from source but unique
// and adds key-value pair into Ctx.
func (b *BugLog) With(k string, v interface{}, tags ...string) Ctx {
	return b.branch().addKV(k, v).addTags(tags...)
}

func (b *BugLog) addKV(k string, v interface{}, tags ...string) *BugLog {
	if before, ok := b.Fields[k]; ok {
		before.Value = v
		b.Fields[k] = before
	} else {
		b.Fields[k] = Field{
			Key:   k,
			Value: v,
		}
	}
	return b
}

// WithTags returns a new instance of BugLog from source but unique
// and sets giving tags.
func (b *BugLog) WithTags(tags ...string) Ctx {
	if len(tags) == 0 {
		return b
	}

	return b.branch().addTags(tags...)
}

func (b *BugLog) addTags(tags ...string) *BugLog {
	b.Tags = append(b.Tags, tags...)
	return b
}

// FromTitle returns a new instance of BugLog from source but unique
// and sets giving title.
func (b *BugLog) WithTitle(title string, v ...interface{}) Ctx {
	if len(v) != 0 {
		title = fmt.Sprintf(title, v...)
	}

	br := b.branch()
	br.Title = title
	return br
}

// Red logs giving status message at giving time with RedLvl.
func (b *BugLog) Red(msg string, vals ...interface{}) Ctx {
	if len(vals) != 0 {
		msg = fmt.Sprintf(msg, vals...)
	}

	return b.logStatus(Status{
		Message: msg,
		Level:   RedLvl,
		Graph:   GetMethodGraph(3),
	})
}

// Yellow logs giving status message at giving time with YellowLvl.
func (b *BugLog) Yellow(msg string, vals ...interface{}) Ctx {
	if len(vals) != 0 {
		msg = fmt.Sprintf(msg, vals...)
	}

	return b.logStatus(Status{
		Message: msg,
		Level:   YellowLvl,
		Graph:   GetMethodGraph(3),
	})
}

// Error logs giving status message at giving time with ErrorLvl.
func (b *BugLog) Error(err error, msg string, vals ...interface{}) Ctx {
	if len(vals) != 0 {
		msg = fmt.Sprintf(msg, vals...)
	}

	return b.logStatus(Status{
		Err:     err,
		Message: msg,
		Level:   ErrorLvl,
		Graph:   GetMethodGraph(3),
	})
}

// Info logs giving status message at giving time with InfoLvl.
func (b *BugLog) Info(msg string, vals ...interface{}) Ctx {
	if len(vals) != 0 {
		msg = fmt.Sprintf(msg, vals...)
	}

	return b.logStatus(Status{
		Message: msg,
		Level:   InfoLvl,
		Graph:   GetMethodGraph(3),
	})
}

func (b *BugLog) logStatus(s Status) Ctx {
	blog := *b
	blog.Status = s
	blog.bugs = nil
	blog.computes = nil
	blog.Metrics = make([]Metric, len(b.computes))

	if blog.Title == "" {
		blog.Title = s.Graph.By.Function
	}

	// add the metrics into their respective spots.
	for index, compute := range b.computes {
		blog.Metrics[index] = compute.Compute()
	}

	if b.bugs != nil {
		if err := b.bugs.Recv(blog); err != nil {
			log.Printf("error logging: %+s", err)
		}
	}

	hl.ml.Lock()
	defer hl.ml.Unlock()
	if hl.hls != nil {
		if err := hl.hls.Recv(blog); err != nil {
			log.Printf("error logging: %+s", err)
		}
	}

	return b
}

// branch duplicates BugLog and copies appropriate
// dataset over to ensure uniqueness of values to source.
func (b *BugLog) branch() *BugLog {
	br := *b
	br.From = time.Now()
	br.Fields = make(map[string]Field)
	br.Signature = randName(20)

	br.computes = make([]Compute, len(b.computes))
	copy(br.computes, b.computes)

	br.Tags = make([]string, len(b.Tags))
	copy(br.Tags, b.Tags)

	for k, v := range b.Fields {
		br.Fields[k] = v
	}

	return &br
}

//**************************************************************
// internal methods and impl
//**************************************************************

func makeLocation(d int) Location {
	var loc Location
	loc.Function, loc.File, loc.Line = GetMethod(d)
	return loc
}

// We omit vowels from the set of available characters to reduce the chances
// of "bad words" being formed.
var alphanums = []rune("bcdfghjklmnpqrstvwxz0123456789")

// String generates a random alphanumeric string, without vowels, which is n
// characters long.  This will panic if n is less than zero.
func randName(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = alphanums[rand.Intn(len(alphanums))]
	}
	return string(b)
}

type sortedFields []Field

func (s sortedFields) Len() int {
	return len(s)
}

func (s sortedFields) Less(i, j int) bool {
	return s[i].Key < s[j].Key
}

func (s sortedFields) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
