package entity

import (
	"bytes"
	"encoding/base64"
	"strconv"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

type WgKey wgtypes.Key

func (k WgKey) MarshalEasyJSON() ([]byte, error) {
	if k.IsEmpty() {
		return []byte(strconv.Quote("")), nil
	}

	return []byte(strconv.Quote(k.String())), nil
}

func (k *WgKey) UnmarshalEasyJSON(data []byte) error {
	d, err := strconv.Unquote(string(data))
	if err != nil {
		return err
	}

	if d == "" {
		*k = WgKey(wgtypes.Key{})
		return nil
	}

	key, err := wgtypes.ParseKey(d)
	if err != nil {
		return err
	}

	*k = WgKey(key)

	return nil
}

func (k WgKey) String() string {
	return base64.StdEncoding.EncodeToString(k[:])
}

func (k WgKey) IsEmpty() bool {
	return bytes.Equal(k[:], make([]byte, len(k)))
}
