package jsonstream

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
)

func Ingest[T any](ctx context.Context, data io.Reader, doer func(key string, chunk *T) error) error {
	decoder := json.NewDecoder(data)

	if _, err := expectToken(decoder, "{"); err != nil {
		return fmt.Errorf("unexpected token: %w", err)
	}

	for decoder.More() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		key, err := expectToken(decoder, "")
		if err != nil {
			break
		}

		var chunk T
		if err := decoder.Decode(&chunk); err != nil {
			return fmt.Errorf("failed to decode chunk: %w", err)
		}

		if err := doer(key, &chunk); err != nil {
			return fmt.Errorf("failed to process chunk: %w", err)
		}
	}

	if _, err := expectToken(decoder, "}"); err != nil {
		return fmt.Errorf("unexpected token: %w", err)
	}

	return nil
}

func expectToken(decoder *json.Decoder, val string) (string, error) {
	token, err := decoder.Token()
	if err != nil {
		return "", fmt.Errorf("failed to get token: %w", err)
	}

	switch tt := token.(type) {
	case json.Delim:
		if val != "" && tt.String() != val {
			return "", fmt.Errorf("invalid token, expected %s got: %s", val, tt.String())
		}

		return tt.String(), nil
	case string:
		if val != "" && tt != val {
			return "", fmt.Errorf("invalid token, expected %s got: %s", val, tt)
		}

		return tt, nil
	default:
		return "", fmt.Errorf("invalid token type")
	}
}
