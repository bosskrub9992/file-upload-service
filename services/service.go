package services

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"slices"

	"github.com/bosskrub9992/file-upload-service/configs"
	"github.com/bosskrub9992/file-upload-service/responses"

	"github.com/gabriel-vasile/mimetype"
)

type ObjectStorer interface {
	Upload(ctx context.Context, src []byte, bucket, remoteObjectName, srcContentType string) error
}

type EmailSender interface {
	SendUploadSuccess(to string) error
}

type Service struct {
	cfg          *configs.Config
	objectStorer ObjectStorer
	emailSender  EmailSender
}

func New(cfg *configs.Config, objectStorer ObjectStorer, emailSender EmailSender) *Service {
	return &Service{
		cfg:          cfg,
		objectStorer: objectStorer,
		emailSender:  emailSender,
	}
}

func (s Service) Upload(ctx context.Context, file []byte, fileName string) (sendEmailErrorCh chan error, err error) {
	fileExtension := filepath.Ext(fileName)
	fileContentType := detectFileContentType(file)

	if err := validateFileContent(s.cfg.UploadFile.AllowFileExtension, fileContentType, fileExtension); err != nil {
		slog.Error(err.Error())
		return nil, responses.ErrValidateFailed
	}

	if err := s.objectStorer.Upload(ctx, file, s.cfg.UploadFile.Bucket, fileName, fileContentType); err != nil {
		slog.Error(err.Error())
		return nil, responses.ErrAPIFailed
	}

	sendEmailErrorCh = make(chan error)
	go func() {
		if err := s.emailSender.SendUploadSuccess("to@email.com"); err != nil {
			slog.Error(err.Error())
			sendEmailErrorCh <- err
		}
	}()

	return sendEmailErrorCh, nil
}

func detectFileContentType(file []byte) string {
	mtype := mimetype.Detect(file)
	if mtype != nil {
		return mtype.String()
	}
	return ""
}

func validateFileContent(
	allowContentTypeToFileExtension map[string][]string,
	requestContentType, requestFileExtension string,
) error {
	allowFileExtensions, found := allowContentTypeToFileExtension[requestContentType]
	if !found {
		return fmt.Errorf("not allow file content type: '%s'", requestContentType)
	}
	if !slices.Contains(allowFileExtensions, requestFileExtension) {
		return fmt.Errorf("not allow file extension: '%s' for content type: '%s'",
			requestFileExtension,
			requestContentType,
		)
	}
	return nil
}
