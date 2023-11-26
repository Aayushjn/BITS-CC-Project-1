package com.springboot.TravelManagementSystem.services;

import com.springboot.TravelManagementSystem.exceptions.BookingException;
import com.springboot.TravelManagementSystem.exceptions.InvalidCredentialException;
import com.springboot.TravelManagementSystem.models.Booking;
import com.springboot.TravelManagementSystem.models.BookingStatus;
import com.springboot.TravelManagementSystem.models.CurrentUserLoginSession;
import com.springboot.TravelManagementSystem.models.User;
import com.springboot.TravelManagementSystem.repository.BookingRepository;
import com.springboot.TravelManagementSystem.repository.SessionRepository;
import com.springboot.TravelManagementSystem.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Service
public class BookingServicesProvider implements BookingServices {

    @Autowired
    UserRepository userRepository;

    @Autowired
    SessionRepository sessionRepository;

    @Autowired
    BookingRepository bookingRepository;

    @Override
    public Booking makeBooking(Booking bookings, String authKey) throws BookingException, InvalidCredentialException {

        Optional<CurrentUserLoginSession> session = sessionRepository.findByAuthKey(authKey);
        if(!session.isPresent()) throw new InvalidCredentialException("Please Login First !..");

        Optional<User> user=userRepository.findById(session.get().getUserId());
        bookings.setUser(user.get());
        bookings.setStatus(BookingStatus.Booked);
        return bookingRepository.save(bookings);

    }

    @Override
    public Booking cancelBooking(Integer bookingId) throws BookingException {
        Booking booking=null;
        Optional<Booking> book=bookingRepository.findById(bookingId);
        if(book.isPresent())
        {
            booking=book.get();
            User user=booking.getUser();
            user.getBookings().remove(booking);
            booking.setStatus(BookingStatus.Cancelled);
            return bookingRepository.save(booking);
        }
        else {
            throw new BookingException("Booking Doesnot Exist!...");
        }
    }

    @Override
    public List<Booking> viewBooking(Integer userId) throws BookingException {
        User user=null;
        Optional<User> userOpt=userRepository.findByUserId(userId);
        if(userOpt.isPresent())
        {
            user=userOpt.get();
            List<Booking> bookings=user.getBookings();
            if(bookings.isEmpty())
            {
                throw new BookingException("NO BOOKING EXIST !...");
            }
            return bookings;
        }
        else
        {
            throw new BookingException("No User Exist !..");
        }
    }

    @Override
    public List<Booking> viewAllBooking(String auhKey) throws BookingException {
        Optional<CurrentUserLoginSession> currUser=sessionRepository.findByAuthKey(auhKey);
        String userType=userRepository.findById(currUser.get().getUserId()).get().getUserType();
        List<Booking> bookings=null;
        if(userType.equalsIgnoreCase("user"))
        {
            throw new BookingException("UnAuthorized Request");
        }
        else {
            bookings=bookingRepository.findAll();
            if(bookings.isEmpty())
            {
                throw new BookingException("NO BOOKING AVALIABLE!!!!");
            }
            else {
                return bookings;
            }
        }
    }
}
