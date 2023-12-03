package com.springboot.TravelManagementSystem.repository;

import com.springboot.TravelManagementSystem.exceptions.InvalidCredentialException;
import com.springboot.TravelManagementSystem.models.CurrentUserLoginSession;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface SessionRepository extends JpaRepository<CurrentUserLoginSession,Integer> {
    public Optional<CurrentUserLoginSession> findByUserId(Integer userId) throws InvalidCredentialException;

    @Query("select c from CurrentUserLoginSession c where c.authKey=?1")
    public Optional<CurrentUserLoginSession> findByAuthKey(String key);
}
