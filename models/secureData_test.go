package models

import (
	"bytes"
	"crypto/rand"
	"testing"

	"golang.org/x/crypto/nacl/box"
)

func mustKey(t *testing.T) (privkey, pubkey *[32]byte, valid bool) {
	t.Helper()
	var err error
	privkey, pubkey, err = box.GenerateKey(rand.Reader)
	if err != nil {
		t.Errorf("Error generating keys: %v", err)
		valid = false
	}
	return
}

func TestSecureData(t *testing.T) {
	ourPrivateKey, ourPublicKey, valid := mustKey(t)
	if !valid {
		return
	}
	payload := []byte("Hello, World")
	msg := &SecureData{}
	if err := msg.Seal(*ourPublicKey, payload); err != nil {
		t.Errorf("%v", err)
		return
	} else {
		t.Logf("Sealed message `Hello, World`")
	}
	out, err := msg.Open(*ourPrivateKey)
	if err != nil {
		t.Errorf("%v", err)
		return
	} else {
		t.Logf("Opened sealed message")
	}
	if !bytes.Equal(out, payload) {
		t.Errorf("Expected %s, got %s", string(payload), string(out))
	} else {
		t.Logf("Sealed message intact")
	}
	msg.Nonce[0] = msg.Nonce[0] ^ byte(0)
	_, err = msg.Open(*ourPrivateKey)
	if err != Corrupt {
		t.Errorf("corruptViaNonce: Expected error %v, not %v", Corrupt, err)
	} else {
		t.Logf("corruptViaNonce: Got expected error %v", err)
	}
	msg.Nonce[0] = msg.Nonce[0] ^ byte(0)
	msg.Key[0] = msg.Key[0] ^ byte(0)
	_, err = msg.Open(*ourPrivateKey)
	if err != Corrupt {
		t.Errorf("corruptViaKey: Expected error %v, not %v", Corrupt, err)
	} else {
		t.Logf("corruptViaKey: Got expected error %v", err)
	}
	msg.Key[0] = msg.Key[0] ^ byte(0)
	msg.Payload[0] = msg.Payload[0] ^ byte(0)
	_, err = msg.Open(*ourPrivateKey)
	if err != Corrupt {
		t.Errorf("corruptViaPayload: Expected error %v, not %v", Corrupt, err)
	} else {
		t.Logf("corruptViaPayload: Got expected error %v", err)
	}
	msg.Payload[0] = msg.Payload[0] ^ byte(0)
	msg.Payload = msg.Payload[1:]
	_, err = msg.Open(*ourPrivateKey)
	if err != Corrupt {
		t.Errorf("corruptViaPayload: Expected error %v, not %v", Corrupt, err)
	} else {
		t.Logf("corruptViaPayload: Got expected error %v", err)
	}
	var nonce []byte
	nonce, msg.Nonce = msg.Nonce, msg.Nonce[1:]
	_, err = msg.Open(*ourPrivateKey)
	if err != BadNonce {
		t.Errorf("badNonce: Expected error %v, not %v", BadNonce, err)
	} else {
		t.Logf("badNonce: Got expected error %v", err)
	}
	msg.Nonce, msg.Key = nonce, msg.Key[1:]
	_, err = msg.Open(*ourPrivateKey)
	if err != BadKey {
		t.Errorf("badKey: Expected error %v, not %v", BadKey, err)
	} else {
		t.Logf("badKey: Got expected error %v", err)
	}
}
