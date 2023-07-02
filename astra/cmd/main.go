package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/andibalo/ramein/astra/internal/config"
	"github.com/andibalo/ramein/astra/internal/db"
	"github.com/andibalo/ramein/astra/internal/logger"
	"github.com/andibalo/ramein/astra/internal/model"
	"github.com/andibalo/ramein/astra/internal/repository"
	"github.com/andibalo/ramein/commons/kafka"
	"github.com/centrifugal/centrifuge"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gocql/gocql"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type clientMessage struct {
	Timestamp int64  `json:"timestamp"`
	Input     string `json:"input"`
}

func handleLog(e centrifuge.LogEntry) {
	log.Printf("%s: %v", e.Message, e.Fields)
}

type connectData struct {
	Email string `json:"email"`
}

type contextKey int

var ginContextKey contextKey

// GinContextToContextMiddleware - at the resolver level we only have access
// to context.Context inside centrifuge, but we need the gin context. So we
// create a gin middleware to add its context to the context.Context used by
// centrifuge websocket server.
func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GinContextFromContext - we recover the gin context from the context.Context
// struct where we added it just above
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ginContextKey)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}
	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

// Finally we can use gin context in the auth middleware of centrifuge.
func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// We get gin ctx from context.Context struct.
		//gc, err := GinContextFromContext(ctx)
		//if err != nil {
		//	fmt.Printf("Failed to retrieve gin context")
		//	fmt.Print(err.Error())
		//	return
		//}
		//// And now we can access gin session.
		//s := sessions.Default(gc)
		//username := s.Get("user").(string)
		//if username != "" {
		//	fmt.Printf("Successful websocket auth for user %s\n", username)
		//} else {
		//	fmt.Printf("Failed websocket auth for user %s\n", username)
		//	return
		//}
		newCtx := centrifuge.SetCredentials(ctx, &centrifuge.Credentials{
			UserID: "42",
			Info:   []byte(`{"name": "Alexander"}`),
		})
		r = r.WithContext(newCtx)
		h.ServeHTTP(w, r)
	})
}

func main() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	l := logger.InitLogger()

	cfg := config.InitConfig(l)

	session, err := db.InitDB(cfg)
	if err != nil {
		panic(err)
	}

	err = db.InitKeyspaceAndTables(cfg, session)
	if err != nil {
		panic(err)
	}

	defer session.Close()

	pendingMessagesTopicSyncProducer, err := kafka.NewSyncProducer(
		cfg.KafkaHosts(),
		cfg.KafkaPendingMessagesTopic(),
		kafka.WithLogger(l),
	)

	if err != nil {
		l.Error("error init pending_messages kafka producer", zap.Error(err))
	}

	defer pendingMessagesTopicSyncProducer.Close()

	messageRepo := repository.NewMessageRepository(session, cfg.Logger())

	node, _ := centrifuge.New(centrifuge.Config{
		LogLevel:   centrifuge.LogLevelDebug,
		LogHandler: handleLog,
	})

	node.OnConnecting(func(ctx context.Context, event centrifuge.ConnectEvent) (centrifuge.ConnectReply, error) {
		// Let's include user email into connect reply, so we can display username in chat.
		// This is an optional step, actually.
		cred, ok := centrifuge.GetCredentials(ctx)
		if !ok {
			return centrifuge.ConnectReply{}, centrifuge.DisconnectServerError
		}
		data, _ := json.Marshal(connectData{
			Email: cred.UserID,
		})
		return centrifuge.ConnectReply{
			Data: data,
		}, nil
	})

	node.OnConnect(func(client *centrifuge.Client) {
		transport := client.Transport()
		log.Printf("[user %s] connected via %s with protocol: %s", client.UserID(), transport.Name(), transport.Protocol())

		// Event handler should not block, so start separate goroutine to
		// periodically send messages to client.
		go func() {
			for {
				select {
				case <-client.Context().Done():
					return
				case <-time.After(5 * time.Second):
					err := client.Send([]byte(`{"time": "` + strconv.FormatInt(time.Now().Unix(), 10) + `"}`))
					if err != nil {
						if err == io.EOF {
							return
						}
						log.Printf("error sending message: %s", err)
					}
				}
			}
		}()

		client.OnRefresh(func(e centrifuge.RefreshEvent, cb centrifuge.RefreshCallback) {
			log.Printf("[user %s] connection is going to expire, refreshing", client.UserID())

			cb(centrifuge.RefreshReply{
				ExpireAt: time.Now().Unix() + 60,
			}, nil)
		})

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			log.Printf("[user %s] subscribes on %s", client.UserID(), e.Channel)

			//if !channelSubscribeAllowed(e.Channel) {
			//	cb(centrifuge.SubscribeReply{}, centrifuge.ErrorPermissionDenied)
			//	return
			//}

			cb(centrifuge.SubscribeReply{
				Options: centrifuge.SubscribeOptions{
					EnableRecovery: true,
					EmitPresence:   true,
					EmitJoinLeave:  true,
					PushJoinLeave:  true,
					Data:           []byte(`{"msg": "welcome"}`),
				},
			}, nil)
		})

		client.OnPublish(func(e centrifuge.PublishEvent, cb centrifuge.PublishCallback) {
			log.Printf("[user %s] publishes into channel %s: %s", client.UserID(), e.Channel, string(e.Data))

			// channel authorization
			if !client.IsSubscribed(e.Channel) {
				cb(centrifuge.PublishReply{}, centrifuge.ErrorPermissionDenied)
				return
			}

			// write message to db

			// sent ack to client

			// check active user connection in with Presence method

			// send msg to active clients, sent to queue for non-active clients
			var msg clientMessage
			err := json.Unmarshal(e.Data, &msg)
			if err != nil {
				cb(centrifuge.PublishReply{}, centrifuge.ErrorBadRequest)
				return
			}
			msg.Timestamp = time.Now().Unix()
			data, _ := json.Marshal(msg)

			result, err := node.Publish(
				e.Channel, data,
				centrifuge.WithHistory(300, time.Minute),
				centrifuge.WithClientInfo(e.ClientInfo),
			)

			cb(centrifuge.PublishReply{Result: &result}, err)
		})

		client.OnPresence(func(e centrifuge.PresenceEvent, cb centrifuge.PresenceCallback) {
			log.Printf("[user %s] calls presence on %s", client.UserID(), e.Channel)

			if !client.IsSubscribed(e.Channel) {
				cb(centrifuge.PresenceReply{}, centrifuge.ErrorPermissionDenied)
				return
			}
			cb(centrifuge.PresenceReply{}, nil)
		})

		client.OnUnsubscribe(func(e centrifuge.UnsubscribeEvent) {
			log.Printf("[user %s] unsubscribed from %s: %s", client.UserID(), e.Channel, e.Reason)
		})

		client.OnAlive(func() {
			log.Printf("[user %s] connection is still active", client.UserID())
		})

		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			log.Printf("[user %s] disconnected: %s", client.UserID(), e.Reason)
		})
	})

	// We also start a separate goroutine for centrifuge itself, since we
	// still need to run gin web server.
	go func() {
		if err := node.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	r := gin.Default()
	// Here we tell gin to use the middleware we created just above
	r.Use(GinContextToContextMiddleware())

	r.GET("/connection/websocket", gin.WrapH(authMiddleware(centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}))))

	r.GET("/chat", func(c *gin.Context) {
		s := sessions.Default(c)
		if s.Get("user") != nil {
			c.HTML(200, "chat.html", gin.H{})
		} else {
			c.JSON(403, gin.H{
				"message": "Not logged in!",
			})
		}
		c.Abort()
	})

	r.GET("/test-insert", func(c *gin.Context) {

		conversationID, _ := gocql.RandomUUID()
		messageID, _ := gocql.RandomUUID()

		m := model.Message{
			ConversationID:    conversationID,
			MessageID:         messageID,
			ConversationName:  "test",
			FromUserID:        "user-id",
			FromUserNumber:    "0929024",
			FromUserFirstName: "andi",
			FromUserLastName:  "balo",
			FromUserEmail:     "andialo214@gmail.com",
			TextContent:       "hello",
			SentAt:            time.Now(),
			CreatedBy:         "",
			CreatedAt:         time.Now(),
			UpdatedBy:         "",
			UpdatedAt:         nil,
			DeletedBy:         "",
			DeletedAt:         nil,
		}

		err = messageRepo.SaveMessage(m)

		if err != nil {
			c.String(500, "fail")
		}

		c.String(200, "success")
	})

	_ = r.Run() // listen and serve on 0.0.0.0:8080
}
