package PSUtils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ed25519"
	"golang.org/x/crypto/ssh"
	mathrand "math/rand"
	"strings"
)

type KeyPair struct {
	PrivateKey string
	AuthKey    string
	Password string
	Type       string // rsa/des/ecdsa/ed25519
	// KeySize    int    // for rsa/des: 1024/2048/4096
	// Cipher     string // for rsa/des/ecdsa: des/3des/aes-128/AES-192/aes-256
	// Curve      string // for ecdsa/ed25519
}

// default keysize: 2048
// default cipher:  aes-256
func GenerateRsaSshKeyPair(password string, keySize int) (*KeyPair, error) {
	size := keySize
	if keySize != 1024 && keySize != 2048 && keySize != 4096 {
		size = 2048
	}
	cipher := x509.PEMCipherAES256
	privateKey, err := rsa.GenerateKey(rand.Reader, size)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	if password != "" {
		privateKeyBlock, err = x509.EncryptPEMBlock(rand.Reader, privateKeyBlock.Type, privateKeyBlock.Bytes,
			[]byte(password), cipher)
		if err != nil {
			return nil, err
		}
	}

	privateKeyString := string(pem.EncodeToMemory(privateKeyBlock))
	publicRsaKey, err := ssh.NewPublicKey(privateKey.Public())
	if err != nil {
		return nil, err
	}
	authorizedKeyString := string(ssh.MarshalAuthorizedKey(publicRsaKey))
	return &KeyPair{
		PrivateKey: privateKeyString,
		AuthKey:    authorizedKeyString,
		Type:       "rsa",
		Password: password,
	}, nil
}

func GenerateED25519SshKeyPair() (*KeyPair, error) {
	pubKey, privKey, _ := ed25519.GenerateKey(rand.Reader)
	publicKey, _ := ssh.NewPublicKey(pubKey)

	pemKey := &pem.Block{
		Type:  "OPENSSH PRIVATE KEY",
		Bytes: marshalED25519PrivateKey(privKey),
	}
	privateKey := pem.EncodeToMemory(pemKey)
	authorizedKey := ssh.MarshalAuthorizedKey(publicKey)

	return &KeyPair{
		PrivateKey: string(privateKey),
		AuthKey:    string(authorizedKey),
		Password:   "",
		Type:       "ed25519",
	}, nil
}

func GenerateAuthorizedKeyFrom(privateKey, password string) (string, error) {
	var signer ssh.Signer
	var err error
	if password != "" && strings.Contains(password, "ENCRYPTED") {
		signer, err = ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(password))
		if err != nil {
			return "", err
		}
	} else {
		signer, err = ssh.ParsePrivateKey([]byte(privateKey))
		if err != nil {
			return "", err
		}
	}

	return string(ssh.MarshalAuthorizedKey(signer.PublicKey())), nil
}

/* Writes ed25519 private keys into the new OpenSSH private key format.
I have no idea why this isn't implemented anywhere yet, you can do seemingly
everything except write it to disk in the OpenSSH private key format. */
func marshalED25519PrivateKey(key ed25519.PrivateKey) []byte {
	// Add our key header (followed by a null byte)
	magic := append([]byte("openssh-key-v1"), 0)

	var w struct {
		CipherName   string
		KdfName      string
		KdfOpts      string
		NumKeys      uint32
		PubKey       []byte
		PrivKeyBlock []byte
	}

	// Fill out the private key fields
	pk1 := struct {
		Check1  uint32
		Check2  uint32
		Keytype string
		Pub     []byte
		Priv    []byte
		Comment string
		Pad     []byte `ssh:"rest"`
	}{}

	// Set our check ints
	ci := mathrand.Uint32()
	pk1.Check1 = ci
	pk1.Check2 = ci

	// Set our key type
	pk1.Keytype = ssh.KeyAlgoED25519

	// Add the pubkey to the optionally-encrypted block
	pk, ok := key.Public().(ed25519.PublicKey)
	if !ok {
		//fmt.Fprintln(os.Stderr, "ed25519.PublicKey type assertion failed on an ed25519 public key. This should never ever happen.")
		return nil
	}
	pubKey := []byte(pk)
	pk1.Pub = pubKey

	// Add our private key
	pk1.Priv = []byte(key)

	// Might be useful to put something in here at some point
	pk1.Comment = ""

	// Add some padding to match the encryption block size within PrivKeyBlock (without Pad field)
	// 8 doesn't match the documentation, but that's what ssh-keygen uses for unencrypted keys. *shrug*
	bs := 8
	blockLen := len(ssh.Marshal(pk1))
	padLen := (bs - (blockLen % bs)) % bs
	pk1.Pad = make([]byte, padLen)

	// Padding is a sequence of bytes like: 1, 2, 3...
	for i := 0; i < padLen; i++ {
		pk1.Pad[i] = byte(i + 1)
	}

	// Generate the pubkey prefix "\0\0\0\nssh-ed25519\0\0\0 "
	prefix := []byte{0x0, 0x0, 0x0, 0x0b}
	prefix = append(prefix, []byte(ssh.KeyAlgoED25519)...)
	prefix = append(prefix, []byte{0x0, 0x0, 0x0, 0x20}...)

	// Only going to support unencrypted keys for now
	w.CipherName = "none"
	w.KdfName = "none"
	w.KdfOpts = ""
	w.NumKeys = 1
	w.PubKey = append(prefix, pubKey...)
	w.PrivKeyBlock = ssh.Marshal(pk1)

	magic = append(magic, ssh.Marshal(w)...)

	return magic
}

