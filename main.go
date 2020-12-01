package main

import (
	"bufio"
	"crypto/x509"
	"encoding/base64"
	"os"
	"regexp"

	cms "github.com/github/ietf-cms"
)

func main() {
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		println("usage:", os.Args[0], "<filename>")
		os.Exit(1)
	}

	sd, err := cms.ParseSignedData(signatureBytes())
	if err != nil {
		println("error: failed to parse signature")
		os.Exit(1)
	}

	chains, err := sd.VerifyDetached(messageBytes(), verifyOpts())
	if err != nil {
		switch t := err.(type) {
		case x509.CertificateInvalidError:
			if t.Reason == x509.Expired {
				println("signing certificate has expired")
				os.Exit(2)
			}
		}

		println("invalid signature")
		os.Exit(2)
	}

	cert := chains[0][0][0]
	subj := cert.Subject.String()

	println("valid signature from", subj)
}

func verifyOpts() x509.VerifyOptions {
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM([]byte(root))
	roots.AppendCertsFromPEM([]byte(intermediate))

	return x509.VerifyOptions{
		Roots: roots,
	}
}

func signatureBytes() []byte {
	asciiSignature := readBetween(os.Args[1], "S/MIME Cryptographic Signature", "------=_Part")
	rawSignature, _ := base64.StdEncoding.DecodeString(string(asciiSignature))
	return rawSignature
}

func messageBytes() []byte {
	return readBetween(os.Args[1], "------=_Part_", "------=_Part_")
}

func readBetween(filename, startAtPattern, endAtPattern string) []byte {
	file, _ := os.Open(filename)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var insideInterestingBlock bool
	var output string

	for scanner.Scan() {
		endMatched, _ := regexp.MatchString(endAtPattern, scanner.Text())
		if insideInterestingBlock && endMatched {
			break
		}

		if insideInterestingBlock {
			output += scanner.Text()
			output += "\r\n"
		}

		startMatched, _ := regexp.MatchString(startAtPattern, scanner.Text())
		if startMatched {
			insideInterestingBlock = true
		}

	}

	outputBytes := []byte(output)

	if len(outputBytes) == 0 {
		println("error: reading file", filename, "failed")
		os.Exit(1)
	}

	return outputBytes[:len(outputBytes)-2]
}
