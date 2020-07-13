package routes

import (
	"crypto"
	"crypto/ed25519"
	crand "crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	mrand "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.io/hashgraph/stable-coin/pb"
)

var hederaClient *hedera.Client
var hederaTopicID hedera.ConsensusTopicID
var adminPrivateKey ed25519.PrivateKey

type transactionResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

func init() {
	_ = godotenv.Load()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, NoColor: false})

	// initialize the hedera client

	network, err := networkFromFile(os.Getenv("HEDERA_NETWORK"))
	if err != nil {
		panic(err)
	}

	hederaClient = hedera.NewClient(network)

	// configure the operator, or payer of transaction fees

	operatorId, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
	if err != nil {
		panic(err)
	}

	operatorKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("OPERATOR_KEY"))
	if err != nil {
		panic(err)
	}

	hederaClient.SetOperator(operatorId, operatorKey)

	// parse the admin key for token operations on the network

	adminHederaPrivateKey, err := hedera.Ed25519PrivateKeyFromString(os.Getenv("ISSUER_PRIVATE_KEY"))
	if err != nil {
		panic(err)
	}

	adminPrivateKey = adminHederaPrivateKey.Bytes()

	// ensure that we have a topic ID

	topicIDString := os.Getenv("TOPIC_ID")

	if len(topicIDString) == 0 {
		// no topic ID specified, create a new topic
		id, err := hedera.NewConsensusTopicCreateTransaction().Execute(hederaClient)
		if err != nil {
			panic(err)
		}

		receipt, err := id.GetReceipt(hederaClient)
		if err != nil {
			panic(err)
		}

		hederaTopicID = receipt.GetConsensusTopicID()

		log.Info().Msgf("no $TOPIC_ID specified, created a new topic %s, set this as TOPIC_ID in the environment and re-run the application", hederaTopicID)

		// the application requires TOPIC_ID to be specified as an environment variable
		os.Exit(0)
	} else {
		hederaTopicID, err = hedera.TopicIDFromString(topicIDString)
		if err != nil {
			panic(err)
		}
	}
}

func SendRawTransaction(c echo.Context) error {
	var req struct {
		Primitive string `json:"primitive"`
	}

	err := c.Bind(&req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, transactionResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	primitive, err := base64.StdEncoding.DecodeString(req.Primitive)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, transactionResponse{
			Status:  false,
			Message: err.Error(),
		})
	}

	err = sendRaw(primitive)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, transactionResponse{
			Status:  false,
			Message: err.Error(),
		})
	} else {
		return c.JSON(http.StatusAccepted, transactionResponse{
			Status:  true,
			Message: "raw transaction request sent",
		})
	}
}

func sendTransaction(v proto.Message, p *pb.Primitive) error {
	var err error
	p.Header, err = makePrimitiveHeader(v)
	if err != nil {
		return err
	}

	messageBytes, err := proto.Marshal(p)
	if err != nil {
		return err
	}

	err = sendRaw(messageBytes)
	if err != nil {
		return err
	}
	return nil
}

func sendRaw(raw []byte) error {
	for {
		_, err := hedera.NewConsensusMessageSubmitTransaction().
			SetMessage(raw).
			SetTopicID(hederaTopicID).
			Execute(hederaClient)

		if err != nil {
			if strings.Contains(err.Error(), "server closed the stream without sending trailers") {
				// resubmit
				log.Warn().Msg("server closed stream - resubmitting")
			} else if strings.Contains(err.Error(), "DUPLICATE_TRANSACTION") {
				// resubmit
				log.Warn().Msg("duplicate transaction id - resubmitting")
			} else {
				log.Error().Err(err)
				return err
			}
		} else {
			return err
		}
	}
}

func makePrimitiveHeader(v proto.Message) (*pb.PrimitiveHeader, error) {
	message, err := proto.Marshal(v)
	if err != nil {
		return nil, err
	}

	var nonce uint64 = mrand.Uint64()
	message = append(message, []byte(strconv.FormatUint(nonce, 10))[:]...)

	signature, err := adminPrivateKey.Sign(crand.Reader, message, crypto.Hash(0))
	if err != nil {
		return nil, err
	}

	header := pb.PrimitiveHeader{
		PublicKey: hex.EncodeToString(adminPrivateKey.Public().(ed25519.PublicKey)),
		Random:    uint64(nonce),
		Signature: signature,
	}

	return &header, nil
}

// FIXME: Contribute this back to hedera-sdk-go
func networkFromFile(filename string) (map[string]hedera.AccountID, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var v struct {
		Network map[string]string `json:"network"`
	}

	err = json.Unmarshal(fileData, &v)
	if err != nil {
		return nil, err
	}

	network := map[string]hedera.AccountID{}

	for address, nodeId := range v.Network {
		nodeIdValue, err := hedera.AccountIDFromString(nodeId)
		if err != nil {
			return nil, err
		}

		network[address] = nodeIdValue
	}

	return network, nil
}
