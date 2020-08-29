package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p-core/host"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"os"
	"time"
)

func chatInputLoop(ctx context.Context, h host.Host, topic *pubsub.Topic, donec chan struct{}) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		msgId := make([]byte, 10)
		_, err := rand.Read(msgId)
		defer func() {
			fmt.Fprintln(os.Stderr, err)
		}()
		if err != nil {
			continue
		}
		now := time.Now().Unix()
		req := &Request{
			Type: Request_SEND_MESSAGE.Enum(),
			SendMessage: &SendMessage{
				Id:      msgId,
				Data:    []byte(msg),
				Created: &now,
			},
		}
		msgBytes, err := req.Marshal()
		if err != nil {
			continue
		}
		err = topic.Publish(ctx, msgBytes)
		if err != nil {
			continue
		}
	}
	donec <- struct{}{}
}
