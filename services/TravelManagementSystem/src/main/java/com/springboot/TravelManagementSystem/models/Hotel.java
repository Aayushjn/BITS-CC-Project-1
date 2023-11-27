package com.springboot.TravelManagementSystem.models;


import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Entity
@NoArgsConstructor
public class Hotel {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Integer hotelId;

    @NotBlank @NotBlank @NotEmpty
    private String hotelName;

    @Enumerated(EnumType.STRING)
    private HotelType hotelType=HotelType.AC;

    @NotBlank @NotBlank @NotEmpty
    private String hotelDescription;

    @NotBlank @NotBlank @NotEmpty
    private String address;

    private Integer rent;

    @Enumerated(EnumType.STRING)
    private HotelStatus status=HotelStatus.NOT_BOOKED;

    @JsonIgnore
    @ManyToOne(cascade = CascadeType.ALL)
    private Packages packages;

}
