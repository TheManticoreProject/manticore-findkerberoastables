package main

import (
	"fmt"

	"github.com/TheManticoreProject/FindKerberoastables/core"

	"github.com/TheManticoreProject/Manticore/logger"
	"github.com/TheManticoreProject/Manticore/network/ldap"
	"github.com/TheManticoreProject/Manticore/windows/credentials"

	"github.com/TheManticoreProject/goopts/parser"
)

var (
	// Configuration
	debug     bool
	printSPNs bool

	// Authentication
	authDomain   string
	authUsername string
	authPassword string
	authHashes   string

	// LDAP Connection Settings
	domainController string
	ldapPort         int
	useLdaps         bool
	useKerberos      bool
)

func parseArgs() {
	ap := parser.ArgumentsParser{
		Banner: "FindKerberoastables - by Remi GASCOU (Podalirius) @ TheManticoreProject - v1.0.0",
		Options: parser.ArgumentsParserOptions{
			ShowBannerOnHelp: true,
			ShowBannerOnRun:  true,
		},
	}

	// Configuration flags
	group_config, err := ap.NewArgumentGroup("Configuration")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_config.NewBoolArgument(&debug, "", "--debug", false, "Debug mode.")
		group_config.NewBoolArgument(&printSPNs, "-s", "--print-spns", false, "Print SPNs.")
	}
	// LDAP Connection Settings
	group_ldapSettings, err := ap.NewArgumentGroup("LDAP Connection Settings")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_ldapSettings.NewStringArgument(&domainController, "-dc", "--dc-ip", "", true, "IP Address of the domain controller or KDC (Key Distribution Center) for Kerberos. If omitted, it will use the domain part (FQDN) specified in the identity parameter.")
		group_ldapSettings.NewTcpPortArgument(&ldapPort, "-lp", "--ldap-port", 389, false, "Port number to connect to LDAP server.")
		group_ldapSettings.NewBoolArgument(&useLdaps, "-L", "--use-ldaps", false, "Use LDAPS instead of LDAP.")
		group_ldapSettings.NewBoolArgument(&useKerberos, "-k", "--use-kerberos", false, "Use Kerberos instead of NTLM.")
	}
	// Authentication flags
	group_auth, err := ap.NewArgumentGroup("Authentication")
	if err != nil {
		fmt.Printf("[error] Error creating ArgumentGroup: %s\n", err)
	} else {
		group_auth.NewStringArgument(&authDomain, "-d", "--domain", "", true, "Active Directory domain to authenticate to.")
		group_auth.NewStringArgument(&authUsername, "-u", "--username", "", true, "User to authenticate as.")
		group_auth.NewStringArgument(&authPassword, "-p", "--password", "", false, "Password to authenticate with.")
		group_auth.NewStringArgument(&authHashes, "-H", "--hashes", "", false, "NT/LM hashes, format is LMhash:NThash.")
	}

	ap.Parse()
}

func main() {
	parseArgs()

	if useLdaps && ldapPort == 389 {
		ldapPort = 636
	}

	creds, err := credentials.NewCredentials(authDomain, authUsername, authPassword, authHashes)
	if err != nil {
		fmt.Printf("[error] Error creating credentials: %s\n", err)
		return
	}

	ldapSession, err := ldap.NewSession(domainController, ldapPort, creds, useLdaps, useKerberos)
	if err != nil {
		logger.Warn(fmt.Sprintf("%s\n", err))
		return
	}
	success, err := ldapSession.Connect()
	if !success {
		logger.Warn(fmt.Sprintf("%s\n", err))
		return
	}
	defer ldapSession.Close()

	kerberoastables, err := core.GetKerberoastables(ldapSession, authDomain)
	if err != nil {
		logger.Warn(fmt.Sprintf("%s\n", err))
		return
	}

	lenKerberoastables := len(kerberoastables)
	if debug {
		logger.Debug(fmt.Sprintf("Found %d kerberoastable accounts.", lenKerberoastables))
	}
	kerberoastable_id := 0
	for dn, spns := range kerberoastables {
		kerberoastable_id += 1
		isLastDN := kerberoastable_id == lenKerberoastables
		if isLastDN {
			logger.Print(fmt.Sprintf("└── %s\n", dn))
		} else {
			logger.Print(fmt.Sprintf("├── %s\n", dn))
		}
		childIndent := "│   "
		if isLastDN {
			childIndent = "    "
		}
		if !printSPNs {
			continue
		}
		for spn_id, spn := range spns {
			if spn_id < len(spns)-1 {
				logger.Print(fmt.Sprintf("%s├── %s\n", childIndent, spn))
			} else {
				logger.Print(fmt.Sprintf("%s└── %s\n", childIndent, spn))
			}
		}
	}

	fmt.Println("Done")
}
