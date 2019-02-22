package json_test

import (
	gjson "encoding/json"
	"testing"

	"github.com/gokit/history"
	"github.com/gokit/history/encoders/json"
	"github.com/stretchr/testify/require"
)

type WritableBuffer struct {
	Data []byte
}

func (wb *WritableBuffer) Write(l history.Level, content interface{}) error {
	if bu, ok := content.([]byte); ok {
		wb.Data = append(wb.Data, bu...)
	}
	return nil
}

func TestJSONEncoding(t *testing.T) {
	var w = new(WritableBuffer)
	w.Data = make([]byte, 0, 512)

	var section = history.WithEntry(3)
	section.Int("id", 2)
	section.String("name", "beta")
	section.Bool("flag", true)

	var err = section.WriteERROR(json.MessageObject("Digesting details"), w)
	require.NoError(t, err)
	require.NotEmpty(t, w.Data)

	var message = string(w.Data)
	require.Contains(t, message, `Digesting details`)
	require.Contains(t, message, `"flag": true`)
	require.Contains(t, message, `"name": "beta"`)
	require.Contains(t, message, `"id": 2`)
}

func TestGetJSON(t *testing.T) {
	t.Run("basic list", func(t *testing.T) {
		event := json.List()
		event.AddString("thunder")
		event.AddInt(234)
		require.Equal(t, "[\"thunder\", 234]", event.Message())
	})

	t.Run("basic fields", func(t *testing.T) {
		event := json.MessageObject("My log")
		event.String("name", "thunder")
		event.Int("id", 234)
		require.Equal(t, "{\"message\": \"My log\", \"name\": \"thunder\", \"id\": 234}", event.Message())
	})

	t.Run("with Object fields", func(t *testing.T) {
		event := json.MessageObject("My log")
		event.String("name", "thunder")
		event.Int("id", 234)

		var jsid, err = gjson.Marshal(map[string]interface{}{"id": 23})
		require.NoError(t, err)
		require.NotEmpty(t, jsid)

		event.Bytes("data", jsid)
		require.Equal(t, "{\"message\": \"My log\", \"name\": \"thunder\", \"id\": 234, \"data\": {\"id\":23}}", event.Message())
	})

	t.Run("with Entry fields", func(t *testing.T) {
		event := json.MessageObject("My log")
		event.String("name", "thunder")
		event.Int("id", 234)
		event.ObjectFor("data", func(event history.Encoder) {
			event.Int("id", 23)
		})
		require.Equal(t, "{\"message\": \"My log\", \"name\": \"thunder\", \"id\": 234, \"data\": {\"id\": 23}}", event.Message())
	})

	t.Run("with bytes fields", func(t *testing.T) {
		event := json.MessageObject("My log")
		event.String("name", "thunder")
		event.Int("id", 234)
		event.Bytes("data", []byte("{\"id\": 23}"))
		require.Equal(t, "{\"message\": \"My log\", \"name\": \"thunder\", \"id\": 234, \"data\": {\"id\": 23}}", event.Message())
	})

	t.Run("using context fields", func(t *testing.T) {
		event := json.MessageObjectWithEmbed("My log", "data", nil)
		event.String("name", "thunder")
		event.Int("id", 234)
		require.Equal(t, "{\"message\": \"My log\", \"data\": {\"name\": \"thunder\", \"id\": 234}}", event.Message())
	})

	t.Run("using context fields with hook", func(t *testing.T) {
		event := json.MessageObjectWithEmbed("My log", "data", func(event history.Encoder) {
			event.Bool("w", true)
		})

		event.String("name", "thunder")
		event.Int("id", 234)
		require.Equal(t, "{\"message\": \"My log\", \"w\": true, \"data\": {\"name\": \"thunder\", \"id\": 234}}", event.Message())
	})
}

func BenchmarkGetJSON(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	b.Run("with basic fields", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := b.N; i > 0; i-- {
			event := json.MessageObject("My log")
			event.String("name", "thunder")
			event.Int("id", 234)
			event.Message()
		}
	})

	b.Run("with Entry fields", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := b.N; i > 0; i-- {
			event := json.MessageObject("My log")
			event.String("name", "thunder")
			event.Int("id", 234)
			event.ObjectFor("data", func(event history.Encoder) {
				event.Int("id", 23)
			})
			event.Message()
		}
	})

	b.Run("with bytes fields", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		bu := []byte("{\"id\": 23}")
		for i := b.N; i > 0; i-- {
			event := json.MessageObject("My log")
			event.String("name", "thunder")
			event.Int("id", 234)
			event.Bytes("data", bu)
			event.Message()
		}
	})
}

func BenchmarkJSONAsEntries(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	b.Run("entries with no write", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := b.N; i > 0; i-- {
			var section = history.WithEntry(3)
			section.Int("id", 2)
			section.String("name", "beta")
			section.Bool("flag", true)
			section.ListFor("heights", func(encoder history.Encoder) {
				encoder.AddInt(20)
				encoder.AddInt(10)
				encoder.AddInt(50)

				encoder.AddObject(func(encoder history.Encoder) {
					section.String("name", "beta")
					section.Bool("flag", true)
				})
			})
		}
	})

	b.Run("entries with write", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var w = new(WritableBuffer)
		w.Data = make([]byte, 0, 512)

		for i := b.N; i > 0; i-- {
			var section = history.WithEntry(3)
			section.Int("id", 2)
			section.String("name", "beta")
			section.Bool("flag", true)
			section.WriteINFO(json.Object(), w)
		}

	})

	b.Run("prepared entries with write", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var w WritableBuffer
		w.Data = make([]byte, 0, 512)

		var section = history.WithEntry(3)
		section.Int("id", 2)
		section.String("name", "beta")
		section.Bool("flag", true)

		for i := b.N; i > 0; i-- {
			section.WriteINFO(json.Object(), &w)
		}

	})

	b.Run("prepared entries with error write", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var w WritableBuffer
		w.Data = make([]byte, 0, 512)

		var section = history.WithEntry(3)
		section.Int("id", 2)
		section.String("name", "beta")
		section.Bool("flag", true)

		for i := b.N; i > 0; i-- {
			section.WriteERROR(json.Object(), &w)
		}

	})

	b.Run("prepared entries with add and write", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var w WritableBuffer
		w.Data = make([]byte, 0, 512)

		var section = history.WithEntry(3)
		section.Int("id", 2)
		section.String("name", "beta")
		section.Bool("flag", true)

		for i := b.N; i > 0; i-- {
			section.With(func(encoder history.Encoder) {
				encoder.Int("day", 1)
			}).WriteINFO(json.Object(), &w)
		}
	})

	b.Run("prepared entries with add and error write", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var w WritableBuffer
		w.Data = make([]byte, 0, 512)

		var section = history.WithEntry(3)
		section.Int("id", 2)
		section.String("name", "beta")
		section.Bool("flag", true)

		for i := b.N; i > 0; i-- {
			section.With(func(encoder history.Encoder) {
				encoder.Int("day", 1)
			}).WriteERROR(json.Object(), &w)
		}
	})
}

func BenchmarkReplayJSONEncoding(b *testing.B) {
	b.Run("simple-replay", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var w = new(WritableBuffer)
		for i := 0; i < b.N; i++ {
			var section = history.Replayable()
			section = section.Int("id", 2)
			section = section.String("name", "beta")
			section = section.Bool("flag", true)
			section.WriteERROR(json.Object(), w)
		}
	})

	b.Run("prepared-replay", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var section = history.Replayable()
		section = section.Int("id", 2)
		section = section.String("name", "beta")
		section = section.Bool("flag", true)

		var w = new(WritableBuffer)
		for i := 0; i < b.N; i++ {
			section.WriteERROR(json.Object(), w)
		}
	})

	b.Run("multi-layer-replay", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		var w = new(WritableBuffer)
		for i := 0; i < b.N; i++ {
			var section = history.Replayable()
			section = section.Int("id", 2)
			section = section.String("name", "beta")
			section = section.Bool("flag", true)

			section = section.With(func(encoder history.Encoder) {
				encoder.Int("age", 22)
				encoder.Int("worker_id", 22)
				encoder.Int("rank_id", 22)

				encoder.ListFor("sets", func(encoder history.Encoder) {
					encoder.AddInt(1)
					encoder.AddInt(2)
					encoder.AddInt(3)
				})
			})

			section.WriteERROR(json.Object(), w)
		}
	})
}
