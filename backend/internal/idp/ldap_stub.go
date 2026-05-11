// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package idp

import (
	"crypto/tls"
	"net"

	"go.uber.org/zap"

	"github.com/kkops/backend/internal/config"
)

// StartLDAPListener listens for LDAP connections (stub: accepts and closes; TODO: bind/search).
func StartLDAPListener(cfg *config.Config, log *zap.Logger) {
	if cfg == nil || !cfg.IdP.LDAP.Enabled {
		return
	}
	addr := cfg.IdP.LDAP.ListenAddr
	if addr == "" {
		addr = ":1389"
	}
	var (
		ln  net.Listener
		err error
	)
	if cfg.IdP.LDAP.TLSCert != "" && cfg.IdP.LDAP.TLSKey != "" {
		cert, cerr := tls.LoadX509KeyPair(cfg.IdP.LDAP.TLSCert, cfg.IdP.LDAP.TLSKey)
		if cerr != nil {
			log.Error("ldap tls cert load failed", zap.Error(cerr))
			return
		}
		ln, err = tls.Listen("tcp", addr, &tls.Config{Certificates: []tls.Certificate{cert}})
	} else {
		ln, err = net.Listen("tcp", addr)
	}
	if err != nil {
		log.Error("ldap listen failed", zap.String("addr", addr), zap.Error(err))
		return
	}
	log.Info("LDAP IdP scaffold listening", zap.String("addr", addr))
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				log.Debug("ldap stub connection accepted", zap.String("remote", c.RemoteAddr().String()))
				// TODO: parse BER LDAP bind request and respond with LDAP Result (use github.com/jimlambrt/gldap or similar).
			}(conn)
		}
	}()
}
