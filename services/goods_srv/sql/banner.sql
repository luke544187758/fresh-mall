/*
 Navicat Premium Data Transfer

 Source Server         : mysql8
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : 127.0.0.1:13306
 Source Schema         : mall-goods-srv

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 12/05/2022 03:47:01
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for banner
-- ----------------------------
DROP TABLE IF EXISTS `banner`;
CREATE TABLE `banner` (
  `id` int NOT NULL AUTO_INCREMENT,
  `add_time` datetime NOT NULL,
  `is_deleted` tinyint(1) NOT NULL,
  `update_time` datetime NOT NULL,
  `image` varchar(200) COLLATE utf8mb4_general_ci NOT NULL,
  `url` varchar(200) COLLATE utf8mb4_general_ci NOT NULL,
  `index` int NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- ----------------------------
-- Records of banner
-- ----------------------------
BEGIN;
INSERT INTO `banner` (`id`, `add_time`, `is_deleted`, `update_time`, `image`, `url`, `index`) VALUES (1, '2022-05-09 13:58:20', 0, '2022-05-09 13:58:20', 'https://m.360buyimg.com/babel/jfs/t1/192589/37/23892/129366/6268b747E59790a75/a43836d0c1a67b8d.png', '', 1);
INSERT INTO `banner` (`id`, `add_time`, `is_deleted`, `update_time`, `image`, `url`, `index`) VALUES (2, '2022-05-09 13:58:20', 0, '2022-05-09 13:58:20', 'https://m.360buyimg.com/babel/jfs/t1/94895/22/17698/97464/626765a0Ed99bbeed/b649f5b79fdc1591.png', '', 2);
INSERT INTO `banner` (`id`, `add_time`, `is_deleted`, `update_time`, `image`, `url`, `index`) VALUES (3, '2022-05-09 13:58:20', 0, '2022-05-09 13:58:20', 'https://m.360buyimg.com/babel/jfs/t1/67819/18/17433/105519/626a7226E2e94f956/a926a540ad891a03.png', '', 3);
INSERT INTO `banner` (`id`, `add_time`, `is_deleted`, `update_time`, `image`, `url`, `index`) VALUES (4, '2022-05-09 13:58:20', 0, '2022-05-09 13:58:20', 'https://m.360buyimg.com/babel/jfs/t1/91139/5/27610/70951/626d404bE3431e4a8/3a9033ddb4ed94e4.jpg', '', 4);
INSERT INTO `banner` (`id`, `add_time`, `is_deleted`, `update_time`, `image`, `url`, `index`) VALUES (5, '2022-05-09 13:58:20', 0, '2022-05-09 13:58:20', 'https://m.360buyimg.com/babel/jfs/t1/61121/31/17593/90766/626c08cdEe7921f15/8494669c17dd900f.png', '', 5);
INSERT INTO `banner` (`id`, `add_time`, `is_deleted`, `update_time`, `image`, `url`, `index`) VALUES (6, '2022-05-09 13:58:20', 0, '2022-05-09 13:58:20', 'https://m.360buyimg.com/babel/jfs/t1/189857/7/24347/123508/626b5422E2235f20d/f8e1fb5c2fb70695.jpg', '', 6);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
