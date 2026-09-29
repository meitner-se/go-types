package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimestampUnmarshalText(t *testing.T) {
	t.Run("RFC3339", func(t *testing.T) {
		const raw = "2026-09-01T00:00:00Z"
		want, err := time.Parse(time.RFC3339, raw)
		require.NoError(t, err)

		var ts Timestamp
		require.NoError(t, ts.UnmarshalText([]byte(raw)))
		assert.True(t, ts.IsDefined())
		assert.False(t, ts.IsNil())
		assert.True(t, ts.Timestamp().Equal(want))
	})

	t.Run("RFC3339Nano", func(t *testing.T) {
		const raw = "2026-09-01T00:00:00.123456789Z"
		want, err := time.Parse(time.RFC3339Nano, raw)
		require.NoError(t, err)

		var ts Timestamp
		require.NoError(t, ts.UnmarshalText([]byte(raw)))
		assert.True(t, ts.IsDefined())
		assert.False(t, ts.IsNil())
		assert.True(t, ts.Timestamp().Equal(want))
	})

	t.Run("offset", func(t *testing.T) {
		const raw = "2026-09-01T02:00:00+02:00"
		var ts Timestamp
		require.NoError(t, ts.UnmarshalText([]byte(raw)))
		assert.True(t, ts.IsDefined())
		assert.False(t, ts.IsNil())
		assert.True(t, ts.Timestamp().Equal(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)))
	})

	t.Run("empty", func(t *testing.T) {
		var ts Timestamp
		require.NoError(t, ts.UnmarshalText([]byte{}))
		assert.False(t, ts.IsDefined())
		assert.True(t, ts.IsNil())
	})

	t.Run("malformed", func(t *testing.T) {
		var ts Timestamp
		err := ts.UnmarshalText([]byte("notatimestamp"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "notatimestamp")
		assert.False(t, ts.IsDefined(), "invalid input must not look defined")
		assert.True(t, ts.IsNil())
	})
}

func TestTimestampUnmarshalTextFlagParityWithJSON(t *testing.T) {
	const raw = "2026-09-01T00:00:00Z"

	var fromJSON Timestamp
	require.NoError(t, json.Unmarshal([]byte(`"`+raw+`"`), &fromJSON))

	var fromText Timestamp
	require.NoError(t, fromText.UnmarshalText([]byte(raw)))

	assert.Equal(t, fromJSON.IsDefined(), fromText.IsDefined())
	assert.Equal(t, fromJSON.IsNil(), fromText.IsNil())
	assert.True(t, fromJSON.Equal(fromText))

	var omitted Timestamp // JSON-omitted field stays at the zero value
	var emptyText Timestamp
	require.NoError(t, emptyText.UnmarshalText(nil))
	assert.Equal(t, omitted.IsDefined(), emptyText.IsDefined())
	assert.Equal(t, omitted.IsNil(), emptyText.IsNil())
}

func TestTimestampUnmarshalParam(t *testing.T) {
	const raw = "2026-09-01T00:00:00Z"
	var ts Timestamp
	require.NoError(t, ts.UnmarshalParam(raw))
	assert.True(t, ts.IsDefined())
	assert.False(t, ts.IsNil())
	assert.Equal(t, raw, ts.String())

	require.NoError(t, ts.UnmarshalParam(""))
	assert.False(t, ts.IsDefined())
	assert.True(t, ts.IsNil())

	err := ts.UnmarshalParam("notatimestamp")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "notatimestamp")
}

func TestTimestampMarshalText(t *testing.T) {
	ts := NewTimestamp(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))
	text, err := ts.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "2026-09-01T00:00:00Z", string(text))

	var roundtrip Timestamp
	require.NoError(t, roundtrip.UnmarshalText(text))
	assert.True(t, ts.Equal(roundtrip))

	empty, err := NewTimestampUndefined().MarshalText()
	require.NoError(t, err)
	assert.Empty(t, empty)
}

func TestDateUnmarshalText(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var d Date
		require.NoError(t, d.UnmarshalText([]byte("2026-09-01")))
		assert.True(t, d.IsDefined())
		assert.False(t, d.IsNil())
		assert.Equal(t, "2026-09-01", d.String())
	})

	t.Run("empty", func(t *testing.T) {
		var d Date
		require.NoError(t, d.UnmarshalText([]byte{}))
		assert.False(t, d.IsDefined())
		assert.True(t, d.IsNil())
	})

	t.Run("malformed", func(t *testing.T) {
		var d Date
		err := d.UnmarshalText([]byte("notadate"))
		require.Error(t, err)
		assert.Contains(t, err.Error(), "notadate")
		assert.False(t, d.IsDefined())
		assert.True(t, d.IsNil())
	})
}

func TestDateUnmarshalTextFlagParityWithJSON(t *testing.T) {
	var fromJSON Date
	require.NoError(t, json.Unmarshal([]byte(`"2026-09-01"`), &fromJSON))

	var fromText Date
	require.NoError(t, fromText.UnmarshalText([]byte("2026-09-01")))

	assert.Equal(t, fromJSON.IsDefined(), fromText.IsDefined())
	assert.Equal(t, fromJSON.IsNil(), fromText.IsNil())
	assert.Equal(t, fromJSON.String(), fromText.String())
}

func TestDateUnmarshalParam(t *testing.T) {
	var d Date
	require.NoError(t, d.UnmarshalParam("2026-09-01"))
	assert.True(t, d.IsDefined())
	assert.Equal(t, "2026-09-01", d.String())

	require.NoError(t, d.UnmarshalParam(""))
	assert.False(t, d.IsDefined())
	assert.True(t, d.IsNil())

	err := d.UnmarshalParam("notadate")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "notadate")
}

func TestDateMarshalText(t *testing.T) {
	d := NewDate(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC))
	text, err := d.MarshalText()
	require.NoError(t, err)
	assert.Equal(t, "2026-09-01", string(text))

	empty, err := NewDateUndefined().MarshalText()
	require.NoError(t, err)
	assert.Empty(t, empty)
}
