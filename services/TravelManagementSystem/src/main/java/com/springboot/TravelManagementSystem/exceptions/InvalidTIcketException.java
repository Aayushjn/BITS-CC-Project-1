package com.springboot.TravelManagementSystem.exceptions;

public class InvalidTIcketException extends Exception{
    public InvalidTIcketException() {
        super();
    }

    public InvalidTIcketException(String message) {
        super(message);
    }
}
