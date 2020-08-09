/*
 Navicat Premium Data Transfer

 Source Server         : Ali-sz
 Source Server Type    : MySQL
 Source Server Version : 50730
 Source Host           : 120.78.138.140:3306
 Source Schema         : tpayment

 Target Server Type    : MySQL
 Target Server Version : 50730
 File Encoding         : 65001

 Date: 09/08/2020 17:35:16
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for acquirer
-- ----------------------------
DROP TABLE IF EXISTS `acquirer`;
CREATE TABLE `acquirer` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `name` varchar(32) DEFAULT NULL,
  `addition` varchar(255) DEFAULT NULL,
  `config_file_url` varchar(512) DEFAULT NULL,
  `agency_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for agency
-- ----------------------------
DROP TABLE IF EXISTS `agency`;
CREATE TABLE `agency` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `name` varchar(64) DEFAULT NULL,
  `tel` varchar(64) DEFAULT NULL,
  `addr` varchar(128) DEFAULT NULL,
  `email` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for agency_user_associate
-- ----------------------------
DROP TABLE IF EXISTS `agency_user_associate`;
CREATE TABLE `agency_user_associate` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `agency_id` bigint(20) unsigned DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for device_in_merchant
-- ----------------------------
DROP TABLE IF EXISTS `device_in_merchant`;
CREATE TABLE `device_in_merchant` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `device_id` bigint(20) DEFAULT NULL,
  `merchant_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for entry_type
-- ----------------------------
DROP TABLE IF EXISTS `entry_type`;
CREATE TABLE `entry_type` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `name` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for merchant
-- ----------------------------
DROP TABLE IF EXISTS `merchant`;
CREATE TABLE `merchant` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `name` varchar(64) DEFAULT NULL,
  `tel` varchar(64) DEFAULT NULL,
  `addr` varchar(128) DEFAULT NULL,
  `agency_id` bigint(20) DEFAULT NULL,
  `email` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for merchant_user_associate
-- ----------------------------
DROP TABLE IF EXISTS `merchant_user_associate`;
CREATE TABLE `merchant_user_associate` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `store_id` int(10) unsigned DEFAULT NULL,
  `merchant_id` bigint(20) DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `role` varchar(16) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for payment_method
-- ----------------------------
DROP TABLE IF EXISTS `payment_method`;
CREATE TABLE `payment_method` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `name` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for payment_setting_in_device
-- ----------------------------
DROP TABLE IF EXISTS `payment_setting_in_device`;
CREATE TABLE `payment_setting_in_device` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `merchant_device_id` bigint(20) DEFAULT NULL,
  `payment_methods` json DEFAULT NULL,
  `entry_types` json DEFAULT NULL,
  `payment_types` json DEFAULT NULL,
  `acquirer_id` bigint(20) DEFAULT NULL,
  `mid` varchar(32) DEFAULT NULL,
  `tid` varchar(32) DEFAULT NULL,
  `addition` varchar(512) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for payment_type
-- ----------------------------
DROP TABLE IF EXISTS `payment_type`;
CREATE TABLE `payment_type` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `name` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for tms_app
-- ----------------------------
DROP TABLE IF EXISTS `tms_app`;
CREATE TABLE `tms_app` (
  `name` varchar(32) COLLATE utf8_bin NOT NULL,
  `package_id` varchar(64) COLLATE utf8_bin NOT NULL,
  `description` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `agency_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1505150659148678223 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tms_app_file
-- ----------------------------
DROP TABLE IF EXISTS `tms_app_file`;
CREATE TABLE `tms_app_file` (
  `version_name` varchar(32) COLLATE utf8_bin DEFAULT NULL,
  `version_code` int(11) DEFAULT NULL,
  `update_description` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `file_name` varchar(128) COLLATE utf8_bin DEFAULT NULL,
  `file_url` varchar(255) COLLATE utf8_bin NOT NULL,
  `decode_status` varchar(16) COLLATE utf8_bin DEFAULT NULL,
  `decode_fail_msg` varchar(255) COLLATE utf8_bin DEFAULT NULL,
  `app_id` bigint(20) NOT NULL,
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `appid` (`app_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1505154471485801841 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tms_app_in_device
-- ----------------------------
DROP TABLE IF EXISTS `tms_app_in_device`;
CREATE TABLE `tms_app_in_device` (
  `name` varchar(126) COLLATE utf8_bin DEFAULT NULL,
  `package_id` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `version_name` varchar(64) COLLATE utf8_bin DEFAULT NULL,
  `version_code` int(11) DEFAULT NULL,
  `status` varchar(24) COLLATE utf8_bin DEFAULT NULL,
  `app_id` bigint(20) DEFAULT NULL,
  `app_file_id` bigint(20) DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `external_id` bigint(20) DEFAULT NULL,
  `external_id_type` varchar(16) COLLATE utf8_bin DEFAULT '0',
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `d` (`external_id`,`external_id_type`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1505221047908102140 DEFAULT CHARSET=utf8 COLLATE=utf8_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for tms_batch_update
-- ----------------------------
DROP TABLE IF EXISTS `tms_batch_update`;
CREATE TABLE `tms_batch_update` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `module_id` bigint(20) unsigned DEFAULT NULL,
  `store_id` int(10) unsigned DEFAULT NULL,
  `custom_attributes` json DEFAULT NULL,
  `created_by` int(10) unsigned DEFAULT NULL,
  `updated_by` int(10) unsigned DEFAULT NULL,
  `deleted_by` int(10) unsigned DEFAULT NULL,
  `account_id` bigint(20) DEFAULT NULL,
  `status` int(11) DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `tags` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `device_models` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `update_fail_msg` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tms_device
-- ----------------------------
DROP TABLE IF EXISTS `tms_device`;
CREATE TABLE `tms_device` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `agency_id` int(10) unsigned DEFAULT NULL,
  `device_sn` varchar(126) COLLATE utf8_bin NOT NULL,
  `device_csn` varchar(126) COLLATE utf8_bin DEFAULT NULL COMMENT '第三方ID',
  `device_model` varchar(16) COLLATE utf8_bin DEFAULT '0',
  `alias` json DEFAULT NULL,
  `reboot_mode` varchar(16) COLLATE utf8_bin NOT NULL,
  `reboot_time` varchar(6) COLLATE utf8_bin NOT NULL,
  `reboot_day_in_week` int(11) DEFAULT NULL,
  `reboot_day_in_month` int(11) DEFAULT NULL,
  `location_lat` varchar(16) COLLATE utf8_bin DEFAULT NULL,
  `location_lon` varchar(16) COLLATE utf8_bin DEFAULT NULL,
  `push_token` varchar(64) COLLATE utf8_bin DEFAULT NULL COMMENT '推送ID',
  `battery` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`,`device_sn`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1505156075807081547 DEFAULT CHARSET=utf8 COLLATE=utf8_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for tms_device_and_tag_mid
-- ----------------------------
DROP TABLE IF EXISTS `tms_device_and_tag_mid`;
CREATE TABLE `tms_device_and_tag_mid` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `device_id` bigint(20) NOT NULL,
  `tag_id` bigint(20) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  KEY `device_id` (`device_id`,`tag_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=749 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for tms_model
-- ----------------------------
DROP TABLE IF EXISTS `tms_model`;
CREATE TABLE `tms_model` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(32) COLLATE utf8_bin NOT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `module_id` bigint(20) unsigned DEFAULT NULL,
  `store_id` int(10) unsigned DEFAULT NULL,
  `custom_attributes` json DEFAULT NULL,
  `created_by` int(10) unsigned DEFAULT NULL,
  `updated_by` int(10) unsigned DEFAULT NULL,
  `deleted_by` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for tms_tags
-- ----------------------------
DROP TABLE IF EXISTS `tms_tags`;
CREATE TABLE `tms_tags` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `agency_id` bigint(20) unsigned DEFAULT NULL,
  `name` varchar(64) COLLATE utf8_bin NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=30 DEFAULT CHARSET=utf8 COLLATE=utf8_bin ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Table structure for tms_upload_file
-- ----------------------------
DROP TABLE IF EXISTS `tms_upload_file`;
CREATE TABLE `tms_upload_file` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `device_sn` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `file_name` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `file_url` varchar(256) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `agency_id` bigint(20) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `device_sn` (`device_sn`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=734762 DEFAULT CHARSET=utf8 COLLATE=utf8_bin;

-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `agency_id` bigint(20) unsigned DEFAULT NULL,
  `email` varchar(64) NOT NULL,
  `pwd` varchar(64) DEFAULT NULL,
  `name` varchar(64) DEFAULT NULL,
  `role` varchar(16) DEFAULT NULL,
  `active` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`,`email`) USING BTREE,
  KEY `email` (`email`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=62 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_app_id
-- ----------------------------
DROP TABLE IF EXISTS `user_app_id`;
CREATE TABLE `user_app_id` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `store_id` int(10) unsigned DEFAULT NULL,
  `app_id` varchar(64) DEFAULT NULL,
  `app_secret` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_role
-- ----------------------------
DROP TABLE IF EXISTS `user_role`;
CREATE TABLE `user_role` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `store_id` int(10) unsigned DEFAULT NULL,
  `name` varchar(16) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for user_token
-- ----------------------------
DROP TABLE IF EXISTS `user_token`;
CREATE TABLE `user_token` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `created_at` datetime(6) DEFAULT NULL,
  `updated_at` datetime(6) DEFAULT NULL,
  `deleted_at` datetime(6) DEFAULT NULL,
  `store_id` int(10) unsigned DEFAULT NULL,
  `user_id` bigint(20) DEFAULT NULL,
  `app_id` bigint(20) DEFAULT NULL,
  `token` varchar(64) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`,`app_id`),
  KEY `token` (`token`) USING HASH
) ENGINE=InnoDB AUTO_INCREMENT=77 DEFAULT CHARSET=utf8;

SET FOREIGN_KEY_CHECKS = 1;
