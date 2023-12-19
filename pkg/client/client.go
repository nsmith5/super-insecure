package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func Register(ctx context.Context, username string) error {
	resp, err := http.Post(
		`http://localhost:8080/register`,
		`application/json`,
		strings.NewReader(fmt.Sprintf(`{"username": %q}`, username)),
	)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}

	return nil
}

func ItemGet(ctx context.Context, username, item string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, `http://localhost:8080/items/`+item, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Insecure "+username)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", errors.New(string(body))
	}

	var value struct {
		Value string
	}
	err = json.NewDecoder(resp.Body).Decode(&value)
	if err != nil {
		return "", err
	}

	return value.Value, nil
}

func ItemSet(ctx context.Context, username, item, value string) error {
	body := fmt.Sprintf(`{"value": %q}`, value)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		`http://localhost:8080/items/`+item,
		strings.NewReader(body),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Insecure "+username)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}

	return nil
}

func ItemDelete(ctx context.Context, username, item string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, `http://localhost:8080/items/`+item, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Insecure "+username)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}

	return nil
}
