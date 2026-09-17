CREATE DATABASE  IF NOT EXISTS `zero_order` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `zero_order`;
-- MySQL dump 10.13  Distrib 8.0.41, for Win64 (x86_64)
--
-- Host: 127.0.0.1    Database: zero_order
-- ------------------------------------------------------
-- Server version	9.2.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `order_item`
--

DROP TABLE IF EXISTS `order_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `order_item` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '明细表自身主键',
  `order_no` bigint NOT NULL COMMENT '关联订单号',
  `goods_id` varchar(64) NOT NULL COMMENT '商品ID（来自goods服务）',
  `goods_name` varchar(128) NOT NULL COMMENT '下单时商品名称快照',
  `goods_cover` varchar(256) NOT NULL COMMENT '下单时商品封面快照',
  `price_cent` bigint NOT NULL DEFAULT '0' COMMENT '下单时单件价格(分)',
  `num` int NOT NULL COMMENT '购买数量',
  `total_cent` bigint NOT NULL DEFAULT '0' COMMENT '本条商品小计金额(分)',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_order_no` (`order_no`)
) ENGINE=InnoDB AUTO_INCREMENT=27 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单明细表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `order_logistics`
--

DROP TABLE IF EXISTS `order_logistics`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `order_logistics` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '物流表自身主键',
  `order_id` bigint NOT NULL COMMENT '关联订单主键id',
  `receiver_name` varchar(32) NOT NULL COMMENT '收件人姓名快照',
  `receiver_phone` varchar(32) NOT NULL COMMENT '收件手机号快照',
  `receiver_address` varchar(256) NOT NULL COMMENT '完整收货地址快照',
  `logistics_name` varchar(32) DEFAULT NULL COMMENT '快递公司名称',
  `tracking_no` varchar(64) DEFAULT NULL COMMENT '快递运单号',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_order_no` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单物流表';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `orders`
--

DROP TABLE IF EXISTS `orders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `orders` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '自增主键id',
  `order_no` bigint NOT NULL COMMENT '雪花id',
  `user_id` varchar(64) NOT NULL COMMENT '下单用户ID（来自user服务）',
  `total_cent` bigint NOT NULL COMMENT '订单原价',
  `pay_cent` bigint NOT NULL DEFAULT '0' COMMENT '实付金额(分，优惠后)',
  `status` tinyint NOT NULL DEFAULT '0' COMMENT '0待支付 1已支付 2已发货 3已完成 4已取消 5已退款',
  `pay_type` tinyint DEFAULT NULL COMMENT '支付方式 1余额 2微信 3支付宝',
  `pay_time` datetime DEFAULT NULL COMMENT '支付时间',
  `deliver_time` datetime DEFAULT NULL COMMENT '发货时间',
  `receive_time` datetime DEFAULT NULL COMMENT '确认收货时间',
  `cancel_time` datetime DEFAULT NULL COMMENT '订单取消时间',
  `remark` varchar(256) DEFAULT NULL COMMENT '用户下单备注',
  `expire_time` datetime NOT NULL COMMENT '订单超时关闭时间，一般30分钟',
  `deleted_at` datetime DEFAULT NULL COMMENT '软删除',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_orderId` (`order_no`),
  KEY `idx_userId_create` (`user_id`,`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='订单主表';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-09-17 19:55:14
