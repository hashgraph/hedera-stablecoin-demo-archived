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
            FileWriter mixFileWriter = new FileWriter("stabl-test/mix_" + fileSuffix + ".csv");
            PrintWriter mixPrintWriter = new PrintWriter(mixFileWriter);

            String[] users = new String[iterations];
            Ed25519PrivateKey[] privateKeys = new Ed25519PrivateKey[iterations];

            System.out.println("Generating user names and keys");
            for (int i = 0; i < iterations; i++) {
                users[i] = fileSuffix + "_user_" + i;
                privateKeys[i] = Ed25519PrivateKey.generate();
            }

            for (int i = 0; i < iterations; i++) {
                if ((i % modulo) == 0) {
                    System.out.println(fileSuffix + "-" + i + ".");
                }
                joinPrintWriter.println(privateKeys[i].publicKey.toString() + "," + users[i]);
                buyPrintWriter.println(users[i] + "," + (10_000 + random.nextInt(10_000)));
                burnPrintWriter.println(Primitives.burnPrimitive(privateKeys[i], privateKeys[i].publicKey));

                String toAddress = users[random.nextInt(iterations)];
                int quantity = random.nextInt(100) + 1;
                sendPrintWriter.println(Primitives.sendPrimitive(toAddress, privateKeys[i], privateKeys[i].publicKey, quantity));

                //operation, key,user,amount,primitive
                String lineFormat = "%s,%s,%s,%d,%s";
                String join = String.format(lineFormat, "j", privateKeys[i].publicKey.toString(), users[i], 0, "");
                mixPrintWriter.println(join);
                String mint = String.format(lineFormat, "m", "", users[i], (10_000 + random.nextInt(10_000)), "");
                mixPrintWriter.println(mint);

                for (int transferIndex = 0; transferIndex < 9; transferIndex++) {
                    toAddress = users[random.nextInt(iterations)];
                    String sendPrimitive = Primitives.sendPrimitive(toAddress, privateKeys[i], privateKeys[i].publicKey, random.nextInt(100) + 1);
                    String transfer = String.format(lineFormat, "t", "", "", 0, sendPrimitive);
                    mixPrintWriter.println(transfer);
                }
            }
            joinPrintWriter.close();
            buyPrintWriter.close();
            burnPrintWriter.close();
            sendPrintWriter.close();
            mixPrintWriter.close();
            System.out.println("");
        } catch (UnsupportedEncodingException e) {
            System.out.println(e.getMessage());
        } catch (IOException e) {
            e.printStackTrace();
        }
    }
}
