package com.hedera.stabl.test.model;

public class RESTBalance {
    public long balance;
    public boolean frozen;

    public RESTBalance() {
    }

    public RESTBalance(long balance, boolean frozen) {
        this.balance = balance;
        this.frozen = false;
    }

    public long getBalance() {
        return this.balance;
    }
    public void setBalance(long balance) {
        this.balance = balance;
    }
    public boolean getFrozen() {
        return this.frozen;
    }
    public void setFrozen(boolean frozen) {
        this.frozen = frozen;
    }
}
