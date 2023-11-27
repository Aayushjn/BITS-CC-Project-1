package com.springboot.TravelManagementSystem.controller;

import com.springboot.TravelManagementSystem.exceptions.PackageException;
import com.springboot.TravelManagementSystem.models.Packages;
import com.springboot.TravelManagementSystem.services.PackageService;
import com.springboot.TravelManagementSystem.services.PackageServiceProvider;
import jakarta.validation.Valid;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.w3c.dom.stylesheets.LinkStyle;

import java.util.List;

@RestController
@RequestMapping("/packages")
@CrossOrigin(origins = "*")
public class PackageController {

    @Autowired
    PackageServiceProvider packageServiceProvider;

    @GetMapping("")
    public ResponseEntity<List<Packages>> getAllPackages() throws PackageException
    {
        List<Packages> pkgList=packageServiceProvider.getAllPackages();
        return new ResponseEntity<List<Packages>>(pkgList, HttpStatus.OK);
    }

    @PostMapping("")
    public ResponseEntity<Packages> addNewPackage(@Valid @RequestBody Packages packages,@RequestParam(value = "key") String authKey) throws PackageException
    {
        Packages createPackage=packageServiceProvider.createPackage(packages,authKey);
        return new ResponseEntity<Packages>(createPackage,HttpStatus.CREATED);
    }

    @PutMapping("")
    public ResponseEntity<Packages> updatePackage(@Valid @RequestBody Packages packages,@RequestParam (value ="key") String authKey) throws PackageException
    {
        Packages updatePackage= packageServiceProvider.updatePackage(packages,authKey);
        return new ResponseEntity<Packages>(updatePackage,HttpStatus.OK);
    }

    @DeleteMapping("")
    public ResponseEntity<Packages> deletePackage(@RequestParam(value = "PackageId") Integer pkgId , @RequestParam(value = "key") String authKey) throws PackageException
    {
        Packages pkg=packageServiceProvider.deletePackage(pkgId,authKey);
        return new ResponseEntity<Packages>(pkg,HttpStatus.OK);
    }

}
