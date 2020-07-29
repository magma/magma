/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
//
//go:generate protoc -I testcfg --go_out=plugins=grpc,paths=source_relative:testcfg testcfg/test_configs.proto
//
// Package test provides test for gateway managed configuration (mconfig)
package test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"

	"magma/gateway/mconfig"
	"magma/gateway/mconfig/test/testcfg"
	orcprotos "magma/orc8r/lib/go/protos"
)

// JSON to recreate scenario with "static" mconfig file
const (
	testMmconfigJsonV1 = `{
	  "configs_by_key": {
		"service1": {
		  "@type": "type.googleapis.com/magma.mconfig.Service1Config",
		  "str1": "bla",
		  "str2": "192.168.128.1",
		  "uint1": 1,
		  "uint2": 2,
		  "strarr": []
		},
		"service2": {
		  "@type": "type.googleapis.com/magma.mconfig.Service2Config",
		  "str21": "bla bla",
		  "str22": "abcdefg",
		  "uint21": 1,
		  "uint22": 2,
		  "float1": 3.14,
		  "bool1": true
		},
		"does_not_exist_1": {
		  "@type": "type.googleapis.com/magma.mconfig.DoesNotExist",
		  "bla": 1,
		  "blaBla": 1,
		  "logLevel": "INFO"
		}
	  }
	}`
	lenStr22BS         = "Лучше меньше, да лучше"
	testMmconfigJsonV2 = `{
	  "configs_by_key": {
		"service1": {
		  "@type": "type.googleapis.com/magma.mconfig.Service1Config",
		  "str1": "hello world",
		  "str2": "::1",
		  "uint1": 3,
		  "uint2": 4,
		  "strarr": []
		},
		"service2": {
		  "@type": "type.googleapis.com/magma.mconfig.Service2Config",
		  "str21": "Лучше меньше, да лучше",
		  "str22": "aBcDeFg",
		  "uint21": 21,
		  "uint22": 22,
		  "float1": 14.3,
		  "bool1": false
		}
	  }
	}`
)

func TestGatewayMconfigRefresh(t *testing.T) {
	mconfig.StopRefreshTicker() // stop non-test config refresh

	// Create tmp mconfig test file
	tmpfile, err := ioutil.TempFile("", mconfig.MconfigFileName)
	if err != nil {
		t.Fatal(err)
	}
	// Write V1 style marshaled configs
	if _, err = tmpfile.Write([]byte(testMmconfigJsonV1)); err != nil {
		t.Fatal(err)
	}
	mcpath := tmpfile.Name()
	if err = tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
	t.Logf("Created gateway config file: %s", mcpath)

	// Start configs refresh ticker
	ticker := time.NewTicker(time.Millisecond * 50)
	go func() {
		for {
			<-ticker.C
			refreshErr := mconfig.RefreshConfigsFrom(mcpath)
			if refreshErr != nil {
				t.Error(refreshErr)
			}
		}
	}()

	// Cleanup @ the end of the test
	defer func() {
		ticker.Stop()
		time.Sleep(time.Millisecond * 20)
		t.Logf("Remove temporary gateway config file: %s", mcpath)
		os.Remove(mcpath)
	}()

	time.Sleep(time.Millisecond * 120)
	s1cfg := &testcfg.Service1Config{}
	err = mconfig.GetServiceConfigs("service1", s1cfg)
	if err != nil {
		t.Fatal(err)
	}
	expectedStr2 := "192.168.128.1"
	if s1cfg.Str2 != expectedStr2 {
		t.Fatalf("service1 String2 Mismatch %s != %s", s1cfg.Str2, expectedStr2)
	}
	mc := mconfig.GetGatewayConfigs()
	expectedStr2 = "192.123.155.0"
	s1cfg.Str2 = expectedStr2
	mc.ConfigsByKey["service1"], err = ptypes.MarshalAny(s1cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Test marshaling of new configs
	mc.ConfigsByKey["service2"], err = ptypes.MarshalAny(
		&testcfg.Service2Config{
			Str21: "str21",
			Str22: "str22",
		})

	marshaled, err := orcprotos.MarshalMconfig(mc)
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(mcpath, marshaled, os.ModePerm)
	if err != nil {
		t.Fatal(err)

	}
	s1cfg = &testcfg.Service1Config{}

	// Wait for refresh
	time.Sleep(time.Millisecond * 120)
	err = mconfig.GetServiceConfigs("service1", s1cfg)
	if err != nil {
		t.Fatal(err)
	}
	if s1cfg.Str2 != expectedStr2 {
		t.Fatalf("service1 String2 Mismatch %s != %s", s1cfg.Str2, expectedStr2)
	}

	s2cfg := &testcfg.Service2Config{}
	err = mconfig.GetServiceConfigs("service2", s2cfg)
	if err != nil {
		t.Fatal(err)
	}
	if s2cfg.Str21 != "str21" {
		t.Fatalf(
			"service2 Configs Str21 Mismatch %s != str21", s2cfg.Str21)
	}
	if s2cfg.Str22 != "str22" {
		t.Fatalf(
			"service2 Configs Str22 Mismatch %s != str22", s2cfg.Str22)
	}

	// Test API's type enforcement/safety
	err = mconfig.GetServiceConfigs("service1", s2cfg)
	if err == nil {
		t.Fatal("Expected Error Getting service1 configs into Service2Config type")
	}
	s2cfg = &testcfg.Service2Config{}

	// test V2 - 'configsByKey' encoding version
	if err = ioutil.WriteFile(mcpath, []byte(testMmconfigJsonV2), os.ModePerm); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond * 120)
	err = mconfig.GetServiceConfigs("service2", s2cfg)
	if err != nil {
		t.Fatal(err)
	}
	if s2cfg.GetStr21() != lenStr22BS {
		t.Fatalf("service2 Configs Str21 Mismatch %s != %s", s2cfg.GetStr21(), lenStr22BS)
	}
}
