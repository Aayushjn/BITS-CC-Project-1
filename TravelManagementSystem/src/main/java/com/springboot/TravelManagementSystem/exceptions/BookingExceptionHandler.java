package com.springboot.TravelManagementSystem.exceptions;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.context.request.WebRequest;

import java.time.LocalDateTime;

public class BookingExceptionHandler {

    @ExceptionHandler(BookingException.class)
    public ResponseEntity<ErrorDetails> bookingExceptionHandler(BookingException be, WebRequest req)
    {
        ErrorDetails bookException=new ErrorDetails();
        bookException.setTimestamp(LocalDateTime.now());
        bookException.setMessage(be.getMessage());
        bookException.setDescription(req.getDescription(false));

        return new ResponseEntity<ErrorDetails>(bookException, HttpStatus.BAD_REQUEST);
    }

}
