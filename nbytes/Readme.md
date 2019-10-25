ByteStream
---------
Bytes provides a delimited stream reader and writer which will encoded a sequence of giving byte set 
with a giving character set as ending, a delimiter sort of to indicate to it's reader that this is the end 
of this set. It escapes the delimiter if it appears within the byte sequence to ensure preservation. 

It is useful when multiplexing multiple streams of bytes over a connection where we wish to send 
multiple messages without the use of size headers where we prefix the size of giving stream before sequence 
of bytes, this becomes useful when you have memory restrictions and can't know total size of incoming bytes 
unless fully read out which is not efficient has you may get to use up memory just trying to read up all available
data which maybe more larger than available memory.

By using sequenced based delimiting we can solve this and still adequately convey the ending and beginning 
of a new message stream, which becomes powerful when handling larger and never-ending incoming data.


## Install

```bash
go get -u github.com/influx6/npkg/nbytes
```


## Example


- Multiplexed Streaming

```go
var dest bytes.Buffer
writer := &mb.DelimitedStreamWriter{
	Dest:      &dest,
	Escape:    []byte(":/"),
	Delimiter: []byte("//"),
}

sentences := []string{
	"I went into park stream all alone before the isle lands.",
	"Isle lands of YOR, before the dream verse began we found the diskin.",
	"Break fast in bed, love and eternality for ever",
	"Awaiting the ending seen of waiting for you?",
	"Done be such a waste!",
	"{\"log\":\"token\", \"centry\":\"20\"}",
}

for _, sentence := range sentences {
	written, err := writer.Write([]byte(sentence))
	streamWritten, err := writer.End()
}

reader := &mb.DelimitedStreamReader{
	Src:       bytes.NewReader(dest.Bytes()),
	Escape:    []byte(":/"),
	Delimiter: []byte("//"),
}

for index, sentence := range sentences {
	res := make([]byte, len(sentence))
	read, err := reader.Read(res)
	
	// if we are end of stream segment, continue
	if err != nil && err == ErrEOS {
		continue
	}
}

```

- DelimitedStreamWriter

```go
var dest bytes.Buffer
writer := &mb.DelimitedStreamWriter{
	Dest:      &dest,
	Escape:    []byte(":/"),
	Delimiter: []byte("//"),
}

written, err := writer.Write([]byte("Wondering out the clouds"))
written, err := writer.Write([]byte("of endless streams beyond the shore"))
totalWritten, err := writer.End()
```


- DelimitedStreamWriter

```go
var dest bytes.Buffer
writer := &mb.DelimitedStreamWriter{
	Dest:      &dest,
	Escape:    []byte(":/"),
	Delimiter: []byte("//"),
}

written, err := writer.Write([]byte("Wondering out the clouds"))
written, err := writer.Write([]byte("of endless streams beyond the shore"))
totalWritten, err := writer.End()
```

-- DelimitedStreamReader 

```go
reader := &mb.DelimitedStreamReader{
	Src:       strings.NewReader("Wondering out the :/:///clouds of endless :///streams beyond the shore//"),
	Escape:    []byte(":/"),
	Delimiter: []byte("//"),
}

res := make([]byte, len(spec.In))
read, err := reader.Read(res)
```
