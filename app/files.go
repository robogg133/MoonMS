package app

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const (
	BANNED_ACCOUNTS_PATH = "banned-players.json"
	OPS_PATH             = "ops.json"
	WHITELIST_PATH       = "whitelist.json"
)

type WhitelistEntry struct {
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name,omitempty"`
}

type BanEntry struct {
	UUID    string `json:"uuid,omitempty"`
	Name    string `json:"name,omitempty"`
	IP      string `json:"ip,omitempty"`
	Created string `json:"created"`
	Source  string `json:"source"`
	Expires string `json:"expires"`
	Reason  string `json:"reason,omitempty"`
}

type OPEntry struct {
	UUID              string `json:"uuid"`
	Name              string `json:"name"`
	Level             uint8  `json:"level"`
	BypassPlayerLimit bool   `json:"bypassesPlayerLimit"`
}

func (s *Server) loadFiles() error {

	s.Bans.lock.Lock()
	if err := readFile(filepath.Join(s.Config.AcessFolder, BANNED_ACCOUNTS_PATH), &s.Bans.list); err != nil {
		return err
	}
	s.Bans.lock.Unlock()

	s.OPs.lock.Lock()
	if err := readFile(filepath.Join(s.Config.AcessFolder, OPS_PATH), &s.OPs.list); err != nil {
		return err
	}
	s.OPs.lock.Unlock()

	s.Whitelisteds.lock.Lock()
	if err := readFile(filepath.Join(s.Config.AcessFolder, WHITELIST_PATH), &s.Whitelisteds.list); err != nil {
		return err
	}
	s.Whitelisteds.lock.Unlock()

	s.Bans.lock.RLock()
	s.ban.lock.Lock()
	for _, v := range s.Bans.list {
		if v.Name != "" {
			s.ban.check[v.UUID] = &v
		}

		if v.IP != "" {
			s.ban.check[v.IP] = &v
		}
	}
	s.Bans.lock.RLock()
	s.ban.lock.Unlock()

	s.OPs.lock.RLock()
	s.op.lock.Lock()
	for _, v := range s.OPs.list {
		s.op.check[v.UUID] = &v
	}
	s.OPs.lock.RUnlock()
	s.op.lock.Unlock()

	s.Whitelisteds.lock.RLock()
	s.whitelist.lock.Lock()
	for _, v := range s.Whitelisteds.list {
		s.whitelist.check[v.UUID] = true
	}
	s.Whitelisteds.lock.RUnlock()
	s.whitelist.lock.Unlock()

	return nil
}

func (s *Server) basicFiles() error {

	if err := createFile(filepath.Join(s.Config.AcessFolder, BANNED_ACCOUNTS_PATH)); err != nil {
		return err
	}

	s.LogDebug("created banned-players file")

	if err := createFile(filepath.Join(s.Config.AcessFolder, OPS_PATH)); err != nil {
		return err
	}
	s.LogDebug("created ops file")

	if err := createFile(filepath.Join(s.Config.AcessFolder, WHITELIST_PATH)); err != nil {
		return err
	}

	s.LogDebug("created whitelist file")

	if err := os.MkdirAll(s.Config.PluginsFolder, 0777); err != nil {
		return err
	}

	s.LogDebug("created plugins folder")

	return nil
}

func readFile(src string, dst any) error {
	b, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return nil
	}

	if err := json.Unmarshal(b, dst); err != nil {
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
