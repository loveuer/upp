package cache

import "encoding/json"

type scanner struct {
	err error
	bs  []byte
}

func (s *scanner) Scan(model any) error {
	if s.err != nil {
		return s.err
	}

	return json.Unmarshal(s.bs, model)
}

func newScanner(bs []byte, err error) *scanner {
	return &scanner{bs: bs, err: err}
}
