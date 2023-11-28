//go:build wasm

package plan

import (
	"errors"
)

func DownloadUpgrade(dstRoot, url, daemonName string) error {
	return errors.New("upgrade module not supported for wasm")
}

func EnsureBinary(path string) error {
	return errors.New("upgrade module not supported for wasm")
}

func DownloadURLWithChecksum(url string) (string, error) {
	return "", errors.New("upgrade module not supported for wasm")
}

func ValidateIsURLWithChecksum(urlStr string) error {
	return errors.New("upgrade module not supported for wasm")
}
