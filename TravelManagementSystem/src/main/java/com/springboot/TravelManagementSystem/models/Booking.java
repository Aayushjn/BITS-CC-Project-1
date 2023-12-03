package com.springboot.TravelManagementSystem.models;

import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import lombok.Data;

import java.time.LocalDateTime;

@Data
@Entity @Table(name="Booking")
public class Booking {

    @Id @GeneratedValue(strategy = GenerationType.AUTO)
    private Integer bookingId;

    @NotBlank @NotBlank @NotEmpty
    private String description;

    @Enumerated(EnumType.STRING)
    private BookingType bookingType=BookingType.National;

    @Enumerated(EnumType.STRING)
    private BookingStatus status=BookingStatus.Not_Booked;

    @NotBlank @NotBlank @NotEmpty
    private String bookingTitle;
    //private LocalDateTime date;

    @JsonIgnore
    @ManyToOne(cascade = CascadeType.ALL)
    @JoinColumn(name = "user_id",referencedColumnName = "userId")
    private User user;

}
