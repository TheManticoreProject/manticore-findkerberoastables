package core

import (
	"fmt"

	"github.com/TheManticoreProject/Manticore/network/ldap"
)

func GetKerberoastables(ldapSession *ldap.Session) (map[string][]string, error) {

	results := map[string][]string{}

	query := "(&"
	query += "(|"
	query += "(objectClass=computer)"
	query += "(objectClass=person)"
	query += "(objectClass=user)"
	query += ")"
	query += "(servicePrincipalName=*)"
	query += ")"
	searchResults, err := ldapSession.QueryWholeSubtree("", query, []string{})
	if err != nil {
		return results, fmt.Errorf("error performing LDAP search: %s", err)
	}

	for _, entry := range searchResults {
		results[entry.DN] = entry.GetAttributeValues("servicePrincipalName")
	}

	return results, nil
}
