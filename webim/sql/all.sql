-- MySQL dump 10.13  Distrib 5.7.24, for osx10.13 (x86_64)
--
-- Host: localhost    Database: custmchat
-- ------------------------------------------------------
-- Server version	5.7.24

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
-- Table structure for table `agent`
--

DROP TABLE IF EXISTS `agent`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `agent` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '租户id',
  `group_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `role_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `avatar` varchar(250) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `username` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '用户名',
  `real_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `nick_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `hashed_password` char(60) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'hashed密码',
  `job_number` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '工号',
  `serve_limit` int(11) NOT NULL DEFAULT '0',
  `is_online` tinyint(1) NOT NULL DEFAULT '0',
  `ranking` int(11) NOT NULL DEFAULT '1' COMMENT '对话排序',
  `email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '邮件',
  `mobile` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '手机号',
  `qq_num` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `signature` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `wechat` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `is_admin` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否管理员',
  `perms_range_type` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限范围',
  `account_status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '账号状态',
  `create_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `update_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
  `deleted_at` datetime(6) DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `mobile` (`mobile`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_ent` (`ent_id`),
  KEY `idx_group` (`group_id`),
  KEY `idx_role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `agent`
--

LOCK TABLES `agent` WRITE;
/*!40000 ALTER TABLE `agent` DISABLE KEYS */;
INSERT INTO `agent` VALUES ('bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','bgrg80l5jj83bqe154g0','bgrg80l5jj83bqe154gg','','2550418657@qq.com','','','$2a$10$EruorL4FShhjB14nmztrZ.KetzLurgrUipxg4EL5fbnFzkrTYHuo.','',1,0,0,'2550418657@qq.com','18868905690','','','offline','',1,'all','valid','2019-01-10 16:36:18.684521','2019-01-10 16:37:30.498321',NULL);
/*!40000 ALTER TABLE `agent` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `agent_group`
--

DROP TABLE IF EXISTS `agent_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `agent_group` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ent_id_name` (`ent_id`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `agent_group`
--

LOCK TABLES `agent_group` WRITE;
/*!40000 ALTER TABLE `agent_group` DISABLE KEYS */;
INSERT INTO `agent_group` VALUES ('bgrg80l5jj83bqe154g0','bgrg80l5jj83bqe154fg','超级管理员','超管');
/*!40000 ALTER TABLE `agent_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `agent_group_relation`
--

DROP TABLE IF EXISTS `agent_group_relation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `agent_group_relation` (
  `group_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `uid` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  UNIQUE KEY `uid` (`uid`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `agent_group_relation`
--

LOCK TABLES `agent_group_relation` WRITE;
/*!40000 ALTER TABLE `agent_group_relation` DISABLE KEYS */;
INSERT INTO `agent_group_relation` VALUES ('bgrg80l5jj83bqe154g0','bgrg80l5jj83bqe154h0');
/*!40000 ALTER TABLE `agent_group_relation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `agent_statistic`
--

DROP TABLE IF EXISTS `agent_statistic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `agent_statistic` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `agent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `conversation_count` int(10) unsigned NOT NULL DEFAULT '0',
  `message_count` int(10) unsigned NOT NULL DEFAULT '0',
  `avg_conversation_duration` float NOT NULL DEFAULT '0',
  `avg_first_resp_duration` float NOT NULL DEFAULT '0',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  KEY `ent_id` (`ent_id`),
  KEY `agent_id` (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `agent_statistic`
--

LOCK TABLES `agent_statistic` WRITE;
/*!40000 ALTER TABLE `agent_statistic` DISABLE KEYS */;
/*!40000 ALTER TABLE `agent_statistic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `allocation_rule`
--

DROP TABLE IF EXISTS `allocation_rule`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `allocation_rule` (
  `id` char(20) NOT NULL,
  `ent_id` char(20) NOT NULL,
  `rule_type` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `allocation_rule`
--

LOCK TABLES `allocation_rule` WRITE;
/*!40000 ALTER TABLE `allocation_rule` DISABLE KEYS */;
/*!40000 ALTER TABLE `allocation_rule` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `automatic_message`
--

DROP TABLE IF EXISTS `automatic_message`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `automatic_message` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `channel_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '渠道类型(web, sdk, ...)',
  `msg_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '推广消息，企业欢迎消息，...',
  `msg_content` text COLLATE utf8mb4_unicode_ci,
  `after_seconds` int(11) NOT NULL DEFAULT '0' COMMENT '多长时间(秒)之后发送',
  `enabled` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否启用',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_ent` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `automatic_message`
--

LOCK TABLES `automatic_message` WRITE;
/*!40000 ALTER TABLE `automatic_message` DISABLE KEYS */;
/*!40000 ALTER TABLE `automatic_message` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `conversation`
--

DROP TABLE IF EXISTS `conversation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `conversation` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `trace_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `agent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `agent_msg_count` int(10) unsigned NOT NULL DEFAULT '0',
  `agent_type` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `msg_count` int(10) unsigned NOT NULL DEFAULT '0',
  `title` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `client_first_send_time` datetime(6) DEFAULT NULL,
  `client_msg_count` int(10) unsigned NOT NULL DEFAULT '0',
  `duration` int(10) unsigned NOT NULL DEFAULT '0',
  `first_msg_created_at` datetime(6) DEFAULT NULL,
  `first_response_wait_time` int(11) DEFAULT '0',
  `last_msg_content` varchar(2000) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `last_msg_created_at` datetime(6) DEFAULT NULL,
  `quality_grade` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT '',
  `summary` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `update_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  `ended_at` datetime(6) DEFAULT NULL,
  `ended_by` char(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_ent` (`ent_id`),
  KEY `idx_agent` (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `conversation`
--

LOCK TABLES `conversation` WRITE;
/*!40000 ALTER TABLE `conversation` DISABLE KEYS */;
/*!40000 ALTER TABLE `conversation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `conversation_quality`
--

DROP TABLE IF EXISTS `conversation_quality`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `conversation_quality` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `grade` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `visitor_msg_count` int(11) NOT NULL DEFAULT '0',
  `agent_msg_count` int(11) NOT NULL DEFAULT '0',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '启用/禁用',
  UNIQUE KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `conversation_quality`
--

LOCK TABLES `conversation_quality` WRITE;
/*!40000 ALTER TABLE `conversation_quality` DISABLE KEYS */;
/*!40000 ALTER TABLE `conversation_quality` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `conversation_statistic`
--

DROP TABLE IF EXISTS `conversation_statistic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `conversation_statistic` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `conversation_count` int(10) unsigned NOT NULL DEFAULT '0',
  `effective_conversation_count` int(10) unsigned NOT NULL DEFAULT '0',
  `message_count` int(10) unsigned NOT NULL DEFAULT '0',
  `avg_resp_duration` float NOT NULL DEFAULT '0',
  `avg_conversation_duration` float NOT NULL DEFAULT '0',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `conversation_statistic`
--

LOCK TABLES `conversation_statistic` WRITE;
/*!40000 ALTER TABLE `conversation_statistic` DISABLE KEYS */;
/*!40000 ALTER TABLE `conversation_statistic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `conversation_transfer`
--

DROP TABLE IF EXISTS `conversation_transfer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `conversation_transfer` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `duration` int(11) NOT NULL DEFAULT '30',
  `transfer_target` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `target_type` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '启用/禁用',
  UNIQUE KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `conversation_transfer`
--

LOCK TABLES `conversation_transfer` WRITE;
/*!40000 ALTER TABLE `conversation_transfer` DISABLE KEYS */;
/*!40000 ALTER TABLE `conversation_transfer` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ending_conversation`
--

DROP TABLE IF EXISTS `ending_conversation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ending_conversation` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `no_message_duration` int(11) NOT NULL DEFAULT '-1',
  `offline_duration` int(11) NOT NULL DEFAULT '-1',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '启用/禁用',
  UNIQUE KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ending_conversation`
--

LOCK TABLES `ending_conversation` WRITE;
/*!40000 ALTER TABLE `ending_conversation` DISABLE KEYS */;
/*!40000 ALTER TABLE `ending_conversation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ending_message`
--

DROP TABLE IF EXISTS `ending_message`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ending_message` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `platform` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'web, sdk, wechat, weibo',
  `agent` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `system` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '启用',
  `prompt` tinyint(1) NOT NULL DEFAULT '0',
  UNIQUE KEY `ent_id` (`ent_id`,`platform`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ending_message`
--

LOCK TABLES `ending_message` WRITE;
/*!40000 ALTER TABLE `ending_message` DISABLE KEYS */;
/*!40000 ALTER TABLE `ending_message` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ent_app`
--

DROP TABLE IF EXISTS `ent_app`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ent_app` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `app_name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `create_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `update_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ent_id` (`ent_id`,`app_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ent_app`
--

LOCK TABLES `ent_app` WRITE;
/*!40000 ALTER TABLE `ent_app` DISABLE KEYS */;
INSERT INTO `ent_app` VALUES ('bgrg80l5jj83bqe154hg','bgrg80l5jj83bqe154fg','integrate_settings','2019-01-10 16:36:18.688231','2019-01-10 16:36:18.688231'),('bgrg80l5jj83bqe154i0','bgrg80l5jj83bqe154fg','visitor','2019-01-10 16:36:18.688231','2019-01-10 16:36:18.688231'),('bgrg80l5jj83bqe154ig','bgrg80l5jj83bqe154fg','conversation','2019-01-10 16:36:18.688231','2019-01-10 16:36:18.688231'),('bgrg80l5jj83bqe154j0','bgrg80l5jj83bqe154fg','history_conversation','2019-01-10 16:36:18.688231','2019-01-10 16:36:18.688231'),('bgrg80l5jj83bqe154jg','bgrg80l5jj83bqe154fg','customer','2019-01-10 16:36:18.688231','2019-01-10 16:36:18.688231'),('bgrg80l5jj83bqe154k0','bgrg80l5jj83bqe154fg','data_report','2019-01-10 16:36:18.688231','2019-01-10 16:36:18.688231'),('bgrg80l5jj83bqe154kg','bgrg80l5jj83bqe154fg','ent_info','2019-01-10 16:36:18.688231','2019-01-10 16:36:18.688231'),('bgrg80l5jj83bqe154l0','bgrg80l5jj83bqe154fg','agent_settings','2019-01-10 16:36:18.688231','2019-01-10 16:36:18.688231');
/*!40000 ALTER TABLE `ent_app` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `enterprise`
--

DROP TABLE IF EXISTS `enterprise`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `enterprise` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企(事)业单位名称',
  `avatar` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `industry` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '所属行业',
  `email` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '电子邮件',
  `mobile` varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '联系电话',
  `description` text COLLATE utf8mb4_unicode_ci COMMENT '描述',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `owner` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '所属人',
  `plan` int(11) NOT NULL COMMENT '企业版本',
  `agent_num` int(11) NOT NULL COMMENT '坐席数',
  `trial_status` int(11) NOT NULL COMMENT '试用状态',
  `expiration_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '过期时间',
  `last_activated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '客服最后活跃时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  KEY `idx_ent_name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `enterprise`
--

LOCK TABLES `enterprise` WRITE;
/*!40000 ALTER TABLE `enterprise` DISABLE KEYS */;
INSERT INTO `enterprise` VALUES ('bgrg80l5jj83bqe154fg','test ent 003','','','2550418657@qq.com','18868905690','','2019-01-10 16:36:18.683882','2550418657@qq.com',1,1,2,'2019-02-09 16:36:18.683882','2019-01-10 16:36:18.683882');
/*!40000 ALTER TABLE `enterprise` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `evaluation`
--

DROP TABLE IF EXISTS `evaluation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `evaluation` (
  `ent_id` char(20) NOT NULL,
  `agent_id` char(20) NOT NULL,
  `conv_id` char(20) NOT NULL COMMENT '对话id',
  `level` tinyint(4) NOT NULL,
  `content` varchar(255) NOT NULL DEFAULT '',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  UNIQUE KEY `conv_id` (`conv_id`),
  KEY `idx_ent` (`ent_id`),
  KEY `idx_agent` (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `evaluation`
--

LOCK TABLES `evaluation` WRITE;
/*!40000 ALTER TABLE `evaluation` DISABLE KEYS */;
/*!40000 ALTER TABLE `evaluation` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invitation_logs`
--

DROP TABLE IF EXISTS `invitation_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invitation_logs` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企业ID',
  `trace_id` char(20) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  `look_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '外观ID',
  `look_config_index` int(11) NOT NULL COMMENT '所属规则中外观配置中的索引',
  `mech_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '机制ID',
  `mech_config_index` int(11) NOT NULL COMMENT '所属规则中外观配置中的索引',
  `is_accepted` tinyint(1) NOT NULL DEFAULT '0',
  `conversation_id` char(20) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_ent_id_created_at` (`ent_id`,`created_at`),
  KEY `idx_ent_id_look_id_created_at` (`ent_id`,`look_id`,`created_at`),
  KEY `idx_ent_id_mech_id_created_at` (`ent_id`,`mech_id`,`created_at`),
  KEY `idx_ent_id_trace_id_created_at` (`ent_id`,`trace_id`,`created_at`),
  KEY `idx_ent_id_look_id_config_index_created_at` (`ent_id`,`look_id`,`look_config_index`,`created_at`),
  KEY `idx_ent_id_mech_id_config_index_created_at` (`ent_id`,`mech_id`,`mech_config_index`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invitation_logs`
--

LOCK TABLES `invitation_logs` WRITE;
/*!40000 ALTER TABLE `invitation_logs` DISABLE KEYS */;
/*!40000 ALTER TABLE `invitation_logs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invitation_looks`
--

DROP TABLE IF EXISTS `invitation_looks`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invitation_looks` (
  `rule_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企业ID',
  `name` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `target` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '平台类型',
  `begin_on` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '开始时间',
  `expire_on` datetime(6) DEFAULT NULL COMMENT '过期时间,如果为null,则为永久',
  `enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用',
  `style_config` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '样式配置',
  `is_any` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否任意匹配',
  `version` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '版本号',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  PRIMARY KEY (`rule_id`),
  KEY `idx_enterprise_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invitation_looks`
--

LOCK TABLES `invitation_looks` WRITE;
/*!40000 ALTER TABLE `invitation_looks` DISABLE KEYS */;
/*!40000 ALTER TABLE `invitation_looks` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invitation_looks_priority`
--

DROP TABLE IF EXISTS `invitation_looks_priority`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invitation_looks_priority` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `priority` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '优先级配置',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  PRIMARY KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invitation_looks_priority`
--

LOCK TABLES `invitation_looks_priority` WRITE;
/*!40000 ALTER TABLE `invitation_looks_priority` DISABLE KEYS */;
/*!40000 ALTER TABLE `invitation_looks_priority` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invitation_looks_rules`
--

DROP TABLE IF EXISTS `invitation_looks_rules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invitation_looks_rules` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企业ID',
  `rule` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '规则配置(条件)',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invitation_looks_rules`
--

LOCK TABLES `invitation_looks_rules` WRITE;
/*!40000 ALTER TABLE `invitation_looks_rules` DISABLE KEYS */;
/*!40000 ALTER TABLE `invitation_looks_rules` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invitation_mechs`
--

DROP TABLE IF EXISTS `invitation_mechs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invitation_mechs` (
  `rule_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '规则ID',
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企业ID',
  `name` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `target` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '平台类型',
  `begin_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '开始时间',
  `expire_at` datetime(6) DEFAULT NULL COMMENT '过期时间,如果为null,则为永久',
  `enabled` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否启用',
  `style_config` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '样式配置',
  `is_any` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否任意匹配',
  `version` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '版本号',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  PRIMARY KEY (`rule_id`),
  KEY `idx_ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invitation_mechs`
--

LOCK TABLES `invitation_mechs` WRITE;
/*!40000 ALTER TABLE `invitation_mechs` DISABLE KEYS */;
/*!40000 ALTER TABLE `invitation_mechs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invitation_mechs_priority`
--

DROP TABLE IF EXISTS `invitation_mechs_priority`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invitation_mechs_priority` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企业ID',
  `priority` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '优先级配置',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  PRIMARY KEY (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invitation_mechs_priority`
--

LOCK TABLES `invitation_mechs_priority` WRITE;
/*!40000 ALTER TABLE `invitation_mechs_priority` DISABLE KEYS */;
/*!40000 ALTER TABLE `invitation_mechs_priority` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `invitation_mechs_rules`
--

DROP TABLE IF EXISTS `invitation_mechs_rules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `invitation_mechs_rules` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '规则ID',
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企业ID',
  `rule` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '规则配置(条件)',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `idx_enterprise_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `invitation_mechs_rules`
--

LOCK TABLES `invitation_mechs_rules` WRITE;
/*!40000 ALTER TABLE `invitation_mechs_rules` DISABLE KEYS */;
/*!40000 ALTER TABLE `invitation_mechs_rules` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `leave_message`
--

DROP TABLE IF EXISTS `leave_message`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `leave_message` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci,
  `status` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'unhandled',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `leave_message`
--

LOCK TABLES `leave_message` WRITE;
/*!40000 ALTER TABLE `leave_message` DISABLE KEYS */;
/*!40000 ALTER TABLE `leave_message` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `leave_message_config`
--

DROP TABLE IF EXISTS `leave_message_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `leave_message_config` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `introduction` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `show_visitor_name` tinyint(1) NOT NULL DEFAULT '0',
  `show_telephone` tinyint(1) NOT NULL DEFAULT '0',
  `show_email` tinyint(1) NOT NULL DEFAULT '0',
  `show_wechat` tinyint(1) NOT NULL DEFAULT '0',
  `show_qq` tinyint(1) NOT NULL DEFAULT '0',
  `auto_create_category` tinyint(1) NOT NULL DEFAULT '0',
  `fill_contact` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT 'single',
  `use_default_content` tinyint(1) NOT NULL DEFAULT '0',
  `default_content` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  UNIQUE KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `leave_message_config`
--

LOCK TABLES `leave_message_config` WRITE;
/*!40000 ALTER TABLE `leave_message_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `leave_message_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `login_limit`
--

DROP TABLE IF EXISTS `login_limit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `login_limit` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0',
  `group_ids` text COLLATE utf8mb4_unicode_ci,
  `city_list` text COLLATE utf8mb4_unicode_ci,
  `allowed_ip_list` text COLLATE utf8mb4_unicode_ci,
  UNIQUE KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `login_limit`
--

LOCK TABLES `login_limit` WRITE;
/*!40000 ALTER TABLE `login_limit` DISABLE KEYS */;
/*!40000 ALTER TABLE `login_limit` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `login_records`
--

DROP TABLE IF EXISTS `login_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `login_records` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `agent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `login_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '登录时间',
  `login_client` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '登录客户端',
  `login_ip` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '登录ip',
  `device_info` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '登录设备信息',
  PRIMARY KEY (`id`),
  KEY `idx_ent` (`ent_id`),
  KEY `idx_agent` (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `login_records`
--

LOCK TABLES `login_records` WRITE;
/*!40000 ALTER TABLE `login_records` DISABLE KEYS */;
INSERT INTO `login_records` VALUES ('bgrg8ul5jj83bqe15570','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-10 16:38:18.205491','web','::1','PostmanRuntime/7.4.0');
/*!40000 ALTER TABLE `login_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `message`
--

DROP TABLE IF EXISTS `message`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `message` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `trace_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `agent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `conversation_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `from_type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `content_type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime(6) DEFAULT NULL,
  `read_time` datetime(6) DEFAULT NULL,
  `content` text COLLATE utf8mb4_unicode_ci,
  `msg_type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `extra` text COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`),
  KEY `idx_conversation_id` (`conversation_id`),
  KEY `msg_find_by_ent_and_track` (`ent_id`,`trace_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `message`
--

LOCK TABLES `message` WRITE;
/*!40000 ALTER TABLE `message` DISABLE KEYS */;
/*!40000 ALTER TABLE `message` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `message_beep`
--

DROP TABLE IF EXISTS `message_beep`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `message_beep` (
  `agent_id` char(20) NOT NULL,
  `client_type` char(20) NOT NULL,
  `beep_type` varchar(20) NOT NULL,
  `new_conversation` tinyint(1) NOT NULL DEFAULT '0',
  `new_message` tinyint(1) NOT NULL DEFAULT '0',
  `conversation_transfer_in` tinyint(1) NOT NULL DEFAULT '0',
  `conversation_transfer_out` tinyint(1) NOT NULL DEFAULT '0',
  `colleague_conversation` tinyint(1) NOT NULL DEFAULT '0',
  UNIQUE KEY `agent_id` (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `message_beep`
--

LOCK TABLES `message_beep` WRITE;
/*!40000 ALTER TABLE `message_beep` DISABLE KEYS */;
/*!40000 ALTER TABLE `message_beep` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `perm`
--

DROP TABLE IF EXISTS `perm`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `perm` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企业id',
  `app_name` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品名称',
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限名',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ent_id` (`ent_id`,`app_name`,`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `perm`
--

LOCK TABLES `perm` WRITE;
/*!40000 ALTER TABLE `perm` DISABLE KEYS */;
INSERT INTO `perm` VALUES ('bgrg80l5jj83bqe154lg','bgrg80l5jj83bqe154fg','data_report','use_data_report','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154m0','bgrg80l5jj83bqe154fg','data_report','export_data','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154mg','bgrg80l5jj83bqe154fg','ent_info','check_update_ent_info','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154n0','bgrg80l5jj83bqe154fg','ent_info','check_ent_agent_group_info','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154ng','bgrg80l5jj83bqe154fg','ent_info','create_update_delete_agent_group','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154o0','bgrg80l5jj83bqe154fg','ent_info','create_update_delete_agent_account','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154og','bgrg80l5jj83bqe154fg','ent_info','check_update_visitor_limit','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154p0','bgrg80l5jj83bqe154fg','ent_info','check_update_ent_role','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154pg','bgrg80l5jj83bqe154fg','ent_info','check_update_ent_safety','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154q0','bgrg80l5jj83bqe154fg','agent_settings','check_manage_blacklist','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154qg','bgrg80l5jj83bqe154fg','agent_settings','check_update_visitor_queue','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154r0','bgrg80l5jj83bqe154fg','agent_settings','check_update_prechat_form','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154rg','bgrg80l5jj83bqe154fg','agent_settings','check_update_tag','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154s0','bgrg80l5jj83bqe154fg','agent_settings','check_update_quick_reply','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154sg','bgrg80l5jj83bqe154fg','agent_settings','check_update_auto_message','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154t0','bgrg80l5jj83bqe154fg','agent_settings','check_update_agent_allocation','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154tg','bgrg80l5jj83bqe154fg','agent_settings','check_update_conversation_invite','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154u0','bgrg80l5jj83bqe154fg','agent_settings','check_update_conversation_rule','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154ug','bgrg80l5jj83bqe154fg','agent_settings','check_update_agent_evaluation','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154v0','bgrg80l5jj83bqe154fg','integrate_settings','check_update_website_plugin','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe154vg','bgrg80l5jj83bqe154fg','integrate_settings','check_update_chat_link','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe15500','bgrg80l5jj83bqe154fg','visitor','update_agent_status','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe1550g','bgrg80l5jj83bqe154fg','visitor','invite_visitor','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe15510','bgrg80l5jj83bqe154fg','visitor','add_delete_visitor_tag','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe1551g','bgrg80l5jj83bqe154fg','visitor','add_update_visitor_card','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe15520','bgrg80l5jj83bqe154fg','visitor','block_unblock_visitor','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe1552g','bgrg80l5jj83bqe154fg','conversation','transfer_others_conversation','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe15530','bgrg80l5jj83bqe154fg','conversation','others_transfer_conversation','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe1553g','bgrg80l5jj83bqe154fg','conversation','end_others_conversation','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe15540','bgrg80l5jj83bqe154fg','history_conversation','check_others_conversation','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe1554g','bgrg80l5jj83bqe154fg','history_conversation','export_conversation','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe15550','bgrg80l5jj83bqe154fg','history_conversation','update_conversation_summary','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe1555g','bgrg80l5jj83bqe154fg','customer','use_customer_app','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe15560','bgrg80l5jj83bqe154fg','customer','check_update_customer_card','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813'),('bgrg80l5jj83bqe1556g','bgrg80l5jj83bqe154fg','customer','check_update_customer_list','2019-01-10 16:36:18.688813','2019-01-10 16:36:18.688813');
/*!40000 ALTER TABLE `perm` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `perms_range_groups`
--

DROP TABLE IF EXISTS `perms_range_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `perms_range_groups` (
  `agent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `group_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  UNIQUE KEY `agent_id` (`agent_id`,`group_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `perms_range_groups`
--

LOCK TABLES `perms_range_groups` WRITE;
/*!40000 ALTER TABLE `perms_range_groups` DISABLE KEYS */;
/*!40000 ALTER TABLE `perms_range_groups` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `plan`
--

DROP TABLE IF EXISTS `plan`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `plan` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `plan_type` tinyint(4) NOT NULL DEFAULT '0',
  `agent_serve_limit` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `plan_type` (`plan_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `plan`
--

LOCK TABLES `plan` WRITE;
/*!40000 ALTER TABLE `plan` DISABLE KEYS */;
INSERT INTO `plan` VALUES ('bgpep8d5jj88eme2mleg',1,2),('bgpepst5jj88etu1qkqg',2,10),('bgpeq4d5jj88f8jf2qq0',3,50);
/*!40000 ALTER TABLE `plan` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `prechat_form`
--

DROP TABLE IF EXISTS `prechat_form`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `prechat_form` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '启用/禁用',
  `title` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `form_fields` text COLLATE utf8mb4_unicode_ci,
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_ent` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `prechat_form`
--

LOCK TABLES `prechat_form` WRITE;
/*!40000 ALTER TABLE `prechat_form` DISABLE KEYS */;
/*!40000 ALTER TABLE `prechat_form` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `queue_config`
--

DROP TABLE IF EXISTS `queue_config`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `queue_config` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `queue_size` int(11) NOT NULL DEFAULT '0' COMMENT '排队访客数量',
  `description` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '排队提示文案',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '启用/禁用',
  UNIQUE KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `queue_config`
--

LOCK TABLES `queue_config` WRITE;
/*!40000 ALTER TABLE `queue_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `queue_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `quickreply_group`
--

DROP TABLE IF EXISTS `quickreply_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `quickreply_group` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `created_by` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `quickreply_group`
--

LOCK TABLES `quickreply_group` WRITE;
/*!40000 ALTER TABLE `quickreply_group` DISABLE KEYS */;
/*!40000 ALTER TABLE `quickreply_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `quickreply_item`
--

DROP TABLE IF EXISTS `quickreply_item`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `quickreply_item` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `quickreply_group_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci,
  `created_by` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `quickreply_item`
--

LOCK TABLES `quickreply_item` WRITE;
/*!40000 ALTER TABLE `quickreply_item` DISABLE KEYS */;
/*!40000 ALTER TABLE `quickreply_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role`
--

DROP TABLE IF EXISTS `role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_ent` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role`
--

LOCK TABLES `role` WRITE;
/*!40000 ALTER TABLE `role` DISABLE KEYS */;
INSERT INTO `role` VALUES ('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154fg','超管');
/*!40000 ALTER TABLE `role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_perm`
--

DROP TABLE IF EXISTS `role_perm`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `role_perm` (
  `role_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `perm_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '权限id',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
  UNIQUE KEY `role_id` (`role_id`,`perm_id`),
  KEY `perm_id` (`perm_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_perm`
--

LOCK TABLES `role_perm` WRITE;
/*!40000 ALTER TABLE `role_perm` DISABLE KEYS */;
INSERT INTO `role_perm` VALUES ('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154lg','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154m0','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154mg','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154n0','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154ng','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154o0','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154og','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154p0','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154pg','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154q0','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154qg','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154r0','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154rg','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154s0','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154sg','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154t0','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154tg','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154u0','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154ug','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154v0','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154vg','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe15500','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe1550g','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe15510','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe1551g','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe15520','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe1552g','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe15530','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe1553g','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe15540','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe1554g','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe15550','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe1555g','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe15560','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588'),('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe1556g','2019-01-10 16:36:18.693588','2019-01-10 16:36:18.693588');
/*!40000 ALTER TABLE `role_perm` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `send_file`
--

DROP TABLE IF EXISTS `send_file`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `send_file` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT '0',
  UNIQUE KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `send_file`
--

LOCK TABLES `send_file` WRITE;
/*!40000 ALTER TABLE `send_file` DISABLE KEYS */;
/*!40000 ALTER TABLE `send_file` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_group`
--

DROP TABLE IF EXISTS `user_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `user_group` (
  `id` char(20) NOT NULL,
  `ent_id` char(20) NOT NULL,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_group`
--

LOCK TABLES `user_group` WRITE;
/*!40000 ALTER TABLE `user_group` DISABLE KEYS */;
/*!40000 ALTER TABLE `user_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `visit`
--

DROP TABLE IF EXISTS `visit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `visit` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企业id',
  `trace_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `visit_page_cnt` int(11) NOT NULL DEFAULT '1' COMMENT '访问页数',
  `residence_time_sec` int(11) NOT NULL DEFAULT '1' COMMENT '停留秒数',
  `browser_family` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `browser_version_string` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `browser_version` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `os_category` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `os_family` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `os_version_string` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `os_version` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `platform` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `ua_string` text COLLATE utf8mb4_unicode_ci NOT NULL,
  `ip` varchar(50) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '',
  `country` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `province` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `city` varchar(30) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `isp` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `first_page_source` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '首次来路页source',
  `first_page_source_keyword` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '首次来路页keyword',
  `first_page_source_domain` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '首次来路页domain',
  `first_page_source_url` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '首次来路页url',
  `first_page_title` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '首次着陆页title',
  `first_page_domain` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '首次着陆页domain',
  `first_page_url` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '首次着陆页url',
  `latest_title` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '最新着陆页title',
  `latest_url` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '最新着陆页url',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_enterprise_id_trace_id` (`ent_id`,`trace_id`),
  KEY `idx_enterprise_id_created_on` (`ent_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `visit`
--

LOCK TABLES `visit` WRITE;
/*!40000 ALTER TABLE `visit` DISABLE KEYS */;
/*!40000 ALTER TABLE `visit` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `visit_blacklist`
--

DROP TABLE IF EXISTS `visit_blacklist`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `visit_blacklist` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企业id',
  `trace_id` char(20) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
  `visit_id` char(20) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '',
  `agent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '坐席id',
  `conv_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '对话id',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ent_id` (`ent_id`,`trace_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `visit_blacklist`
--

LOCK TABLES `visit_blacklist` WRITE;
/*!40000 ALTER TABLE `visit_blacklist` DISABLE KEYS */;
/*!40000 ALTER TABLE `visit_blacklist` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `visit_page`
--

DROP TABLE IF EXISTS `visit_page`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `visit_page` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企业id',
  `visit_id` char(20) CHARACTER SET ascii COLLATE ascii_bin NOT NULL COMMENT '访问id',
  `ip` varchar(50) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT '',
  `source` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '首次来路页source',
  `source_keyword` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '首次来路页keyword',
  `source_domain` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '首次来路页domain',
  `source_url` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '首次来路页url',
  `title` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '首次着陆页title',
  `domain` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '首次着陆页domain',
  `url` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '首次着陆页url',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
  PRIMARY KEY (`id`),
  KEY `idx_ent_id_visit_id_create_at` (`ent_id`,`visit_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `visit_page`
--

LOCK TABLES `visit_page` WRITE;
/*!40000 ALTER TABLE `visit_page` DISABLE KEYS */;
/*!40000 ALTER TABLE `visit_page` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `visitor`
--

DROP TABLE IF EXISTS `visitor`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `visitor` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企业id',
  `trace_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT 'trace id',
  `name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `age` int(11) NOT NULL DEFAULT '0',
  `gender` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `avatar` text COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '头像',
  `mobile` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `weibo` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `wechat` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `email` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `qq_num` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `address` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `remark` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '备注',
  `visit_cnt` int(11) NOT NULL DEFAULT '1' COMMENT '累计访问次数',
  `visit_page_cnt` int(11) NOT NULL DEFAULT '1' COMMENT '累计访问页数',
  `residence_time_sec` int(11) NOT NULL DEFAULT '1' COMMENT '累计停留时长',
  `last_visit_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '最近访问id',
  `visited_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '最近访问时间戳',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `udx_ent_id_trace_id` (`ent_id`,`trace_id`),
  KEY `idx_ent_id_created_at` (`ent_id`,`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `visitor`
--

LOCK TABLES `visitor` WRITE;
/*!40000 ALTER TABLE `visitor` DISABLE KEYS */;
/*!40000 ALTER TABLE `visitor` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `visitor_statistic`
--

DROP TABLE IF EXISTS `visitor_statistic`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `visitor_statistic` (
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `visitor_count` int(10) unsigned NOT NULL DEFAULT '0',
  `visit_num` int(10) unsigned NOT NULL DEFAULT '0',
  `page_views` int(10) unsigned NOT NULL DEFAULT '0',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `visitor_statistic`
--

LOCK TABLES `visitor_statistic` WRITE;
/*!40000 ALTER TABLE `visitor_statistic` DISABLE KEYS */;
/*!40000 ALTER TABLE `visitor_statistic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `visitor_tag`
--

DROP TABLE IF EXISTS `visitor_tag`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `visitor_tag` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `creator` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `color` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `use_count` int(11) NOT NULL DEFAULT '0',
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_ent` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `visitor_tag`
--

LOCK TABLES `visitor_tag` WRITE;
/*!40000 ALTER TABLE `visitor_tag` DISABLE KEYS */;
INSERT INTO `visitor_tag` VALUES ('bgrgh9t5jj83bqe15580','bgrg80l5jj83bqe154fg','bgrg80l5jj83bqe154h0','test tag update','浅蓝',0,'2019-01-10 16:56:07.364274','2019-01-10 17:00:40.680643');
/*!40000 ALTER TABLE `visitor_tag` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `visitor_tag_relation`
--

DROP TABLE IF EXISTS `visitor_tag_relation`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `visitor_tag_relation` (
  `visitor_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `tag_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  UNIQUE KEY `visitor_id` (`visitor_id`,`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `visitor_tag_relation`
--

LOCK TABLES `visitor_tag_relation` WRITE;
/*!40000 ALTER TABLE `visitor_tag_relation` DISABLE KEYS */;
/*!40000 ALTER TABLE `visitor_tag_relation` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2019-01-11 10:12:37
