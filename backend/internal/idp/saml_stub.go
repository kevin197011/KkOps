// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package idp

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/kkops/backend/internal/config"
)

// RegisterSAMLRoutes registers SAML scaffold routes when enabled (TODO: full SP integration).
func RegisterSAMLRoutes(r *gin.Engine, cfg *config.Config) {
	if cfg == nil || !cfg.IdP.SAML.Enabled {
		return
	}
	g := r.Group("/saml")
	{
		g.GET("/metadata", samlMetadata(cfg))
		// IdP-initiated style POST — returns HTML that POSTs a placeholder SAMLResponse (TODO: real ACS URL + canonicalization).
		g.GET("/idp-initiated-stub", samlIDPInitiatedStub(cfg))
	}
}

func samlMetadata(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		xml := `<?xml version="1.0" encoding="UTF-8"?>
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" entityID="` + issuerBase(cfg) + `/saml">
  <IDPSSODescriptor protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol">
    <SingleSignOnService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="` + issuerBase(cfg) + `/saml/idp-initiated-stub"/>
  </IDPSSODescriptor>
</EntityDescriptor>`
		c.Data(http.StatusOK, "application/xml", []byte(xml))
	}
}

func issuerBase(cfg *config.Config) string {
	if cfg != nil && cfg.IdP.Issuer != "" {
		return cfg.IdP.Issuer
	}
	return "http://localhost:8080"
}

func samlIDPInitiatedStub(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := getPrivateKey()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "no signing key"})
			return
		}
		_ = key
		assertion := `<samlp:Response xmlns:samlp="urn:oasis:names:tc:SAML:2.0:protocol" ID="_stub" Version="2.0" IssueInstant="` +
			time.Now().UTC().Format(time.RFC3339) + `" Destination="https://example.com/acs">
  <Issuer xmlns="urn:oasis:names:tc:SAML:2.0:assertion">` + issuerBase(cfg) + `</Issuer>
  <Status><StatusCode Value="urn:oasis:names:tc:SAML:2.0:status:Success"/></Status>
</samlp:Response>`
		encoded := base64.StdEncoding.EncodeToString([]byte(assertion))
		html := `<!DOCTYPE html><html><body onload="document.forms[0].submit()">
<form method="post" action="https://example.com/acs">
<input type="hidden" name="SAMLResponse" value="` + encoded + `"/>
<!-- TODO: RSA-SHA256 XML signature (exclusive C14N); placeholder response only -->
</form></body></html>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
	}
}
