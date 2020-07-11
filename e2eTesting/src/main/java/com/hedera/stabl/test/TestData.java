package com.hedera.stabl.test;

import java.io.IOException;

public class TestData {
    private static String[] fileSuffixes = {"a","b","c","d","e","f","g","h","i","j"};
    private static int iterations = 1_000_000;

    public static void main(String[] args) throws IOException, InterruptedException {
        if (args.length >= 1) {
            if (args[0].toUpperCase().equals("E2E")) {
                // e2e localhost:3128
                End2End.e2eTest(args[1], args[2]); // this will loop indefinitely
            } else {
                iterations = Integer.parseInt(args[0]);
            }
        }
        for (String fileSuffix : fileSuffixes) {
            FileGen fileGen = new FileGen();
            fileGen.iterations = iterations;
            fileGen.fileSuffix = fileSuffix;

            Thread thread = new Thread(fileGen);
            thread.start();
        }
    }
}
