package com.hedera.stabl.test;

import com.google.protobuf.ByteString;
import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PrivateKey;
import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PublicKey;
import com.hedera.mirror.api.proto.ConsensusTopicResponse;
import com.hedera.mirror.api.proto.Timestamp;
import stabl.Join;
import stabl.Primitive;

import java.time.Instant;
import java.util.Random;

import static com.hedera.stabl.test.Primitives.primitiveHeader;

public class SimulatorThread implements Runnable {
    public io.grpc.stub.StreamObserver<ConsensusTopicResponse> responseObserver = null;

    public SimulatorThread(io.grpc.stub.StreamObserver<ConsensusTopicResponse> responseObserver) {
        this.responseObserver = responseObserver;
    }

    public void run() {
        Random random = new Random();
        int nanos = random.nextInt(Integer.MAX_VALUE);
        System.out.println("Starting Thread " + nanos);
        boolean keepGoing = true;
        while (keepGoing) {
            Ed25519PrivateKey privateKey = Ed25519PrivateKey.generate();
            Ed25519PublicKey publicKey = privateKey.publicKey;

            int randomValue1 = random.nextInt(Integer.MAX_VALUE);
            long time = Instant.now().getEpochSecond();
            Join join = Join.newBuilder()
                .setAddress(publicKey.toString())
                .setUsername("user_" + Integer.toString(nanos) + Integer.toString(randomValue1))
                .build();

            try {
                nanos += 1;
                Primitive.Builder primitive = Primitive.newBuilder()
                    .setHeader(primitiveHeader(privateKey, join.toByteArray(), publicKey.toString()))
                    .setJoin(join);

                Timestamp timestamp = Timestamp.newBuilder().setSeconds(Instant.now().getEpochSecond()).setNanos(nanos).build();
                ConsensusTopicResponse consensusTopicResponse = ConsensusTopicResponse.newBuilder()
                    .setMessage(ByteString.copyFrom(primitive.build().toByteArray()))
                    .setConsensusTimestamp(timestamp)
                    .setRunningHash(ByteString.copyFromUtf8("runningHash"))
                    .setSequenceNumber(1)
                    .build();

                this.responseObserver.onNext(consensusTopicResponse);

                Thread.sleep(1);

            } catch (Exception e) {
                keepGoing = false;
            }

        }
    }
}
