package types

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//nolint:lll
func TestRichText(t *testing.T) {
	t.Run("Unmarshal", func(t *testing.T) {
		tt := []struct {
			content, expected string
		}{
			{
				content:  "<p>Paragraf</p><h1>Stor rubrik</h1><h2>Medium rubrik</h2><h3>Liten rubrik</h3><ul><li><p>punkt</p></li><li><p>lista</p></li></ul><ol><li><p>numrerad </p></li><li><p>lista</p><pre><code>Kodblock</code></pre></li></ol><p><strong>FET </strong><em>KURSIV </em> <u>UNDERSTRYKNING</u> <s>GENOMSTRUKEN</s> <code>KODD</code> <a target=\\\"_blank\\\" rel=\\\"noopener noreferrer nofollow\\\" href=\\\"https://google.com\\\">länk till google</a></p><p></p><img src=\\\"https://fileserver.develop.meitner.se/v1/file/efe378e5-e263-438c-841c-07ab20c60bc0.png\\\">",
				expected: "Paragraf\n\nStor rubrik\n\nMedium rubrik\n\nLiten rubrik\n\npunkt\n\nlista\n\nnumrerad \n\nlista\n\nKodblock\n\nFET KURSIV  UNDERSTRYKNING GENOMSTRUKEN KODD länk till google\n\n",
			},
			{
				content:  "<ul><li><p>punkt</p></li><li><p>lista</p></li></ul>",
				expected: "punkt\n\nlista",
			},
			{
				content:  "<p>hej</p><p>på dig</p><p></p>",
				expected: "hej\n\npå dig",
			},
		}

		for _, tc := range tt {
			t.Run(tc.content, func(t *testing.T) {
				input := fmt.Sprintf(`{"content":"%s"}`, tc.content)

				var output RichText
				err := json.Unmarshal([]byte(input), &output)
				require.NoError(t, err)

				text, err := output.Text()
				require.NoError(t, err)

				assert.Equal(t, tc.expected, text)
			})
		}
	})
}

func TestGob(t *testing.T) {
	person := struct {
		FirstName String
		LastName  String
	}{
		FirstName: NewString("John"),
		LastName:  NewString("Doe"),
	}

	var buf bytes.Buffer

	err := gob.NewEncoder(&buf).Encode(person)
	require.NoError(t, err)

	var decoded struct {
		FirstName String
		LastName  String
	}

	err = gob.NewDecoder(&buf).Decode(&decoded)
	require.NoError(t, err)

	assert.Equal(t, person, decoded)
}

func TestTimestamp(t *testing.T) {
	t.Run("StartOfDay", func(t *testing.T) {
		currentTime := time.Now().UTC()

		tz := currentTime.Location()

		timestamp := NewTimestamp(currentTime).StartOfDay(tz).Timestamp()

		assert.Equal(t, currentTime.Year(), timestamp.Year())
		assert.Equal(t, currentTime.Month(), timestamp.Month())
		assert.Equal(t, currentTime.Day(), timestamp.Day())
		assert.Equal(t, 0, timestamp.Hour())
		assert.Equal(t, 0, timestamp.Minute())
		assert.Equal(t, 0, timestamp.Second())
		assert.Equal(t, 0, timestamp.Nanosecond())
	})

	t.Run("EndOfDay", func(t *testing.T) {
		currentTime := time.Now().UTC()

		tz := currentTime.Location()

		timestamp := NewTimestamp(currentTime).EndOfDay(tz).Timestamp()

		assert.Equal(t, currentTime.Year(), timestamp.Year())
		assert.Equal(t, currentTime.Month(), timestamp.Month())
		assert.Equal(t, currentTime.Day(), timestamp.Day())
		assert.Equal(t, 23, timestamp.Hour())
		assert.Equal(t, 59, timestamp.Minute())
		assert.Equal(t, 59, timestamp.Second())
		assert.Equal(t, 0, timestamp.Nanosecond())
	})
}
