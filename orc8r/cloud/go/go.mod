// Copyright (c) Facebook, Inc. and its affiliates.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
module magma/orc8r/cloud/go

replace (
	magma/orc8r/lib/go => ../../lib/go
	magma/orc8r/lib/go/protos => ../../lib/go/protos
)

require (
	github.com/Masterminds/squirrel v1.1.1-0.20190513200039-d13326f0be73
	github.com/aws/aws-sdk-go v1.19.6
	github.com/coreos/go-systemd v0.0.0-20181031085051-9002847aa142
	github.com/facebookincubator/ent v0.0.0-20191128071424-29c7b0a0d805
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/go-openapi/analysis v0.18.0 // indirect
	github.com/go-openapi/errors v0.18.0
	github.com/go-openapi/jsonpointer v0.18.0 // indirect
	github.com/go-openapi/jsonreference v0.18.0 // indirect
	github.com/go-openapi/loads v0.18.0 // indirect
	github.com/go-openapi/spec v0.18.0 // indirect
	github.com/go-openapi/strfmt v0.18.0
	github.com/go-openapi/swag v0.18.0
	github.com/go-openapi/validate v0.18.0
	github.com/go-sql-driver/mysql v1.4.1-0.20190510102335-877a9775f068
	github.com/go-swagger/go-swagger v0.18.0
	github.com/go-swagger/scan-repo-boundary v0.0.0-20180623220736-973b3573c013 // indirect
	github.com/godbus/dbus v0.0.0-20181101234600-2ff6f7ffd60f // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.3.3
	github.com/golang/snappy v0.0.0-20180518054509-2e65f85255db // indirect
	github.com/google/uuid v1.1.1
	github.com/gorilla/handlers v1.4.0 // indirect
	github.com/hpcloud/tail v1.0.0
	github.com/labstack/echo v0.0.0-20181123063414-c54d9e8eed6c
	github.com/labstack/gommon v0.2.8 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0
	github.com/lib/pq v1.2.0
	github.com/mattn/go-sqlite3 v1.11.0
	github.com/olivere/elastic/v7 v7.0.6
	github.com/onsi/ginkgo v1.7.0 // indirect
	github.com/onsi/gomega v1.4.3 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/alertmanager v0.17.0
	github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/prometheus/client_model v0.2.0
	github.com/prometheus/common v0.2.0
	github.com/prometheus/procfs v0.0.0-20190117184657-bf6a532e95b1
	github.com/prometheus/prometheus v0.0.0-20190115164134-b639fe140c1f
	github.com/spf13/cobra v0.0.3
	github.com/spf13/viper v1.3.1 // indirect
	github.com/stretchr/testify v1.4.0
	github.com/thoas/go-funk v0.4.0
	github.com/toqueteos/webbrowser v1.1.0 // indirect
	github.com/vektra/mockery v0.0.0-20181123154057-e78b021dcbb5
	golang.org/x/lint v0.0.0-20190409202823-959b441ac422
	golang.org/x/net v0.0.0-20200202094626-16171245cfb2
	golang.org/x/tools v0.0.0-20191012152004-8de300cfc20a
	google.golang.org/api v0.3.1 // indirect
	google.golang.org/grpc v1.27.1
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0
	gopkg.in/yaml.v2 v2.2.8
	magma/orc8r/lib/go v0.0.0
	magma/orc8r/lib/go/protos v0.0.0
)

go 1.12
