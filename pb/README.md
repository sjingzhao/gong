### 定义pb源文件

```protobuf
syntax = "proto3";  // proto 版本
option go_package = "./example";  // 关键：指定生成的 Go 代码输出到当前目录的 example 文件夹（包名也为 example）
package example;    // proto 包名（可与 go_package 一致或不同）

// 定义 gRPC 服务
service UserService {
  rpc GetUser(UserRequest) returns (User);
  rpc AddUser(User) returns (UserRequest);
}

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
}

message UserRequest {
  int32 user_id = 1;
}
```

### 生成pb文件

```shell
protoc --go_out=. --go-grpc_out=. usergrpc.proto
```

> ✅ 此时 pb/ 目录下会生成 user.pb.go 和 user_grpc.pb.go

### 创建具体服务

```
// UserServiceImpl 实现 example.UserServiceServer 接口
type UserServiceImpl struct {
	example.UnimplementedUserServiceServer
}

func (s *UserServiceImpl) GetUser(ctx context.Context, req *example.UserRequest) (*example.User, error) {
	fmt.Println("GetUser", req)
	return &example.User{
		Id:    req.UserId,
		Name:  "User_" + string(rune('A'+req.UserId%26)),
		Email: "user" + string(req.UserId) + "@example.com",
	}, nil
}

func (s *UserServiceImpl) AddUser(ctx context.Context, req *example.User) (*example.UserRequest, error) {
	fmt.Println("AddUser", req)
	return &example.UserRequest{UserId: req.Id + 1}, nil
}
```

### 调用服务

```
client, err := grpc.NewClient(":8088", grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
    log.Fatal(err)
}
defer func() {
    _ = client.Close()
}()

// 获取 UserServiceClient
userClient := example.NewUserServiceClient(client)

// 调用 unary RPC
resp, err := userClient.GetUser(context.Background(), &example.UserRequest{UserId: 42})
if err != nil {
    log.Fatal("GetUser error:", err)
}
log.Printf("Unary response: %+v", resp)

// 调用 streaming RPC
stream, err := userClient.AddUser(context.Background(), &example.User{Id: 1, Name: "x", Email: "e"})
if err != nil {
    log.Fatal("StreamUsers error:", err)
}

log.Printf("Stream message: %+v", stream)
```