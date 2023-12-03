package com.springboot.TravelManagementSystem.models;


import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import jakarta.validation.constraints.*;
import lombok.Data;

import java.util.List;

@Data
@Entity @Table(name = "Packages")
public class Packages {

    @Id @GeneratedValue(strategy = GenerationType.AUTO)
    private Integer packageId;

    @Size(min=3, message = "Package name shouldn't be less than 3 characters") @NotNull @NotBlank @NotEmpty
    private String packageName;

    @Size(min = 5,message = "Description shouldn't be less than 5 character ") @NotNull @NotBlank @NotEmpty
    private String packageDescription;

    @Enumerated(EnumType.ORDINAL)
    private PackageType packageType;

    @NotNull @Min(0)
    private Integer packageCost;

    @NotNull @NotBlank @NotEmpty @Size(min=3,message = "Payment Details shouldnot be less than 3 characters")
    private String paymentDetails;


    //ticket one
    @JsonIgnore
    @OneToMany(mappedBy = "packages",cascade = CascadeType.ALL)
    private List<Ticket> ticketDetails;

    //hotel one
    @OneToMany(mappedBy = "packages",cascade = CascadeType.ALL)
    private List<Hotel> hotelDetails;

}
