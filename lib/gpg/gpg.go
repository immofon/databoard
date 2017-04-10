package gpg

import (
	"crypto"
	"io"

	"github.com/juju/errors"

	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/packet"
)

var (
	ErrNotCurrectPassphrase = errors.New("not currect passphrase")
)

func SymmetricEncrypt(ciphertext io.Writer, passphrase []byte) (plaintext io.WriteCloser, err error) {
	plaintext, err = openpgp.SymmetricallyEncrypt(ciphertext, passphrase, nil, &packet.Config{
		DefaultHash:   crypto.SHA256,
		DefaultCipher: packet.CipherAES256,
	})
	if err != nil {
		return nil, errors.Annotate(err, "openpgp.SymmetricallyEncrypt")
	}

	return plaintext, nil
}

var emptyKeyRing = new(openpgp.EntityList)

func SymmetricDecrypt(ciphertext io.Reader, passphrase []byte) (plaintext io.Reader, err error) {
	md, err := openpgp.ReadMessage(ciphertext, emptyKeyRing, func() openpgp.PromptFunction {
		var isFirst = true
		return func(_ []openpgp.Key, _ bool) ([]byte, error) {
			if !isFirst {
				return nil, errors.Trace(ErrNotCurrectPassphrase)
			}
			isFirst = false
			return passphrase, nil
		}
	}(), &packet.Config{
		DefaultHash:   crypto.SHA256,
		DefaultCipher: packet.CipherAES256,
	})

	if err != nil {
		return nil, errors.Annotate(err, "openpgp.ReadMessage")
	}

	return md.UnverifiedBody, nil
}
