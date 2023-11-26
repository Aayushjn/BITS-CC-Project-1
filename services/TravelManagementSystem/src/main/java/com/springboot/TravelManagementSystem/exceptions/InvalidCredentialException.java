package com.springboot.TravelManagementSystem.exceptions;

public class InvalidCredentialException extends Exception{
    public InvalidCredentialException() {
        super();
    }

    public InvalidCredentialException(String message) {
        super(message);
    }
}
