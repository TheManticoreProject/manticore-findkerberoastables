package core

import (
	"fmt"
	"strings"

	"github.com/TheManticoreProject/Manticore/network/ldap"
)

func domainToDN(domain string) string {
	parts := strings.Split(domain, ".")
	var dnParts []string
	for _, part := range parts {
		dnParts = append(dnParts, fmt.Sprintf("DC=%s", part))
	}
	return strings.Join(dnParts, ",")
}

func GetKerberoastables(ldapSession *ldap.Session, domain string) (map[string][]string, error) {

	results := map[string][]string{}

	baseDN := domainToDN(domain)

	query := "(&"
	query += "(|"
	query += "(objectClass=computer)"
	query += "(objectClass=person)"
	query += "(objectClass=user)"
	query += ")"
	query += "(servicePrincipalName=*)"
	query += ")"
	searchResults, err := ldapSession.QueryWholeSubtree(baseDN, query, []string{})
	if err != nil {
		return results, fmt.Errorf("error performing LDAP search: %s", err)
	}

	for _, entry := range searchResults {
		results[entry.DN] = entry.GetAttributeValues("servicePrincipalName")
	}

	return results, nil
}
