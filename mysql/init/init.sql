SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

DROP TABLE IF EXISTS `Book`;
CREATE TABLE `Book` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Title` varchar(256) NOT NULL,
  `ISBN` varchar(256) NOT NULL,
  `Author` varchar(256) NOT NULL,
  PRIMARY KEY (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `LendHistory`;
CREATE TABLE `LendHistory` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `LendId` int(11) NOT NULL,
  `Status` varchar(256) NOT NULL,
  `Date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' ON UPDATE current_timestamp(),
  PRIMARY KEY (`Id`),
  KEY `LendId` (`LendId`),
  CONSTRAINT `LendHistory_ibfk_1` FOREIGN KEY (`LendId`) REFERENCES `UserLendBook` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `Team`;
CREATE TABLE `Team` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Name` varchar(256) NOT NULL,
  `Owner` int(11) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `Owner` (`Owner`),
  CONSTRAINT `Team_ibfk_1` FOREIGN KEY (`Owner`) REFERENCES `User` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `User`;
CREATE TABLE `User` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `Email` varchar(256) NOT NULL,
  `Name` varchar(256) NOT NULL,
  `Password` varchar(512) NOT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE KEY `Email` (`Email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `UserBook`;
CREATE TABLE `UserBook` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UserId` int(11) NOT NULL,
  `BookId` int(11) NOT NULL,
  `State` varchar(256) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `UserId` (`UserId`),
  KEY `BookId` (`BookId`),
  CONSTRAINT `UserBook_ibfk_1` FOREIGN KEY (`UserId`) REFERENCES `User` (`Id`),
  CONSTRAINT `UserBook_ibfk_2` FOREIGN KEY (`BookId`) REFERENCES `Book` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `UserLendBook`;
CREATE TABLE `UserLendBook` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `OwnerId` int(11) NOT NULL,
  `BorrowedId` int(11) NOT NULL,
  `BookId` int(11) NOT NULL,
  `Status` varchar(256) NOT NULL,
  PRIMARY KEY (`Id`),
  UNIQUE KEY `OwnerId_BorrowedId_BookId` (`OwnerId`,`BorrowedId`,`BookId`),
  KEY `BorrowedId` (`BorrowedId`),
  KEY `BookId` (`BookId`),
  CONSTRAINT `UserLendBook_ibfk_1` FOREIGN KEY (`OwnerId`) REFERENCES `User` (`Id`),
  CONSTRAINT `UserLendBook_ibfk_2` FOREIGN KEY (`BorrowedId`) REFERENCES `User` (`Id`),
  CONSTRAINT `UserLendBook_ibfk_3` FOREIGN KEY (`BookId`) REFERENCES `Book` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


DROP TABLE IF EXISTS `UserTeam`;
CREATE TABLE `UserTeam` (
  `Id` int(11) NOT NULL AUTO_INCREMENT,
  `UserId` int(11) NOT NULL,
  `TeamId` int(11) NOT NULL,
  PRIMARY KEY (`Id`),
  KEY `UserId` (`UserId`),
  KEY `TeamId` (`TeamId`),
  CONSTRAINT `UserTeam_ibfk_1` FOREIGN KEY (`UserId`) REFERENCES `User` (`Id`),
  CONSTRAINT `UserTeam_ibfk_2` FOREIGN KEY (`TeamId`) REFERENCES `Team` (`Id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
