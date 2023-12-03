package com.springboot.TravelManagementSystem.repository;

import com.springboot.TravelManagementSystem.models.Packages;
import com.springboot.TravelManagementSystem.models.Ticket;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface TicketRepository extends JpaRepository<Ticket,Integer> {
    public List<Ticket> findByPackages(Packages packages);

}
