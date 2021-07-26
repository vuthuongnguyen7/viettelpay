package viettelpay

import (
	"compress/gzip"
	"encoding/json"
	"io"
)

func MarshalGzipJSON(w io.Writer, data interface{}) error {
	gz := gzip.NewWriter(w)

	if err := json.NewEncoder(gz).Encode(data); err != nil {
		return err
	}

	return gz.Close()
}

func UnmarshalGzipJSON(r io.Reader, data interface{}) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	if err = json.NewDecoder(gz).Decode(data); err != nil {
		return err
	}

	return gz.Close()
}
