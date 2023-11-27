package com.springboot.TravelManagementSystem.services;

import com.springboot.TravelManagementSystem.exceptions.InvalidCredentialException;
import com.springboot.TravelManagementSystem.exceptions.InvalidTIcketException;
import com.springboot.TravelManagementSystem.exceptions.PackageException;
import com.springboot.TravelManagementSystem.models.CurrentUserLoginSession;
import com.springboot.TravelManagementSystem.models.Packages;
import com.springboot.TravelManagementSystem.models.Ticket;
import com.springboot.TravelManagementSystem.models.User;
import com.springboot.TravelManagementSystem.repository.PackageRepository;
import com.springboot.TravelManagementSystem.repository.SessionRepository;
import com.springboot.TravelManagementSystem.repository.TicketRepository;
import com.springboot.TravelManagementSystem.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Service
public class TicketServiceHandler implements TicketService{

    @Autowired
    TicketRepository ticketRepository;

    @Autowired
    PackageRepository packageRepository;

    @Autowired
    SessionRepository sessionRepository;

    @Autowired
    UserRepository userRepository;


    @Override
    public Ticket createTicket(Ticket ticket, String authKey, Integer packageId) throws InvalidTIcketException, InvalidCredentialException, PackageException {

        Optional<CurrentUserLoginSession> culs=sessionRepository.findByAuthKey(authKey);

        if(!culs.isPresent())
        {
            throw new InvalidCredentialException("Login In First");
        }

        Optional<User> user=userRepository.findById(culs.get().getUserId());
        Optional<Packages> pkgOpt=packageRepository.findById(packageId);

        if(pkgOpt.isPresent())
        {
            ticket.setPackages(pkgOpt.get());
        }
        Ticket ticketCreated=ticketRepository.save(ticket);
        return ticketCreated;
    }

    @Override
    public Ticket getTicket(Integer ticketId) throws InvalidTIcketException {
        Optional<Ticket> ticket=ticketRepository.findById(ticketId);

        if(ticket.isPresent()==false) throw new InvalidTIcketException("No such ticket exist");

        return ticket.get();
    }

    @Override
    public Ticket deleteTicket(Integer ticketId) throws InvalidTIcketException {

        Ticket ticket=null;

        Optional<Ticket> user=ticketRepository.findById(ticketId);

        if(user.isPresent())
        {
            ticket=user.get();
            ticket.setTicketStatus(false);
            ticketRepository.save(ticket);
        }
        else
        {
            throw new InvalidTIcketException("Ticket Doesnot Exist");
        }
        return ticket;
    }

    @Override
    public List<Ticket> getAllTicket(Integer packageId, String key) throws InvalidTIcketException {

        List<Ticket> tickets=null;

        Optional<CurrentUserLoginSession> culs=sessionRepository.findByAuthKey(key);

        String userType=userRepository.findById(culs.get().getUserId()).get().getUserType();
        if(userType.equalsIgnoreCase("user") && packageId!=null )
        {
            Optional<Packages> pkg=packageRepository.findById(packageId);
            Packages packages=pkg.get();
            tickets=ticketRepository.findByPackages(packages);
        }
        else if(userType.equalsIgnoreCase("admin") && packageId==null)
        {
            tickets=ticketRepository.findAll();
        } else if (userType.equalsIgnoreCase("admin") && packageId!=null) {
            Optional<Packages> pkg=packageRepository.findById(packageId);
            Packages packages=pkg.get();
            tickets=ticketRepository.findByPackages(packages);
        }
        if(tickets.isEmpty() || tickets==null)
        {
            throw new InvalidTIcketException("Tickets Not Avaliable");
        }

        return tickets;
    }
}
