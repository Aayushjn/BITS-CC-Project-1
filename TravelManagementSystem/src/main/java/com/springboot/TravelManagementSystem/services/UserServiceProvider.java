package com.springboot.TravelManagementSystem.services;

import com.springboot.TravelManagementSystem.models.CurrentUserLoginSession;
import com.springboot.TravelManagementSystem.models.SessionDTO;
import com.springboot.TravelManagementSystem.models.User;
import com.springboot.TravelManagementSystem.models.UserDTO;
import com.springboot.TravelManagementSystem.repository.SessionRepository;
import com.springboot.TravelManagementSystem.repository.SessionRepository;
import com.springboot.TravelManagementSystem.repository.UserRepository;
import jakarta.transaction.Transactional;
import org.aspectj.weaver.AnnotationNameValuePair;
import com.springboot.TravelManagementSystem.exceptions.AlreadyExistsException;
import com.springboot.TravelManagementSystem.exceptions.InvalidCredentialException;
import org.hibernate.context.spi.CurrentSessionContext;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;
import java.util.Optional;
import java.util.UUID;

@Service
@Transactional
public class UserServiceProvider implements UserService{

    @Autowired
    UserRepository userRepo;

    @Autowired
    SessionRepository sessionRepo;

    @Override
    public User userSignup(User user) throws AlreadyExistsException {
        Optional<User> option=userRepo.findByEmail(user.getEmail());
        if(option.isPresent())
        {
            throw new AlreadyExistsException("USER ALREADY REGISTERED");
        }
        userRepo.save(user);
        return user;
    }

    @Override
    public SessionDTO userLogin(UserDTO user) throws InvalidCredentialException {
        Optional<User> option=userRepo.findByEmail(user.getEmail());
        if(!option.isPresent())
        {
            throw new InvalidCredentialException("USER DOESN'T EXIST !");
        }
        Optional<CurrentUserLoginSession> checkSession=sessionRepo.findByUserId(option.get().getUserId());
        if(checkSession.isPresent())
        {
            throw new InvalidCredentialException("USER ALREADY LOGGED IN !");
        }
        if(!(option.get().getEmail().equals(user.getEmail())))
        {
            throw new InvalidCredentialException("Invalid Email Address");
        }
        else if(!(option.get().getPassword().equals(user.getPassword())))
        {
            throw new InvalidCredentialException("Invalid Password");
        }
        else if(!(option.get().getPassword().equals(user.getPassword()) && option.get().getEmail().equals(user.getEmail())))
        {
            throw new InvalidCredentialException("Invalid Credentials");
        }

        SessionDTO sessionDTO=new SessionDTO();
        CurrentUserLoginSession culs=new CurrentUserLoginSession();
        String authKey= UUID.randomUUID().toString();
        culs.setAuthKey(authKey);
        culs.setSessionStartTime(LocalDateTime.now());
        sessionDTO.setAuthKey(culs.getAuthKey());
        sessionDTO.setSessionStartTime(culs.getSessionStartTime());
        culs.setUserId(option.get().getUserId());
        sessionRepo.save(culs);
        return sessionDTO;
    }

    @Override
    public String userLogout(String authKey) throws InvalidCredentialException {
        Optional<CurrentUserLoginSession> culs=sessionRepo.findByAuthKey(authKey);
        if(!culs.isPresent())
        {
            throw new InvalidCredentialException("User Doesnot logged in with the key : " + authKey);
        }
        sessionRepo.delete(culs.get());
        return " Logout Successfully ! ";
    }

    @Override
    public boolean updateUser(User user) throws InvalidCredentialException {
      Optional<User> checkUser=userRepo.findByEmail(user.getEmail());
        if(!checkUser.isPresent()) throw new InvalidCredentialException("User Doesn't exist with id " + user.getEmail());
        User u=checkUser.get();
        u.setName(user.getName());
        u.setAddress(user.getAddress());
        u.setMobile(user.getMobile());
        userRepo.save(u);
        return u!=null;
    }

    @Override
    public User getUser(String authKey) throws InvalidCredentialException {
        Optional<CurrentUserLoginSession> culs=sessionRepo.findByAuthKey(authKey);
        if(!culs.isPresent())
        {
            throw new InvalidCredentialException("Invalid Authentication Key");
        }
        Optional<User> user=userRepo.findById(culs.get().getUserId());
        return user.get();
    }

    @Override
    public User deleteUser(Integer userId, String authKey) throws InvalidCredentialException {
        Optional<CurrentUserLoginSession> culs=sessionRepo.findByAuthKey(authKey);
        if(!culs.isPresent())
        {
            throw new InvalidCredentialException("Invalid Authentication Key");
        }
        String userType= userRepo.findById(culs.get().getUserId()).get().getUserType();
        if(userType.equalsIgnoreCase("user"))
        {
            throw new InvalidCredentialException("Unauthorized Request!!!");
        }
        Optional<User> user=userRepo.findById(userId);
        if(!user.isPresent()) throw new InvalidCredentialException("User Doesnt Exist with id : " + userId);
        userRepo.delete(user.get());
        return user.get();
    }

    @Override
    public User userAdmin(String userEmail, String passcode) throws InvalidCredentialException {
        if(!passcode.equals("admin"))
        {
            throw new InvalidCredentialException("Invalid passcode !");
        }
        else if(userEmail.equals(null))
        {
            throw new InvalidCredentialException("Invalid Email Address");
        }
        Optional<User> user=userRepo.findByEmail(userEmail);
        if(!user.isPresent())

        {
            throw new InvalidCredentialException("User Doesnt exist with id " + userEmail);
        }
        user.get().setUserType("admin");
        userRepo.save(user.get());
        return user.get();
    }
}


