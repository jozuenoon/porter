package jsonstream_test

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"porter"
	"porter/pkg/jsonstream"
	"porter/testdata"
)

func TestIngestor(t *testing.T) {
	t.Run("ingestor single", func(t *testing.T) {
		stream := toStream(single)

		var out []*porter.Port

		var okey string

		err := jsonstream.Ingest(context.Background(), stream, func(key string, port *porter.Port) error {
			okey = key

			out = append(out, port)

			return nil
		})
		require.NoError(t, err)

		require.Len(t, out, 1)

		require.Equal(t, "PLGDN", okey)
		require.Equal(t, "Gdansk", out[0].Name)
		require.Equal(t, "Gdansk", out[0].City)
		require.Equal(t, "Poland", out[0].Country)
		require.Equal(t, "Pomeranian Voivodeship", out[0].Province)
		require.Equal(t, "Europe/Warsaw", out[0].Timezone)
		require.Equal(t, "PLGDN", out[0].Unlocs[0])
		require.Equal(t, "45511", out[0].Code)

		c0, err := decimal.NewFromString("18.6466384")
		require.NoError(t, err)

		c1, err := decimal.NewFromString("54.35202520000001")
		require.NoError(t, err)

		require.True(t, c0.Equal(out[0].Coordinates[0]))
		require.True(t, c1.Equal(out[0].Coordinates[1]))
	})

	t.Run("ingestor 4 ports", func(t *testing.T) {
		stream := toStream(ports4)

		var out []*porter.Port

		err := jsonstream.Ingest(context.Background(), stream, func(_ string, port *porter.Port) error {
			out = append(out, port)

			return nil
		})
		require.NoError(t, err)

		require.Len(t, out, 4)

		require.Equal(t, "Gdansk", out[0].Name)
		require.Equal(t, "Gdynia", out[1].Name)
		require.Equal(t, "Swinoujscie", out[2].Name)
		require.Equal(t, "Szczecin", out[3].Name)
	})

	t.Run("ingest all", func(t *testing.T) {
		stream := toStream(testdata.Ports)

		var out []*porter.Port

		err := jsonstream.Ingest(context.Background(), stream, func(_ string, port *porter.Port) error {
			out = append(out, port)

			return nil
		})
		require.NoError(t, err)

		require.Len(t, out, 1632)
	})

	t.Run("ingest empty", func(t *testing.T) {
		stream := toStream("{}")

		var out []*porter.Port

		err := jsonstream.Ingest(context.Background(), stream, func(_ string, port *porter.Port) error {
			out = append(out, port)

			return nil
		})
		require.NoError(t, err)

		require.Len(t, out, 0)
	})

	t.Run("ingest wrong structure", func(t *testing.T) {
		stream := toStream(`{"config": [1,2,3,4]}`)

		var out []*porter.Port

		err := jsonstream.Ingest(context.Background(), stream, func(_ string, port *porter.Port) error {
			out = append(out, port)

			return nil
		})
		require.Error(t, err)
		require.Len(t, out, 0)
	})
}

func BenchmarkIngestAll(b *testing.B) {
	stream := toStream(testdata.Ports)

	for i := 0; i < b.N; i++ {
		err := jsonstream.Ingest(context.Background(), stream, func(_ string, port *porter.Port) error {
			return nil
		})
		require.NoError(b, err)

		_, _ = stream.Seek(0, io.SeekStart)
	}
}

func toStream(data string) io.ReadSeeker {
	return bytes.NewReader([]byte(data))
}

const single = `
{
"PLGDN": {
  "name": "Gdansk",
  "city": "Gdansk",
  "country": "Poland",
  "alias": [],
  "regions": [],
  "coordinates": [
    18.6466384,
    54.35202520000001
  ],
  "province": "Pomeranian Voivodeship",
  "timezone": "Europe/Warsaw",
  "unlocs": [
    "PLGDN"
  ],
  "code": "45511"
}
}
`

const ports4 = `
{
"PLGDN": {
  "name": "Gdansk",
  "city": "Gdansk",
  "country": "Poland",
  "alias": [],
  "regions": [],
  "coordinates": [
    18.6466384,
    54.35202520000001
  ],
  "province": "Pomeranian Voivodeship",
  "timezone": "Europe/Warsaw",
  "unlocs": [
    "PLGDN"
  ],
  "code": "45511"
},
"PLGDY": {
  "name": "Gdynia",
  "coordinates": [
    18.55,
    54.5
  ],
  "city": "Gdynia",
  "province": "Pomorskie",
  "country": "Poland",
  "alias": [],
  "regions": [],
  "timezone": "Europe/Warsaw",
  "unlocs": [
    "PLGDY"
  ],
  "code": "45505"
},
"PLSWI": {
  "name": "Swinoujscie",
  "coordinates": [
    14.25,
    53.9
  ],
  "city": "Swinoujscie",
  "country": "Poland",
  "alias": [],
  "regions": [],
  "province": "West Pomeranian Voivodeship",
  "timezone": "Europe/Warsaw",
  "unlocs": [
    "PLSWI"
  ],
  "code": "45512"
},
"PLSZZ": {
  "name": "Szczecin",
  "city": "Szczecin",
  "country": "Poland",
  "alias": [],
  "regions": [],
  "coordinates": [
    14.5528116,
    53.4285438
  ],
  "province": "West Pomeranian Voivodeship",
  "timezone": "Europe/Warsaw",
  "unlocs": [
    "PLSZZ"
  ],
  "code": "45507"
}
}
`
