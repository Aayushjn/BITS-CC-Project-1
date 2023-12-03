package com.springboot.TravelManagementSystem.services;

import com.springboot.TravelManagementSystem.exceptions.HotelException;
import com.springboot.TravelManagementSystem.exceptions.PackageException;
import com.springboot.TravelManagementSystem.models.Hotel;
import com.springboot.TravelManagementSystem.models.Packages;

import java.util.List;

public interface HotelService {

    public Hotel addHotel(Integer pkgId,Hotel hotel,String authKey) throws HotelException;

    public Hotel deleteHotel(Integer hotelId,String authKey) throws HotelException;

    public Hotel updateHotel(Hotel hotel,String authKey) throws HotelException;

    public List<Hotel> getAllHotel() throws HotelException;


}
