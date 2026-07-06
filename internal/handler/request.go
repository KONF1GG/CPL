package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type CreateVMRequest struct {
	Name   string `json:"name"`
	CPU    int    `json:"cpu"`
	RamMB  int    `json:"ram_mb"`
	DiskGB int    `json:"disk_gb"`
}

func decodeJSON(r *http.Request, dst any) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return errors.New("request body must contain a single JSON value")
	}
	return nil
}

func parseID(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0, err
	}
	if id == 0 || id > uint64(^uint(0)) {
		return 0, fmt.Errorf("invalid id %q", s)
	}
	return uint(id), nil
}
