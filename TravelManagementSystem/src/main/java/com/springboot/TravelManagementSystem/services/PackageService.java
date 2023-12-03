package com.springboot.TravelManagementSystem.services;

import com.springboot.TravelManagementSystem.exceptions.PackageException;
import com.springboot.TravelManagementSystem.models.Packages;

import java.util.List;

public interface PackageService {

    //admin functionality

    public Packages createPackage(Packages pkg, String authKey) throws PackageException;

    // admin functionality
    public Packages updatePackage(Packages pkg,String authKey) throws PackageException;

    //admin functionality
    public Packages deletePackage(Integer pkgId,String authKey) throws PackageException;

    public List<Packages> getAllPackages() throws PackageException;

}
