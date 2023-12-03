package com.springboot.TravelManagementSystem.models;

import lombok.Data;
import org.springframework.format.annotation.DateTimeFormat;

import java.time.LocalDateTime;

@Data
public class SessionDTO {
    private String authKey;
    private LocalDateTime sessionStartTime;
}
