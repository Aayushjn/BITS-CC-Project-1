package com.springboot.TravelManagementSystem.repository;

import com.springboot.TravelManagementSystem.models.Booking;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface BookingRepository extends JpaRepository<Booking,Integer> {



}
