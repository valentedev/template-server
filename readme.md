Disabling directory listings of /static/
$ find ./ui/static -type d -exec touch {}/index.html \;

Generate TLS certificates - go to /your-project/tls-folder and run:
go run /opt/homebrew/Cellar/go/1.21.5/libexec/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost
