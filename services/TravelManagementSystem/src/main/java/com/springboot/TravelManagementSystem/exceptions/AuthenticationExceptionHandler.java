package com.springboot.TravelManagementSystem.exceptions;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.context.request.WebRequest;

import java.time.LocalDateTime;

@ControllerAdvice
public class AuthenticationExceptionHandler {
    @ExceptionHandler(InvalidCredentialException.class)
    public ResponseEntity<ErrorDetails> InvalidCredentialExceptionHandler(InvalidCredentialException ice, WebRequest req){
            ErrorDetails authEx=new ErrorDetails();
            authEx.setTimestamp(LocalDateTime.now());
            authEx.setMessage(ice.getMessage());
            authEx.setDescription(req.getDescription(false));

            return new ResponseEntity<ErrorDetails>(authEx, HttpStatus.BAD_REQUEST);
    }
    @ExceptionHandler(AlreadyExistsException.class)
    public ResponseEntity<ErrorDetails> DuplicateSignupExceptionHandler(AlreadyExistsException aee, WebRequest req){
        ErrorDetails authEx=new ErrorDetails();
        authEx.setTimestamp(LocalDateTime.now());
        authEx.setMessage(aee.getMessage());
        authEx.setDescription(req.getDescription(false));

        return new ResponseEntity<ErrorDetails>(authEx, HttpStatus.BAD_REQUEST);
    }
}
