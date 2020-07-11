package com.hedera.stabl.test.model;

public class RESTMintTo {
    public String address;
    public String quantity;

    public RESTMintTo() {
    }

    public RESTMintTo(String address, String quantity) {
        this.address = address;
        this.quantity = quantity;
    }

    public String getQuantity() {
        return this.quantity;
    }
    public void setQuantity(String quantity) {
        this.quantity = quantity;
    }
    public String getAddress() {
        return this.address;
    }
    public void setAddress(String address) {
        this.address = address;
    }
}
