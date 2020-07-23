package com.hedera.stabl.test;

import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PrivateKey;
import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PublicKey;

import java.io.FileWriter;
import java.io.IOException;
import java.io.PrintWriter;
import java.io.UnsupportedEncodingException;
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
//            FileWriter mixJoinFileWriter = new FileWriter("stabl-test/mix-join_" + fileSuffix + ".csv");
//            PrintWriter mixJoinPrintWriter = new PrintWriter(mixJoinFileWriter);
//            FileWriter mixMintFileWriter = new FileWriter("stabl-test/mix-mint_" + fileSuffix + ".csv");
//            PrintWriter mixMintPrintWriter = new PrintWriter(mixMintFileWriter);
//            FileWriter mixTransferFileWriter = new FileWriter("stabl-test/mix-transfer_" + fileSuffix + ".csv");
//            PrintWriter mixTransferPrintWriter = new PrintWriter(mixTransferFileWriter);

            String[] users = new String[iterations];
            Ed25519PrivateKey[] privateKeys = new Ed25519PrivateKey[iterations];

            for (int i = 0; i < iterations; i++) {
                if ((i % modulo) == 0) {
                    System.out.println(fileSuffix + "-" + i + ".");
                }
                users[i] = fileSuffix + "_user_" + i;
                privateKeys[i] = Ed25519PrivateKey.generate();

                joinPrintWriter.println(privateKeys[i].publicKey.toString() + "," + users[i]);
                buyPrintWriter.println(users[i] + "," + (100_000 + random.nextInt(100_000)));
                burnPrintWriter.println(Primitives.burnPrimitive(privateKeys[i], privateKeys[i].publicKey));

                if (i >= 10) {
                    for (int xferCount = 0; xferCount < 100; xferCount++) {
                        // randomly pick to/from
                        int midpoint = i / 2;
                        String toAddress = users[random.nextInt(midpoint)]; // target any of the previously created users
                        int toIndex = random.nextInt(midpoint) + midpoint;
                        int quantity = random.nextInt(10) + 1;
                        sendPrintWriter.println(Primitives.sendPrimitive(toAddress, privateKeys[toIndex], privateKeys[toIndex].publicKey, quantity));
                    }
                }
//                mixPrintWriter.println(join);

//                if (i >= 5) { // start minting five users ago
//                    String mint = String.format(lineFormat, "m", "", users[i-5], (10_000 + random.nextInt(10_000)), "");
//                    mixPrintWriter.println(mint);
//                }
//                if (i >= 10) { // no point transferring to self, need some users to transfer to
//                    for (int transferIndex = 0; transferIndex < 9; transferIndex++) {
//                        toAddress = users[random.nextInt(i-9)]; // randomly pick a user from the already created users
//                        Ed25519PrivateKey fromPrivateKey = privateKeys[i-10]; // transfer from ten users ago
//                        String sendPrimitive = Primitives.sendPrimitive(toAddress, fromPrivateKey, fromPrivateKey.publicKey, random.nextInt(100) + 1);
//                        String transfer = String.format(lineFormat, "t", "", "", 0, sendPrimitive);
//                        mixPrintWriter.println(transfer);
//                    }
//                }
            }
            joinPrintWriter.close();
            buyPrintWriter.close();
            burnPrintWriter.close();
            sendPrintWriter.close();
//            mixPrintWriter.close();
            System.out.println("");
        } catch (UnsupportedEncodingException e) {
            System.out.println(e.getMessage());
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
