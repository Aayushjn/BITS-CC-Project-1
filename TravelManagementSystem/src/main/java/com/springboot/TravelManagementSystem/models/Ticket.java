package com.springboot.TravelManagementSystem.models;


import com.fasterxml.jackson.annotation.JsonIgnore;
import jakarta.persistence.*;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotEmpty;
import jakarta.validation.constraints.NotNull;
import lombok.Data;
import lombok.NoArgsConstructor;

@Data
@Entity
@NoArgsConstructor
@Table(name = "Ticket")
public class Ticket {

    @Id
    @GeneratedValue(strategy = GenerationType.AUTO)
    private Integer ticketId;

    @NotNull @NotEmpty @NotBlank
    private String ticketDescription;

    private boolean ticketStatus;

    @JsonIgnore
    @ManyToOne(cascade = CascadeType.ALL)
    @JoinColumn(name = "packageId")
    private Packages packages;

}
