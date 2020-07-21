package com.hedera.stabl.test;

import com.google.common.primitives.Bytes;
import com.google.protobuf.ByteString;
import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PrivateKey;
import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PublicKey;
import stabl.Burn;
import stabl.Primitive;
import stabl.PrimitiveHeader;
import stabl.Transfer;
import java.util.Random;
import java.io.UnsupportedEncodingException;
import java.util.Base64;

public class Primitives {
    private static final Random random = new Random();

    public static String burnPrimitive(Ed25519PrivateKey privateKey, Ed25519PublicKey publicKey) throws UnsupportedEncodingException {

        long amount = random.nextInt(1000) + 100;
        Burn.Builder burn = Burn.newBuilder()
                .setAmount(amount);

        Primitive.Builder primitive = Primitive.newBuilder()
                .setHeader(primitiveHeader(privateKey, burn.build().toByteArray(), publicKey.toString()))
                .setBurn(burn);

        return Base64.getEncoder().encodeToString(primitive.build().toByteArray());
    }

    public static String sendPrimitive(String toAddress, Ed25519PrivateKey privateKey, Ed25519PublicKey publicKey, int quantity) throws UnsupportedEncodingException {

        Transfer.Builder transfer = Transfer.newBuilder()
                .setToAddress(toAddress)
                .setQuantity(quantity);

        Primitive.Builder primitive = Primitive.newBuilder()
                .setHeader(primitiveHeader(privateKey, transfer.build().toByteArray(), publicKey.toString()))
                .setTransfer(transfer);

        return Base64.getEncoder().encodeToString(primitive.build().toByteArray());
    }

    public static PrimitiveHeader primitiveHeader(Ed25519PrivateKey privateKey, byte[] primitive, String publicKey) throws UnsupportedEncodingException {
        long rand = random.nextInt() + 1;
        if (rand < 0) { rand = rand * -1;}
        String randomString = Long.toString(rand);
        if (randomString.length() > 19) {
            System.out.println("too long");
            System.exit(0);
        }

        byte[] signMe = Bytes.concat(primitive, randomString.getBytes("UTF8"));
        byte[] signature = privateKey.sign(signMe);

        PrimitiveHeader.Builder header = PrimitiveHeader.newBuilder()
                .setRandom(rand)
                .setSignature(ByteString.copyFrom(signature))
                .setPublicKey(publicKey);

        return header.build();
    }
}
