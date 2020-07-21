package com.hedera.stabl.test;

import com.google.protobuf.ByteString;
import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PrivateKey;
import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PublicKey;
import com.hedera.mirror.api.proto.ConsensusServiceGrpc;
import com.hedera.mirror.api.proto.ConsensusTopicResponse;
import com.hedera.mirror.api.proto.ConsensusTopicQuery;
import com.hedera.mirror.api.proto.Timestamp;
import io.grpc.*;
import io.grpc.stub.ServerCallStreamObserver;
import stabl.Join;
import stabl.Primitive;

import java.io.IOException;
import java.time.Instant;
import java.util.Random;
import java.util.concurrent.TimeUnit;

import static com.hedera.stabl.test.Primitives.primitiveHeader;

public class MirrorSimulator {
    private Server server;

    private void start(int threadCount) throws IOException {
        /* The port on which the server should run */
        int port = 5600;
        server = ServerBuilder.forPort(port)
            .addService(new MirrorConsensusTopicService(threadCount))
            .build()
            .start();
        System.out.println("Server started, listening on " + port);
        Runtime.getRuntime().addShutdownHook(new Thread() {
            @Override
            public void run() {
                // Use stderr here since the logger may have been reset by its JVM shutdown hook.
                System.err.println("*** shutting down gRPC server since JVM is shutting down");
                try {
                    MirrorSimulator.this.stop();
                } catch (InterruptedException e) {
                    e.printStackTrace(System.err);
                }
                System.err.println("*** server shut down");
            }
        });
    }

    private void stop() throws InterruptedException {
        if (server != null) {
            server.shutdown().awaitTermination(30, TimeUnit.SECONDS);
        }
    }

    /**
     * Await termination on the main thread since the grpc library uses daemon threads.
     */
    private void blockUntilShutdown() throws InterruptedException {
        if (server != null) {
            server.awaitTermination();
        }
    }
    /**
     * Main method.
     */
    public static void run(int threadCount) throws Exception {
        final MirrorSimulator server = new MirrorSimulator();
        server.start(threadCount);
        server.blockUntilShutdown();
    }

    private static class MirrorConsensusTopicService extends ConsensusServiceGrpc.ConsensusServiceImplBase {

        private int threadCount;

        public MirrorConsensusTopicService(int threadCount) {
            this.threadCount = threadCount;
        }

        @Override
        public void subscribeTopic(ConsensusTopicQuery request,
                              io.grpc.stub.StreamObserver<ConsensusTopicResponse> responseObserver) {
            for (int i=0; i < threadCount; i++) {
                Thread t = new Thread(new SimulatorThread(responseObserver));
                t.start();
            }

            boolean keepGoing = true;
            while (keepGoing) {
                ServerCallStreamObserver observer = (ServerCallStreamObserver)responseObserver;
                if (observer.isCancelled()) {
                    System.out.println("client disconnected");
                    keepGoing = false;
                } else {
                    try {
                        Thread.sleep(1000);
                    } catch (InterruptedException e) {
                        e.printStackTrace();
                    }
                }
            }
        }
    }


}
