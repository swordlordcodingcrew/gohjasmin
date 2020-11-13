module gohjasmin

go 1.14

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/mattn/go-sqlite3 v1.14.4
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/sirupsen/logrus v1.7.0 // indirect
	github.com/spf13/cobra v1.1.1 // indirect
	github.com/spf13/viper v1.7.1 // indirect
	golang.org/x/crypto v0.0.0-20201112155050-0c6587e931a9 // indirect
	internal/gohjaslib v1.0.0
)

replace internal/gohjaslib => ./internal/gohjaslib
