package com.springboot.TravelManagementSystem.models;

import jakarta.persistence.*;
import jakarta.validation.constraints.Email;
import lombok.Data;

import java.time.LocalDateTime;

@Data
@Entity @Table(name="User_Session")
public class CurrentUserLoginSession {
    @Id @GeneratedValue(strategy = GenerationType.AUTO)
    private Integer id;
    private Integer userId;
    private String authKey;
    private LocalDateTime sessionStartTime;
}
