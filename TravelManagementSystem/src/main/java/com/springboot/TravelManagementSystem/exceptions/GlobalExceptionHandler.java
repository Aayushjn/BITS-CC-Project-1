package com.springboot.TravelManagementSystem.exceptions;

import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.MethodArgumentNotValidException;
import org.springframework.web.bind.annotation.ControllerAdvice;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.context.request.WebRequest;

import java.nio.file.AccessDeniedException;
import java.time.LocalDateTime;

@ControllerAdvice
public class GlobalExceptionHandler {


    @ExceptionHandler(MethodArgumentNotValidException.class)
    public ResponseEntity<ErrorDetails> MethodArgumentNotValidException(MethodArgumentNotValidException maexp)
    {
        ErrorDetails exp=new ErrorDetails();
        exp.setTimestamp(LocalDateTime.now());
        exp.setMessage("Validation Error");
        exp.setDescription(maexp.getBindingResult().getFieldError().getDefaultMessage());
        return new ResponseEntity<ErrorDetails>(exp,HttpStatus.BAD_REQUEST);
    }


    @ExceptionHandler(AccessDeniedException.class)
    public ResponseEntity<ErrorDetails> AccessDeniedExceptionHandler(AccessDeniedException ade,WebRequest req)
    {
        ErrorDetails exp=new ErrorDetails();
        exp.setTimestamp(LocalDateTime.now());
        exp.setMessage(ade.getMessage());
        exp.setDescription(req.getDescription(false));
        return new ResponseEntity<ErrorDetails>(exp,HttpStatus.FORBIDDEN);
    }

    @ExceptionHandler(InvalidTIcketException.class)
    public ResponseEntity<ErrorDetails> InvalidTIcketExceptionHandler(InvalidTIcketException ite,WebRequest req)
    {
        ErrorDetails exp=new ErrorDetails();
        exp.setTimestamp(LocalDateTime.now());
        exp.setMessage(ite.getMessage());
        exp.setDescription(req.getDescription(false));
        return new ResponseEntity<ErrorDetails>(exp,HttpStatus.NO_CONTENT);
    }


    @ExceptionHandler(PackageException.class)
    public ResponseEntity<ErrorDetails> packageExceptionHandler(PackageException pe,WebRequest req)
    {
        ErrorDetails pkExp=new ErrorDetails();
        pkExp.setTimestamp(LocalDateTime.now());
        pkExp.setMessage(pe.getMessage());
        pkExp.setDescription(req.getDescription(false));

        return new ResponseEntity<ErrorDetails>(pkExp,HttpStatus.BAD_REQUEST);

    }

    @ExceptionHandler(Exception.class)
    public ResponseEntity<ErrorDetails> OtherExceptionHandler(Exception exp, WebRequest req)
    {
        ErrorDetails authEx=new ErrorDetails();
        authEx.setTimestamp(LocalDateTime.now());
        authEx.setMessage(exp.getMessage());
        authEx.setDescription(req.getDescription(false));

        return new ResponseEntity<ErrorDetails>(authEx, HttpStatus.BAD_REQUEST);
    }
}
