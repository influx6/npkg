package nframes

import (
	"path"
	"runtime"
	"strings"

	"github.com/influx6/npkg"
)

// series of possible levels
const (
	FATAL     Level = 0x10
	ERROR           = 0x8
	DEBUG           = 0x4
	WARNING         = 0x2
	INFO            = 0x1
	ALL             = INFO | WARNING | DEBUG | ERROR | FATAL
	STACKABLE       = DEBUG | ERROR | FATAL
)

//**************************************************************
// Level
//**************************************************************

// Level defines a int type which represent the a giving level of entry for a giving entry.
type Level uint8

// Text2Level returns Level value for the giving string.
//
// It returns ALL as a default value if it does not know the level string.
func Text2Level(lvl string) Level {
	switch strings.ToLower(lvl) {
	case "FATAL", "fatal":
		return FATAL
	case "warning", "WARNING":
		return WARNING
	case "debug", "DEBUG":
		return DEBUG
	case "error", "ERROR":
		return ERROR
	case "info", "INFO":
		return INFO
	}
	return ALL
}

// String returns the string version of the Level.
func (l Level) String() string {
	switch l {
	case FATAL:
		return "FATAL"
	case DEBUG:
		return "DEBUG"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case INFO:
		return "INFO"
	}
	return "UNKNOWN"
}

//************************************************************
// Stack Frames
//************************************************************

// GetFrameDetails uses runtime.CallersFrames instead of runtime.FuncForPC
// which in go1.12 misses certain frames, to ensure backward compatibility
// with the previous version, this method is written to provide alternative
// setup that uses the recommended way of go1.12.
func GetFrameDetails(skip int, size int) []FrameDetail {
	var frames = make([]uintptr, size)
	var written = runtime.Callers(skip, frames)
	if written == 0 {
		return nil
	}

	var details = make([]FrameDetail, 0, len(frames))

	frames = frames[:written]
	var rframes = runtime.CallersFrames(frames)
	for {
		frame, more := rframes.Next()
		if !more {
			break
		}

		var detail FrameDetail
		detail.File = frame.File
		detail.Line = frame.Line
		detail.Method = frame.Function
		detail.FileName, detail.Package = fileToPackageAndFilename(frame.File)

		details = append(details, detail)
	}
	return details
}

// GetFrames returns a slice of stack frames for a giving size, skipping the provided
// `skip` count.
func GetFrames(skip int, size int) []uintptr {
	var frames = make([]uintptr, size)
	var written = runtime.Callers(skip, frames)
	return frames[:written]
}

// Frames is a slice of pointer uints.
type Frames []uintptr

// Details returns a slice of FrameDetails describing with a snapshot
// of giving stack frame pointer details.
func (f Frames) Details() []FrameDetail {
	var details = make([]FrameDetail, len(f))
	for ind, ptr := range f {
		details[ind] = Frame(ptr).Detail()
	}
	return details
}

// Encode encodes all Frames within slice into provided object encoder with keyname "_stack_frames".
func (f Frames) Encode(encoder npkg.ObjectEncoder) error {
	return encoder.ListFor("_stack_frames", f.EncodeList)
}

// EncodeList encodes all Frames within slice into provided list encoder.
func (f Frames) EncodeList(encoder npkg.ListEncoder) error {
	for _, frame := range f {
		var fr = Frame(frame)
		if err := encoder.AddObject(fr); err != nil {
			return err
		}
	}
	return nil
}

// Frame represents a program counter inside a stack frame.
// For historical reasons if Frame is interpreted as a uintptr
// its value represents the program counter + 1.
type Frame uintptr

// FrameDetail represent the snapshot description of a
// Frame pointer.
type FrameDetail struct {
	Line     int
	Method   string
	File     string
	Package  string
	FileName string
}

const srcSub = "/src/"

// EncodeObject encodes giving frame into provided encoder.
func (f Frame) EncodeObject(encode npkg.ObjectEncoder) error {
	fn := runtime.FuncForPC(f.Pc())
	if fn == nil {
		return nil
	}

	if err := encode.String("method", fn.Name()); err != nil {
		return err
	}

	var file, line = fn.FileLine(f.Pc())
	if line >= 0 {
		if err := encode.Int("line", line); err != nil {
			return err
		}
	}

	if file != "" && file != "???" {
		if err := encode.String("file", file); err != nil {
			return err
		}

		var fileName, pkgName = fileToPackageAndFilename(file)
		if err := encode.String("file_name", fileName); err != nil {
			return err
		}
		if err := encode.String("package", pkgName); err != nil {
			return err
		}
	}
	return nil
}

const winSlash = '\\'

func toSlash(s string) string {
	for index, item := range s {
		if item == winSlash {
			s = s[:index] + "/" + s[index+1:]
		}
	}
	return s
}

func fileToPackageAndFilename(file string) (filename, pkg string) {
	if runtime.GOOS == "windows" {
		file = toSlash(file)
	}

	var pkgIndex = strings.Index(file, srcSub)
	if pkgIndex != -1 {
		var pkgFileBase = file[pkgIndex+5:]
		if lastSlash := strings.LastIndex(pkgFileBase, "/"); lastSlash != -1 {
			filename = pkgFileBase[lastSlash+1:]
			pkg = pkgFileBase[:lastSlash]
		}
	}
	return
}

// Pc returns the program counter for this frame;
// multiple frames may have the same PC value.
func (f Frame) Pc() uintptr { return uintptr(f) - 1 }

// Detail returns the detail for giving Frame pointer.
func (f Frame) Detail() FrameDetail {
	var detail = FrameDetail{
		Method: "Unknown",
		File:   "...",
		Line:   -1,
	}

	fn := runtime.FuncForPC(f.Pc())
	if fn == nil {
		return detail
	}

	detail.Method = fn.Name()
	var file, line = fn.FileLine(f.Pc())
	if line >= 0 {
		detail.Line = line
	}
	if file != "" && file != "???" {
		detail.File = file
	}
	if detail.File != "..." {
		pkgPieces := strings.SplitAfter(detail.File, "/src/")
		var pkgFileBase string
		if len(pkgPieces) > 1 {
			pkgFileBase = pkgPieces[1]
		}

		detail.Package = path.Dir(pkgFileBase)
		detail.FileName = path.Base(pkgFileBase)
	}
	return detail
}

// File returns the full path to the file that contains the
// function for this Frame's pc.
func (f Frame) File() string {
	fn := runtime.FuncForPC(f.Pc())
	if fn == nil {
		return "unknown"
	}
	file, _ := fn.FileLine(f.Pc())
	return file
}

// Line returns the line number of source code of the
// function for this Frame's pc.
func (f Frame) Line() int {
	fn := runtime.FuncForPC(f.Pc())
	if fn == nil {
		return 0
	}
	_, line := fn.FileLine(f.Pc())
	return line
}

// Name returns the name of this function, if known.
func (f Frame) Name() string {
	fn := runtime.FuncForPC(f.Pc())
	if fn == nil {
		return "unknown"
	}
	return fn.Name()
}
