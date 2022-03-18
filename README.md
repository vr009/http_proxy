# http_proxy

This is a simple MITM - http proxy server with minimal tls support.

## To start run:
`./build.sh`

This script generates the new CA, which is needed to be installed to other trusted CA of host system.

If you don't want to have CA as a trusted one you can use curl command with the `--insecure` flag or use `--cacert <path-to-ca>`.
