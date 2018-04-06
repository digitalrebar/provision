package models

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/box"
)

// SecureData is used to store and send access controlled Param values
// to the locations they are needed at.  SecureData uses a simple
// encryption mechanism based on the NACL Box API (as implemented by
// libsodium, golang.org/x/crypto/nacl/box, and many others)
type SecureData struct {
	// Key is a 32 byte long curve25519 public key.
	Key []byte
	// Nonce must be 24 bytes of cryptographically random numbers
	Nonce []byte
	// Payload is the encrypted payload.
	Payload []byte
}

var (
	BadKey   = errors.New("Key must be 32 bytes long")
	BadNonce = errors.New("Nonce must be 24 bytes long")
	Corrupt  = errors.New("SecureData corrupted")
)

// Validate makes sure that the lengths we expect for the Key and
// Nonce are correct
func (s *SecureData) Validate() error {
	if len(s.Key) != 32 {
		return BadKey
	}
	if len(s.Nonce) != 24 {
		return BadNonce
	}
	return nil
}

// Seal takes curve25519 public key advertised by where the payload
// should be stored, and fills in the SecureData with the data
// required for the Open operation to succeed.
func (s *SecureData) Seal(peerPublicKey [32]byte, data []byte) error {
	ourPublicKey, ourPrivateKey, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return fmt.Errorf("Error generating ephemeral local keys: %v", err)
	}
	s.Key = ourPublicKey[:]
	nonce := [24]byte{}
	_, err = io.ReadFull(rand.Reader, nonce[:])
	if err != nil {
		return fmt.Errorf("Error generating nonce: %v", err)
	}
	s.Nonce = nonce[:]
	box.Seal(s.Payload, data, &nonce, &peerPublicKey, ourPrivateKey)
	return nil
}

// Open opens a sealed SecureData item.
func (s *SecureData) Open(targetPrivateKey [32]byte) ([]byte, error) {
	err := s.Validate()
	if err != nil {
		return nil, err
	}
	peerPublicKey := [32]byte{}
	copy(peerPublicKey[:], s.Key)
	nonce := [24]byte{}
	copy(nonce[:], s.Nonce)
	res := []byte{}
	_, opened := box.Open(res, s.Payload, &nonce, &peerPublicKey, &targetPrivateKey)
	if !opened {
		return nil, Corrupt
	}
	return res, nil
}
