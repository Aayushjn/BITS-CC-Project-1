package com.springboot.TravelManagementSystem.services;

import com.springboot.TravelManagementSystem.exceptions.PackageException;
import com.springboot.TravelManagementSystem.models.CurrentUserLoginSession;
import com.springboot.TravelManagementSystem.models.Hotel;
import com.springboot.TravelManagementSystem.models.Packages;
import com.springboot.TravelManagementSystem.repository.HotelRepository;
import com.springboot.TravelManagementSystem.repository.PackageRepository;
import com.springboot.TravelManagementSystem.repository.SessionRepository;
import com.springboot.TravelManagementSystem.repository.UserRepository;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import java.util.List;
import java.util.Optional;

@Service
public class PackageServiceProvider implements PackageService{

//    @Autowired
//    HotelRepository hotelRepository;

    @Autowired
    PackageRepository packageRepository;

    @Autowired
    SessionRepository sessionRepository;

    @Autowired
    UserRepository userRepository;


    @Override
    public Packages createPackage(Packages pkg, String authKey) throws PackageException {

        Optional<CurrentUserLoginSession> culs=sessionRepository.findByAuthKey(authKey);
        if(!culs.isPresent())
        {
            throw new PackageException("Invalid Authentication key");
        }
        String userType=userRepository.findById(culs.get().getUserId()).get().getUserType();
        Optional<Packages> user=packageRepository.findById(culs.get().getUserId());
        if(userType.equalsIgnoreCase("user"))
        {
            throw new PackageException("Unauthorized Request");
        }
        else if(user.isPresent())
        {
            throw new PackageException("Package Already Exist with id " + pkg.getPackageId());
        }
        //hotel have to save the created package
        System.out.println(pkg);
//        List<Hotel> hotelList=pkg.getHotelDetails();
//        System.out.println(hotelList);
//        for(Hotel hotel :hotelList)
//        {
//            hotel.setPackages(pkg);
//        }
        Packages packagesCreate=packageRepository.save(pkg);

        return packagesCreate;
    }

    @Override
    public Packages updatePackage(Packages pkg, String authKey) throws PackageException {
        Optional<CurrentUserLoginSession> culs=sessionRepository.findByAuthKey(authKey);
        if(!culs.isPresent())
        {
            throw new PackageException("Invalid Authentication key");
        }
        String userType=userRepository.findById(culs.get().getUserId()).get().getUserType();
        Optional<Packages> getPkg=packageRepository.findById(pkg.getPackageId());
        if(userType.equalsIgnoreCase("user")) {
            throw new PackageException("Unauthorized Request");
        }
        else if(!getPkg.isPresent())
        {
            throw new PackageException("Package Not Present with id " + pkg.getPackageId());
        }
        Packages packagesUpdated=packageRepository.save(pkg);

        return packagesUpdated;
    }

    @Override
    public Packages deletePackage(Integer pkgId, String authKey) throws PackageException {
        Optional<CurrentUserLoginSession> culs=sessionRepository.findByAuthKey(authKey);
        if(!culs.isPresent())
        {
            throw new PackageException("Invalid Authentication key");
        }
        String userType=userRepository.findById(culs.get().getUserId()).get().getUserType();
        if(userType.equalsIgnoreCase("user"))
        {
            throw new PackageException("Unauthorized Request");
        }
        Optional<Packages> pkg=packageRepository.findById(pkgId);
        if(!pkg.isPresent())
        {
            throw new PackageException("Package doesnot exist with package id " + pkgId);
        }
        packageRepository.delete(pkg.get());
        return pkg.get();

    }

    @Override
    public List<Packages> getAllPackages() throws PackageException {
        List<Packages> pkgList=packageRepository.findAll();
        if(pkgList.isEmpty())
        {
            throw new PackageException("Package Not Avaliable");
        }
        return pkgList;
    }
}
