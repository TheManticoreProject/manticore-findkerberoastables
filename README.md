![](./.github/banner.png)

<p align="center">
      A tool to extract and display the list of Kerberoastable users and computers (service principal names) from an Active Directory environment over LDAP.
      <br>
      <a href="https://github.com/TheManticoreProject/FindKerberoastables/actions/workflows/release.yaml" title="Build"><img alt="Build and Release" src="https://github.com/TheManticoreProject/FindKerberoastables/actions/workflows/release.yaml/badge.svg"></a>
      <img alt="GitHub release (latest by date)" src="https://img.shields.io/github/v/release/TheManticoreProject/FindKerberoastables">
      <img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/TheManticoreProject/FindKerberoastables">
      <a href="https://twitter.com/intent/follow?screen_name=podalirius_" title="Follow"><img src="https://img.shields.io/twitter/follow/podalirius_?label=Podalirius&style=social"></a>
      <a href="https://www.youtube.com/c/Podalirius_?sub_confirmation=1" title="Subscribe"><img alt="YouTube Channel Subscribers" src="https://img.shields.io/youtube/channel/subscribers/UCF_x5O7CSfr82AfNVTKOv_A?style=social"></a>
      <br>
</p>

## Features

- [x] Extract the list of Kerberoastable users from an Active Directory.

## Usage

Run the tool with `-h` to display the help message:

```
$ ./FindKerberoastables -h
FindKerberoastables - by Remi GASCOU (Podalirius) @ TheManticoreProject - v1.0.0

Usage: FindKerberoastables --domain <string> --username <string> [--password <string>] [--hashes <string>] [--debug] [--print-spns] --dc-ip <string> [--ldap-port <tcp port>] [--use-ldaps] [--use-kerberos]

  Authentication:
    -d, --domain <string>   Active Directory domain to authenticate to.
    -u, --username <string> User to authenticate as.
    -p, --password <string> Password to authenticate with. (default: "")
    -H, --hashes <string>   NT/LM hashes, format is LMhash:NThash. (default: "")

  Configuration:
    --debug          Debug mode. (default: false)
    -s, --print-spns Print SPNs. (default: false)

  LDAP Connection Settings:
    -dc, --dc-ip <string>       IP Address of the domain controller or KDC (Key Distribution Center) for Kerberos. If omitted, it will use the domain part (FQDN) specified in the identity parameter.
    -lp, --ldap-port <tcp port> Port number to connect to LDAP server. (default: 389)
    -L, --use-ldaps             Use LDAPS instead of LDAP. (default: false)
    -k, --use-kerberos          Use Kerberos instead of NTLM. (default: false)
```

### Basic Usage

Extract kerberoastable accounts from an Active Directory:

```
$ ./FindKerberoastables -d "EXAMPLE.local" -u "Administrator" -p "Password123!" -dc "10.0.0.1"
FindKerberoastables - by Remi GASCOU (Podalirius) @ TheManticoreProject - v1.0.0

├── CN=WEB-SERVER,CN=Computers,DC=EXAMPLE,DC=local
└── CN=SQL-SERVER,CN=Computers,DC=EXAMPLE,DC=local
Done
```

### Display Service Principal Names (SPNs)

Use the `-s` or `--print-spns` flag to display the SPNs for each kerberoastable account:

```
$ ./FindKerberoastables -d "EXAMPLE.local" -u "Administrator" -p "Password123!" -dc "10.0.0.1" -s
FindKerberoastables - by Remi GASCOU (Podalirius) @ TheManticoreProject - v1.0.0

├── CN=WEB-SERVER,CN=Computers,DC=EXAMPLE,DC=local
│   ├── HTTP/WEB-SERVER
│   └── HTTP/WEB-SERVER.EXAMPLE.local
└── CN=SQL-SERVER,CN=Computers,DC=EXAMPLE,DC=local
    ├── MSSQLSvc/SQL-SERVER
    └── MSSQLSvc/SQL-SERVER.EXAMPLE.local:1433
Done
```

### Using LDAPS

For environments requiring LDAPS (encrypted LDAP over TLS), use the `-L` flag. The tool automatically selects port 636 when LDAPS is enabled:

```
$ ./FindKerberoastables -d "EXAMPLE.local" -u "Administrator" -p "Password123!" -dc "10.0.0.1" -L
```

### Using Pass-the-Hash

Authenticate using NTLM hash instead of plaintext password:

```
$ ./FindKerberoastables -d "EXAMPLE.local" -u "Administrator" -H "aad3b435b51404eeaad3b435b51404ee:5f4dcc3b5aa765d61d8327deb882cf99" -dc "10.0.0.1" -L
```

### Using Kerberos Authentication

Use Kerberos for authentication instead of NTLM:

```
$ ./FindKerberoastables -d "EXAMPLE.local" -u "Administrator" -p "Password123!" -dc "10.0.0.1" -L -k
```

## Demonstration

The following demonstrates finding kerberoastable accounts in an Active Directory domain:

```
$ ./FindKerberoastables -d "MANTICORE.local" -u "Administrator" -p "MyPassword123!" -dc "192.168.1.10" -L -s
FindKerberoastables - by Remi GASCOU (Podalirius) @ TheManticoreProject - v1.0.0

[2026-08-27 10h42m35s] DEBUG: Found 2 kerberoastable accounts.
├── CN=MANTICORE-DC1,OU=Domain Controllers,DC=MANTICORE,DC=local
│   ├── TERMSRV/MANTICORE-DC1
│   ├── ldap/MANTICORE-DC1.MANTICORE.local
│   ├── DNS/MANTICORE-DC1.MANTICORE.local
│   ├── HOST/MANTICORE-DC1
│   ├── HOST/MANTICORE-DC1.MANTICORE.local
│   └── GC/MANTICORE-DC1.MANTICORE.local/MANTICORE.local
└── CN=krbtgt,CN=Users,DC=MANTICORE,DC=local
    └── kadmin/changepw
Done
```
## Contributing

Pull requests are welcome. Feel free to open an issue if you want to add other features.

## Credits

- [Remi GASCOU (Podalirius)](https://github.com/Podalirius) for the creation of the [FindKerberoastables](https://github.com/TheManticoreProject/FindKerberoastables) project.

