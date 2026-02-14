package app

import (
	"os"
	"path/filepath"
)

const (
	BANNED_IPS_PATH      = "access/banned-ips.json"
	BANNED_ACCOUNTS_PATH = "access/banned-accounts.json"
	OPS_PATH             = "access/ops.json"
	WHITELIST_PATH       = "access/whitelist.json"
)

func (s *Server) basicFiles() error {
	if err := createFile(BANNED_IPS_PATH); err != nil {
		return err
	}
	s.LogDebug("creating ip file")

	if err := createFile(BANNED_ACCOUNTS_PATH); err != nil {
		return err
	}

	if err := createFile(BANNED_IPS_PATH); err != nil {
		return err
	}

	if err := createFile(OPS_PATH); err != nil {
		return err
	}

	if err := createFile(WHITELIST_PATH); err != nil {
		return err
	}

	return nil
}

func createFile(s string) error {
	_, err := os.Stat(s)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Dir(s), 0755)
			if err != nil {
				return err
			}
			f, err := os.Create(s)
			if err != nil {

				return err
			}
			f.Close()
		} else {
			return err
		}
	}
	return nil
}
