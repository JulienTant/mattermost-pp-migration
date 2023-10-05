build:
	# build go app for linux (amd64 and arm64), darwin (amd64 and arm64), and windows (amd64)
	GOOS=linux GOARCH=amd64 go build -o bin/mattermost-pp-migration-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build -o bin/mattermost-pp-migration-linux-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/mattermost-pp-migration-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -o bin/mattermost-pp-migration-darwin-arm64 main.go
	GOOS=windows GOARCH=amd64 go build -o bin/mattermost-pp-migration-windows-amd64.exe main.go