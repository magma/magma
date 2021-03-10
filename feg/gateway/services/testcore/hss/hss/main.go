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

// This starts the home subscriber server (hss) service.
package main

import (
	"context"
	"flag"
	"log"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/registry"
	"magma/feg/gateway/services/testcore/hss/servicers"
	"magma/feg/gateway/services/testcore/hss/storage"
	"magma/gateway/streamer"
	"magma/orc8r/lib/go/service"
)

func main() {
	flag.Parse()
	srv, err := service.NewServiceWithOptions(registry.ModuleName, registry.MOCK_HSS)
	if err != nil {
		log.Fatalf("Error creating hss service: %s", err)
	}
	config, err := servicers.GetHSSConfig()
	if err != nil {
		log.Printf("Error getting hss config: %s", err)
	}
	store := storage.NewMemorySubscriberStore()
	servicer, err := servicers.NewHomeSubscriberServer(store, config)
	if err != nil {
		log.Fatalf("Error creating home subscriber server: %s", err)
	}
	protos.RegisterHSSConfiguratorServer(srv.GrpcServer, servicer)

	if config.StreamSubscribers {
		streamerClient := streamer.NewStreamerClient(registry.Get())
		l := storage.NewSubscriberListener(store)
		if err = streamerClient.AddListener(l); err != nil {
			log.Printf("Failed to start subscriber streaming: %s", err.Error())
		} else {
			go streamerClient.Stream(l)
		}
	}

	subscribers, err := servicers.GetConfiguredSubscribers()
	if err != nil {
		log.Printf("Could not fetch preconfigured subscribers: %s", err)
	} else {
		// Add preconfigured subscribers
		for _, sub := range subscribers {
			_, err = servicer.AddSubscriber(context.Background(), sub)
			if err != nil {
				log.Printf("Error adding subscriber: %s", err)
			}
		}
	}
	// Start diameter server
	startedChan := make(chan string, 1)
	go func() {
		log.Printf("Starting home subscriber server with configs:\n\t%+v", *servicer.Config)
		err := servicer.Start(startedChan) // blocks
		log.Fatal(err)
	}()
	localAddr := <-startedChan
	log.Printf("Started home subscriber server @ %s", localAddr)

	// Run the service
	err = srv.Run()
	if err != nil {
		log.Fatalf("Error running hss service: %s", err)
	}
}
