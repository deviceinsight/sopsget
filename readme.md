# Sopsget

A tool which checks your .sops.yaml for missing fingerprints with `gpg` and automatically fetches them from `keys.openpgpg.org`

## Usage

* `sopsget` \
checks if a .sops.yaml is present in your current directory and uses that to processing
*  `sopsget ./folder/file` \
checks processes the specified file

Constraint:
the file has to have the `.yaml` or `.yml` ending

##  Build

### Build locally:

* `go build`
* make the built binary available on your path