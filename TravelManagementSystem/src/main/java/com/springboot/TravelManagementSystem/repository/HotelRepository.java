package com.springboot.TravelManagementSystem.repository;

import com.springboot.TravelManagementSystem.models.Hotel;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface HotelRepository extends JpaRepository<Hotel,Integer> {


}
