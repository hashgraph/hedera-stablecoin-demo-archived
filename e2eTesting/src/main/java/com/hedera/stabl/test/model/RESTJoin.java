package com.hedera.stabl.test.model;

public class RESTJoin {
    public String username;
    public String publicKey;

    public RESTJoin() {
    }

    public RESTJoin(String username, String publicKey) {
        this.username = username;
        this.publicKey = publicKey;
    }

    public String getPublicKey() {
        return this.publicKey;
    }
    public void setPublicKey(String publicKey) {
        this.publicKey = publicKey;
    }
    public String getUsername() {
        return this.username;
    }

    public void setUsername(String username) {
        this.username = username;
    }
}
