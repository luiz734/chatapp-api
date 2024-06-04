git reset --hard &&\
chmod +x run.sh &&\
rm go.sum go.mod &&\
go mod init chatapp-api &&\
go mod tidy &&\
direnv allow &&\
go run .
