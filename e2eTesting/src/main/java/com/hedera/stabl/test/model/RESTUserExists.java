package com.hedera.stabl.test.model;

public class RESTUserExists {
    public boolean exists;

    public RESTUserExists() {
    }

    public RESTUserExists(boolean exists) {
        this.exists = exists;
    }

    public boolean getExists() {
        return this.exists;
    }
    public void setExists(boolean exists) {
        this.exists = exists;
    }
}
