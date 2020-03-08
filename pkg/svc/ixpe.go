package svc

//go:generate mkdir -p ../proto/generated
//go:generate sh -c "protoc --go_out=paths=source_relative,plugins=grpc:../proto/generated -I=../proto ../proto/*.proto"

import (
	"errors"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v6"
	iPXEBuilder "github.com/pojntfx/ipxebuilderd/pkg/proto/generated"
	"github.com/pojntfx/ipxebuilderd/pkg/workers"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/bloom42/libs/rz-go/log"
)

const (
	totalCompileSteps = 2247
)

// IPXEBuilder manages iPXEs.
type IPXEBuilder struct {
	iPXEBuilder.UnimplementedIPXEBuilderServer
	Builder                        *workers.Builder
	S3BucketName, S3HostPortPublic string
	S3Secure                       bool
	s3Client                       *minio.Client
}

func (i *IPXEBuilder) getURL(objectName string) (string, error) {
	// TODO: Specify if minio is secure or not (this should also be set in the `Connect` func with the SSL param and be false by default)
	u, err := url.Parse("http://")
	if i.S3Secure {
		u, err = url.Parse("https://")
	}
	if err != nil {
		return "", err
	}

	u.Path = path.Join(u.Path, i.S3HostPortPublic, i.S3BucketName, objectName)

	return u.String(), nil
}

// Extract extracts the iPXE source code.
func (i *IPXEBuilder) Extract() error {
	return i.Builder.Extract()
}

// ConnectToS3 connects to S3.
func (i *IPXEBuilder) ConnectToS3(s3HostPort, s3AccessKey, s3SecretKey string) error {
	policy := "{\"Version\": \"2012-10-17\",\"Statement\": [{\"Action\": [\"s3:GetObject\"],\"Effect\": \"Allow\",\"Principal\": {\"AWS\": [\"*\"]},\"Resource\": [\"arn:aws:s3:::" + i.S3BucketName + "/*\"],\"Sid\": \"\"}]}"

	minioClient, err := minio.New(s3HostPort, s3AccessKey, s3SecretKey, i.S3Secure)

	if err := minioClient.MakeBucket(i.S3BucketName, "us-east-1"); err != nil {
		exists, err := minioClient.BucketExists(i.S3BucketName)
		if err != nil && !exists {
			return err
		}
	}

	if err := minioClient.SetBucketPolicy(i.S3BucketName, policy); err != nil {
		return err
	}

	i.s3Client = minioClient

	return err
}

func (i *IPXEBuilder) upload(path string) error {
	return nil
}

// Create creates an iPXE.
func (i *IPXEBuilder) Create(req *iPXEBuilder.IPXE, srv iPXEBuilder.IPXEBuilder_CreateServer) error {
	delta := int64(totalCompileSteps)
	stdoutChan, stderrChan, doneChan, errChan := make(chan string), make(chan string), make(chan string), make(chan error)

	go i.Builder.Build(req.GetScript(), req.GetPlatform(), req.GetDriver(), req.GetExtension(), stdoutChan, stderrChan, doneChan, errChan)

	log.Info("Building iPXE")

	for {
		select {
		case <-stdoutChan:
			delta--
			if err := srv.Send(&iPXEBuilder.IPXEStatus{
				Delta: delta,
				URL:   "",
			}); err != nil {
				return err
			}
		case stderrErr := <-stderrChan:
			if !strings.Contains(stderrErr, "ar:") { // `ar` prints to `stderr` for some reason
				return errors.New(stderrErr)
			}

			delta--
			if err := srv.Send(&iPXEBuilder.IPXEStatus{
				Delta: delta,
				URL:   "",
			}); err != nil {
				return err
			}
		case outPath := <-doneChan:
			rest, rawObjectName := filepath.Split(outPath)
			platform := filepath.Base(rest)
			objectName := path.Join(uuid.NewV4().String(), platform, rawObjectName)
			contentType := "application/octet-stream"

			_, err := i.s3Client.FPutObject(i.S3BucketName, objectName, outPath, minio.PutObjectOptions{ContentType: contentType})
			if err != nil {
				return err
			}

			url, err := i.getURL(objectName)
			if err != nil {
				return err
			}

			if err := srv.Send(&iPXEBuilder.IPXEStatus{
				Delta: 0,
				URL:   url,
			}); err != nil {
				return err
			}

			return nil
		case err := <-errChan:
			return err
		}
	}
}
