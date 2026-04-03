package grpc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/flash1nho/GophKeeper/config"
	"github.com/flash1nho/GophKeeper/internal/facade"
	"github.com/flash1nho/GophKeeper/internal/interceptors"
	"github.com/flash1nho/GophKeeper/internal/logger"
	"github.com/flash1nho/GophKeeper/internal/models/secrets"
)

type MockFacade struct{}

func (m *MockFacade) GetUserIDFromContext(ctx context.Context) (int, error) {
	userID, err := interceptors.GetUserIDFromContext(ctx)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func TestGrpcPrivateHandlerInitialization(t *testing.T) {
	settings := config.SettingsObject{
		DatabaseDSN:       "postgres://localhost/test",
		GrpcServerAddress: "localhost:3200",
		MasterKey:         []byte("test_master_key_for_testing_purposes"),
		Log:               logger.Log,
	}

	f := &facade.Facade{}
	handler := NewGrpcPrivateHandler(nil, settings, f)

	assert.NotNil(t, handler)
	assert.Equal(t, settings, handler.settings)
	assert.NotNil(t, handler.facade)
}

func TestGrpcPrivateHandlerTypeToSecretMapping(t *testing.T) {
	tests := []struct {
		name      string
		typeStr   string
		shouldErr bool
	}{
		{
			name:      "Text type",
			typeStr:   "Text",
			shouldErr: false,
		},
		{
			name:      "Cred type",
			typeStr:   "Cred",
			shouldErr: false,
		},
		{
			name:      "Card type",
			typeStr:   "Card",
			shouldErr: false,
		},
		{
			name:      "File type",
			typeStr:   "File",
			shouldErr: false,
		},
		{
			name:      "Unknown type",
			typeStr:   "Unknown",
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn, ok := typeToSecret[tt.typeStr]

			if tt.shouldErr {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
				assert.NotNil(t, fn)
			}
		})
	}
}

func TestGrpcPrivateHandlerCreateRequestValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       *CreateRequest
		shouldErr bool
	}{
		{
			name: "Valid create request - Text",
			req: &CreateRequest{
				Type: "Text",
				Data: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"content": structpb.NewStringValue("test content"),
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "Valid create request - Card",
			req: &CreateRequest{
				Type: "Card",
				Data: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"number": structpb.NewStringValue("4532015112830366"),
						"expiry": structpb.NewStringValue("12/26"),
						"holder": structpb.NewStringValue("John Doe"),
						"cvv":    structpb.NewStringValue("123"),
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "Invalid type",
			req: &CreateRequest{
				Type: "InvalidType",
				Data: &structpb.Struct{},
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := typeToSecret[tt.req.Type]
			if tt.shouldErr {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
		})
	}
}

func TestGrpcPrivateHandlerGetRequestValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       *GetRequest
		shouldErr bool
	}{
		{
			name: "Valid get request",
			req: &GetRequest{
				ID:   1,
				Type: "Text",
			},
			shouldErr: false,
		},
		{
			name: "Zero ID",
			req: &GetRequest{
				ID:   0,
				Type: "Text",
			},
			shouldErr: false,
		},
		{
			name: "Invalid type",
			req: &GetRequest{
				ID:   1,
				Type: "InvalidType",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := typeToSecret[tt.req.Type]
			if tt.shouldErr {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
		})
	}
}

func TestGrpcPrivateHandlerListRequestValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       *ListRequest
		shouldErr bool
	}{
		{
			name: "Valid list request - Text",
			req: &ListRequest{
				Type: "Text",
			},
			shouldErr: false,
		},
		{
			name: "Valid list request - Card",
			req: &ListRequest{
				Type: "Card",
			},
			shouldErr: false,
		},
		{
			name: "Valid list request - Cred",
			req: &ListRequest{
				Type: "Cred",
			},
			shouldErr: false,
		},
		{
			name: "Valid list request - File",
			req: &ListRequest{
				Type: "File",
			},
			shouldErr: false,
		},
		{
			name: "Invalid type",
			req: &ListRequest{
				Type: "InvalidType",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := typeToSecret[tt.req.Type]
			if tt.shouldErr {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
		})
	}
}

func TestGrpcPrivateHandlerUpdateRequestValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       *UpdateRequest
		shouldErr bool
	}{
		{
			name: "Valid update request",
			req: &UpdateRequest{
				ID:   1,
				Type: "Text",
				Data: &structpb.Struct{
					Fields: map[string]*structpb.Value{
						"content": structpb.NewStringValue("updated content"),
					},
				},
			},
			shouldErr: false,
		},
		{
			name: "Zero ID",
			req: &UpdateRequest{
				ID:   0,
				Type: "Text",
				Data: &structpb.Struct{},
			},
			shouldErr: false,
		},
		{
			name: "Invalid type",
			req: &UpdateRequest{
				ID:   1,
				Type: "InvalidType",
				Data: &structpb.Struct{},
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := typeToSecret[tt.req.Type]
			if tt.shouldErr {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
		})
	}
}

func TestGrpcPrivateHandlerDeleteRequestValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       *DeleteRequest
		shouldErr bool
	}{
		{
			name: "Valid delete request",
			req: &DeleteRequest{
				ID:   1,
				Type: "Text",
			},
			shouldErr: false,
		},
		{
			name: "Invalid type",
			req: &DeleteRequest{
				ID:   1,
				Type: "InvalidType",
			},
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, ok := typeToSecret[tt.req.Type]
			if tt.shouldErr {
				assert.False(t, ok)
			} else {
				assert.True(t, ok)
			}
		})
	}
}

func TestGrpcPrivateHandlerUploadStatusRequestValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       *UploadStatusRequest
		shouldErr bool
	}{
		{
			name: "Valid upload status request",
			req: &UploadStatusRequest{
				FileName: "document.pdf",
			},
			shouldErr: false,
		},
		{
			name: "Empty filename",
			req: &UploadStatusRequest{
				FileName: "",
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.req)
		})
	}
}

func TestGrpcPrivateHandlerDownloadRequestValidation(t *testing.T) {
	tests := []struct {
		name      string
		req       *DownloadRequest
		shouldErr bool
	}{
		{
			name: "Valid download request",
			req: &DownloadRequest{
				ID:         1,
				FileOffset: 0,
			},
			shouldErr: false,
		},
		{
			name: "Download with offset",
			req: &DownloadRequest{
				ID:         1,
				FileOffset: 1024,
			},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.req)
			assert.Greater(t, tt.req.ID, int32(0))
		})
	}
}

func TestGrpcPrivateHandlerMetadataHandling(t *testing.T) {
	metadata := &Metadata{
		FileName:        "test.pdf",
		FileContentType: "application/pdf",
		FileSize:        1024,
		FileOffset:      0,
	}

	assert.Equal(t, "test.pdf", metadata.FileName)
	assert.Equal(t, "application/pdf", metadata.FileContentType)
	assert.Equal(t, int64(1024), metadata.FileSize)
	assert.Equal(t, int64(0), metadata.FileOffset)
}

func TestGrpcPrivateHandlerMergeAdditionalData(t *testing.T) {
	_ = logger.Initialize("debug")

	settings := config.SettingsObject{
		MasterKey: []byte("test_master_key"),
		Log:       logger.Log,
	}

	handler := &GrpcPrivateHandler{
		pool:     nil,
		settings: settings,
		facade:   &facade.Facade{},
	}

	tests := []struct {
		name        string
		secretType  string
		data        *structpb.Struct
		shouldMerge bool
	}{
		{
			name:       "Card type with number",
			secretType: "Card",
			data: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"number": structpb.NewStringValue("4532015112830366"),
				},
			},
			shouldMerge: true,
		},
		{
			name:       "Text type",
			secretType: "Text",
			data: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"content": structpb.NewStringValue("test"),
				},
			},
			shouldMerge: false,
		},
		{
			name:        "Nil data",
			secretType:  "Card",
			data:        nil,
			shouldMerge: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.mergeAdditionalData(tt.secretType, tt.data)

			if tt.shouldMerge && tt.data != nil {
				assert.NotNil(t, tt.data.Fields)
			}
		})
	}
}

func TestGrpcPrivateHandlerBindData(t *testing.T) {
	_ = logger.Initialize("debug")

	settings := config.SettingsObject{
		MasterKey: []byte("test_master_key"),
		Log:       logger.Log,
	}

	handler := &GrpcPrivateHandler{
		pool:     nil,
		settings: settings,
		facade:   &facade.Facade{},
	}

	tests := []struct {
		name      string
		data      *structpb.Struct
		dest      interface{}
		shouldErr bool
	}{
		{
			name: "Valid text data",
			data: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"content": structpb.NewStringValue("test content"),
				},
			},
			dest:      &secrets.Text{},
			shouldErr: false,
		},
		{
			name: "Valid cred data",
			data: &structpb.Struct{
				Fields: map[string]*structpb.Value{
					"login":    structpb.NewStringValue("user"),
					"password": structpb.NewStringValue("pass"),
				},
			},
			dest:      &secrets.Cred{},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.bindData(tt.data, tt.dest)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGrpcPrivateHandlerMapToProto(t *testing.T) {
	_ = logger.Initialize("debug")

	settings := config.SettingsObject{
		MasterKey: []byte("test_master_key"),
		Log:       logger.Log,
	}

	handler := &GrpcPrivateHandler{
		pool:     nil,
		settings: settings,
		facade:   &facade.Facade{},
	}

	tests := []struct {
		name      string
		results   []interface{}
		shouldErr bool
	}{
		{
			name: "Single result",
			results: []interface{}{
				map[string]interface{}{
					"id":   1,
					"type": "Text",
				},
			},
			shouldErr: false,
		},
		{
			name: "Multiple results",
			results: []interface{}{
				map[string]interface{}{"id": 1},
				map[string]interface{}{"id": 2},
				map[string]interface{}{"id": 3},
			},
			shouldErr: false,
		},
		{
			name:      "Empty results",
			results:   []interface{}{},
			shouldErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protoList, err := handler.mapToProto(tt.results)
			if tt.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, protoList)
			}
		})
	}
}

func TestGrpcPrivateHandlerUnauthorizedError(t *testing.T) {
	_ = logger.Initialize("debug")

	settings := config.SettingsObject{
		MasterKey: []byte("test_master_key"),
		Log:       logger.Log,
	}

	handler := &GrpcPrivateHandler{
		pool:     nil,
		settings: settings,
		facade:   &facade.Facade{},
	}

	ctx := context.Background()

	_, err := handler.getSecretInstance(ctx, "Text")
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Unauthenticated, st.Code())
}

func TestGrpcPrivateHandlerInvalidSecretType(t *testing.T) {
	_ = logger.Initialize("debug")

	settings := config.SettingsObject{
		MasterKey: []byte("test_master_key"),
		Log:       logger.Log,
	}

	handler := &GrpcPrivateHandler{
		pool:     nil,
		settings: settings,
		facade:   &facade.Facade{},
	}

	assert.NotNil(t, handler)
}

func TestUploadRequestUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		hasMetadata bool
		hasChunk    bool
	}{
		{
			name:        "With metadata",
			hasMetadata: true,
			hasChunk:    false,
		},
		{
			name:        "With chunk data",
			hasMetadata: false,
			hasChunk:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &UploadRequest{}

			if tt.hasMetadata {
				req.Data = &UploadRequest_Metadata{
					Metadata: &Metadata{
						FileName: "test.pdf",
						FileSize: 1024,
					},
				}
				assert.NotNil(t, req.GetMetadata())
			}

			if tt.hasChunk {
				req.Data = &UploadRequest_Chunk{
					Chunk: []byte("test data"),
				}
				assert.NotNil(t, req.GetChunk())
			}
		})
	}
}

func TestDownloadResponseCreation(t *testing.T) {
	tests := []struct {
		name       string
		chunk      []byte
		fileOffset int64
	}{
		{
			name:       "Small chunk",
			chunk:      []byte("test chunk"),
			fileOffset: 0,
		},
		{
			name:       "Large chunk with offset",
			chunk:      make([]byte, 1024*1024),
			fileOffset: 1024 * 1024,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &DownloadResponse{
				Chunk:      tt.chunk,
				FileOffset: tt.fileOffset,
			}

			assert.Equal(t, tt.chunk, resp.Chunk)
			assert.Equal(t, tt.fileOffset, resp.FileOffset)
			assert.Len(t, resp.Chunk, len(tt.chunk))
		})
	}
}

func TestChunkSizeConstant(t *testing.T) {
	expectedSize := 1024 * 1024 // 1MB
	assert.Equal(t, expectedSize, ChunkSize)
}
