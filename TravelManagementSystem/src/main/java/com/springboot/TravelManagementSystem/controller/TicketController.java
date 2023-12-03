package com.springboot.TravelManagementSystem.controller;

import com.springboot.TravelManagementSystem.exceptions.InvalidCredentialException;
import com.springboot.TravelManagementSystem.exceptions.InvalidTIcketException;
import com.springboot.TravelManagementSystem.exceptions.PackageException;
import com.springboot.TravelManagementSystem.models.Ticket;
import com.springboot.TravelManagementSystem.services.TicketServiceHandler;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/tickets")
@CrossOrigin(origins = "*")
public class TicketController {

    @Autowired
    TicketServiceHandler ticketServiceHandler;

    @GetMapping("/all")
    public ResponseEntity<List<Ticket>> getAllTickets(@RequestParam(required = false) Integer packageId,@RequestParam(required = false) String key) throws InvalidTIcketException
    {
        List<Ticket> tickets=ticketServiceHandler.getAllTicket(packageId,key);

        return new ResponseEntity<List<Ticket>>(tickets, HttpStatus.OK);
    }

    @GetMapping("/")
    public ResponseEntity<Ticket> getTicket(@RequestParam Integer ticketId) throws InvalidTIcketException
    {
        Ticket ticket=ticketServiceHandler.getTicket(ticketId);
        System.out.println(ticket.getTicketId());
        //Map<String,Integer> map=new HashMap<>();
        //map.put("Ticket_ID",ticket.getTicketId());
        return new ResponseEntity<Ticket>(ticket,HttpStatus.FOUND);
    }

    @PostMapping("/")
    public ResponseEntity<Ticket> createTicket(@Valid @RequestBody Ticket ticket,@RequestParam("key") String key,@RequestParam(required = false,defaultValue = "0") Integer packageId) throws InvalidTIcketException , InvalidCredentialException, PackageException
    {
        Ticket createdTicket=null;
        createdTicket=ticketServiceHandler.createTicket(ticket,key,packageId);
        return new ResponseEntity<Ticket>(createdTicket,HttpStatus.CREATED);
    }

    @DeleteMapping("/<ticketId>")
    public ResponseEntity<Map<String, Integer>> cancelTicket(@RequestParam Integer ticketId) throws InvalidTIcketException
    {
        Ticket ticket=null;
        Map<String, Integer> map=new HashMap<>();
        ticket=ticketServiceHandler.deleteTicket(ticketId);
        map.put("Ticket_Cancelled_ID",ticket.getTicketId());
        return new ResponseEntity<Map<String ,Integer>>(map,HttpStatus.OK);
    }

}
