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

 Date: 21/07/2022 21:27:50
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for ordergoods
-- ----------------------------
DROP TABLE IF EXISTS `ordergoods`;
CREATE TABLE `ordergoods` (
  `id` bigint NOT NULL,
  `add_time` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  `update_time` datetime NOT NULL,
  `order` bigint NOT NULL,
  `goods` bigint NOT NULL,
  `goods_name` varchar(20) COLLATE utf8mb4_general_ci NOT NULL,
  `goods_image` varchar(200) COLLATE utf8mb4_general_ci NOT NULL,
  `goods_price` decimal(10,5) NOT NULL,
  `nums` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
