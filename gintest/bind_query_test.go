package gintest

import (
	"encoding"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/meitner-se/go-types"
)

// timestampQuery is the publicapis-shaped query struct from INF-9812.
type timestampQuery struct {
	ChangedAfter types.Timestamp `form:"changedAfter"`
}

type dateQuery struct {
	On types.Date `form:"on"`
}

// uuidQuery is the publicapis-shaped query struct from INF-9831.
type uuidQuery struct {
	SchoolID types.UUID `form:"schoolID"`
}

func bindQuery(t *testing.T, rawQuery string, dest any) error {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := httptest.NewRequest(http.MethodGet, "/?"+rawQuery, nil)
	c.Request = req
	return c.ShouldBindQuery(dest)
}

// TestTimestampShouldBindQueryBareRFC3339 is the INF-9812 reproduction:
// publicapis calls gin's ShouldBindQuery on a struct field of type types.Timestamp.
// Partners send a bare RFC3339 value (?changedAfter=2026-09-01T00:00:00Z), not a
// JSON-quoted string.
func TestTimestampShouldBindQueryBareRFC3339(t *testing.T) {
	const raw = "2026-09-01T00:00:00Z"
	want, err := time.Parse(time.RFC3339, raw)
	require.NoError(t, err)

	var q timestampQuery
	err = bindQuery(t, "changedAfter="+raw, &q)
	require.NoError(t, err, "gin %s ShouldBindQuery of bare RFC3339 into types.Timestamp", gin.Version)
	assert.True(t, q.ChangedAfter.IsDefined())
	assert.False(t, q.ChangedAfter.IsNil())
	assert.True(t, q.ChangedAfter.Timestamp().Equal(want))
}

func TestTimestampShouldBindQueryRFC3339Nano(t *testing.T) {
	const raw = "2026-09-01T00:00:00.123456789Z"
	want, err := time.Parse(time.RFC3339Nano, raw)
	require.NoError(t, err)

	var q timestampQuery
	err = bindQuery(t, "changedAfter="+raw, &q)
	require.NoError(t, err)
	assert.True(t, q.ChangedAfter.IsDefined())
	assert.False(t, q.ChangedAfter.IsNil())
	assert.True(t, q.ChangedAfter.Timestamp().Equal(want))
}

func TestTimestampShouldBindQueryOffset(t *testing.T) {
	const raw = "2026-09-01T02:00:00+02:00"
	want, err := time.Parse(time.RFC3339, raw)
	require.NoError(t, err)

	var q timestampQuery
	// '+' is a space in query strings; partners must send %2B for an offset.
	err = bindQuery(t, "changedAfter=2026-09-01T02:00:00%2B02:00", &q)
	require.NoError(t, err)
	assert.True(t, q.ChangedAfter.IsDefined())
	assert.False(t, q.ChangedAfter.IsNil())
	assert.True(t, q.ChangedAfter.Timestamp().Equal(want))
	assert.True(t, q.ChangedAfter.Timestamp().Equal(time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)))
}

func TestTimestampShouldBindQueryEmptyAndAbsent(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var q timestampQuery
		err := bindQuery(t, "changedAfter=", &q)
		require.NoError(t, err)
		assert.False(t, q.ChangedAfter.IsDefined())
		assert.True(t, q.ChangedAfter.IsNil())
	})

	t.Run("absent", func(t *testing.T) {
		var q timestampQuery
		err := bindQuery(t, "", &q)
		require.NoError(t, err)
		assert.False(t, q.ChangedAfter.IsDefined())
		assert.True(t, q.ChangedAfter.IsNil())
	})
}

func TestTimestampShouldBindQueryMalformed(t *testing.T) {
	var q timestampQuery
	err := bindQuery(t, "changedAfter=notatimestamp", &q)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "notatimestamp")
}

func TestTimestampQueryFlagParityWithJSON(t *testing.T) {
	const raw = "2026-09-01T00:00:00Z"

	var fromJSON types.Timestamp
	require.NoError(t, json.Unmarshal([]byte(`"`+raw+`"`), &fromJSON))

	var fromQuery timestampQuery
	require.NoError(t, bindQuery(t, "changedAfter="+raw, &fromQuery))

	assert.Equal(t, fromJSON.IsDefined(), fromQuery.ChangedAfter.IsDefined())
	assert.Equal(t, fromJSON.IsNil(), fromQuery.ChangedAfter.IsNil())
	assert.True(t, fromJSON.Equal(fromQuery.ChangedAfter))
}

func TestDateShouldBindQuery(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		var q dateQuery
		err := bindQuery(t, "on=2026-09-01", &q)
		require.NoError(t, err)
		assert.True(t, q.On.IsDefined())
		assert.False(t, q.On.IsNil())
		assert.Equal(t, "2026-09-01", q.On.String())
	})

	t.Run("empty", func(t *testing.T) {
		var q dateQuery
		err := bindQuery(t, "on=", &q)
		require.NoError(t, err)
		assert.False(t, q.On.IsDefined())
		assert.True(t, q.On.IsNil())
	})

	t.Run("malformed", func(t *testing.T) {
		var q dateQuery
		err := bindQuery(t, "on=notadate", &q)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "notadate")
	})

	t.Run("flag parity with JSON", func(t *testing.T) {
		var fromJSON types.Date
		require.NoError(t, json.Unmarshal([]byte(`"2026-09-01"`), &fromJSON))

		var fromQuery dateQuery
		require.NoError(t, bindQuery(t, "on=2026-09-01", &fromQuery))

		assert.Equal(t, fromJSON.IsDefined(), fromQuery.On.IsDefined())
		assert.Equal(t, fromJSON.IsNil(), fromQuery.On.IsNil())
		assert.Equal(t, fromJSON.String(), fromQuery.On.String())
	})
}

// TestUUIDShouldBindQueryBareCanonical is the INF-9831 reproduction:
// publicapis calls gin's ShouldBindQuery on a struct field of type types.UUID.
// Partners send a bare hyphenated UUID (?schoolID=550e8400-e29b-41d4-a716-446655440000),
// not a JSON-quoted string.
func TestUUIDShouldBindQueryBareCanonical(t *testing.T) {
	const raw = "550e8400-e29b-41d4-a716-446655440000"
	want := uuid.MustParse(raw)

	var q uuidQuery
	err := bindQuery(t, "schoolID="+raw, &q)
	require.NoError(t, err, "gin %s ShouldBindQuery of bare UUID into types.UUID", gin.Version)
	assert.True(t, q.SchoolID.IsDefined())
	assert.False(t, q.SchoolID.IsNil())
	assert.Equal(t, want, q.SchoolID.UUID())
}

func TestUUIDShouldBindQueryEmptyAndAbsent(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		var q uuidQuery
		err := bindQuery(t, "schoolID=", &q)
		require.NoError(t, err)
		assert.False(t, q.SchoolID.IsDefined())
		assert.True(t, q.SchoolID.IsNil())
	})

	t.Run("absent", func(t *testing.T) {
		var q uuidQuery
		err := bindQuery(t, "", &q)
		require.NoError(t, err)
		assert.False(t, q.SchoolID.IsDefined())
		assert.True(t, q.SchoolID.IsNil())
	})
}

func TestUUIDShouldBindQueryMalformed(t *testing.T) {
	var q uuidQuery
	err := bindQuery(t, "schoolID=not-a-uuid", &q)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid UUID")
}

func TestUUIDQueryFlagParityWithJSON(t *testing.T) {
	const raw = "550e8400-e29b-41d4-a716-446655440000"

	var fromJSON types.UUID
	require.NoError(t, json.Unmarshal([]byte(`"`+raw+`"`), &fromJSON))

	var fromQuery uuidQuery
	require.NoError(t, bindQuery(t, "schoolID="+raw, &fromQuery))

	assert.Equal(t, fromJSON.IsDefined(), fromQuery.SchoolID.IsDefined())
	assert.Equal(t, fromJSON.IsNil(), fromQuery.SchoolID.IsNil())
	assert.Equal(t, fromJSON.UUID(), fromQuery.SchoolID.UUID())
}

// TestGinV1120StructKindBindingPath records which interfaces gin v1.12.0
// actually consults for struct-kind fields. This decides the production
// implementation: TextUnmarshaler is only used when the form tag names
// parser=encoding.TextUnmarshaler; the default path is BindUnmarshaler
// then json.Unmarshal of the raw query value.
func TestGinV1120StructKindBindingPath(t *testing.T) {
	require.Equal(t, "v1.12.0", gin.Version)

	t.Run("default path does not call TextUnmarshaler", func(t *testing.T) {
		var dest struct {
			Value textOnly `form:"value"`
		}
		err := bindQuery(t, "value=hello", &dest)
		require.Error(t, err, "gin v1.12.0 should fall through to json.Unmarshal for a TextUnmarshaler-only struct")
		assert.False(t, dest.Value.called)
	})

	t.Run("parser tag calls TextUnmarshaler", func(t *testing.T) {
		var dest struct {
			Value textOnly `form:"value,parser=encoding.TextUnmarshaler"`
		}
		err := bindQuery(t, "value=hello", &dest)
		require.NoError(t, err)
		assert.True(t, dest.Value.called)
		assert.Equal(t, "hello", dest.Value.got)
	})

	t.Run("default path calls BindUnmarshaler", func(t *testing.T) {
		var dest struct {
			Value paramOnly `form:"value"`
		}
		err := bindQuery(t, "value=hello", &dest)
		require.NoError(t, err)
		assert.True(t, dest.Value.called)
		assert.Equal(t, "hello", dest.Value.got)
	})
}

type textOnly struct {
	called bool
	got    string
}

func (t *textOnly) UnmarshalText(text []byte) error {
	t.called = true
	t.got = string(text)
	return nil
}

type paramOnly struct {
	called bool
	got    string
}

func (p *paramOnly) UnmarshalParam(param string) error {
	p.called = true
	p.got = param
	return nil
}

var (
	_ encoding.TextUnmarshaler = (*textOnly)(nil)
	_ binding.BindUnmarshaler  = (*paramOnly)(nil)
)

// TestQueryParamTypesAudit notes which go-types already bind from a bare
// query string through gin v1.12.0's json.Unmarshal fallback (values that
// happen to be valid JSON) and which do not.
func TestQueryParamTypesAudit(t *testing.T) {
	t.Run("Int64 limit=100", func(t *testing.T) {
		var q struct {
			Limit types.Int64 `form:"limit"`
		}
		err := bindQuery(t, "limit=100", &q)
		require.NoError(t, err)
		assert.True(t, q.Limit.IsDefined())
		assert.False(t, q.Limit.IsNil())
		assert.Equal(t, int64(100), q.Limit.Int64())
	})

	t.Run("Bool flag=true", func(t *testing.T) {
		var q struct {
			Flag types.Bool `form:"flag"`
		}
		err := bindQuery(t, "flag=true", &q)
		require.NoError(t, err)
		assert.True(t, q.Flag.IsDefined())
		assert.False(t, q.Flag.IsNil())
		assert.True(t, q.Flag.Bool())
	})

	t.Run("Float64 x=1.5", func(t *testing.T) {
		var q struct {
			X types.Float64 `form:"x"`
		}
		err := bindQuery(t, "x=1.5", &q)
		require.NoError(t, err)
		assert.True(t, q.X.IsDefined())
		assert.False(t, q.X.IsNil())
		assert.Equal(t, 1.5, q.X.Float64())
	})

	t.Run("UUID bare hyphenated", func(t *testing.T) {
		id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		var q struct {
			ID types.UUID `form:"id"`
		}
		err := bindQuery(t, "id="+id.String(), &q)
		require.NoError(t, err)
		assert.True(t, q.ID.IsDefined())
		assert.False(t, q.ID.IsNil())
		assert.Equal(t, id, q.ID.UUID())
	})
}
