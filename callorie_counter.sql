-- MySQL dump 10.13  Distrib 8.0.28, for Linux (x86_64)
--
-- Host: 127.0.0.1    Database: callorie_counter
-- ------------------------------------------------------
-- Server version	8.0.28-0ubuntu0.20.04.3

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `food`
--

DROP TABLE IF EXISTS `food`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `food` (
  `id` int NOT NULL AUTO_INCREMENT,
  `product_name` varchar(100) DEFAULT NULL,
  `callories` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=37 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `food`
--

LOCK TABLES `food` WRITE;
/*!40000 ALTER TABLE `food` DISABLE KEYS */;
INSERT INTO `food` VALUES (1,'яичница жаренная',152),(2,'чай',2),(3,'хлебец',14),(4,'яичница жаренная с наггетсами',248),(5,'яичница жаренная с наггетсами и сыром',280),(6,'макароны с наггетсами',274),(7,'гречневые хлопья с наггетсами',180),(8,'гречка с наггетсами',247),(9,'наггетсы 3шт',144),(10,'наггетсы 2шт',96),(11,'кусок пиццы',200),(12,'рис с наггетсами',189),(13,'паста',125),(14,'макароны со свининой',185),(15,'гречка',37),(16,'гречневые хлопья',29),(17,'рис',39),(18,'чебупицца',125),(20,'шоколадка горькая',27),(21,'печенька на фруктозе',50),(22,'пол кусочка пиццы',100),(23,'булочка с корицей',233),(24,'конфета с нугой и карамелью',51),(25,'пончик с сахарной пудрой',255),(26,'сарделька',220),(28,'молочный шоколад',45),(29,'отчет за денб',0),(30,'куриные катлеты',150),(31,'конфета',50),(33,'шоколад горький',27),(34,'макароны',40),(35,'мюсли батончик',100),(36,'молоко',3);
/*!40000 ALTER TABLE `food` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lunchs`
--

DROP TABLE IF EXISTS `lunchs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lunchs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` varchar(100) NOT NULL,
  `date` datetime DEFAULT NULL,
  `callories` int DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lunchs`
--

LOCK TABLES `lunchs` WRITE;
/*!40000 ALTER TABLE `lunchs` DISABLE KEYS */;
INSERT INTO `lunchs` VALUES (5,'1459568357','2022-02-14 17:30:48',16),(6,'1459568357','2022-02-14 20:19:15',940),(7,'1459568357','2022-02-14 20:23:16',51),(8,'1459568357','2022-02-15 06:29:25',248),(9,'1459568357','2022-02-15 10:42:29',286),(10,'1459568357','2022-02-15 15:21:14',413),(11,'1459568357','2022-02-15 21:51:38',350),(12,'1459568357','2022-02-16 11:38:54',552),(13,'1459568357','2022-02-16 18:37:16',352),(14,'1459568357','2022-02-16 20:13:38',326),(15,'1459568357','2022-02-17 08:38:47',359),(17,'1459568357','2022-02-18 17:34:16',409),(18,'1459568357','2022-02-18 17:37:23',192),(19,'1459568357','2022-02-18 17:38:15',100),(20,'1568063222','2022-02-19 00:40:56',16),(21,'1568063222','2022-02-19 00:52:20',5);
/*!40000 ALTER TABLE `lunchs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users_status`
--

DROP TABLE IF EXISTS `users_status`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users_status` (
  `user_id` varchar(100) NOT NULL,
  `waiting` varchar(100) NOT NULL,
  `product_name` varchar(100) DEFAULT NULL,
  `callories` varchar(100) DEFAULT NULL,
  `current_callories` int DEFAULT NULL,
  PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users_status`
--

LOCK TABLES `users_status` WRITE;
/*!40000 ALTER TABLE `users_status` DISABLE KEYS */;
INSERT INTO `users_status` VALUES ('1459568357','no_waiting','','',0),('1568063222','no_waiting','','',0);
/*!40000 ALTER TABLE `users_status` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-02-21  8:30:02
