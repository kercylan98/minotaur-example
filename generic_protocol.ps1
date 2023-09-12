$FolderPath = "./protocol/protocol"
if (-Not (Test-Path -Path $FolderPath)) {
    New-Item -Path $FolderPath -ItemType Directory
}

protoc --proto_path . ./protocol/message_system.proto --go_out=protocol/protocol --go-grpc_out=protocol/protocol/ protocol/*.proto