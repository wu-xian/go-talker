genny -imp "github.com/golang/protobuf/proto" -in .\generic.go.1 -out MessageWarpper.pb.go gen "T=MessageWarpper"
genny -imp "github.com/golang/protobuf/proto" -in .\generic.go.1 -out MessageLoginRequest.pb.go gen "T=MessageLoginRequest"
genny -imp "github.com/golang/protobuf/proto" -in .\generic.go.1 -out MessageLoginResponse.pb.go gen "T=MessageLoginResponse"
genny -imp "github.com/golang/protobuf/proto" -in .\generic.go.1 -out MessageLogoutRequest.pb.go gen "T=MessageLogoutRequest"
genny -imp "github.com/golang/protobuf/proto" -in .\generic.go.1 -out MessageLogoutResponse.pb.go gen "T=MessageLogoutResponse"
genny -imp "github.com/golang/protobuf/proto" -in .\generic.go.1 -out MessageClientLogin.pb.go gen "T=MessageClientLogin"
genny -imp "github.com/golang/protobuf/proto" -in .\generic.go.1 -out MessageClientLogout.pb.go gen "T=MessageClientLogout"
genny -imp "github.com/golang/protobuf/proto" -in .\generic.go.1 -out MessageClientReceived.pb.go gen "T=MessageClientReceived"
genny -imp "github.com/golang/protobuf/proto" -in .\generic.go.1 -out MessageClientSend.pb.go gen "T=MessageClientSend"
protoc --go_out=.\ .\message.proto