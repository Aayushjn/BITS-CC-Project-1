package com.springboot.TravelManagementSystem.exceptions;

public class AlreadyExistsException extends Exception{
    public AlreadyExistsException() {
        super();
    }

    public AlreadyExistsException(String message) {
        super(message);
    }
}
