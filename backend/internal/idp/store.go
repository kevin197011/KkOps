// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package idp

import (
	"sync"
	"time"
)

type authCodeEntry struct {
	UserID      uint
	ClientID    string
	RedirectURI string
	Scope       string
	State       string
	Expiry      time.Time
}

var (
	codeStore   = make(map[string]*authCodeEntry)
	codeStoreMu sync.RWMutex
)

func saveCode(code string, e *authCodeEntry) {
	e.Expiry = time.Now().Add(5 * time.Minute)
	codeStoreMu.Lock()
	defer codeStoreMu.Unlock()
	codeStore[code] = e
}

func consumeCode(code string) (*authCodeEntry, bool) {
	codeStoreMu.Lock()
	defer codeStoreMu.Unlock()
	e, ok := codeStore[code]
	if !ok || e == nil {
		return nil, false
	}
	if time.Now().After(e.Expiry) {
		delete(codeStore, code)
		return nil, false
	}
	delete(codeStore, code)
	return e, true
}
