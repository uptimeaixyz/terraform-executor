// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: api/proto/executor.proto

package executor

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Executor_AppendCode_FullMethodName      = "/executor.Executor/AppendCode"
	Executor_Plan_FullMethodName            = "/executor.Executor/Plan"
	Executor_Apply_FullMethodName           = "/executor.Executor/Apply"
	Executor_Destroy_FullMethodName         = "/executor.Executor/Destroy"
	Executor_GetStateList_FullMethodName    = "/executor.Executor/GetStateList"
	Executor_ClearCode_FullMethodName       = "/executor.Executor/ClearCode"
	Executor_CreateContext_FullMethodName   = "/executor.Executor/CreateContext"
	Executor_DeleteContext_FullMethodName   = "/executor.Executor/DeleteContext"
	Executor_CreateWorkspace_FullMethodName = "/executor.Executor/CreateWorkspace"
	Executor_DeleteWorkspace_FullMethodName = "/executor.Executor/DeleteWorkspace"
	Executor_AddProviders_FullMethodName    = "/executor.Executor/AddProviders"
	Executor_AddSecretEnv_FullMethodName    = "/executor.Executor/AddSecretEnv"
	Executor_AddSecretVar_FullMethodName    = "/executor.Executor/AddSecretVar"
	Executor_ClearProviders_FullMethodName  = "/executor.Executor/ClearProviders"
	Executor_ClearWorkspace_FullMethodName  = "/executor.Executor/ClearWorkspace"
	Executor_ClearSecretVars_FullMethodName = "/executor.Executor/ClearSecretVars"
)

// ExecutorClient is the client API for Executor service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// The Executor service definition.
type ExecutorClient interface {
	// Appends code to the Terraform configuration.
	AppendCode(ctx context.Context, in *AppendCodeRequest, opts ...grpc.CallOption) (*AppendCodeResponse, error)
	// Generates a Terraform plan and returns the result.
	Plan(ctx context.Context, in *PlanRequest, opts ...grpc.CallOption) (*PlanResponse, error)
	// Applies the Terraform plan and returns the execution result.
	Apply(ctx context.Context, in *ApplyRequest, opts ...grpc.CallOption) (*ApplyResponse, error)
	// Destroys the Terraform-managed infrastructure and returns the result.
	Destroy(ctx context.Context, in *DestroyRequest, opts ...grpc.CallOption) (*DestroyResponse, error)
	// Retrieves the Terraform state list.
	GetStateList(ctx context.Context, in *GetStateListRequest, opts ...grpc.CallOption) (*GetStateListResponse, error)
	// Clears the Terraform files.
	ClearCode(ctx context.Context, in *ClearCodeRequest, opts ...grpc.CallOption) (*ClearCodeResponse, error)
	// Creates a new context.
	CreateContext(ctx context.Context, in *CreateContextRequest, opts ...grpc.CallOption) (*CreateContextResponse, error)
	// Deletes a context.
	DeleteContext(ctx context.Context, in *DeleteContextRequest, opts ...grpc.CallOption) (*DeleteContextResponse, error)
	// Creates a workspace within a context.
	CreateWorkspace(ctx context.Context, in *CreateWorkspaceRequest, opts ...grpc.CallOption) (*CreateWorkspaceResponse, error)
	// Deletes a workspace within a context.
	DeleteWorkspace(ctx context.Context, in *DeleteWorkspaceRequest, opts ...grpc.CallOption) (*DeleteWorkspaceResponse, error)
	// Adds providers to the Terraform configuration.
	AddProviders(ctx context.Context, in *AddProvidersRequest, opts ...grpc.CallOption) (*AddProvidersResponse, error)
	// Adds a secret env to the Terraform configuration.
	AddSecretEnv(ctx context.Context, in *AddSecretEnvRequest, opts ...grpc.CallOption) (*AddSecretEnvResponse, error)
	// Adds a secret variable to the Terraform configuration.
	AddSecretVar(ctx context.Context, in *AddSecretVarRequest, opts ...grpc.CallOption) (*AddSecretVarResponse, error)
	// Clears the providers from the Terraform configuration.
	ClearProviders(ctx context.Context, in *ClearProvidersRequest, opts ...grpc.CallOption) (*ClearProvidersResponse, error)
	// Clears the workspace.
	ClearWorkspace(ctx context.Context, in *ClearWorkspaceRequest, opts ...grpc.CallOption) (*ClearWorkspaceResponse, error)
	// Clears the secret vars from the Terraform configuration.
	ClearSecretVars(ctx context.Context, in *ClearSecretVarsRequest, opts ...grpc.CallOption) (*ClearSecretVarsResponse, error)
}

type executorClient struct {
	cc grpc.ClientConnInterface
}

func NewExecutorClient(cc grpc.ClientConnInterface) ExecutorClient {
	return &executorClient{cc}
}

func (c *executorClient) AppendCode(ctx context.Context, in *AppendCodeRequest, opts ...grpc.CallOption) (*AppendCodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AppendCodeResponse)
	err := c.cc.Invoke(ctx, Executor_AppendCode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) Plan(ctx context.Context, in *PlanRequest, opts ...grpc.CallOption) (*PlanResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PlanResponse)
	err := c.cc.Invoke(ctx, Executor_Plan_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) Apply(ctx context.Context, in *ApplyRequest, opts ...grpc.CallOption) (*ApplyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ApplyResponse)
	err := c.cc.Invoke(ctx, Executor_Apply_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) Destroy(ctx context.Context, in *DestroyRequest, opts ...grpc.CallOption) (*DestroyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DestroyResponse)
	err := c.cc.Invoke(ctx, Executor_Destroy_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) GetStateList(ctx context.Context, in *GetStateListRequest, opts ...grpc.CallOption) (*GetStateListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetStateListResponse)
	err := c.cc.Invoke(ctx, Executor_GetStateList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) ClearCode(ctx context.Context, in *ClearCodeRequest, opts ...grpc.CallOption) (*ClearCodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClearCodeResponse)
	err := c.cc.Invoke(ctx, Executor_ClearCode_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) CreateContext(ctx context.Context, in *CreateContextRequest, opts ...grpc.CallOption) (*CreateContextResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateContextResponse)
	err := c.cc.Invoke(ctx, Executor_CreateContext_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) DeleteContext(ctx context.Context, in *DeleteContextRequest, opts ...grpc.CallOption) (*DeleteContextResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteContextResponse)
	err := c.cc.Invoke(ctx, Executor_DeleteContext_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) CreateWorkspace(ctx context.Context, in *CreateWorkspaceRequest, opts ...grpc.CallOption) (*CreateWorkspaceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateWorkspaceResponse)
	err := c.cc.Invoke(ctx, Executor_CreateWorkspace_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) DeleteWorkspace(ctx context.Context, in *DeleteWorkspaceRequest, opts ...grpc.CallOption) (*DeleteWorkspaceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteWorkspaceResponse)
	err := c.cc.Invoke(ctx, Executor_DeleteWorkspace_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) AddProviders(ctx context.Context, in *AddProvidersRequest, opts ...grpc.CallOption) (*AddProvidersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddProvidersResponse)
	err := c.cc.Invoke(ctx, Executor_AddProviders_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) AddSecretEnv(ctx context.Context, in *AddSecretEnvRequest, opts ...grpc.CallOption) (*AddSecretEnvResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddSecretEnvResponse)
	err := c.cc.Invoke(ctx, Executor_AddSecretEnv_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) AddSecretVar(ctx context.Context, in *AddSecretVarRequest, opts ...grpc.CallOption) (*AddSecretVarResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddSecretVarResponse)
	err := c.cc.Invoke(ctx, Executor_AddSecretVar_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) ClearProviders(ctx context.Context, in *ClearProvidersRequest, opts ...grpc.CallOption) (*ClearProvidersResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClearProvidersResponse)
	err := c.cc.Invoke(ctx, Executor_ClearProviders_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) ClearWorkspace(ctx context.Context, in *ClearWorkspaceRequest, opts ...grpc.CallOption) (*ClearWorkspaceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClearWorkspaceResponse)
	err := c.cc.Invoke(ctx, Executor_ClearWorkspace_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *executorClient) ClearSecretVars(ctx context.Context, in *ClearSecretVarsRequest, opts ...grpc.CallOption) (*ClearSecretVarsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClearSecretVarsResponse)
	err := c.cc.Invoke(ctx, Executor_ClearSecretVars_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ExecutorServer is the server API for Executor service.
// All implementations must embed UnimplementedExecutorServer
// for forward compatibility.
//
// The Executor service definition.
type ExecutorServer interface {
	// Appends code to the Terraform configuration.
	AppendCode(context.Context, *AppendCodeRequest) (*AppendCodeResponse, error)
	// Generates a Terraform plan and returns the result.
	Plan(context.Context, *PlanRequest) (*PlanResponse, error)
	// Applies the Terraform plan and returns the execution result.
	Apply(context.Context, *ApplyRequest) (*ApplyResponse, error)
	// Destroys the Terraform-managed infrastructure and returns the result.
	Destroy(context.Context, *DestroyRequest) (*DestroyResponse, error)
	// Retrieves the Terraform state list.
	GetStateList(context.Context, *GetStateListRequest) (*GetStateListResponse, error)
	// Clears the Terraform files.
	ClearCode(context.Context, *ClearCodeRequest) (*ClearCodeResponse, error)
	// Creates a new context.
	CreateContext(context.Context, *CreateContextRequest) (*CreateContextResponse, error)
	// Deletes a context.
	DeleteContext(context.Context, *DeleteContextRequest) (*DeleteContextResponse, error)
	// Creates a workspace within a context.
	CreateWorkspace(context.Context, *CreateWorkspaceRequest) (*CreateWorkspaceResponse, error)
	// Deletes a workspace within a context.
	DeleteWorkspace(context.Context, *DeleteWorkspaceRequest) (*DeleteWorkspaceResponse, error)
	// Adds providers to the Terraform configuration.
	AddProviders(context.Context, *AddProvidersRequest) (*AddProvidersResponse, error)
	// Adds a secret env to the Terraform configuration.
	AddSecretEnv(context.Context, *AddSecretEnvRequest) (*AddSecretEnvResponse, error)
	// Adds a secret variable to the Terraform configuration.
	AddSecretVar(context.Context, *AddSecretVarRequest) (*AddSecretVarResponse, error)
	// Clears the providers from the Terraform configuration.
	ClearProviders(context.Context, *ClearProvidersRequest) (*ClearProvidersResponse, error)
	// Clears the workspace.
	ClearWorkspace(context.Context, *ClearWorkspaceRequest) (*ClearWorkspaceResponse, error)
	// Clears the secret vars from the Terraform configuration.
	ClearSecretVars(context.Context, *ClearSecretVarsRequest) (*ClearSecretVarsResponse, error)
	mustEmbedUnimplementedExecutorServer()
}

// UnimplementedExecutorServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedExecutorServer struct{}

func (UnimplementedExecutorServer) AppendCode(context.Context, *AppendCodeRequest) (*AppendCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendCode not implemented")
}
func (UnimplementedExecutorServer) Plan(context.Context, *PlanRequest) (*PlanResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Plan not implemented")
}
func (UnimplementedExecutorServer) Apply(context.Context, *ApplyRequest) (*ApplyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Apply not implemented")
}
func (UnimplementedExecutorServer) Destroy(context.Context, *DestroyRequest) (*DestroyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Destroy not implemented")
}
func (UnimplementedExecutorServer) GetStateList(context.Context, *GetStateListRequest) (*GetStateListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStateList not implemented")
}
func (UnimplementedExecutorServer) ClearCode(context.Context, *ClearCodeRequest) (*ClearCodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearCode not implemented")
}
func (UnimplementedExecutorServer) CreateContext(context.Context, *CreateContextRequest) (*CreateContextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateContext not implemented")
}
func (UnimplementedExecutorServer) DeleteContext(context.Context, *DeleteContextRequest) (*DeleteContextResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteContext not implemented")
}
func (UnimplementedExecutorServer) CreateWorkspace(context.Context, *CreateWorkspaceRequest) (*CreateWorkspaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWorkspace not implemented")
}
func (UnimplementedExecutorServer) DeleteWorkspace(context.Context, *DeleteWorkspaceRequest) (*DeleteWorkspaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteWorkspace not implemented")
}
func (UnimplementedExecutorServer) AddProviders(context.Context, *AddProvidersRequest) (*AddProvidersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProviders not implemented")
}
func (UnimplementedExecutorServer) AddSecretEnv(context.Context, *AddSecretEnvRequest) (*AddSecretEnvResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSecretEnv not implemented")
}
func (UnimplementedExecutorServer) AddSecretVar(context.Context, *AddSecretVarRequest) (*AddSecretVarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddSecretVar not implemented")
}
func (UnimplementedExecutorServer) ClearProviders(context.Context, *ClearProvidersRequest) (*ClearProvidersResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearProviders not implemented")
}
func (UnimplementedExecutorServer) ClearWorkspace(context.Context, *ClearWorkspaceRequest) (*ClearWorkspaceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearWorkspace not implemented")
}
func (UnimplementedExecutorServer) ClearSecretVars(context.Context, *ClearSecretVarsRequest) (*ClearSecretVarsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClearSecretVars not implemented")
}
func (UnimplementedExecutorServer) mustEmbedUnimplementedExecutorServer() {}
func (UnimplementedExecutorServer) testEmbeddedByValue()                  {}

// UnsafeExecutorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ExecutorServer will
// result in compilation errors.
type UnsafeExecutorServer interface {
	mustEmbedUnimplementedExecutorServer()
}

func RegisterExecutorServer(s grpc.ServiceRegistrar, srv ExecutorServer) {
	// If the following call pancis, it indicates UnimplementedExecutorServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Executor_ServiceDesc, srv)
}

func _Executor_AppendCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).AppendCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_AppendCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).AppendCode(ctx, req.(*AppendCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_Plan_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlanRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).Plan(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_Plan_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).Plan(ctx, req.(*PlanRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_Apply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ApplyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).Apply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_Apply_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).Apply(ctx, req.(*ApplyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_Destroy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DestroyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).Destroy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_Destroy_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).Destroy(ctx, req.(*DestroyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_GetStateList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStateListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).GetStateList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_GetStateList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).GetStateList(ctx, req.(*GetStateListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_ClearCode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).ClearCode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_ClearCode_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).ClearCode(ctx, req.(*ClearCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_CreateContext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateContextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).CreateContext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_CreateContext_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).CreateContext(ctx, req.(*CreateContextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_DeleteContext_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteContextRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).DeleteContext(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_DeleteContext_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).DeleteContext(ctx, req.(*DeleteContextRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_CreateWorkspace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWorkspaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).CreateWorkspace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_CreateWorkspace_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).CreateWorkspace(ctx, req.(*CreateWorkspaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_DeleteWorkspace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteWorkspaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).DeleteWorkspace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_DeleteWorkspace_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).DeleteWorkspace(ctx, req.(*DeleteWorkspaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_AddProviders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddProvidersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).AddProviders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_AddProviders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).AddProviders(ctx, req.(*AddProvidersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_AddSecretEnv_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddSecretEnvRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).AddSecretEnv(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_AddSecretEnv_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).AddSecretEnv(ctx, req.(*AddSecretEnvRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_AddSecretVar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddSecretVarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).AddSecretVar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_AddSecretVar_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).AddSecretVar(ctx, req.(*AddSecretVarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_ClearProviders_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearProvidersRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).ClearProviders(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_ClearProviders_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).ClearProviders(ctx, req.(*ClearProvidersRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_ClearWorkspace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearWorkspaceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).ClearWorkspace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_ClearWorkspace_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).ClearWorkspace(ctx, req.(*ClearWorkspaceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Executor_ClearSecretVars_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClearSecretVarsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ExecutorServer).ClearSecretVars(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Executor_ClearSecretVars_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ExecutorServer).ClearSecretVars(ctx, req.(*ClearSecretVarsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Executor_ServiceDesc is the grpc.ServiceDesc for Executor service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Executor_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "executor.Executor",
	HandlerType: (*ExecutorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppendCode",
			Handler:    _Executor_AppendCode_Handler,
		},
		{
			MethodName: "Plan",
			Handler:    _Executor_Plan_Handler,
		},
		{
			MethodName: "Apply",
			Handler:    _Executor_Apply_Handler,
		},
		{
			MethodName: "Destroy",
			Handler:    _Executor_Destroy_Handler,
		},
		{
			MethodName: "GetStateList",
			Handler:    _Executor_GetStateList_Handler,
		},
		{
			MethodName: "ClearCode",
			Handler:    _Executor_ClearCode_Handler,
		},
		{
			MethodName: "CreateContext",
			Handler:    _Executor_CreateContext_Handler,
		},
		{
			MethodName: "DeleteContext",
			Handler:    _Executor_DeleteContext_Handler,
		},
		{
			MethodName: "CreateWorkspace",
			Handler:    _Executor_CreateWorkspace_Handler,
		},
		{
			MethodName: "DeleteWorkspace",
			Handler:    _Executor_DeleteWorkspace_Handler,
		},
		{
			MethodName: "AddProviders",
			Handler:    _Executor_AddProviders_Handler,
		},
		{
			MethodName: "AddSecretEnv",
			Handler:    _Executor_AddSecretEnv_Handler,
		},
		{
			MethodName: "AddSecretVar",
			Handler:    _Executor_AddSecretVar_Handler,
		},
		{
			MethodName: "ClearProviders",
			Handler:    _Executor_ClearProviders_Handler,
		},
		{
			MethodName: "ClearWorkspace",
			Handler:    _Executor_ClearWorkspace_Handler,
		},
		{
			MethodName: "ClearSecretVars",
			Handler:    _Executor_ClearSecretVars_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/proto/executor.proto",
}
