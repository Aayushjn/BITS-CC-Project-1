package com.springboot.TravelManagementSystem.controller;

import com.springboot.TravelManagementSystem.exceptions.BookingException;
import com.springboot.TravelManagementSystem.exceptions.InvalidCredentialException;
import com.springboot.TravelManagementSystem.models.Booking;
import com.springboot.TravelManagementSystem.services.BookingServices;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/bookings")
@CrossOrigin(origins = "*")
public class BookingController {

    @Autowired
    BookingServices bookingServices;

    @DeleteMapping("/")
    public ResponseEntity<Map<String,String>> cancelBooking(@RequestParam Integer bookingId) throws BookingException
    {
        Booking booking=null;
        booking=bookingServices.cancelBooking(bookingId);
        Map<String,String> map=new HashMap<>();
        map.put("status ", String.valueOf(booking.getStatus()));
        return new ResponseEntity<Map<String,String>>(map, HttpStatus.OK);
    }

    @PostMapping("/")
    public ResponseEntity<Map<String,Integer>> makeBookings(@Valid @RequestBody Booking booking, @RequestParam String key) throws BookingException, InvalidCredentialException{
        Booking createBooking=null;
        //System.out.println(booking);
        createBooking=bookingServices.makeBooking(booking,key);
        Map<String,Integer> map=new HashMap<>();
        map.put("Booking Done with id ", createBooking.getBookingId());
        return new ResponseEntity<Map<String,Integer>>(map, HttpStatus.CREATED);
    }

    @GetMapping("/")
    public ResponseEntity<List<Booking>> viewBooking(@RequestParam(required = false) Integer userId) throws BookingException
    {
        List<Booking> booking=null;
        booking=bookingServices.viewBooking(userId);
        return new ResponseEntity<List<Booking>>(booking,HttpStatus.OK);
    }

    @GetMapping("/all")
    public ResponseEntity<List<Booking>> viewAllBooking(@RequestParam(required = false) String authKey) throws BookingException
    {
        List<Booking> bookings=null;
        bookings=bookingServices.viewAllBooking(authKey);
        return new ResponseEntity<List<Booking>>(bookings,HttpStatus.OK);
    }

}
