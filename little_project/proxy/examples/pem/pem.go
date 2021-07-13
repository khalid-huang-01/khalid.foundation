package main

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"log"
	"os"
)

func main() {
//	var pubPEMData = []byte(`
//-----BEGIN PUBLIC KEY-----
//MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
//WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
//CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x
//qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
//yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY
//nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
//6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q
//TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/
//a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9
//PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ
//yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUk
//AIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ==
//-----END PUBLIC KEY-----
//`)


	var crtPEMData = []byte(`
-----BEGIN CERTIFICATE-----
MIIBWjCCAQCgAwIBAgICBAAwCgYIKoZIzj0EAwIwEzERMA8GA1UEAxMIS3ViZUVk
Z2UwIBcNMjEwNDE2MDExOTEzWhgPMjEyMTAzMjMwMTE5MTNaMBMxETAPBgNVBAMT
CEt1YmVFZGdlMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEErbUyo8qMGfhG5ql
beslWs7oZxVHV0UhWEqJtEI1aIOy4+yGkf5FPpUP5tjXBhg9BMTRp8njAflm+hvf
xGhjjKNCMEAwDgYDVR0PAQH/BAQDAgKkMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggr
BgEFBQcDAjAPBgNVHRMBAf8EBTADAQH/MAoGCCqGSM49BAMCA0gAMEUCIGjAgFkT
BVlKEeuQUKRajUOcI+W+dGNTT41RUKkrM/faAiEA1E+OAgx0Td4hpUKFoBK2/jW1
3axIlyt4XiTOfe3QrJI=
-----END CERTIFICATE-----
`)
	// block.Type的值就是我们上面的内容里面的BEGIN的内容

	block, _ := pem.Decode(crtPEMData)
	if block == nil {
		log.Fatal("failed to decode PEM block containing public key")
	}
	log.Printf("%+v", block)

	if err := pem.Encode(os.Stdout, block); err != nil {
		log.Fatal(err)
	}

	m := pem.EncodeToMemory(block) // []byte

	fmt.Printf("m=[%s]\n", m)

	sEnc := base64.StdEncoding.EncodeToString(block.Bytes)
	fmt.Println(sEnc)
}
