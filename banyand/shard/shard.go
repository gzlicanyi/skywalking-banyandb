// Licensed to Apache Software Foundation (ASF) under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Apache Software Foundation (ASF) licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package shard

import (
	"time"

	"github.com/apache/skywalking-banyandb/banyand/internal/bus"
	"github.com/apache/skywalking-banyandb/banyand/storage"
	"github.com/apache/skywalking-banyandb/pkg/logger"
	"github.com/apache/skywalking-banyandb/pkg/run"
)

var (
	_ bus.MessageListener    = (*Shard)(nil)
	_ run.PreRunner          = (*Shard)(nil)
	_ storage.DataSubscriber = (*Shard)(nil)
	_ storage.DataPublisher  = (*Shard)(nil)
)

type Shard struct {
	log       *logger.Logger
	publisher bus.Publisher
}

func (s Shard) ComponentName() string {
	return "shard"
}

func (s *Shard) Pub(publisher bus.Publisher) error {
	s.publisher = publisher
	return nil
}

func (s *Shard) Sub(subscriber bus.Subscriber) error {
	return subscriber.Subscribe(storage.TraceRaw, s)
}

func (s *Shard) PreRun() error {
	s.log = logger.GetLogger("shard")
	s.log.Info("pre running")
	return nil
}

func (s *Shard) Name() string {
	return "shard"
}

func (s Shard) Rev(message bus.Message) {
	s.log.Info("rev", logger.Any("msg", message.Data()))
	_ = s.publisher.Publish(storage.TraceSharded, bus.NewMessage(bus.MessageID(time.Now().UnixNano()), "sharded message"))
}
