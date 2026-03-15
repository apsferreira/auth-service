package messaging

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// OTPRequestedEvent — JSON marshaling / unmarshaling
//
// Publisher.PublishOTPRequested wraps an OTPRequestedEvent and serialises it
// to JSON before placing it on the wire. These tests verify that the struct
// round-trips correctly without requiring a live RabbitMQ connection.
// ---------------------------------------------------------------------------

func TestOTPRequestedEvent_MarshalJSON(t *testing.T) {
	event := OTPRequestedEvent{
		Email:       "user@example.com",
		ServiceName: "Libri",
	}

	data, err := json.Marshal(event)
	require.NoError(t, err)

	assert.Contains(t, string(data), `"email":"user@example.com"`)
	assert.Contains(t, string(data), `"service_name":"Libri"`)
}

func TestOTPRequestedEvent_UnmarshalJSON(t *testing.T) {
	raw := `{"email":"ana@example.com","service_name":"Nitro"}`

	var event OTPRequestedEvent
	err := json.Unmarshal([]byte(raw), &event)

	require.NoError(t, err)
	assert.Equal(t, "ana@example.com", event.Email)
	assert.Equal(t, "Nitro", event.ServiceName)
}

func TestOTPRequestedEvent_RoundTrip(t *testing.T) {
	original := OTPRequestedEvent{
		Email:       "roundtrip@example.com",
		ServiceName: "Focus Hub",
	}

	data, err := json.Marshal(original)
	require.NoError(t, err)

	var decoded OTPRequestedEvent
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, original.Email, decoded.Email)
	assert.Equal(t, original.ServiceName, decoded.ServiceName)
}

func TestOTPRequestedEvent_EmptyServiceName(t *testing.T) {
	// Marshaling with an empty ServiceName must still produce valid JSON —
	// the consumer on the other end is responsible for any validation.
	event := OTPRequestedEvent{Email: "x@x.com", ServiceName: ""}

	data, err := json.Marshal(event)
	require.NoError(t, err)

	var decoded OTPRequestedEvent
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, "", decoded.ServiceName)
}

func TestOTPRequestedEvent_UnknownFieldsIgnored(t *testing.T) {
	// Extra fields in the JSON payload must not cause an error (forward
	// compatibility: other services may add fields over time).
	raw := `{"email":"compat@example.com","service_name":"Libri","extra_field":"ignored"}`

	var event OTPRequestedEvent
	err := json.Unmarshal([]byte(raw), &event)

	require.NoError(t, err)
	assert.Equal(t, "compat@example.com", event.Email)
	assert.Equal(t, "Libri", event.ServiceName)
}

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

func TestRoutingKeyConstants(t *testing.T) {
	assert.Equal(t, "otp.requested", RoutingKeyOTPRequested)
	assert.Equal(t, "otp.verified", RoutingKeyOTPVerified)
}

// ---------------------------------------------------------------------------
// NewPublisher — connection failure handling
//
// With no live broker available, NewPublisher must return a non-nil error
// rather than hanging indefinitely. The retry loop in connect() will exhaust
// maxRetries quickly because amqp.Dial returns immediately for an invalid URL.
// ---------------------------------------------------------------------------

func TestNewPublisher_InvalidURL_ReturnsError(t *testing.T) {
	// Use an unreachable URL — amqp.Dial fails fast with an invalid address.
	p, err := NewPublisher("amqp://invalid-host-that-does-not-exist:5672/")

	assert.Nil(t, p)
	assert.Error(t, err)
}
