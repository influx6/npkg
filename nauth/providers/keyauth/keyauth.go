package keyauth

import "crypto"

type KeyAuth struct {
	Secret     string
	PrivateKey crypto.PrivateKey
	PublicKey  crypto.PublicKey
	Signer     crypto.Signer
}
