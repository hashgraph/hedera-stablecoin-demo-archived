package com.hedera.stabl.test;

import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PrivateKey;
import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PublicKey;

import java.io.FileWriter;
import java.io.PrintWriter;
import java.util.Random;

public class FileGen implements Runnable {
    public int iterations = 0;
    public String fileSuffix = "";
    private int modulo;
    private static final Random random = new Random();

    public void run() {
        modulo = iterations / 10;
        System.out.println("Generating '" + fileSuffix + "': ");
        try {
            FileWriter joinFileWriter = new FileWriter("stabl-test/join_" + fileSuffix + ".csv");
            PrintWriter joinPrintWriter = new PrintWriter(joinFileWriter);
            FileWriter buyFileWriter = new FileWriter("stabl-test/buy_" + fileSuffix + ".csv");
            PrintWriter buyPrintWriter = new PrintWriter(buyFileWriter);
            FileWriter burnFileWriter = new FileWriter("stabl-test/burn_" + fileSuffix + ".csv");
            PrintWriter burnPrintWriter = new PrintWriter(burnFileWriter);
            FileWriter sendFileWriter = new FileWriter("stabl-test/send_" + fileSuffix + ".csv");
            PrintWriter sendPrintWriter = new PrintWriter(sendFileWriter);

            for (long i = 0; i < iterations; i++) {
                if ((i % modulo) == 0) {
                    System.out.println(fileSuffix + "-" + i + ".");
                }
                Ed25519PrivateKey privateKey = Ed25519PrivateKey.generate();
                Ed25519PublicKey publicKey = privateKey.publicKey;
                joinPrintWriter.println(publicKey.toString() + "," + fileSuffix + "_user_" + i);
                buyPrintWriter.println(fileSuffix + "_user_" + i + "," + (10_000 + random.nextInt(10_000)));
                burnPrintWriter.println(Primitives.burnPrimitive(privateKey, publicKey));
                String toAddress = fileSuffix + "_user_" + random.nextInt(iterations);
                int quantity = random.nextInt(100) + 1;
                sendPrintWriter.println(Primitives.sendPrimitive(toAddress, privateKey, publicKey, quantity));
            }
            joinPrintWriter.close();
            buyPrintWriter.close();
            burnPrintWriter.close();
            sendPrintWriter.close();
            System.out.println();
        } catch (Exception e) {
            System.out.println(e.getMessage());
        }
    }
}
