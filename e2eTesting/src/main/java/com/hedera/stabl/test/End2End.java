package com.hedera.stabl.test;

import com.hedera.stabl.test.model.*;
import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PrivateKey;
import com.hedera.hashgraph.sdk.crypto.ed25519.Ed25519PublicKey;

import javax.ws.rs.client.Client;
import javax.ws.rs.client.ClientBuilder;
import javax.ws.rs.client.Entity;
import javax.ws.rs.core.MediaType;
import javax.ws.rs.core.Response;
import java.io.File;
import java.io.FileWriter;
import java.io.IOException;
import java.io.UnsupportedEncodingException;
import java.time.Instant;
import java.time.ZoneOffset;
import java.time.format.DateTimeFormatter;
import java.time.format.FormatStyle;
import java.util.ArrayList;
import java.util.List;
import java.util.Locale;
import java.util.Random;

public class End2End {

    static class Stats {
        private float total = 0;
        private float min = -1;
        private float max = -1;

        public float range_0_10 = 0;
        public float range_0_10_total = 0;
        public float range_11_15 = 0;
        public float range_11_15_total = 0;
        public float range_16_20 = 0;
        public float range_16_20_total = 0;
        public float range_20_plus = 0;
        public float range_20_total = 0;
        public float apiErrors = 0;
        public float lastValue = -1;

        private String logFile = "";
        private String operation = "";

        public Stats(String logFile, String operation) {
            this.logFile = logFile;
            this.operation = operation;
        }
        public void setValues (float value) {
            lastValue = value;
            if (value == 100) {
                apiErrors += 1;
            } else {
                total = total + value;
                if ((min == -1) || (min > value)) {
                    min = value;
                }
                if ((max == -1) || (max < value)) {
                    max = value;
                }

                if (value <= 10) {
                    range_0_10 += 1;
                    range_0_10_total += value;
                } else if ((value > 10) && (value <= 15)) {
                    range_11_15 += 1;
                    range_11_15_total += value;
                } else if ((value > 15) && (value <= 20)) {
                    range_16_20 += 1;
                    range_16_20_total += value;
                } else {
                    range_20_plus += 1;
                    range_20_total += value;
                }
            }
        }

        public void printOut(FileWriter fr) throws IOException {
            String output = "";
            output = String.format("%s Min=%.1f, Max=%.1f, 0-10=%.1f, 11-15=%.1f, 16-20=%.1f, 20+=%.1f"
                    , operation, min, max, range_0_10, range_11_15, range_16_20, range_20_plus);
            System.out.println(output);

            float range0_10_average = range_0_10_total;
            if (range_0_10 != 0) {
                range0_10_average = range_0_10_total / range_0_10;
            }
            float range_11_15_average = range_11_15_total;
            if (range_11_15 != 0) {
                range_11_15_average = range_11_15_total / range_11_15;
            }
            float range_16_20_average = range_16_20_total;
            if (range_16_20 != 0) {
                range_16_20_average = range_16_20_total / range_16_20;
            }
            float range_20_plus_average = range_20_plus;
            if (range_20_plus != 0) {
                range_20_plus_average = range_20_total / range_20_plus;
            }

            output = String.format("%-9s, %-3.2f, %-3.2f, %-10.0f, %-12.2f, %-11.0f, %-13.2f, %-11.0f, %-13.2f, %-9.0f, %-11.2f\n"
                    , operation
                    , min
                    , max
                    , range_0_10
                    , range0_10_average
                    , range_11_15
                    , range_11_15_average
                    , range_16_20
                    , range_16_20_average
                    , range_20_plus
                    , range_20_plus_average);
            fr.write(output + "\n");
        }

    }
    private static float timeout = 1000 * 20; // 20s
    private static int queryInterval = 100; // in milliseconds
    private static Client client = ClientBuilder.newClient();
    private static final Random random = new Random();
    private static List<String> usernames = new ArrayList<String>();
    private static List<Ed25519PrivateKey> privateKeys = new ArrayList<Ed25519PrivateKey>();
    private static List<Ed25519PublicKey> publicKeys = new ArrayList<Ed25519PublicKey>();
    private static long cycle = 0;
    private static long initialBalance = 10000;
    private static String sendHost = "";
    private static String queryHost = "";
    private static Stats joinStats;
    private static Stats mintStats;
    private static Stats sendStats;

    private static DateTimeFormatter DATE_TIME_FORMATTER = DateTimeFormatter.ofLocalizedDateTime(FormatStyle.SHORT).withLocale(Locale.US).withZone(ZoneOffset.UTC);

    public static void e2eTest(String sendHost, String queryHost) throws IOException, InterruptedException {
        String testId = Integer.toString(random.nextInt(10000));
        boolean delayTest = false;
        if (End2End.sendHost.startsWith("http://")) {
            End2End.sendHost = sendHost;
        } else {
            End2End.sendHost = "http://" + sendHost;
        }
        if (End2End.queryHost.startsWith("http://")) {
            End2End.queryHost = queryHost;
        } else {
            End2End.queryHost = "http://" + queryHost;
        }
        long filePrefix = Instant.now().getEpochSecond();
        String fileName = filePrefix + "-log.txt";
        String fileNameDetail = filePrefix + "-detail.txt";

        File logDetail = new File(fileNameDetail);
        FileWriter frDetail = new FileWriter(logDetail, true);
        frDetail.write("Median,,,\n");
        frDetail.write("Average,,,\n");
        frDetail.write("Count,,,\n");
        frDetail.write("Timeouts,,,\n");
        frDetail.write(",Join, MintTo, Transfer\n");
        frDetail.close();

        joinStats = new Stats(fileName, "Join");
        mintStats = new Stats(fileName, "Mint");
        sendStats = new Stats(fileName, "Send");

        while (true) {
            System.out.println("");
            Ed25519PrivateKey privateKey = Ed25519PrivateKey.generate();
            Ed25519PublicKey publicKey = privateKey.publicKey;
            String userName = "End2End_" + testId + "_" + Long.toString(cycle);
            String detail = "";

            // join network
            if (join(userName, publicKey)) {
                // mint
                if (mintTo(userName, publicKey)) {
                    if (cycle != 0) {
                        if ( ! send(privateKey, publicKey)) {
                            delayTest = true;
                        }
                    }
                } else {
                    delayTest = true;
                }
                detail += "," + joinStats.lastValue;
                detail += "," + mintStats.lastValue;
                detail += "," + sendStats.lastValue;
                detail += "\n";

                frDetail = new FileWriter(logDetail, true);
                frDetail.write(detail);
                frDetail.close();

                usernames.add(userName);
                privateKeys.add(privateKey);
                publicKeys.add(publicKey);

                // increase cycles
                cycle = cycle + 1;
                // print averages
                if (cycle > 1) {
                    File logFile = new File(fileName);
                    FileWriter fr = new FileWriter(logFile, true);
                    fr.write("Operation, Min, Max, 0-10 count, 0-10 average, 11-15 count, 11-15 average, 16-20 count, 16-20 average, 20- count, 20- average\n");
                    joinStats.printOut(fr);
                    mintStats.printOut(fr);
                    sendStats.printOut(fr);

                    float max = joinStats.max;
                    if (mintStats.max > max) { max = mintStats.max; }
                    if (sendStats.max > max) { max = sendStats.max; }
                    float min = joinStats.min;
                    if (mintStats.min < min) { min = mintStats.min; }
                    if (sendStats.min < min) { min = sendStats.min; }
                    float count_0_10 = joinStats.range_0_10 + mintStats.range_0_10 + sendStats.range_0_10;
                    float total_0_10 = joinStats.range_0_10_total + mintStats.range_0_10_total + sendStats.range_0_10_total;
                    float avg_0_10 = 0;
                    if (count_0_10 != 0) {
                        avg_0_10 = total_0_10 / count_0_10;
                    }
                    float count_11_15 = joinStats.range_11_15 + mintStats.range_11_15 + sendStats.range_11_15;
                    float total_11_15 = joinStats.range_11_15_total + mintStats.range_11_15_total + sendStats.range_11_15_total;
                    float avg_11_15 = 0;
                    if (count_11_15 != 0) {
                        avg_11_15 = total_11_15 / count_11_15;
                    }
                    float count_16_20 = joinStats.range_16_20 + mintStats.range_16_20 + sendStats.range_16_20;
                    float total_16_20 = joinStats.range_16_20_total + mintStats.range_16_20_total + sendStats.range_16_20_total;
                    float avg_16_20 = 0;
                    if (count_16_20 != 0) {
                        avg_16_20 = (float)total_16_20 / count_16_20;
                    }
                    float count_20_plus = joinStats.range_20_plus + mintStats.range_20_plus + sendStats.range_20_plus;
                    float total_20_plus = joinStats.range_20_total + mintStats.range_20_total + sendStats.range_20_total;
                    float avg_20_plus = 0;
                    if (count_20_plus != 0) {
                        avg_20_plus = (float)total_20_plus / count_20_plus;
                    }
                    float count_overall = count_0_10 + count_11_15 + count_16_20 + count_20_plus;
                    float total_overall = total_0_10 + total_11_15 + total_16_20 + total_20_plus;
                    float avg_overall = 0;
                    if (count_overall != 0) {
                        avg_overall = (float)total_overall / count_overall;
                    }
                    String output = String.format("%-9s, %-3.2f, %-3.2f, %-10.0f, %-12.2f, %-11.0f, %-13.2f, %-11.0f, %-13.2f, %-9.0f, %-11.2f\n"
                            , "Summary"
                            , min
                            , max
                            , count_0_10
                            , avg_0_10
                            , count_11_15
                            , avg_11_15
                            , count_16_20
                            , avg_16_20
                            , count_20_plus
                            , avg_20_plus
                    );
                    output += String.format("Overall average: %.2f\n", avg_overall);
                    output += String.format("API Timeouts (>%.0f): Join=%.0f, Buy=%.0f, Send=%.0f\n", timeout / 1000, joinStats.apiErrors, mintStats.apiErrors, sendStats.apiErrors);
                    fr.write(output);
                    fr.close();

                }
            } else {
                delayTest = true;
            }

            if (delayTest) {
                System.out.print("Sleeping 3s after timeout");
                sleep(3);
                delayTest = false;
            } else {
                //sleep(3);
            }
        }
    }

    private static void sleep(int duration) throws InterruptedException {
        for (int i=0; i < duration; i++) {
            System.out.print(".");
            Thread.sleep(1000);
        }
        System.out.println("");
    }
    private static boolean join(String userName, Ed25519PublicKey publicKey) throws InterruptedException {
        RESTJoin restJoin = new RESTJoin(userName, publicKey.toString());
        Instant startTime = Instant.now();
        if (postJoin(restJoin)) {
            float duration = Instant.now().toEpochMilli() - startTime.toEpochMilli();
            System.out.print(String.format("%s Join api (%.2fs)", DATE_TIME_FORMATTER.format(startTime), duration / 1000));

            long queryStart = Instant.now().toEpochMilli();
            long queryCount = 0;
            while (true) {
                queryCount += 1;
                Thread.sleep(queryInterval);
                float callTime = Instant.now().toEpochMilli();
                if (validUser(userName)) {
                    duration = Instant.now().toEpochMilli() - startTime.toEpochMilli();
                    callTime = Instant.now().toEpochMilli() - callTime;
                    joinStats.setValues(duration / 1000);
                    System.out.println(", complete (" + duration / 1000 + "s), user api get call (" + callTime / 1000 + "s), queryCount (" + queryCount + ")");
                    return true;
                } else if (Instant.now().toEpochMilli() - queryStart > timeout) {
                    System.out.println(" Timeout for user " + userName + " after " + timeout / 1000 + " seconds and " + queryCount + " queries");
                    joinStats.setValues(100);
                    return false;
                }
            }
        } else {
            System.out.println(startTime + " Join api failed");
            return false;
        }
    }

    private static boolean mintTo(String userName, Ed25519PublicKey publicKey) throws InterruptedException {
        RESTMintTo restMintTo = new RESTMintTo(userName, Long.toString(initialBalance));
        Instant startTime = Instant.now();
        if (postMintTo(restMintTo)) {
            float duration = Instant.now().toEpochMilli() - startTime.toEpochMilli();
            System.out.print(String.format("%s MintTo api (%.2fs)", DATE_TIME_FORMATTER.format(startTime), duration / 1000));
            long queryStart = Instant.now().toEpochMilli();
            long queryCount = 0;
            while (true) {
                queryCount += 1;
                Thread.sleep(queryInterval);
                float callTime = Instant.now().toEpochMilli();
                if (balance(userName) == initialBalance) {
                    duration = Instant.now().toEpochMilli() - startTime.toEpochMilli();
                    callTime = Instant.now().toEpochMilli() - callTime;
                    mintStats.setValues(duration / 1000);
                    System.out.println(", complete (" + duration / 1000 + "s), user api get call (" + callTime / 1000 + "s), queryCount (" + queryCount + ")");
                    return true;
                } else if (Instant.now().toEpochMilli() - queryStart > timeout) {
                    System.out.println(" Timeout for user " + userName + " minting " + initialBalance + " after " + timeout / 1000 + " seconds and " + queryCount + " queries");
                    mintStats.setValues(100);
                    return false;
                }
            }
        } else {
            System.out.println(startTime + " MintTo api failed");
            return false;
        }
    }

    private static boolean send(Ed25519PrivateKey privateKey, Ed25519PublicKey publicKey) throws InterruptedException, UnsupportedEncodingException {
        int userToSendTo = random.nextInt(usernames.size());
        int quantity = random.nextInt(100) + 1;
        String sendPrimitive = Primitives.sendPrimitive(usernames.get(userToSendTo), privateKey, publicKey, quantity);
        RESTPrimitive restPrimitive = new RESTPrimitive(sendPrimitive);
        Instant startTime = Instant.now();
        long currentBalance = balance(usernames.get(userToSendTo));
        if (postPrimitive(restPrimitive)) {
            float duration = Instant.now().toEpochMilli() - startTime.toEpochMilli();
            System.out.print(String.format("%s Transfer api (%.2fs)", DATE_TIME_FORMATTER.format(startTime), duration / 1000));
            long queryStart = Instant.now().toEpochMilli();
            long queryCount = 0;
            while (true) {
                queryCount += 1;
                Thread.sleep(queryInterval);
                float callTime = Instant.now().toEpochMilli();
                if (balance(usernames.get(userToSendTo)) != currentBalance) {
                    duration = Instant.now().toEpochMilli() - startTime.toEpochMilli();
                    callTime = Instant.now().toEpochMilli() - callTime;
                    sendStats.setValues(duration / 1000);
                    System.out.println(", complete (" + duration / 1000 + "s), user api get call (" + callTime / 1000 + "s), queryCount (" + queryCount + ")");
                    return true;
                } else if (Instant.now().toEpochMilli() - queryStart > timeout) {
                    sendStats.setValues(100);
                    System.out.println(" Timeout from " + publicKey + " to " + usernames.get(userToSendTo) + " transferring " + quantity + " after " + timeout / 1000 + " seconds and " + queryCount + " queries");
                    return false;
                }
            }
        } else {
            System.out.println(startTime + " Transfer api failed");
            return false;
        }
    }

    private static boolean postJoin(RESTJoin restJoin) {
        String url = "/v1/token/join";
        Response response = client.target(sendHost + url).request(MediaType.APPLICATION_JSON).post(Entity.entity(restJoin, MediaType.APPLICATION_JSON));
        return (response.getStatus() == 202);
    }

    private static boolean postMintTo(RESTMintTo restMintTo) {
        String url = "/v1/token/mintTo";
        Response response = client.target(sendHost + url).request(MediaType.APPLICATION_JSON).post(Entity.entity(restMintTo, MediaType.APPLICATION_JSON));
        return (response.getStatus() == 202);
    }

    private static boolean postPrimitive(RESTPrimitive restPrimitive) {
        String url = "/v1/token/transaction";
        Response response = client.target(sendHost + url).request(MediaType.APPLICATION_JSON).post(Entity.entity(restPrimitive, MediaType.APPLICATION_JSON));
        return (response.getStatus() == 202);
    }

    private static boolean validUser(String username) {
        String url = "/v1/token/userExists/";
        RESTUserExists userExists = client.target(queryHost + url + username).request(MediaType.APPLICATION_JSON).get(RESTUserExists.class);
        return userExists.exists;
    }

    private static long balance(String username) {
        String url = "/v1/token/balance/";
        RESTBalance balance = client.target(queryHost + url + username).request(MediaType.APPLICATION_JSON).get(RESTBalance.class);
        return balance.balance;
    }
}
