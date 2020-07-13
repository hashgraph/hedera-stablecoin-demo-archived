package main

import (
	"bytes"
	"crypto/ed25519"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/data"
	"github.io/hashgraph/stable-coin/domain"
	"github.io/hashgraph/stable-coin/mirror/api"
	"github.io/hashgraph/stable-coin/mirror/operation"
	"github.io/hashgraph/stable-coin/mirror/state"
	"github.io/hashgraph/stable-coin/pb"
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

var mirrorClient hedera.MirrorClient
var listenAttempts = 0

func init() {
	_ = godotenv.Load()

	// Configure the logger to be pretty
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})

	// configure log level for mirror from env
	lvl, err := zerolog.ParseLevel(strings.ToLower(os.Getenv("MIRROR_LOG")))
	if err != nil {
		panic(err)
	}

	zerolog.SetGlobalLevel(lvl)
}

func main() {
	var err error
	mirrorClient, err = hedera.NewMirrorClient(os.Getenv("HEDERA_MIRROR_NODE"))
	if err != nil {
		panic(err)
	}

	// add the admin user if missing
	if _, exists := state.User.Load("Admin"); !exists {

		adminPublicKey, err := hedera.Ed25519PublicKeyFromString(os.Getenv("ADMIN_PUBLIC_KEY"))
		if err != nil {
			panic(err)
		}

		state.AddUser("Admin", adminPublicKey.Bytes())
	}

	// start the mirror listener
	err = startListening()
	if err != nil {
		panic(err)
	}

	// now that the mirror client is running in the background, proceed to run the mirror API
	api.Run()
}

func startListening() error {
	topicID, err := hedera.TopicIDFromString(os.Getenv("TOPIC_ID"))
	if err != nil {
		return err
	}

	startTime, err := data.GetLatestOperationConsensus()

	if err == sql.ErrNoRows {
		startTime = time.Now()
	} else if err != nil {
		return err
	}

	_, err = hedera.NewMirrorConsensusTopicQuery().
		SetTopicID(topicID).
		SetStartTime(startTime.Add(1*time.Nanosecond)).
		Subscribe(mirrorClient, func(response hedera.MirrorConsensusTopicResponse) {
			listenAttempts = 0

			err := handle(response)
			if err != nil {
				log.Error().Err(err).
					Uint64("seq", response.SequenceNumber).
					Msg("failed to process message; skipping")
			}
		}, handleSubscribeFail)

	if err != nil {
		// FIXME(sdk): immediate subscribe fails should probably hit the subscribe fail
		handleSubscribeFail(err)
	}

	return nil
}

func handleSubscribeFail(err error) {
	listenAttempts += 1

	delay := time.Duration(listenAttempts) * time.Millisecond * 250

	log.Error().Err(err).
		Msgf("mirror subscribe failed; reconnecting in %s...", delay)

	time.Sleep(delay)

	err = startListening()
	if err != nil {
		panic(err)
	}
}

func handle(response hedera.MirrorConsensusTopicResponse) error {
	// statistics -----------------

	if handleCount == 0 {
		startTime = time.Now().Round(time.Millisecond).UnixNano() / timeDivisor
		blockStartTime = startTime
	}

	handleCount = handleCount + 1
	if handleCount%modulo == 0 {
		runTime := time.Now().Round(time.Millisecond).UnixNano()/timeDivisor - startTime
		blockTime := time.Now().Round(time.Millisecond).UnixNano()/timeDivisor - blockStartTime
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
			float64(runTime)/float64(timeMultiplier),
			handleCount,
			avgTPS,
			blockTPS,
			float64(blockTime)/float64(timeMultiplier))
	}

	// ----------------------------

	// parse the primitive operation wrapper from the response
	var primitive pb.Primitive
	err := proto.Unmarshal(response.Message, &primitive)
	if err != nil {
		return err
	}

	primitiveHederaPublicKey, err := hedera.Ed25519PublicKeyFromString(primitive.Header.PublicKey)
	if err != nil {
		return err
	}

	primitivePublicKey := ed25519.PublicKey(primitiveHederaPublicKey.Bytes())

	var op domain.Operation
	switch primitive.Primitive.(type) {
	case *pb.Primitive_Join:
		//if !bytes.Equal(primitivePublicKey, adminPublicKey) {
		//	// not an administrator; ignore
		//	return nil
		//}

		v := primitive.GetJoin()

		err = verify(primitive.Header, v, primitivePublicKey)
		if err != nil {
			return err
		}

		op, err = operation.Announce(v)

	case *pb.Primitive_MintTo:
		//if !bytes.Equal(primitivePublicKey, adminPublicKey) {
		//	// not an administrator; ignore
		//	return nil
		//}

		v := primitive.GetMintTo()

		err = verify(primitive.Header, v, primitivePublicKey)
		if err != nil {
			return err
		}

		op, err = operation.Mint(v)

	case *pb.Primitive_Transfer:
		v := primitive.GetTransfer()

		err = verify(primitive.Header, v, primitivePublicKey)
		if err != nil {
			return err
		}

		op, err = operation.Transfer(primitivePublicKey, v)

	case *pb.Primitive_Burn:
		v := primitive.GetBurn()

		err = verify(primitive.Header, v, primitivePublicKey)
		if err != nil {
			return err
		}

		op, err = operation.Burn(primitivePublicKey, v)

	case *pb.Primitive_Freeze:
		v := primitive.GetFreeze()
		// get admin key
		var adminPubKeyI interface{}
		var exists bool
		if adminPubKeyI, exists = state.User.Load("Admin"); !exists {
			op = domain.Operation{
				Operation:     domain.OpFreeze,
				Status:        domain.OpStatusFailed,
				StatusMessage: "username Admin does not exist",
			}
		} else {
			// admin user exists, check public keys match
			adminPubKey := []byte(adminPubKeyI.(ed25519.PublicKey))

			if bytes.Compare(adminPubKey, primitiveHederaPublicKey.Bytes()) != 0 {
				primitivePublicKeyHex := hex.EncodeToString(primitiveHederaPublicKey.Bytes())
				adminPubKeyHex := hex.EncodeToString(adminPubKey)
				fmt.Println(adminPubKeyHex)
				fmt.Println(primitivePublicKeyHex)
				op = domain.Operation{
					Operation:     domain.OpFreeze,
					Status:        domain.OpStatusFailed,
					StatusMessage: fmt.Sprintf("invalid admin key `%s`", primitivePublicKeyHex),
				}
			} else {
				err = verify(primitive.Header, v, primitivePublicKey)
				if err != nil {
					return err
				}

				op, err = operation.Freeze(primitivePublicKey, v.Account, true)
			}
		}

	case *pb.Primitive_Unfreeze:
		v := primitive.GetUnfreeze()
		// get admin key
		var adminPubKeyI interface{}
		var exists bool
		if adminPubKeyI, exists = state.User.Load("Admin"); !exists {
			op = domain.Operation{
				Operation:     domain.OpUnFreeze,
				Status:        domain.OpStatusFailed,
				StatusMessage: "username Admin does not exist",
			}
		} else {
			// admin user exists, check public keys match
			adminPubKey := []byte(adminPubKeyI.(ed25519.PublicKey))

			if bytes.Compare(adminPubKey, primitiveHederaPublicKey.Bytes()) != 0 {
				adminPubKeyHex := hex.EncodeToString(primitiveHederaPublicKey.Bytes())
				op = domain.Operation{
					Operation:     domain.OpUnFreeze,
					Status:        domain.OpStatusFailed,
					StatusMessage: fmt.Sprintf("invalid admin key `%s`", adminPubKeyHex),
				}
			} else {
				err = verify(primitive.Header, v, primitivePublicKey)
				if err != nil {
					return err
				}

				op, err = operation.Freeze(primitivePublicKey, v.Account, false)
			}
		}

	case *pb.Primitive_Clawback:
		v := primitive.GetClawback()
		// get admin key
		var adminPubKeyI interface{}
		var exists bool
		if adminPubKeyI, exists = state.User.Load("Admin"); !exists {
			op = domain.Operation{
				Operation:     domain.OpClawback,
				Status:        domain.OpStatusFailed,
				StatusMessage: "username Admin does not exist",
			}
		} else {
			// admin user exists, check public keys match
			adminPubKey := []byte(adminPubKeyI.(ed25519.PublicKey))

			if bytes.Compare(adminPubKey, primitiveHederaPublicKey.Bytes()) != 0 {
				adminPubKeyHex := hex.EncodeToString(primitiveHederaPublicKey.Bytes())
				op = domain.Operation{
					Operation:     domain.OpClawback,
					Status:        domain.OpStatusFailed,
					StatusMessage: fmt.Sprintf("invalid admin key `%s`", adminPubKeyHex),
				}
			} else {
				err = verify(primitive.Header, v, primitivePublicKey)
				if err != nil {
					return err
				}

				op, err = operation.Clawback(v)
			}
		}

	case *pb.Primitive_AdminKeyUpdate:
		v := primitive.GetAdminKeyUpdate()
		// get admin key
		var adminPubKeyI interface{}
		var exists bool
		if adminPubKeyI, exists = state.User.Load("Admin"); !exists {
			op = domain.Operation{
				Operation:     domain.OpAdminKeyUpdate,
				Status:        domain.OpStatusFailed,
				StatusMessage: "username Admin does not exist",
			}
		} else {
			// admin user exists, check public keys match
			adminPubKey := []byte(adminPubKeyI.(ed25519.PublicKey))

			if bytes.Compare(adminPubKey, primitiveHederaPublicKey.Bytes()) != 0 {
				adminPubKeyHex := hex.EncodeToString(primitiveHederaPublicKey.Bytes())
				op = domain.Operation{
					Operation:     domain.OpAdminKeyUpdate,
					Status:        domain.OpStatusFailed,
					StatusMessage: fmt.Sprintf("invalid admin key `%s`", adminPubKeyHex),
				}
			} else {
				err = verify(primitive.Header, v, primitivePublicKey)
				if err != nil {
					return err
				}

				op, err = operation.AdminKeyUpdate(primitivePublicKey, v.NewPublicKey)
			}
		}

	default:
		err = fmt.Errorf("unimplemented operation: %T", primitive.Primitive)
	}

	if err != nil {
		return err
	}

	op.Signature = primitive.Header.Signature
	op.Consensus = response.ConsensusTimeStamp.UnixNano()

	state.AddOperation(op)

	return nil
}

func verify(header *pb.PrimitiveHeader, v proto.Message, publicKey ed25519.PublicKey) error {
	message, err := proto.Marshal(v)
	if err != nil {
		return err
	}

	nonce := header.Random
	nonceBytes := []byte(strconv.FormatUint(nonce, 10))[:]

	message = append(message, nonceBytes...)

	verified := ed25519.Verify(publicKey, message, header.Signature)

	if !verified {
		return errors.New("invalid signature")
	}

	return nil
}
