CREATE DATABASE  IF NOT EXISTS `fibear` /*!40100 DEFAULT CHARACTER SET utf8 */;
USE `fibear`;
-- MySQL dump 10.13  Distrib 5.7.17, for macos10.12 (x86_64)
--
-- Host: 127.0.0.1    Database: fibear
-- ------------------------------------------------------
-- Server version	5.7.19

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `block_orders`
--

DROP TABLE IF EXISTS `block_orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `block_orders` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `user_block_date_id` int(11) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `block_orders`
--

LOCK TABLES `block_orders` WRITE;
/*!40000 ALTER TABLE `block_orders` DISABLE KEYS */;
INSERT INTO `block_orders` VALUES (1,2,1,'WAITING'),(2,2,2,'CONFIRM'),(4,1,4,'FINISH');
/*!40000 ALTER TABLE `block_orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `blocks`
--

DROP TABLE IF EXISTS `blocks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `blocks` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `description` varchar(3000) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  `hour_start` time DEFAULT NULL,
  `hour_end` time DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `blocks`
--

LOCK TABLES `blocks` WRITE;
/*!40000 ALTER TABLE `blocks` DISABLE KEYS */;
INSERT INTO `blocks` VALUES (1,'8h - 12h','Morning Block','08:00:00','12:00:00'),(2,'12h - 16h','Lunch Block','12:00:00','16:00:00'),(3,'16h - 20h','Afternoon Block','16:00:00','20:00:00'),(4,'20h - 24h ','Evening Block','20:00:00','24:00:00'),(5,'24h - 4h','Midnight Block','24:00:00','04:00:00');
/*!40000 ALTER TABLE `blocks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `profiles`
--

DROP TABLE IF EXISTS `profiles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `profiles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `avatar` varchar(200) NOT NULL,
  `firstname` varchar(200) NOT NULL,
  `lastname` varchar(200) NOT NULL,
  `sex` tinyint(1) DEFAULT NULL,
  `birthdate` date DEFAULT NULL,
  `province_id` int(11) NOT NULL,
  `star_rate` double DEFAULT NULL,
  `description` varchar(2000) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `province_id` (`province_id`),
  CONSTRAINT `profiles_ibfk_1` FOREIGN KEY (`id`) REFERENCES `users` (`id`),
  CONSTRAINT `profiles_ibfk_2` FOREIGN KEY (`province_id`) REFERENCES `provinces` (`province_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `profiles`
--

LOCK TABLES `profiles` WRITE;
/*!40000 ALTER TABLE `profiles` DISABLE KEYS */;
INSERT INTO `profiles` VALUES (1,1,'https://scontent.fsgn5-2.fna.fbcdn.net/v/t1.0-9/27867199_602928020051032_8822952925598529831_n.jpg?oh=4e3134a045551ef8120080834497ade7&oe=5B034804','tynk','shit',1,'2016-06-15',1,NULL,'Mình là Tynk , châm ngôn của mình là chiều các anh như chiều vợ hihi'),(2,2,'https://scontent.fsgn5-2.fna.fbcdn.net/v/t1.0-9/27867199_602928020051032_8822952925598529831_n.jpg?oh=4e3134a045551ef8120080834497ade7&oe=5B034804','huynk tynk','eat shit',1,'2015-06-11',1,NULL,'Minh la Tynk kute'),(3,3,'https://scontent.fsgn5-2.fna.fbcdn.net/v/t1.0-9/27867199_602928020051032_8822952925598529831_n.jpg?oh=4e3134a045551ef8120080834497ade7&oe=5B034804','le van tynk','ok',0,'1997-10-02',1,3,'Minh la Tynk , thich may anh dep trai khoai to'),(4,4,'https://scontent.fsgn5-2.fna.fbcdn.net/v/t1.0-9/27867199_602928020051032_8822952925598529831_n.jpg?oh=4e3134a045551ef8120080834497ade7&oe=5B034804','tynk van lee','aa',1,'1997-10-02',1,4,'abcccc');
/*!40000 ALTER TABLE `profiles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `provinces`
--

DROP TABLE IF EXISTS `provinces`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `provinces` (
  `province_id` int(11) NOT NULL AUTO_INCREMENT,
  `province_name` varchar(200) NOT NULL,
  PRIMARY KEY (`province_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `provinces`
--

LOCK TABLES `provinces` WRITE;
/*!40000 ALTER TABLE `provinces` DISABLE KEYS */;
INSERT INTO `provinces` VALUES (1,'Hồ Chí Minh');
/*!40000 ALTER TABLE `provinces` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `reviews`
--

DROP TABLE IF EXISTS `reviews`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `reviews` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_reviewed` int(11) DEFAULT NULL,
  `user_id` int(11) DEFAULT NULL,
  `rate` int(11) DEFAULT NULL,
  `description` varchar(2000) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `reviews`
--

LOCK TABLES `reviews` WRITE;
/*!40000 ALTER TABLE `reviews` DISABLE KEYS */;
INSERT INTO `reviews` VALUES (1,1,3,4,'Em này dáng đẹp , mỗi tội ốm quá xương sườn ngồi sau đâm đau'),(2,2,3,5,'Nhiệt tình , full service , LIKE');
/*!40000 ALTER TABLE `reviews` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `roles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES (1,'user'),(2,'bear'),(3,'moderator'),(4,'admin');
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_block_dates`
--

DROP TABLE IF EXISTS `user_block_dates`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_block_dates` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) DEFAULT NULL,
  `block_id` int(11) DEFAULT NULL,
  `block_date` date DEFAULT NULL,
  `description` varchar(3000) DEFAULT NULL,
  `status` varchar(255) DEFAULT NULL,
  `price` double DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_block_dates`
--

LOCK TABLES `user_block_dates` WRITE;
/*!40000 ALTER TABLE `user_block_dates` DISABLE KEYS */;
INSERT INTO `user_block_dates` VALUES (1,3,4,'2018-03-18','Đi chơi hun nhau ôm iếc ','FREE',32000),(2,3,5,'2018-03-18','Chỉ được ôm thôi ','BUSY',12000),(3,3,2,'2018-03-18','abc','FREE',30000),(4,3,1,'2018-03-18','abc','BUSY',30000);
/*!40000 ALTER TABLE `user_block_dates` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_history_orders`
--

DROP TABLE IF EXISTS `user_history_orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_history_orders` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `block_order_id` int(11) DEFAULT NULL,
  `review_id` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_history_orders`
--

LOCK TABLES `user_history_orders` WRITE;
/*!40000 ALTER TABLE `user_history_orders` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_history_orders` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(200) NOT NULL,
  `password_hash` varchar(200) DEFAULT NULL,
  `salt` varchar(200) NOT NULL,
  `secret_key` varchar(200) NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `phone` varchar(32) NOT NULL,
  `email` varchar(200) NOT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `role_id` int(11) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'abc123','$2a$14$cve8Iy2uDpuDRP1exwf6t.NnGctJ6wK.Cuw3krL7lAQc/V8fK9L6u','1','1','2018-03-09 20:06:13','2018-03-09 20:06:13','01216339739','abc@tynk.com',NULL,1),(2,'abc456','$2a$14$cve8Iy2uDpuDRP1exwf6t.NnGctJ6wK.Cuw3krL7lAQc/V8fK9L6u','1','1','2018-03-12 16:40:48','2018-03-12 16:40:48','09113113','tynkancut@tynk.com',NULL,1),(3,'abc789','$2a$14$cve8Iy2uDpuDRP1exwf6t.NnGctJ6wK.Cuw3krL7lAQc/V8fK9L6u','1','1','2018-03-14 13:16:30','2018-03-14 13:16:30','012162222222','tynkcutelover@tynk.com',NULL,2),(4,'abcxyz','$2a$14$cve8Iy2uDpuDRP1exwf6t.NnGctJ6wK.Cuw3krL7lAQc/V8fK9L6u','1','1','2018-03-25 12:06:15','2018-03-25 12:06:15','11111111','abc@tynk.com',NULL,3);
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-07-29  0:46:52
