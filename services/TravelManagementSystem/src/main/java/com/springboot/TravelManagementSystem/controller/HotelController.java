package com.springboot.TravelManagementSystem.controller;

import com.springboot.TravelManagementSystem.exceptions.HotelException;
import com.springboot.TravelManagementSystem.exceptions.PackageException;
import com.springboot.TravelManagementSystem.models.Hotel;
import com.springboot.TravelManagementSystem.models.Packages;
import com.springboot.TravelManagementSystem.services.HotelService;
import com.springboot.TravelManagementSystem.services.HotelServiceProvider;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/hotel")
@CrossOrigin(origins = "*")
public class HotelController {
    @Autowired
    HotelServiceProvider hotelServiceProvider;


    @GetMapping("")
    public ResponseEntity<List<Hotel>> getAllHotels() throws HotelException
    {
        List<Hotel> hotelList=hotelServiceProvider.getAllHotel();
        return new ResponseEntity<List<Hotel>>(hotelList, HttpStatus.OK);
    }

    @PostMapping("")
    public ResponseEntity<Hotel> addHotel(@RequestBody @Valid Hotel hotel, @RequestParam(value = "key") String key,@RequestParam(value="packageId") Integer packageId ) throws HotelException
    {
        Hotel returnHotel=hotelServiceProvider.addHotel(packageId,hotel,key);
        return new ResponseEntity<Hotel>(returnHotel, HttpStatus.CREATED);
    }

    @DeleteMapping("")
    public ResponseEntity<Hotel> deleteHotel(@RequestParam(value = "hotelId") Integer hotelId,@RequestParam(value = "key") String key) throws HotelException
    {
        Hotel returnHotel=hotelServiceProvider.deleteHotel(hotelId,key);
        return new ResponseEntity<Hotel>(returnHotel,HttpStatus.OK);
    }

    @PutMapping("")
    public ResponseEntity<Hotel> updateHotel(@RequestBody @Valid Hotel hotel,@RequestParam(value = "key") String key) throws HotelException
    {
        Hotel returnHotel=hotelServiceProvider.updateHotel(hotel,key);
        return new ResponseEntity<Hotel>(returnHotel,HttpStatus.OK);
    }

}
