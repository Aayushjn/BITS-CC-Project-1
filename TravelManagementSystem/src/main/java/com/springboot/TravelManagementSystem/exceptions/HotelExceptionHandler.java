package com.springboot.TravelManagementSystem.exceptions;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.context.request.WebRequest;

import java.time.LocalDateTime;

@ControllerAdvice
public class HotelExceptionHandler {

    @ExceptionHandler(HotelException.class)
    public ResponseEntity<ErrorDetails> OtherExceptionHandler(HotelException he, WebRequest req)
    {
        ErrorDetails er=new ErrorDetails();
        er.setTimestamp(LocalDateTime.now());
        er.setMessage(he.getMessage());
        er.setDescription(req.getDescription(false));
        return new ResponseEntity<ErrorDetails>(er, HttpStatus.BAD_REQUEST);
    }


}
