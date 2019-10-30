```sh
go mod init $(basename `pwd`)
go mod tidy
dev_appserver.py app.yaml
open localhost:8080
open localhost:8000
```
