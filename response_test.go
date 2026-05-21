package fetch

import (
	"bytes"
	"testing"
)

func TestNewResponse_RoundTrip(t *testing.T) {
	hdrs := []Header{
		{Key: "Content-Type", Value: "application/json"},
		{Key: "X-Custom", Value: "123"},
	}
	body := []byte(`{"hello":"world"}`)

	resp := NewResponse(200, hdrs, body)

	if resp.Status != 200 {
		t.Errorf("Expected status 200, got %d", resp.Status)
	}

	if !bytes.Equal(resp.Body(), body) {
		t.Errorf("Expected body %s, got %s", body, resp.Body())
	}

	if resp.GetHeader("Content-Type") != "application/json" {
		t.Errorf("Expected header application/json, got %s", resp.GetHeader("Content-Type"))
	}
	
	if resp.GetHeader("x-custom") != "123" {
		t.Errorf("Expected header 123, got %s", resp.GetHeader("x-custom"))
	}
}

func TestNewResponse_EmptyBody(t *testing.T) {
	resp := NewResponse(204, nil, nil)

	if resp.Status != 204 {
		t.Errorf("Expected status 204, got %d", resp.Status)
	}

	b := resp.Body()
	if len(b) != 0 {
		t.Errorf("Expected empty body, got length %d", len(b))
	}
}
