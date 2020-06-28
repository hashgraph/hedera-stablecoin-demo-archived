package main

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/mirror/api"
	"github.io/hashgraph/stable-coin/mirror/operation"
	"github.io/hashgraph/stable-coin/pb"
	"os"
	"strconv"
	"time"
)

// statistics -----------------
// TPS calculations
var handleCount int64 = 0
var modulo int64 = 2000
var startTime int64 = 0
var blockStartTime int64 = 0
var blockTPS float64 = 0
var avgTPS float64 = 0
var timeMultiplier int64 = 100
var timeDivisor int64 = 1e7
// ----------------------------

// TODO: Set from .env
var adminPublicKey ed25519.PublicKey

func init() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})

	// Uncomment for a lot more logging
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

func main() {
	mirrorClient, err := hedera.NewMirrorClient(os.Getenv("HEDERA_MIRROR_NODE"))
	if err != nil {
		panic(err)
	}

	topicID, err := hedera.TopicIDFromString(os.Getenv("TOPIC_ID"))
	if err != nil {
		panic(err)
	}

	_, err = hedera.NewMirrorConsensusTopicQuery().
		SetTopicID(topicID).
		SetStartTime(time.Unix(0, 0)).
		Subscribe(mirrorClient, func(response hedera.MirrorConsensusTopicResponse) {
			err := handle(response)
			if err != nil {
				panic(err)
			}
		}, func(err error) {
			panic(err)
		})

	if err != nil {
		panic(err)
	}

	// now that the mirror client is running in the background, proceed to run the mirror API
	api.Run()
}

func handle(response hedera.MirrorConsensusTopicResponse) error {
	// statistics -----------------

	if handleCount == 0 {
		startTime = time.Now().Round(time.Millisecond).UnixNano() / timeDivisor
		blockStartTime = startTime
	}

	handleCount = handleCount + 1
	if handleCount % modulo == 0 {
		runTime := time.Now().Round(time.Millisecond).UnixNano() / timeDivisor - startTime
		blockTime := time.Now().Round(time.Millisecond).UnixNano() / timeDivisor - blockStartTime
		blockStartTime = time.Now().Round(time.Millisecond).UnixNano() / timeDivisor
		blockTPS = 0
		avgTPS = 0
		if blockTime != 0 {
			blockTPS = float64(timeMultiplier) * float64(modulo) / float64(blockTime)
		}
		if runTime != 0 {
			avgTPS = float64(timeMultiplier) * float64(handleCount) / float64(runTime)
		}
		fmt.Printf("%.1f,%d,%.0f,%.0f,%.1f\n",
			float64(runTime) / float64(timeMultiplier),
			handleCount,
			avgTPS,
			blockTPS,
			float64(blockTime) / float64(timeMultiplier))
	}

	// ----------------------------

	// parse the primitive operation wrapper from the response
	var primitive pb.Primitive
	err := proto.Unmarshal(response.Message, &primitive)
	if err != nil {
		return err
	}

	switch primitive.Primitive.(type) {
	case *pb.Primitive_Join:
		v := primitive.GetJoin()

		//err = verify(primitive.Header, v)
		//if err != nil {
		//	return err
		//}

		err = operation.Announce(v)

	default:
		err = fmt.Errorf("unimplemented operation: %T", primitive.Primitive)
	}

	if err != nil {
		return err
	}

	return nil
}

func verify(header *pb.PrimitiveHeader, v proto.Message) error {
	message, err := proto.Marshal(v)
	if err != nil {
		return err
	}

	nonce := header.Random
	nonceBytes := []byte(strconv.FormatUint(nonce, 10))[:]

	message = append(message, nonceBytes...)

	verified := ed25519.Verify(adminPublicKey, message, header.Signature)

	if !verified {
		return errors.New("invalid signature")
	}

	return nil
}
