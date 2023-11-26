package com.springboot.TravelManagementSystem.services;

import com.springboot.TravelManagementSystem.exceptions.BookingException;
import com.springboot.TravelManagementSystem.exceptions.InvalidCredentialException;
import com.springboot.TravelManagementSystem.models.Booking;

import java.util.List;

public interface BookingServices {

    public Booking makeBooking(Booking bookings,String authKey) throws BookingException , InvalidCredentialException;
    public Booking cancelBooking(Integer bookingId) throws BookingException;

    public List<Booking> viewBooking(Integer userId) throws BookingException;


    //admin can access

    public List<Booking> viewAllBooking(String auhKey) throws BookingException;

}
