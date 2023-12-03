package com.springboot.TravelManagementSystem.controller;

import com.springboot.TravelManagementSystem.exceptions.AlreadyExistsException;
import com.springboot.TravelManagementSystem.exceptions.InvalidCredentialException;
import com.springboot.TravelManagementSystem.models.SessionDTO;
import com.springboot.TravelManagementSystem.models.User;
import com.springboot.TravelManagementSystem.models.UserDTO;
import com.springboot.TravelManagementSystem.services.UserService;
import jakarta.persistence.criteria.CriteriaBuilder;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/users")
@CrossOrigin(origins = "*")
public class UserController {

    @Autowired
    UserService service;

    @PostMapping("/signup")
    public ResponseEntity<String> userSignup(@Valid @RequestBody User user) throws AlreadyExistsException
    {
        service.userSignup(user);
        return new ResponseEntity<String>("Registered Successfully", HttpStatus.OK);
    }

    @PostMapping("/login")
    public ResponseEntity<String> userSignIn(@Valid @RequestBody UserDTO user) throws AlreadyExistsException , InvalidCredentialException
    {
        service.userLogin(user);
        return new ResponseEntity<String>("Log In Successfully",HttpStatus.OK);
    }

    @PostMapping("/logout")
    public ResponseEntity<String> userLogout(@RequestParam(value = "key") String authKey) throws InvalidCredentialException
    {
        service.userLogout(authKey);
        return new ResponseEntity<String>("Logged out successfully !!",HttpStatus.OK);
    }

    @PutMapping("/profile")
    public ResponseEntity<String> updateUser(@Valid @RequestBody User user) throws InvalidCredentialException
    {
        service.updateUser(user);
        return new ResponseEntity<String>("User Updated Successfully",HttpStatus.ACCEPTED);
    }

    @GetMapping("/")
    public ResponseEntity<Map<String,Integer>> getUser(@RequestParam String key) throws InvalidCredentialException
    {
        User user=null;
        Map<String,Integer> map=new HashMap<String,Integer>();
        user=service.getUser(key);
        map.put("User ",user.getUserId());
        //map.put("")
        return new ResponseEntity<Map<String, Integer>>(map,HttpStatus.OK);
    }

    @DeleteMapping("/")
    public ResponseEntity<Void> deleteUser(@RequestParam Integer userId,@RequestParam String authKey) throws InvalidCredentialException
    {
        this.service.deleteUser(userId,authKey);
        return new ResponseEntity<>(HttpStatus.OK);
    }

    @PostMapping("/appoint")
    public ResponseEntity<String> appointAdmin(@RequestParam("email") String email,@RequestParam("code") String passcode) throws InvalidCredentialException
    {
       service.userAdmin(email,passcode);
       return new ResponseEntity<String>("Admin Appointed",HttpStatus.OK);
    }

}


