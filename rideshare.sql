-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1:3308
-- Generation Time: Dec 17, 2021 at 09:12 AM
-- Server version: 10.4.21-MariaDB
-- PHP Version: 7.3.31

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `rideshare`
--

-- --------------------------------------------------------

--
-- Table structure for table `customer`
--

CREATE TABLE `customer` (
  `ID` int(5) NOT NULL,
  `FirstName` varchar(255) NOT NULL,
  `LastName` varchar(255) NOT NULL,
  `MobileNumber` varchar(8) NOT NULL,
  `EmailAddress` varchar(255) NOT NULL,
  `Password` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `customer`
--

INSERT INTO `customer` (`ID`, `FirstName`, `LastName`, `MobileNumber`, `EmailAddress`, `Password`) VALUES
(1, 'John', 'Sussies', '91234530', 'Susan@np.com', 'password'),
(2, 'Willy', 'Wiki', '91234569', 'Willy@np.com', 'password'),
(3, 'Johsnon', 'Sussy', '91234530', 'Johnson@np.com', 'password');

-- --------------------------------------------------------

--
-- Table structure for table `driver`
--

CREATE TABLE `driver` (
  `DriverID` varchar(255) NOT NULL,
  `FirstName` varchar(255) NOT NULL,
  `LastName` varchar(255) NOT NULL,
  `MobileNo` varchar(8) NOT NULL,
  `EmailAddress` varchar(255) NOT NULL,
  `Password` varchar(255) NOT NULL,
  `CarLiscenseNo` varchar(8) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `driver`
--

INSERT INTO `driver` (`DriverID`, `FirstName`, `LastName`, `MobileNo`, `EmailAddress`, `Password`, `CarLiscenseNo`) VALUES
('S1234567A', 'Ethan', 'Wiki', '91234567', 'Ethan@np.com', 'password', 'S123LO'),
('S1234567B', 'Elvan', 'Wiki', '91234570', 'Elvan@np.com', 'password', 'S124LO'),
('S1234567C', 'Azzi', 'Wiki', '91234570', 'Azzi@np.com', 'password', 'S125LO');

-- --------------------------------------------------------

--
-- Table structure for table `trip`
--

CREATE TABLE `trip` (
  `TripID` int(11) NOT NULL,
  `DriverID` varchar(255) NOT NULL,
  `CustomerID` int(11) NOT NULL,
  `PickUpLocation` varchar(255) NOT NULL,
  `DropOffLocation` varchar(255) NOT NULL,
  `PickUpTime` varchar(255) NOT NULL,
  `DropOffTime` varchar(255) NOT NULL,
  `Status` varchar(255) NOT NULL CHECK (`Status` in ('Pending','On The Way','In Transit','Completed','Failed'))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

--
-- Dumping data for table `trip`
--

INSERT INTO `trip` (`TripID`, `DriverID`, `CustomerID`, `PickUpLocation`, `DropOffLocation`, `PickUpTime`, `DropOffTime`, `Status`) VALUES
(1, 'S1234567A', 1, '123456', '123457', '12:30', '13:00', 'Completed'),
(2, 'S1234567A', 2, '123456', '123457', '8:35PM', '8:45PM', 'Completed'),
(3, 'S1234567B', 1, '1234567', '1234568', '22:00', ' ', 'Pending');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `customer`
--
ALTER TABLE `customer`
  ADD PRIMARY KEY (`ID`);

--
-- Indexes for table `driver`
--
ALTER TABLE `driver`
  ADD PRIMARY KEY (`DriverID`),
  ADD UNIQUE KEY `DriverID` (`DriverID`);

--
-- Indexes for table `trip`
--
ALTER TABLE `trip`
  ADD PRIMARY KEY (`TripID`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `customer`
--
ALTER TABLE `customer`
  MODIFY `ID` int(5) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;

--
-- AUTO_INCREMENT for table `trip`
--
ALTER TABLE `trip`
  MODIFY `TripID` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
