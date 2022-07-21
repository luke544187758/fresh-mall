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

 Date: 21/07/2022 21:28:14
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for shoppingcart
-- ----------------------------
DROP TABLE IF EXISTS `shoppingcart`;
CREATE TABLE `shoppingcart` (
  `id` bigint NOT NULL,
  `add_time` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  `update_time` datetime NOT NULL,
  `user` bigint NOT NULL,
  `goods` bigint NOT NULL,
  `nums` int NOT NULL,
  `checked` tinyint(1) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

SET FOREIGN_KEY_CHECKS = 1;
