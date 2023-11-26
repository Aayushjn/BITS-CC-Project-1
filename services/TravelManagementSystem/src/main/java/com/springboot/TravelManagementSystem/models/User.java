package com.springboot.TravelManagementSystem.models;

import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import jakarta.validation.constraints.*;
import lombok.Data;
import jdk.jfr.DataAmount;

import java.util.ArrayList;
import java.util.List;

@Data
@Entity
public class User {

     @Id @GeneratedValue(strategy = GenerationType.AUTO)
     private Integer userId;

     @NotBlank @NotBlank @NotEmpty @Size(min=3,max=20,message = "Name must contain atleast 3 character")
     private String name;

     @NotNull @Pattern(regexp = "^[789][0-9]{9}",message = "Mobile No. should be 10 digits")
     private String mobile;

     @NotBlank @NotBlank @NotEmpty
     private String address;

     @JsonIgnore
     private String userType="User";

     @Email(message = "Invalid Email address.")@Column(unique = true)
     private String email;

     @Pattern(regexp = "[A-Za-z0-9@]{6,15}",message = "Password must be 6 to 15 characters and must have atleast 1 alphabet and 1 number")
     @NotNull @NotBlank @NotEmpty
     private String password;

    @JsonIgnore
    @OneToMany(cascade = CascadeType.ALL,mappedBy = "user")
    private List<Booking> bookings=new ArrayList<>();



}
