package ntrees

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDOMAttributeWriter(t *testing.T) {
	var attrs = NewDOMAttrEncoder("")
	require.NoError(t, attrs.QuotedString("media", "img/svg"))
	require.NoError(t, attrs.QuotedString("serial", "223-232-uyl"))
	require.NoError(t, attrs.Int("counter", 20))
	require.NoError(t, attrs.Float("rate", 20.23))

	require.Equal(t, `media="img/svg" serial="223-232-uyl" counter=20 rate=20.23`, attrs.String())
}

func TestDOMAttributeWriterWithList(t *testing.T) {
	var attrs = NewDOMAttrEncoder("")
	require.NoError(t, attrs.QuotedString("media", "img/svg"))
	require.NoError(t, attrs.QuotedString("serial", "223-232-uyl"))
	require.NoError(t, attrs.Int("counter", 20))
	require.NoError(t, attrs.Float("rate", 20.23))

	require.NoError(t, attrs.List("records", "1", "2", "3", "4"))
	require.Equal(t, `media="img/svg" serial="223-232-uyl" counter=20 rate=20.23 records="1,2,3,4"`, attrs.String())
}

func TestDOMAttributeWriterWithListAndSubAttrs(t *testing.T) {
	var attrs = NewDOMAttrEncoder("")
	require.NoError(t, attrs.QuotedString("media", "img/svg"))
	require.NoError(t, attrs.QuotedString("serial", "223-232-uyl"))
	require.NoError(t, attrs.Int("counter", 20))
	require.NoError(t, attrs.Float("rate", 20.23))

	require.NoError(t, attrs.WithAttr("flags", func(encoder AttrEncoder) error {
		if err := encoder.QuotedString("bb", "1"); err != nil {
			return err
		}
		return encoder.QuotedString("bc", "0")
	}))

	require.NoError(t, attrs.List("records", "1", "2", "3", "4"))
	require.Equal(t, `media="img/svg" serial="223-232-uyl" counter=20 rate=20.23 flags.bb="1" flags.bc="0" records="1,2,3,4"`, attrs.String())
}
