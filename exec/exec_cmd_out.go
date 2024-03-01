package exec

import (
	"bytes"
	"fmt"
	"io"
)

type CommandInOut interface {
	Out() io.Writer
	Err() io.Writer
	Log() io.Writer
	In() io.Reader
}

func NewCommandInOut(out, err, log io.Writer, in io.Reader) CommandInOut {
	return &inOut{
		out: out,
		err: err,
		log: log,
		in:  in,
	}
}

type inOut struct {
	out io.Writer
	err io.Writer
	log io.Writer
	in  io.Reader
}

func (s *inOut) Out() io.Writer {
	return s.out
}

func (s *inOut) Err() io.Writer {
	return s.err
}

func (s *inOut) Log() io.Writer {
	return s.log
}

func (s *inOut) In() io.Reader {
	return s.in
}

//goland:noinspection GoUnusedExportedFunction
func NewBufferedInOut() *BufferedInOut {
	var out bytes.Buffer
	var err bytes.Buffer
	var log bytes.Buffer
	var in bytes.Buffer

	return &BufferedInOut{
		out:       &out,
		err:       &err,
		logBuffer: &log,
		logWriter: nil,
		in:        &in,
	}
}

//goland:noinspection GoUnusedExportedFunction
func NewBufferedInOutWithLog(log io.Writer) *BufferedInOut {
	var out bytes.Buffer
	var err bytes.Buffer
	var in bytes.Buffer

	return &BufferedInOut{
		out:       &out,
		err:       &err,
		logBuffer: getBufferOrNil(log),
		logWriter: log,
		in:        &in,
	}
}

func getBufferOrNil(writer io.Writer) *bytes.Buffer {
	if writer == nil {
		var buffer bytes.Buffer
		return &buffer
	} else {
		return nil
	}
}

type BufferedInOut struct {
	out       *bytes.Buffer
	err       *bytes.Buffer
	logBuffer *bytes.Buffer
	logWriter io.Writer
	in        *bytes.Buffer
}

func (s *BufferedInOut) Out() io.Writer {
	return s.out
}

func (s *BufferedInOut) Err() io.Writer {
	return s.err
}

func (s *BufferedInOut) Log() io.Writer {
	if s.logBuffer != nil {
		return s.logBuffer
	} else {
		return s.logWriter
	}
}

func (s *BufferedInOut) In() io.Reader {
	return s.in
}

func (s *BufferedInOut) GetCombinedOutput() string {
	outStr := s.GetOut()
	errStr := s.GetErr()

	if len(errStr) == 0 {
		return outStr
	} else {
		return fmt.Sprintf("stdout:\n%s\nstderr:\n%s\n", outStr, errStr)
	}
}

func (s *BufferedInOut) GetOut() string {
	return s.out.String()
}

func (s *BufferedInOut) GetErr() string {
	return s.err.String()
}

func (s *BufferedInOut) GetLog() string {
	if s.logBuffer == nil {
		return ""
	}
	return s.logBuffer.String()
}
