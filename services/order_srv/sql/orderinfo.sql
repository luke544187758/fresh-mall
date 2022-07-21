/*
 Navicat Premium Data Transfer

 Source Server         : mysql8
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : 127.0.0.1:13306
 Source Schema         : mall-order-srv

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 21/07/2022 21:28:03
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for orderinfo
-- ----------------------------
DROP TABLE IF EXISTS `orderinfo`;
CREATE TABLE `orderinfo` (
  `id` bigint NOT NULL,
  `add_time` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  `update_time` datetime NOT NULL,
  `user` bigint NOT NULL,
  `order_sn` varchar(30) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `pay_type` varchar(30) COLLATE utf8mb4_general_ci NOT NULL,
  `status` varchar(30) COLLATE utf8mb4_general_ci NOT NULL,
  `trade_no` varchar(100) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `order_mount` float NOT NULL,
  `pay_time` datetime DEFAULT NULL,
  `address` varchar(100) COLLATE utf8mb4_general_ci NOT NULL,
  `signer_name` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
  `signer_mobile` varchar(11) COLLATE utf8mb4_general_ci NOT NULL,
  `remark` varchar(200) COLLATE utf8mb4_general_ci NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `orderinfo_order_sn` (`order_sn`),
  UNIQUE KEY `orderinfo_trade_no` (`trade_no`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
