CREATE TABLE `travel_management`.`user` (
  `user_id` INT NOT NULL,
  `address` VARCHAR(255) NULL DEFAULT NULL,
  `email` VARCHAR(255) NULL DEFAULT NULL,
  `mobile` VARCHAR(255) NOT NULL,
  `name` VARCHAR(25) NULL DEFAULT NULL,
  `password` VARCHAR(255) NOT NULL,
  `user_type` VARCHAR(255) NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE INDEX `email_UNIQUE` (`email` ASC) VISIBLE);

CREATE TABLE `travel_management`.`user_session` (
  `id` INT NOT NULL,
  `auth_key` VARCHAR(45) NULL DEFAULT NULL,
  `session_start_time` DATETIME(6) NULL DEFAULT NULL,
  `user_id` INT NULL DEFAULT NULL,
  PRIMARY KEY (`id`));

CREATE TABLE `travel_management`.`booking` (
  `booking_id` INT NOT NULL,
  `booking_title` VARCHAR(255) NULL DEFAULT NULL,
  `booking_type` VARCHAR(255) NULL DEFAULT NULL,
  `address` VARCHAR(255) NULL DEFAULT NULL,
  `status` VARCHAR(255) NULL DEFAULT NULL,
  `user_id` INT NULL DEFAULT NULL,
  PRIMARY KEY (`booking_id`),
  INDEX `fk_user_idx` (`user_id` ASC) VISIBLE,
  CONSTRAINT `fk_user`
    FOREIGN KEY (`user_id`)
    REFERENCES `travel_management`.`user` (`user_id`)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT);

CREATE TABLE `travel_management`.`packages` (
  `package_id` INT NOT NULL,
  `package_cost` INT NOT NULL,
  `package_description` VARCHAR(255) NOT NULL,
  `package_name` VARCHAR(255) NOT NULL,
  `package_type` TINYINT NULL DEFAULT NULL,
  `packages_details` VARCHAR(255) NOT NULL,
  PRIMARY KEY (`package_id`));

CREATE TABLE `travel_management`.`hotel` (
  `hotel_id` INT NOT NULL,
  `address` VARCHAR(255) NULL DEFAULT NULL,
  `hotel_description` VARCHAR(255) NULL DEFAULT NULL,
  `hotel_name` VARCHAR(255) NULL DEFAULT NULL,
  `hotel_type` VARCHAR(255) NULL DEFAULT NULL,
  `rent` INT NULL DEFAULT NULL,
  `status` VARCHAR(255) NULL DEFAULT NULL,
  `packages_package_id` INT NULL DEFAULT NULL,
  PRIMARY KEY (`hotel_id`),
  INDEX `fk_package_idx` (`packages_package_id` ASC) VISIBLE,
  CONSTRAINT `fk_package`
    FOREIGN KEY (`packages_package_id`)
    REFERENCES `travel_management`.`packages` (`package_id`)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT);

CREATE TABLE `travel_management`.`ticket` (
  `ticket_id` INT NOT NULL,
  `ticket_description` VARCHAR(255) NOT NULL,
  `ticket_status` BIT(1) NULL DEFAULT NULL,
  `package_id` INT NULL DEFAULT NULL,
  PRIMARY KEY (`ticket_id`),
  INDEX `fk_package_idx` (`package_id` ASC) VISIBLE,
  CONSTRAINT `fk1_package`
    FOREIGN KEY (`package_id`)
    REFERENCES `travel_management`.`packages` (`package_id`)
    ON DELETE RESTRICT
    ON UPDATE RESTRICT);



