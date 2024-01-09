build:
	@echo 'Building cmd/web...'
	go build -o ./tmp/web ./cmd/web \
	&& cp -r ./tls /tmp/ \
	&& ./tmp/web