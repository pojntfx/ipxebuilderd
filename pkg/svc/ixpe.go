package svc

//go:generate mkdir -p ../proto/generated
//go:generate sh -c "protoc --go_out=paths=source_relative,plugins=grpc:../proto/generated -I=../proto ../proto/*.proto"

import (
	"context"
	"errors"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v6"
	"github.com/pojntfx/ipxebuilderd/pkg/constants"
	iPXEBuilder "github.com/pojntfx/ipxebuilderd/pkg/proto/generated"
	"github.com/pojntfx/ipxebuilderd/pkg/workers"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/bloom42/libs/rz-go/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	delta := int64(constants.TotalCompileSteps)
	stdoutChan, stderrChan, doneChan, errChan := make(chan string), make(chan string), make(chan string), make(chan error)

	go i.Builder.Build(req.GetScript(), req.GetPlatform(), req.GetDriver(), req.GetExtension(), stdoutChan, stderrChan, doneChan, errChan)

	log.Info("Building iPXE")

	for {
		select {
		case <-stdoutChan:
			delta--
			if err := srv.Send(&iPXEBuilder.IPXEStatus{
				Delta: delta,
				Id:    "",
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
				Id:    "",
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

			if err := srv.Send(&iPXEBuilder.IPXEStatus{
				Delta: 0,
				Id:    objectName,
			}); err != nil {
				return err
			}

			return nil
		case err := <-errChan:
			return err
		}
	}
}

// List returns the iPXEs.
func (i *IPXEBuilder) List(_ context.Context, req *iPXEBuilder.IPXEBuilderListArgs) (*iPXEBuilder.IPXEBuilderListReply, error) {
	doneChan := make(chan struct{})
	defer close(doneChan)

	objects := i.s3Client.ListObjectsV2(i.S3BucketName, "", true, doneChan)

	var ipxes []*iPXEBuilder.IPXEOut
	for object := range objects {
		url, err := i.getURL(object.Key)
		if err != nil {
			return nil, err
		}

		iPXEPath, iPXEFile := path.Split(object.Key)
		iPXEPathParts := strings.Split(iPXEPath, "/")
		iPXEFileParts := strings.Split(iPXEFile, ".")

		iPXE := iPXEBuilder.IPXEOut{
			Id:        object.Key,
			Platform:  iPXEPathParts[1],
			Driver:    iPXEFileParts[0],
			Extension: iPXEFileParts[1],
			URL:       url,
		}

		ipxes = append(ipxes, &iPXE)
	}

	return &iPXEBuilder.IPXEBuilderListReply{
		IPXEs: ipxes,
	}, nil
}

// Get gets one of the iPXEs.
func (i *IPXEBuilder) Get(_ context.Context, req *iPXEBuilder.IPXEId) (*iPXEBuilder.IPXEOut, error) {
	object, err := i.s3Client.StatObject(i.S3BucketName, req.GetId(), minio.StatObjectOptions{})
	if err != nil {
		statusErr := status.Errorf(codes.Unknown, err.Error())
		if err.Error() == "The specified key does not exist." {
			statusErr = status.Errorf(codes.NotFound, "iPXE not found")
		}

		return nil, statusErr
	}

	url, err := i.getURL(object.Key)
	if err != nil {
		return nil, err
	}

	iPXEPath, iPXEFile := path.Split(object.Key)
	iPXEPathParts := strings.Split(iPXEPath, "/")
	iPXEFileParts := strings.Split(iPXEFile, ".")

	return &iPXEBuilder.IPXEOut{
		Id:        object.Key,
		Platform:  iPXEPathParts[1],
		Driver:    iPXEFileParts[0],
		Extension: iPXEFileParts[1],
		URL:       url,
	}, nil
}

// Delete deletes an iPXE.
func (i *IPXEBuilder) Delete(_ context.Context, req *iPXEBuilder.IPXEId) (*iPXEBuilder.IPXEId, error) {
	id := req.GetId()

	if err := i.s3Client.RemoveObject(i.S3BucketName, id); err != nil {
		return &iPXEBuilder.IPXEId{
			Id: id,
		}, err
	}

	return &iPXEBuilder.IPXEId{
		Id: id,
	}, nil
}
