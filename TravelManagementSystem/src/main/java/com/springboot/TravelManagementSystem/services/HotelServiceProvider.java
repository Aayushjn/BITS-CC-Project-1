package com.springboot.TravelManagementSystem.services;

import com.springboot.TravelManagementSystem.exceptions.HotelException;
import com.springboot.TravelManagementSystem.exceptions.PackageException;
import com.springboot.TravelManagementSystem.models.*;
import com.springboot.TravelManagementSystem.repository.HotelRepository;
import com.springboot.TravelManagementSystem.repository.PackageRepository;
import com.springboot.TravelManagementSystem.repository.SessionRepository;
import com.springboot.TravelManagementSystem.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Service
public class HotelServiceProvider implements HotelService {

    @Autowired
    private HotelRepository hotelRepository;

    @Autowired
    private SessionRepository sessionRepository;

    @Autowired
    private UserRepository userRepository;

    @Autowired
    private PackageRepository packageRepository;

    @Override
    public Hotel addHotel(Integer pkgId, Hotel hotel, String authKey) throws HotelException {
        //return null;
        Optional<CurrentUserLoginSession> opt = sessionRepository.findByAuthKey(authKey);

        if (!opt.isPresent()) {
            throw new HotelException("Admin login First .... Current Admin Not Logged in ");
        } else {
            CurrentUserLoginSession currentUserLoginSession = opt.get();
            Optional<User> op2 = userRepository.findById(currentUserLoginSession.getUserId());
            User user = op2.get();
            if (user.getUserType().equalsIgnoreCase("admin")) {
                Optional<Packages> optionalPackages = packageRepository.findById(pkgId);
                if (!optionalPackages.isPresent()) {
                    throw new HotelException("Package with given package Id doesnt exist! ");
                }

                Packages packages = optionalPackages.get();

                hotel.setPackages(packages);
                Hotel returnHotel = hotelRepository.save(hotel);
                hotel.setStatus(HotelStatus.BOOKED);
                return returnHotel;
            } else {
                throw new HotelException("Kindly login as admin");
            }
        }


    }

    @Override
    public Hotel deleteHotel(Integer hotelId, String authKey) throws HotelException {

        Optional<CurrentUserLoginSession> optional = sessionRepository.findByAuthKey(authKey);
        if (!optional.isPresent()) {
            throw new HotelException(" Admin login First .... Current Admin Not Logged in ...");
        } else {
            CurrentUserLoginSession culs = optional.get();
            Optional<User> opt2 = userRepository.findById(culs.getUserId());
            User user = opt2.get();
            if (user.getUserType().equals("admin")) {
                Optional<Hotel> hotelOptional = hotelRepository.findById(hotelId);
                if (!hotelOptional.isPresent()) {
                    throw new HotelException("No hotel is present with this id");
                } else {
                    Hotel hotel = hotelOptional.get();
                    Packages packages = hotel.getPackages();

                    hotel.setPackages(null);
                    hotelRepository.delete(hotel);
                    packageRepository.save(packages);
                    hotel.setStatus(HotelStatus.CANCELLED);
                    return hotel;
                }
            } else {
                throw new HotelException("Kindly login as admin!!");
            }
        }
//        return null;
    }

    @Override
    public Hotel updateHotel(Hotel hotel, String authKey) throws HotelException {

        Optional<CurrentUserLoginSession> optional = sessionRepository.findByAuthKey(authKey);
        if (!optional.isPresent()) {
            throw new HotelException("Admin login First .... Current Admin Not Logged in ");
        } else {
            CurrentUserLoginSession currentUserLoginSession = optional.get();
            Optional<User> op2 = userRepository.findById(currentUserLoginSession.getUserId());
            User user = op2.get();
            if (user.getUserType().equalsIgnoreCase("admin")) {
                Optional<Hotel> optionalHotel = hotelRepository.findById(hotel.getHotelId());
                if (!optionalHotel.isPresent()) {
                    throw new HotelException("No Hotel present with the given details !");
                } else {
                    Hotel returnedHotel = optionalHotel.get();
                    if (hotel.getAddress() != null) {
                        returnedHotel.setAddress(hotel.getAddress());
                    }
                    if (hotel.getHotelDescription() != null) {
                        returnedHotel.setHotelDescription(hotel.getHotelDescription());
                    }
                    if (hotel.getHotelName() != null) {
                        returnedHotel.setHotelName(hotel.getHotelName());
                    }
                    if (hotel.getHotelType() != null) {
                        returnedHotel.setHotelType(hotel.getHotelType());
                    }
                    if (hotel.getRent() != 0) {
                        returnedHotel.setRent(hotel.getRent());
                    }
                    if (hotel.getStatus() != null) {
                        returnedHotel.setStatus(hotel.getStatus());
                    }
                    return hotelRepository.save(returnedHotel);
                }

            } else {
                throw new HotelException("Kindly login as a admin");
            }
        }
    }

    @Override
    public List<Hotel> getAllHotel() throws HotelException {
        List<Hotel> hotelList=hotelRepository.findAll();
        if(hotelList.isEmpty())
        {
            throw new HotelException("Hotel Not Avaliable");
        }
        return hotelList;
    }
}