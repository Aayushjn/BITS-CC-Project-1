package com.springboot.TravelManagementSystem.services;

import com.springboot.TravelManagementSystem.exceptions.InvalidCredentialException;
import com.springboot.TravelManagementSystem.exceptions.InvalidTIcketException;
import com.springboot.TravelManagementSystem.exceptions.PackageException;
import com.springboot.TravelManagementSystem.models.Ticket;

import java.util.List;

public interface TicketService {

    //user or admin can create ticket
    public Ticket createTicket(Ticket ticket,String authKey,Integer packageId) throws InvalidTIcketException , InvalidCredentialException, PackageException;

    //user or admin see ticket

    public Ticket getTicket(Integer ticketId) throws InvalidTIcketException;

    // user or admin delete ticket
    public Ticket deleteTicket(Integer ticketId) throws InvalidTIcketException;

    // admin can see all ticket whereas user can see only their tickets

    public List<Ticket> getAllTicket(Integer packageId, String key) throws InvalidTIcketException;

}
