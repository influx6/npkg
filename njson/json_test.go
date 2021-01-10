package njson_test

import (
	gnjson "encoding/json"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/influx6/npkg"
	"github.com/influx6/npkg/njson"
)

type WritableBuffer struct {
	Data []byte
}

func (wb *WritableBuffer) Write(content []byte) (int, error) {
	wb.Data = append(wb.Data, content...)
	return len(content), nil
}

func TestGetJSON(t *testing.T) {
	t.Run("basic list", func(t *testing.T) {
		event := njson.JSONL()
		event.AddString("thunder")
		event.AddInt(234)
		require.Equal(t, "[\"thunder\", 234]", event.Message())
	})

	t.Run("basic list with bytes", func(t *testing.T) {
		event := njson.JSONL()
		event.AddString("thunder")
		event.AddInt(234)

		var jsid, err = gnjson.Marshal(map[string]interface{}{"id": 23})
		require.NoError(t, err)
		require.NotEmpty(t, jsid)

		event.AddBytes(jsid)

		require.Equal(t, "[\"thunder\", 234, {\"id\":23}]", event.Message())
	})

	t.Run("basic fields", func(t *testing.T) {
		event := njson.MJSON("My log")
		event.String("name", "thunder")
		event.Int("id", 234)
		require.Equal(t, "{\"message\": \"My log\", \"name\": \"thunder\", \"id\": 234}", event.Message())
	})

	t.Run("with Object fields", func(t *testing.T) {
		event := njson.MJSON("My log")
		event.String("name", "thunder")
		event.Int("id", 234)

		var jsid, err = gnjson.Marshal(map[string]interface{}{"id": 23})
		require.NoError(t, err)
		require.NotEmpty(t, jsid)

		event.Bytes("data", jsid)
		require.Equal(t, "{\"message\": \"My log\", \"name\": \"thunder\", \"id\": 234, \"data\": {\"id\":23}}", event.Message())
	})

	t.Run("with Entry fields", func(t *testing.T) {
		event := njson.MJSON("My log")
		event.String("name", "thunder")
		event.Int("id", 234)
		event.ObjectFor("data", func(event npkg.ObjectEncoder) {
			event.Int("id", 23)
		})
		require.Equal(t, "{\"message\": \"My log\", \"name\": \"thunder\", \"id\": 234, \"data\": {\"id\": 23}}", event.Message())
	})

	t.Run("with bytes fields", func(t *testing.T) {
		event := njson.MJSON("My log")
		event.String("name", "thunder")
		event.Int("id", 234)
		event.Bytes("data", []byte("{\"id\": 23}"))
		require.Equal(t, "{\"message\": \"My log\", \"name\": \"thunder\", \"id\": 234, \"data\": {\"id\": 23}}", event.Message())
	})
}

func BenchmarkGetJSON(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()

	b.Run("with basic fields", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := b.N; i > 0; i-- {
			event := njson.MJSON("My log")
			event.String("name", "thunder")
			event.Int("id", 234)
			event.Message()
		}
	})

	b.Run("with Entry fields", func(b *testing.B) {
		b.ResetTimer()
		b.ReportAllocs()

		for i := b.N; i > 0; i-- {
			event := njson.MJSON("My log")
			event.String("name", "thunder")
			event.Int("id", 234)
			event.ObjectFor("data", func(event npkg.ObjectEncoder) {
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
			event := njson.MJSON("My log")
			event.String("name", "thunder")
			event.Int("id", 234)
			event.Bytes("data", bu)
			event.Message()
		}
	})
}
