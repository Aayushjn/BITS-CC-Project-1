package com.springboot.TravelManagementSystem.repository;

import com.springboot.TravelManagementSystem.models.Packages;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface PackageRepository extends JpaRepository<Packages,Integer> {

}
