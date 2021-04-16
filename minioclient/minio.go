package minioclient

import (
	"context"
	"io"
	"time"

	"github.com/minio/minio-go"
)

type MINIOConfig struct {
	endpoint        string
	accessKeyID     string
	secretAccessKey string
}

func NewMINIOConfig(endpoint, accessKeyID, secretAccessKey string) *MINIOConfig {
	return &MINIOConfig{
		endpoint:        endpoint,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
	}
}

func (this_ *MINIOConfig) NewBucket(name string) (*MINIOClient, error) {
	c, err := minio.New(this_.endpoint, this_.accessKeyID, this_.secretAccessKey, false)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	err = c.MakeBucketWithContext(ctx, name, "us-east-1")
	if err != nil {
		exists, errBucketExists := c.BucketExists(name)
		if errBucketExists != nil && !exists {
			return nil, err
		}
	}
	return &MINIOClient{
		c:          c,
		bucketName: name,
	}, nil
}

type MINIOClient struct {
	c          *minio.Client
	bucketName string
}

func (this_ *MINIOClient) GetBaseClient() *minio.Client {
	return this_.c
}

func (this_ *MINIOClient) FGetObject(objectName, filePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return this_.c.FGetObjectWithContext(ctx, this_.bucketName, objectName, filePath, minio.GetObjectOptions{})
}

func (this_ *MINIOClient) GetObject(objectName string) (*minio.Object, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return this_.c.GetObjectWithContext(ctx, this_.bucketName, objectName, minio.GetObjectOptions{})
}

func (this_ *MINIOClient) FPutObject(objectName, filePath string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return this_.c.FPutObjectWithContext(ctx, this_.bucketName, objectName, filePath, minio.PutObjectOptions{})
}
func (this_ *MINIOClient) StatObject(objectName string) (minio.ObjectInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return this_.c.StatObjectWithContext(ctx, this_.bucketName, objectName, minio.StatObjectOptions{})
}

func (this_ *MINIOClient) IsExistsObject(objectName string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	_, e := this_.c.StatObjectWithContext(ctx, this_.bucketName, objectName, minio.StatObjectOptions{})
	if e == nil {
		return true, nil
	}
	err, ok := e.(minio.ErrorResponse)
	if ok {
		if err.StatusCode == 404 {
			return false, nil
		}
	}
	return false, err
}

func (this_ *MINIOClient) PutObject(objectName string, reader io.Reader) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	return this_.c.PutObjectWithContext(ctx, this_.bucketName, objectName, reader, -1, minio.PutObjectOptions{})
}

func (this_ *MINIOClient) RemoveObjects(objectNames ...string) error {
	if len(objectNames) == 0 {
		return nil
	}
	for _, v := range objectNames {
		err := this_.RemoveObject(v)
		if err == nil {
			return err
		}
	}
	return nil
}

func (this_ *MINIOClient) RemoveObject(objectName string) error {
	return this_.c.RemoveObject(this_.bucketName, objectName)
}
