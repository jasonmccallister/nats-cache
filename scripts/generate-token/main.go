package main

import (
	"flag"
	"fmt"
	"os"

	"aidanwoods.dev/go-paseto"
)

func main() {
	subj := flag.String("subject", "", "The subject to use for the token")
	secret := flag.String("secret", "", "The secret to use for generating the token")
	flag.Parse()

	if *subj == "" {
		fmt.Println("subject is required")
		os.Exit(1)
	}

	token := paseto.NewToken()
	token.SetSubject(*subj)

	// if a secret is provided, use it, otherwise generate a new one
	var secretKey paseto.V4AsymmetricSecretKey
	if *secret != "" {
		s, err := paseto.NewV4AsymmetricSecretKeyFromHex(*secret)
		if err != nil {
			fmt.Println("failed to create secret key:", err)
			os.Exit(1)
		}

		secretKey = s
	} else {
		secretKey = paseto.NewV4AsymmetricSecretKey()
	}

	publicKey := secretKey.Public()

	signed := token.V4Sign(secretKey, nil)

	fmt.Println("secret key hex:", secretKey.ExportHex())
	fmt.Println("public key hex:", publicKey.ExportHex())
	fmt.Println("subject:", *subj)
	fmt.Println("signed token:", signed)
}
