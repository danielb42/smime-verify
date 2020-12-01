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

	var expired bool
	if _, err := sd.VerifyDetached(messageBytes(), verifyOpts()); err != nil {
		if e, t := err.(x509.CertificateInvalidError); t && e.Reason == x509.Expired {
			expired = true
		} else {
			println("Signature is INVALID")
			os.Exit(2)
		}
	}

	chains, _ := sd.VerifyDetachedIgnoreExpiry(messageBytes(), verifyOpts())

	cert := chains[0][0][0]
	subj := cert.Subject.String()

	if expired {
		println("Signature is VALID but EXPIRED")
	} else {
		println("Signature is VALID")
	}

	println("Signed by:", subj)
}

func verifyOpts() x509.VerifyOptions {
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM([]byte(DTrust_Root))
	roots.AppendCertsFromPEM([]byte(Telesec_Root))

	intermediates := x509.NewCertPool()
	intermediates.AppendCertsFromPEM([]byte(DTrust_Intermediate))
	intermediates.AppendCertsFromPEM([]byte(Telesec_Intermediate))

	return x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermediates,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageAny},
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
