git reset --hard &&\
rm go.sum go.mod &&\
go mod init chatapp-api &&\
go mod tidy &&\
go run .
