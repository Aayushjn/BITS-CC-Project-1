package com.springboot.TravelManagementSystem.services;

import com.springboot.TravelManagementSystem.exceptions.AlreadyExistsException;
import com.springboot.TravelManagementSystem.exceptions.InvalidCredentialException;
import com.springboot.TravelManagementSystem.models.SessionDTO;
import com.springboot.TravelManagementSystem.models.User;
import com.springboot.TravelManagementSystem.models.UserDTO;

public interface UserService {
    public User userSignup(User user) throws AlreadyExistsException;
    public SessionDTO userLogin(UserDTO user) throws InvalidCredentialException;

    public String userLogout(String authKey) throws InvalidCredentialException;

    public boolean updateUser(User user) throws InvalidCredentialException;

    public User getUser(String authKey) throws InvalidCredentialException;


    //admin use only

    public User deleteUser(Integer userId ,String authKey) throws InvalidCredentialException;

    public User userAdmin(String userEmail,String passcode) throws InvalidCredentialException;

}
