package main

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/openpgp"
)

func VerifySignedKey(sshKey, Signature string) (bool, error) {
	keyRingReader := strings.NewReader(getKey())
	signatureReader := strings.NewReader(Signature)

	verificationTarget := strings.NewReader(sshKey)

	keyring, err := openpgp.ReadArmoredKeyRing(keyRingReader)
	if err != nil {
		fmt.Println("Read Armored Key Ring: " + err.Error())
		return false, err
	}

	_, err = openpgp.CheckArmoredDetachedSignature(keyring, verificationTarget, signatureReader)
	if err != nil {
		fmt.Println("Check Detached Signature: " + err.Error())
		return false, err
	}

	return true, nil

}
