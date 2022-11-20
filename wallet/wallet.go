package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey *ecdsa.PrivateKey
	publicKey  *ecdsa.PublicKey
	address    string
}

func NewWallet() *Wallet {
	key, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	public := key.PublicKey

	// following steps are to generate shorter blockchain address from public key
	// https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses
	h2 := sha256.New()
	h2.Write(public.X.Bytes())
	h2.Write(public.Y.Bytes())
	digest2 := h2.Sum(nil)

	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)

	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])

	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)

	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)

	checksum := digest6[:4]

	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], checksum[:])

	address := base58.Encode(dc8)

	return &Wallet{
		privateKey: key,
		publicKey:  &public,
		address:    address,
	}
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%x%x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) Address() string {
	return w.address
}
