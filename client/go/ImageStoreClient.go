/*
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package client

import (
	pb "ElephantTrunkArch/ImageStore/protobuff/go"
	"crypto/tls"
	"crypto/x509"
	"io"
	"io/ioutil"
	"time"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// GrpcClient structure
type GrpcClient struct {
	is pb.IsClient
}

const (
	chunkSize = 4095 * 1024 // 4 MB
)

// NewImageStoreClient : This is the constructor to initialize the GrpcClient
func NewImageStoreClient(RootCA string, ClientCert string, ClientKey string, hostname, port string) (*GrpcClient, error) {
	addr := hostname + ":" + port
	glog.Infof("Addr: %s", addr)

	// Read certificate binary
	certPEMBlock, err := ioutil.ReadFile(ClientCert)
	if err != nil {
		glog.Errorf("Failed to Read Client Certificate : %s", err)
		return nil, err
	}

	keyPEMBlock, err := ioutil.ReadFile(ClientKey)
	if err != nil {
		glog.Errorf("Failed to Read Client Key : %s", err)
		return nil, err
	}

	// Load the certificates from binary
	certificate, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		glog.Errorf("Failed to Load ClientKey Pair : %s", err)
		return nil, err
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(RootCA)
	if err != nil {
		glog.Errorf("Failed to Read CA Certificate : %s", err)
		return nil, err
	}

	// Append the Certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		glog.Errorf("Failed to Append Certificate")
		return nil, nil
	}

	// Create the TLS credentials for transport
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		glog.Errorf("Did not connect: %v", err)
		return nil, err
	}
	isClient := pb.NewIsClient(conn)
	return &GrpcClient{is: isClient}, nil
}

// Read is a wrapper around gRPC go client implementation for
// Read interface.
func (pClient *GrpcClient) Read(imgHandle string) ([]byte, error) {
	// Set the gRPC timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// gRPC call
	client, err := pClient.is.Read(ctx, &pb.ReadReq{ReadKeyname: imgHandle})
	if err != nil {
		glog.Errorf("Error: %v", err)
		return nil, err
	}
	var blob []byte
	for {
		c, err := client.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			glog.Errorf("Error while receiving: %v", err)
			return nil, err
		}
		blob = append(blob, c.Chunk...)
	}
	return blob, err
}

// Store is a wrapper around gRPC go client implementation for
// Store interface.
func (pClient *GrpcClient) Store(imgFrame []byte, memType string) (string, error) {
	// Set the gRPC timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := pClient.is.Store(ctx)

	chnk := &pb.StoreReq{}
	//Iterating through the ByteArray for every chunkSize
	for currentByte := 0; currentByte < len(imgFrame); currentByte += chunkSize {
		if currentByte+chunkSize > len(imgFrame) {
			chnk.Chunk = imgFrame[currentByte:len(imgFrame)]
		} else {
			chnk.Chunk = imgFrame[currentByte : currentByte+chunkSize]
		}
		if err := resp.Send(chnk); err != nil {
			imgFrame = nil
		}
	}
	chnk.MemoryType = memType
	replymsg, err := resp.CloseAndRecv()
	if err != nil {
		glog.Errorf("Unexpected Error : %s", err)
	}
	return replymsg.StoreKeyname, err
}

// Remove is a wrapper around gRPC go client implementation for
// Remove interface.
func (pClient *GrpcClient) Remove(imgHandle string) (bool, error) {
	// Set the gRPC timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// gRPC call
	_, err := pClient.is.Remove(ctx, &pb.RemoveReq{RemKeyname: imgHandle})
	if err != nil {
		glog.Errorf("Error: %v", err)
		return false, err
	}
	return true, err
}