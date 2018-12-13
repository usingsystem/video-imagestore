"""
Copyright (c) 2018 Intel Corporation.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
"""

# Python grpc client implementation

import logging
import argparse
import hashlib
import time
import sys
import os
from ImageStore.client.py.client import \
    GrpcImageStoreClient

logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s : %(levelname)s : \
                    %(name)s : [%(filename)s] :' +
                    '%(funcName)s : in line : [%(lineno)d] : %(message)s')
log = logging.getLogger("GRPC_TEST")


CA_CERT = "/etc/ssl/grpc_internal/ca_certificate.pem"
IM_CLIENT_KEY = "/etc/ssl/imagestore/imagestore_client_key.pem"
IM_CLIENT_CERT = "/etc/ssl/imagestore/imagestore_client_certificate.pem"

if __name__ == '__main__':

    client = GrpcImageStoreClient(IM_CLIENT_CERT, IM_CLIENT_KEY, CA_CERT,
                                  hostname="localhost")

    # Testing Read("imgHandle") gRPC call
    keyname = "inmem_335afcab"
    config = client.Read(keyname)

    # Testing Store("value") gRPC call
    keyname = client.Store(bytes(0x00), 'inmemory')
    print(keyname)

    # Testing Remove("imgHandle") gRPC call
    client.Remove(keyname)