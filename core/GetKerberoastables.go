package core

import (
	"fmt"

	"github.com/TheManticoreProject/Manticore/network/ldap"
)

func GetKerberoastables(ldapSession *ldap.Session) (map[string][]string, error) {

	results := map[string][]string{}

	success, err := ldapSession.Connect()
	if !success {
		return results, fmt.Errorf("error connecting to LDAP: %s\n", err)
	}

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
		return results, fmt.Errorf("error performing LDAP search: %s\n", err)
	}

	for _, entry := range searchResults {
		results[entry.DN] = entry.GetAttributeValues("servicePrincipalName")
	}

	ldapSession.Close()

	return results, nil
}
