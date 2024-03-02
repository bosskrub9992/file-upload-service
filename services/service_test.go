package services

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/bosskrub9992/file-upload-service/configs"
	"github.com/bosskrub9992/file-upload-service/responses"
	"github.com/bosskrub9992/file-upload-service/services/mocks"

	"github.com/stretchr/testify/mock"
)

func TestService_Upload(t *testing.T) {
	ctx := context.Background()
	cfg := configs.New()
	objectStorer := mocks.NewObjectStorer(t)
	emailSender := mocks.NewEmailSender(t)
	jpegIMG, err := getBytesFromImage("./images/please.jpeg")
	if err != nil {
		t.Error(err)
		return
	}
	type fields struct {
		cfg          *configs.Config
		objectStorer ObjectStorer
		emailSender  EmailSender
	}
	type args struct {
		ctx      context.Context
		file     []byte
		fileName string
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		mockFunc   func()
		assertFunc func(sendEmailErrorCh chan error, err error)
	}{
		{
			name: "validate failed, should respond validate failed",
			fields: fields{
				cfg:          cfg,
				objectStorer: objectStorer,
				emailSender:  emailSender,
			},
			args: args{
				ctx:      ctx,
				file:     jpegIMG,
				fileName: "please.png",
			},
			wantErr:  true,
			mockFunc: func() {},
			assertFunc: func(sendEmailErrorCh chan error, err error) {
				if !errors.As(err, &responses.ErrValidateFailed) {
					t.Error("wrong error response")
				}
			},
		},
		{
			name: "upload failed, should respond api failed",
			fields: fields{
				cfg:          cfg,
				objectStorer: objectStorer,
				emailSender:  emailSender,
			},
			args: args{
				ctx:      ctx,
				file:     jpegIMG,
				fileName: "please.jpeg",
			},
			wantErr: true,
			mockFunc: func() {
				objectStorer.EXPECT().Upload(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(errors.New("mock error occurred")).
					Once()
			},
			assertFunc: func(sendEmailErrorCh chan error, err error) {
				if !errors.As(err, &responses.ErrAPIFailed) {
					t.Error("wrong error response")
				}
			},
		},
		{
			name: "upload success, should send email",
			fields: fields{
				cfg:          cfg,
				objectStorer: objectStorer,
				emailSender:  emailSender,
			},
			args: args{
				ctx:      ctx,
				file:     jpegIMG,
				fileName: "please.jpeg",
			},
			wantErr: false,
			mockFunc: func() {
				objectStorer.EXPECT().Upload(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil).
					Once()
				emailSender.EXPECT().SendUploadSuccess(mock.Anything).
					Return(nil).
					Once()
			},
			assertFunc: func(sendEmailErrorCh chan error, err error) {
				select {
				case <-time.After(1 * time.Second):
					// wait for 1 second to ensure that system will do function sendUploadSuccess
				case sendEmailError := <-sendEmailErrorCh:
					if sendEmailError != nil {
						t.Errorf("should have no error when sending email but got 1: %+v", sendEmailError)
					}
				}
			},
		},
		{
			name: "upload success but send email failed, should have error from sendEmailErrorCh",
			fields: fields{
				cfg:          cfg,
				objectStorer: objectStorer,
				emailSender:  emailSender,
			},
			args: args{
				ctx:      ctx,
				file:     jpegIMG,
				fileName: "please.jpeg",
			},
			wantErr: false,
			mockFunc: func() {
				objectStorer.EXPECT().Upload(ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
					Return(nil).
					Once()
				emailSender.EXPECT().SendUploadSuccess(mock.Anything).
					Return(errors.New("send email failed")).
					Once()
			},
			assertFunc: func(sendEmailErrorCh chan error, err error) {
				select {
				case <-time.After(1 * time.Second):
					// wait for 1 second to ensure that system will do function sendUploadSuccess
				case sendEmailError := <-sendEmailErrorCh:
					if sendEmailError == nil {
						t.Error("should have error but got no error")
					}
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(
				tt.fields.cfg,
				tt.fields.objectStorer,
				tt.fields.emailSender,
			)
			tt.mockFunc()
			sendEmailErrorCh, err := s.Upload(tt.args.ctx, tt.args.file, tt.args.fileName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Upload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.assertFunc(sendEmailErrorCh, err)
		})
	}
}

func Test_validateFileContent(t *testing.T) {
	allowContentTypeToFileExtension := map[string][]string{
		"image/jpeg": {".jpeg", ".jpg", ".JPEG"},
		"image/png":  {".png", ".PNG"},
		"image/heic": {".heic"},
	}
	type args struct {
		allowContentTypeToFileExtension map[string][]string
		requestContentType              string
		requestFileExtension            string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "image/jpeg vs .jpeg",
			args: args{
				allowContentTypeToFileExtension: allowContentTypeToFileExtension,
				requestContentType:              "image/jpeg",
				requestFileExtension:            ".jpeg",
			},
			wantErr: false,
		},
		{
			name: "image/jpeg vs .png",
			args: args{
				allowContentTypeToFileExtension: allowContentTypeToFileExtension,
				requestContentType:              "image/jpeg",
				requestFileExtension:            ".png",
			},
			wantErr: true,
		},
		{
			name: "image/png vs .png",
			args: args{
				allowContentTypeToFileExtension: allowContentTypeToFileExtension,
				requestContentType:              "image/png",
				requestFileExtension:            ".png",
			},
			wantErr: false,
		},
		{
			name: "image/png vs .jpeg",
			args: args{
				allowContentTypeToFileExtension: allowContentTypeToFileExtension,
				requestContentType:              "image/png",
				requestFileExtension:            ".jpeg",
			},
			wantErr: true,
		},
		{
			name: "image/heic vs .heic",
			args: args{
				allowContentTypeToFileExtension: allowContentTypeToFileExtension,
				requestContentType:              "image/heic",
				requestFileExtension:            ".heic",
			},
			wantErr: false,
		},
		{
			name: "image/heic vs .jpeg",
			args: args{
				allowContentTypeToFileExtension: allowContentTypeToFileExtension,
				requestContentType:              "image/heic",
				requestFileExtension:            ".jpeg",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateFileContent(tt.args.allowContentTypeToFileExtension, tt.args.requestContentType, tt.args.requestFileExtension); (err != nil) != tt.wantErr {
				t.Errorf("validateFileContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_detectFileContentType(t *testing.T) {
	heicIMG, err := getBytesFromImage("./images/please.heic")
	if err != nil {
		t.Error(err)
		return
	}
	jpegIMG, err := getBytesFromImage("./images/please.jpeg")
	if err != nil {
		t.Error(err)
		return
	}
	jpgIMG, err := getBytesFromImage("./images/please.jpg")
	if err != nil {
		t.Error(err)
		return
	}
	pngIMG, err := getBytesFromImage("./images/please.png")
	if err != nil {
		t.Error(err)
		return
	}

	type args struct {
		file []byte
	}
	tests := []struct {
		name                string
		args                args
		wantFileContentType string
	}{
		{
			name: "test .jpeg",
			args: args{
				file: jpegIMG,
			},
			wantFileContentType: "image/jpeg",
		},
		{
			name: "test .jpg",
			args: args{
				file: jpgIMG,
			},
			wantFileContentType: "image/jpeg",
		},
		{
			name: "test .png",
			args: args{
				file: pngIMG,
			},
			wantFileContentType: "image/png",
		},
		{
			name: "test .heic",
			args: args{
				file: heicIMG,
			},
			wantFileContentType: "image/heic",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFileContentType := detectFileContentType(tt.args.file)
			if gotFileContentType != tt.wantFileContentType {
				t.Errorf("detectFileContent() gotFileContentType = %v, want %v", gotFileContentType, tt.wantFileContentType)
			}
		})
	}
}

func getBytesFromImage(path string) ([]byte, error) {
	img, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer img.Close()
	srcByte, err := io.ReadAll(img)
	if err != nil {
		return nil, err
	}
	return srcByte, nil
}
