-- MySQL dump 10.13  Distrib 5.7.24, for osx10.13 (x86_64)
--
-- Host: 127.0.0.1    Database: custmchat
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
  `public_email` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `public_telephone` varchar(32) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
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
INSERT INTO `agent` VALUES ('bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','bgrg80l5jj83bqe154g0','bgrg80l5jj83bqe154gg','','2550418657@qq.com','','','$2a$10$0D1CpoQq3hJXg7ct9V0oD.K3Wl8CDhW/fgqYxYMiV3kZ759IgGSQG','',1,0,0,'2550418657@qq.com','18868905690','','','','','unavailable','',1,'all','valid','2019-01-10 16:36:18.684521','2019-01-14 17:16:29.506837',NULL);
/*!40000 ALTER TABLE `agent` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `agent_group`
--

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
INSERT INTO `agent_statistic` VALUES ('123456','567890',1,2,10.32,5.88,'2019-01-24 16:36:47.433713','2019-01-26 00:36:47.433713'),('123456','567890',9,20,100.1,123.123,'2019-01-24 16:36:47.433713','2019-01-26 00:36:47.433713'),('123456','567890',1,2,10.32,5.88,'2019-01-23 16:37:59.446313','2019-01-26 08:37:59.446313'),('123456','567890',9,20,100.1,123.123,'2019-01-23 16:37:59.446313','2019-01-26 08:37:59.446313');
/*!40000 ALTER TABLE `agent_statistic` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `allocation_rule`
--

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
-- Table structure for table `auth_group`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth_group` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(80) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_group`
--

LOCK TABLES `auth_group` WRITE;
/*!40000 ALTER TABLE `auth_group` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_group_permissions`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth_group_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `group_id` int(11) NOT NULL,
  `permission_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_group_permissions_group_id_permission_id_0cd325b0_uniq` (`group_id`,`permission_id`),
  KEY `auth_group_permissio_permission_id_84c5c92e_fk_auth_perm` (`permission_id`),
  CONSTRAINT `auth_group_permissio_permission_id_84c5c92e_fk_auth_perm` FOREIGN KEY (`permission_id`) REFERENCES `auth_permission` (`id`),
  CONSTRAINT `auth_group_permissions_group_id_b120cbf9_fk_auth_group_id` FOREIGN KEY (`group_id`) REFERENCES `auth_group` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_group_permissions`
--

LOCK TABLES `auth_group_permissions` WRITE;
/*!40000 ALTER TABLE `auth_group_permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_group_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_permission`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth_permission` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `content_type_id` int(11) NOT NULL,
  `codename` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_permission_content_type_id_codename_01ab375a_uniq` (`content_type_id`,`codename`),
  CONSTRAINT `auth_permission_content_type_id_2f476e4b_fk_django_co` FOREIGN KEY (`content_type_id`) REFERENCES `django_content_type` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=25 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_permission`
--

LOCK TABLES `auth_permission` WRITE;
/*!40000 ALTER TABLE `auth_permission` DISABLE KEYS */;
INSERT INTO `auth_permission` VALUES (1,'Can add log entry',1,'add_logentry'),(2,'Can change log entry',1,'change_logentry'),(3,'Can delete log entry',1,'delete_logentry'),(4,'Can view log entry',1,'view_logentry'),(5,'Can add permission',2,'add_permission'),(6,'Can change permission',2,'change_permission'),(7,'Can delete permission',2,'delete_permission'),(8,'Can view permission',2,'view_permission'),(9,'Can add group',3,'add_group'),(10,'Can change group',3,'change_group'),(11,'Can delete group',3,'delete_group'),(12,'Can view group',3,'view_group'),(13,'Can add user',4,'add_user'),(14,'Can change user',4,'change_user'),(15,'Can delete user',4,'delete_user'),(16,'Can view user',4,'view_user'),(17,'Can add content type',5,'add_contenttype'),(18,'Can change content type',5,'change_contenttype'),(19,'Can delete content type',5,'delete_contenttype'),(20,'Can view content type',5,'view_contenttype'),(21,'Can add session',6,'add_session'),(22,'Can change session',6,'change_session'),(23,'Can delete session',6,'delete_session'),(24,'Can view session',6,'view_session');
/*!40000 ALTER TABLE `auth_permission` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_user`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `password` varchar(128) NOT NULL,
  `last_login` datetime(6) DEFAULT NULL,
  `is_superuser` tinyint(1) NOT NULL,
  `username` varchar(150) NOT NULL,
  `first_name` varchar(30) NOT NULL,
  `last_name` varchar(150) NOT NULL,
  `email` varchar(254) NOT NULL,
  `is_staff` tinyint(1) NOT NULL,
  `is_active` tinyint(1) NOT NULL,
  `date_joined` datetime(6) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_user`
--

LOCK TABLES `auth_user` WRITE;
/*!40000 ALTER TABLE `auth_user` DISABLE KEYS */;
INSERT INTO `auth_user` VALUES (1,'pbkdf2_sha256$120000$9bd7etX9taCg$jd53C89miyO6S03cDsfOAXJ6KIQkL1c+0TH9l+dRyY8=','2019-03-05 08:37:03.589845',1,'zhiru.chen','','','forfd8960@gmail.com',1,1,'2019-03-05 08:36:45.469291');
/*!40000 ALTER TABLE `auth_user` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_user_groups`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth_user_groups` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `group_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_user_groups_user_id_group_id_94350c0c_uniq` (`user_id`,`group_id`),
  KEY `auth_user_groups_group_id_97559544_fk_auth_group_id` (`group_id`),
  CONSTRAINT `auth_user_groups_group_id_97559544_fk_auth_group_id` FOREIGN KEY (`group_id`) REFERENCES `auth_group` (`id`),
  CONSTRAINT `auth_user_groups_user_id_6a12ed8b_fk_auth_user_id` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_user_groups`
--

LOCK TABLES `auth_user_groups` WRITE;
/*!40000 ALTER TABLE `auth_user_groups` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_user_groups` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `auth_user_user_permissions`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `auth_user_user_permissions` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `permission_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_user_user_permissions_user_id_permission_id_14a6b632_uniq` (`user_id`,`permission_id`),
  KEY `auth_user_user_permi_permission_id_1fbb5f2c_fk_auth_perm` (`permission_id`),
  CONSTRAINT `auth_user_user_permi_permission_id_1fbb5f2c_fk_auth_perm` FOREIGN KEY (`permission_id`) REFERENCES `auth_permission` (`id`),
  CONSTRAINT `auth_user_user_permissions_user_id_a95ead1b_fk_auth_user_id` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `auth_user_user_permissions`
--

LOCK TABLES `auth_user_user_permissions` WRITE;
/*!40000 ALTER TABLE `auth_user_user_permissions` DISABLE KEYS */;
/*!40000 ALTER TABLE `auth_user_user_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `automatic_message`
--

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
  `agent_effective_msg_count` int(11) NOT NULL DEFAULT '0',
  `client_last_send_time` datetime(6) DEFAULT NULL,
  `first_msg_create_time` datetime(6) DEFAULT NULL,
  `eval_content` varchar(500) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `eval_level` int(11) NOT NULL DEFAULT '0',
  `has_summary` tinyint(1) NOT NULL DEFAULT '0',
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
-- Table structure for table `django_admin_log`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `django_admin_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `action_time` datetime(6) NOT NULL,
  `object_id` longtext,
  `object_repr` varchar(200) NOT NULL,
  `action_flag` smallint(5) unsigned NOT NULL,
  `change_message` longtext NOT NULL,
  `content_type_id` int(11) DEFAULT NULL,
  `user_id` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `django_admin_log_content_type_id_c4bce8eb_fk_django_co` (`content_type_id`),
  KEY `django_admin_log_user_id_c564eba6_fk_auth_user_id` (`user_id`),
  CONSTRAINT `django_admin_log_content_type_id_c4bce8eb_fk_django_co` FOREIGN KEY (`content_type_id`) REFERENCES `django_content_type` (`id`),
  CONSTRAINT `django_admin_log_user_id_c564eba6_fk_auth_user_id` FOREIGN KEY (`user_id`) REFERENCES `auth_user` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `django_admin_log`
--

LOCK TABLES `django_admin_log` WRITE;
/*!40000 ALTER TABLE `django_admin_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `django_admin_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `django_content_type`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `django_content_type` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_label` varchar(100) NOT NULL,
  `model` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `django_content_type_app_label_model_76bd3d3b_uniq` (`app_label`,`model`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `django_content_type`
--

LOCK TABLES `django_content_type` WRITE;
/*!40000 ALTER TABLE `django_content_type` DISABLE KEYS */;
INSERT INTO `django_content_type` VALUES (1,'admin','logentry'),(3,'auth','group'),(2,'auth','permission'),(4,'auth','user'),(5,'contenttypes','contenttype'),(7,'sales','enterprise'),(6,'sessions','session');
/*!40000 ALTER TABLE `django_content_type` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `django_migrations`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `django_migrations` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `applied` datetime(6) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `django_migrations`
--

LOCK TABLES `django_migrations` WRITE;
/*!40000 ALTER TABLE `django_migrations` DISABLE KEYS */;
INSERT INTO `django_migrations` VALUES (1,'contenttypes','0001_initial','2019-03-05 08:18:35.910775'),(2,'auth','0001_initial','2019-03-05 08:18:36.227717'),(3,'admin','0001_initial','2019-03-05 08:18:36.309909'),(4,'admin','0002_logentry_remove_auto_add','2019-03-05 08:18:36.322340'),(5,'admin','0003_logentry_add_action_flag_choices','2019-03-05 08:18:36.339290'),(6,'contenttypes','0002_remove_content_type_name','2019-03-05 08:18:36.404621'),(7,'auth','0002_alter_permission_name_max_length','2019-03-05 08:18:36.436041'),(8,'auth','0003_alter_user_email_max_length','2019-03-05 08:18:36.477273'),(9,'auth','0004_alter_user_username_opts','2019-03-05 08:18:36.489470'),(10,'auth','0005_alter_user_last_login_null','2019-03-05 08:18:36.516813'),(11,'auth','0006_require_contenttypes_0002','2019-03-05 08:18:36.521352'),(12,'auth','0007_alter_validators_add_error_messages','2019-03-05 08:18:36.534956'),(13,'auth','0008_alter_user_username_max_length','2019-03-05 08:18:36.564131'),(14,'auth','0009_alter_user_last_name_max_length','2019-03-05 08:18:36.595555'),(15,'sessions','0001_initial','2019-03-05 08:18:36.626394');
/*!40000 ALTER TABLE `django_migrations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `django_session`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `django_session` (
  `session_key` varchar(40) NOT NULL,
  `session_data` longtext NOT NULL,
  `expire_date` datetime(6) NOT NULL,
  PRIMARY KEY (`session_key`),
  KEY `django_session_expire_date_a5c62663` (`expire_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `django_session`
--

LOCK TABLES `django_session` WRITE;
/*!40000 ALTER TABLE `django_session` DISABLE KEYS */;
INSERT INTO `django_session` VALUES ('d7pjyu96wanj17fm7evkicohj02c1y0c','MDU1YzI5MGMzMmQ1N2IyMzgxMGM4NDI1NTk3MDAxYzA5MDBhYmUwYzp7Il9hdXRoX3VzZXJfaWQiOiIxIiwiX2F1dGhfdXNlcl9iYWNrZW5kIjoiZGphbmdvLmNvbnRyaWIuYXV0aC5iYWNrZW5kcy5Nb2RlbEJhY2tlbmQiLCJfYXV0aF91c2VyX2hhc2giOiI4MzRiNTE0YWNkMjM0ZjU5NDUzYmFmOWIzNjBjMGU3ODc2MzM0YTg5In0=','2019-03-19 08:37:03.609745');
/*!40000 ALTER TABLE `django_session` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ending_conversation`
--

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
-- Table structure for table `ent_all_configs`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ent_all_configs` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `config_content` text COLLATE utf8mb4_unicode_ci,
  `create_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `update_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  UNIQUE KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ent_all_configs`
--

LOCK TABLES `ent_all_configs` WRITE;
/*!40000 ALTER TABLE `ent_all_configs` DISABLE KEYS */;
INSERT INTO `ent_all_configs` VALUES ('bh3em3d5jj8bc62oksq0','bgrg80l5jj83bqe154fg','{\"agents_permissions_config\":{\"agents_permissions_level\":1},\"auto_reply_msg_settings\":{\"web\":{\"content\":\"אוכחוכחכחלחלח\",\"count_down\":0,\"status\":\"open\"},\"baidu_bcp\":{\"content\":\"您好，工作人员正在忙，请稍等。\",\"count_down\":0,\"status\":\"open\"},\"mini_program\":{\"content\":\"客服无应答\",\"count_down\":0,\"status\":\"open\"},\"sdk\":{\"content\":\"抱歉让您久等了（客服无应答）\",\"count_down\":0,\"status\":\"close\"},\"toutiao\":{\"content\":\"\",\"count_down\":0,\"status\":\"close\"},\"weibo\":{\"content\":\"抱歉让您久等了（微博客服无应答）\",\"count_down\":0,\"status\":\"open\"},\"weixin\":{\"content\":\"抱歉让您久等了\",\"count_down\":0,\"status\":\"close\"}},\"chat_link_auto_msg_config\":{\"web\":{\"status\":\"open\"}},\"client_waking_auto_msg\":{\"web\":{\"content\":\"顾客有事 等会\",\"count_down\":0,\"status\":\"open\"},\"baidu_bcp\":{\"content\":\"我还在等待你的消息哟~请问还有什么可以帮到您的吗？\",\"count_down\":0,\"status\":\"open\"},\"mini_program\":{\"content\":\"访客无应答\",\"count_down\":0,\"status\":\"open\"},\"sdk\":{\"content\":\"我还在等待你的消息哟~请问还有什么可以帮到您的吗？（顾客无应答）\",\"count_down\":0,\"status\":\"close\"},\"toutiao\":{\"content\":\"\",\"count_down\":0,\"status\":\"close\"},\"weibo\":{\"content\":\"我还在等待你的消息哟~请问还有什么可以帮到您的吗？（微博顾客无应答）\",\"count_down\":0,\"status\":\"open\"},\"weixin\":{\"content\":\"我还在等待你的消息哟~请问还有什么可以帮到您的吗？（微信顾客无应答）\",\"count_down\":0,\"status\":\"open\"}},\"conv_grade_config\":{\"enable\":true,\"first_level\":{\"agent_msg_cnt\":0,\"client_msg_cnt\":19},\"second_level\":{\"agent_msg_cnt\":0,\"client_msg_cnt\":10},\"third_level\":{\"agent_msg_cnt\":0,\"client_msg_cnt\":5}},\"end_conv_expire_config\":{\"mini_program\":3,\"sdk\":3,\"web\":{\"no_msg_end\":3,\"offline_end\":30},\"weibo\":3,\"weixin\":3},\"ending_msg_settings\":{\"web\":{\"agent_ending_message\":\"感谢您的咨询，祝您生活工作顺利！（客服手动结束时）\",\"auto_ending_message\":\"感谢您的咨询，祝您生活工作顺利！（系统自动关闭时给顾客发）\",\"status\":\"open\"},\"baidu_bcp\":{\"agent_ending_message\":\"【系统消息】您好，为了保证服务质量，我们已经结束了对话，期待再次为您服务。\",\"auto_ending_message\":\"【系统消息】您好，由于很久没有收到您的消息，系统自动结束了对话。如果还有需要，欢迎随时联系我们。\",\"status\":\"open\"},\"mini_program\":{\"status\":\"open\"},\"sdk\":{\"agent_ending_message\":\"到此为止，再咨询收费了。（手动结束）\",\"auto_ending_message\":\"您好，由于很久没有收到您的消息，系统自动结束了对话。\\n如果仍有需要，欢迎随时联系我们。\",\"status\":\"open\"},\"toutiao\":{\"agent_ending_message\":\"结束\",\"status\":\"close\"},\"weibo\":{\"agent_ending_message\":\"如有问题，欢迎咨询（微博手动结束）\",\"auto_ending_message\":\"您好，由于很久没有收到您的消息，系统自动结束了对话。\\n如果仍有需要，欢迎随时联系我们。（微博自动结束）\",\"status\":\"open\"},\"weixin\":{\"agent_ending_message\":\"如有问题，欢迎咨询（微信手动结束）\",\"auto_ending_message\":\"您好，由于很久没有收到您的消息，系统自动结束了对话。\\n如果仍有需要，欢迎随时联系我们。（微信自动结束）\",\"status\":\"open\"}},\"invitation_config\":{\"auto\":{\"accept\":{\"countdown\":3,\"status\":\"open\"},\"countdown\":9,\"reject\":{\"countdown\":10,\"type\":\"again\"},\"status\":\"open\"},\"desktop\":{\"actions\":[{\"height\":48,\"id\":\"592qxtz42\",\"position\":{\"bottom\":\"auto\",\"left\":46,\"right\":\"auto\",\"top\":26},\"type\":1,\"width\":122},{\"height\":30,\"id\":\"5b4zpxp7a\",\"position\":{\"bottom\":\"auto\",\"left\":174,\"right\":\"auto\",\"top\":1},\"type\":2,\"width\":25},{\"height\":18,\"id\":\"5b4zpyo2u\",\"position\":{\"bottom\":\"auto\",\"left\":0,\"right\":\"auto\",\"top\":0},\"type\":3,\"width\":103,\"linkType\":1}],\"bgi\":{\"height\":100,\"src\":\"\",\"width\":200},\"position\":{\"bottom\":1,\"side\":0,\"type\":1},\"src\":\"https://s3-qcloud.meiqia.com/pics.meiqia.bucket/7c69d3069c4db8ab872b99e5b374cd15.png\",\"text\":\"    \\n终于等到你\",\"type\":5},\"facade_status\":\"open\",\"manual\":{\"accept\":{\"countdown\":4,\"status\":\"close\"},\"reject\":{\"countdown\":10,\"type\":\"stop\"},\"status\":\"open\"},\"mobile\":{\"actions\":[{\"height\":37,\"id\":\"5954xq4vs\",\"position\":{\"bottom\":\"auto\",\"left\":216,\"right\":\"auto\",\"top\":2},\"type\":2,\"width\":39}],\"bgi\":{\"height\":140,\"src\":\"https://s3-qcloud.meiqia.com/pics.meiqia.bucket/e6a03078a326241a7a6376aa3a72d9c8.png\",\"width\":280},\"position\":{\"type\":2,\"value\":42},\"src\":\"https://s3-qcloud.meiqia.com/pics.meiqia.bucket/845a6a9da8a170132f627573be40e2e6.png\",\"text\":\"您好，当前有专业客服人员在线，让我们来帮助您吧。\",\"type\":2}},\"oauth_settings\":{\"identity_key\":\"\",\"retry_times\":0,\"secret_key\":\"\",\"status\":\"close\",\"success_result\":\"\",\"url\":\"\"},\"promotion_msg_settings\":{\"web\":{\"content\":[\"test msg\",\"test msg\"],\"status\":\"open\",\"stop_after_talk\":false},\"baidu_bcp\":{\"content\":[\"test msg\"],\"status\":\"open\",\"stop_after_talk\":false},\"mini_program\":{\"content\":[\"test msg\",\"test msg\"],\"status\":\"close\",\"stop_after_talk\":false},\"sdk\":{\"content\":[\"test msg\"],\"status\":\"close\",\"stop_after_talk\":false},\"toutiao\":{\"content\":[\"test msg\"],\"status\":\"close\",\"stop_after_talk\":false},\"weibo\":{\"content\":[\"test msg\"],\"status\":\"open\",\"stop_after_talk\":false},\"weixin\":{\"content\":[\"test msg\"],\"status\":\"open\",\"stop_after_talk\":false}},\"queue_settings\":null,\"reserve_clues_config\":{\"enabled\":true,\"fallback\":\"allocate_rule\"},\"robot_settings\":{\"auto_complete_only_query_main_q\":true,\"avatar\":\"https://s3.cn-north-1.amazonaws.com.cn/pics.meiqia.bucket/86d1d383c0f6cf72\",\"correlation_threshold\":4,\"failed_threshold\":3,\"left_msg_cnt\":45175,\"manual_redirect\":false,\"more_like_this_count\":3,\"nickname\":\"爱因机器人\",\"provider\":\"meiqia\",\"response_cant_answer\":[{\"text\":\"抱歉，能否再描述一下您的问题，没太听懂呢~\"},{\"text\":\"抱歉请您换一种方式提问哦，方便小洽正确识别\"}],\"response_eval_useful\":[{\"text\":\"可以\"}],\"response_eval_useless\":[{\"text\":\"不可以\"}],\"response_manual_redirect\":[{\"text\":\"不太懂你说什么，手动转人工\"}],\"response_more_like_this\":[{\"text\":\"\"}],\"response_queueing\":[{\"text\":\"不太听懂你说什么，人工客服正忙，请排队稍后\"}],\"response_redirect\":[{\"text\":\"不太懂你说什么，正在为您转接人工，请稍后~\"}],\"response_reply\":[{\"text\":\"不太懂你说什么，要不给我们留言？\"}],\"rule\":\"robot_first\",\"show_switch\":true,\"signature\":\"\",\"status\":\"open\",\"unmatch_threshold\":3,\"welcome_question_ids\":[705987,698801,684911,682159],\"welcome_text\":\"\\u003cdiv\\u003e欢迎，我是机器人小助手，请问遇到了什么问题吗？\\u003c/div\\u003e\\u003cdiv data-role=\\\"blank-line\\\"\\u003e\\u003cbr/\\u003e\\u003c/div\\u003e\"},\"sales_cloud_config\":null,\"send_file_settings\":{\"widget_status\":\"open\"},\"service_evaluation_config\":{\"agent_invitation\":\"close\",\"agent_visible\":\"close\",\"auto_invitation\":\"open\",\"prompt_text\":\"请您为我的服务做出评价\"},\"standalone_window_config\":{\"background\":{\"color\":\"c4cf48\",\"url\":\"\"},\"desktop\":{\"customer_content\":\"\\u003cdiv\\u003e778\\u003c/div\\u003e\\u003cdiv\\u003e\\u003cimg src=https://s3-qcloud.meiqia.com/pics.meiqia.bucket/6ab41102558be9a21dd500dd6ce39a00.png style=\\\"max-width: 100%;\\\"\\u003e\\u003c/div\\u003e\\u003cdiv\\u003e\\u003cbr\\u003e\\u003c/div\\u003e\\u003cdiv\\u003e\\u003cimg src=https://s3-qcloud.meiqia.com/pics.meiqia.bucket/709e561b79d497069df012e12e44ffd7.jpg style=\\\"max-width: 100%;\\\"\\u003e\\u003c/div\\u003e\\u003cdiv\\u003e\\u003cbr\\u003e\\u003c/div\\u003e\\u003cdiv\\u003e\\u003cbr\\u003e\\u003c/div\\u003e\",\"customer_photo_type\":\"small\",\"theme\":[\"8686a5\",\"white\",\"3ad531\"],\"type\":\"fusion\"},\"mobile\":{\"theme\":[\"573942\",\"white\"],\"type\":\"mustang\"},\"removeBrand\":\"close\",\"ring\":\"open\"},\"ticket_config\":{\"captcha\":\"open\",\"category\":\"close\",\"contactRule\":\"single\",\"defaultTemplate\":\"open\",\"defaultTemplateContent\":\"姓名：\\n上手随时：\",\"email\":\"open\",\"intro\":\"您好，现在为下班时间，请您留下个人信息，我们会在第一时间跟您取得联系。1111111111111111111111111111111111111111111111111\",\"name\":\"open\",\"permission\":\"open\",\"qq\":\"close\",\"tel\":\"open\",\"wechat\":\"open\"},\"timeout_redirect_config\":{\"count_down\":0,\"rules\":[{\"count_down\":0,\"type\":\"group_first\"}],\"status\":\"open\"},\"visitor_visible\":{\"region\":\"open\"},\"web_callback_settings\":{\"callback_switch\":\"open\",\"captcha_switch\":\"open\"},\"welcome_msg_settings\":{\"web\":{\"content\":\"你好，Welcome！伦敦伦敦\",\"status\":\"open\"},\"baidu_bcp\":{\"content\":\"手机百度企业欢迎语\",\"status\":\"open\"},\"mini_program\":{\"content\":\"小程序欢迎语\",\"status\":\"open\"},\"sdk\":{\"content\":\"SDK企业欢迎消息112233\",\"status\":\"open\"},\"toutiao\":{\"content\":\"\",\"status\":\"open\"},\"weibo\":{\"content\":\"你好，请问有什么可以帮到您的吗？（微博企业欢迎语）\",\"status\":\"open\"},\"weixin\":{\"content\":\"你好，美洽企业欢迎咨询\",\"status\":\"open\"}},\"is_activated\":true,\"is_from_baidu_open\":true,\"is_group\":false,\"scheduler_after_client_send_msg\":false}','2019-01-22 18:06:05.271043','2019-01-22 18:26:07.690185');
/*!40000 ALTER TABLE `ent_all_configs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `ent_app`
--

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
-- Table structure for table `ent_plan`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ent_plan` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `plan_type` tinyint(4) NOT NULL DEFAULT '0',
  `trial_status` int(11) NOT NULL DEFAULT '0',
  `agent_serve_limit` int(11) NOT NULL DEFAULT '0',
  `login_agent_limit` int(11) NOT NULL DEFAULT '0',
  `agent_num` int(11) NOT NULL DEFAULT '0',
  `pay_amount` int(11) NOT NULL DEFAULT '0',
  `expiration_time` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `create_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) COMMENT '创建时间',
  `update_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6) COMMENT '修改时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `ent_id` (`ent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `ent_plan`
--

LOCK TABLES `ent_plan` WRITE;
/*!40000 ALTER TABLE `ent_plan` DISABLE KEYS */;
INSERT INTO `ent_plan` VALUES ('bh1fn2d5jj84f1ljc1ug','bgrg80l5jj83bqe154fg',1,1,1,1,1,0,'2019-02-18 10:27:21.000000','2019-01-19 10:27:21.000000','2019-01-19 10:27:21.000000');
/*!40000 ALTER TABLE `ent_plan` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `enterprise`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `enterprise` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '企(事)业单位名称',
  `full_name` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `province` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `city` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `avatar` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `industry` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '所属行业',
  `location` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `address` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `website` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
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
  `contact_mobile` varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `contact_email` varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `contact_qq` varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `contact_wechat` varchar(40) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `contact_signature` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
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
INSERT INTO `enterprise` VALUES ('bgrg80l5jj83bqe154fg','test ent 003','','','','','','','','','2550418657@qq.com','18868905690','','2019-01-10 16:36:18.683882','2550418657@qq.com',1,1,2,'2019-02-09 16:36:18.683882','2019-01-10 16:36:18.683882','','','','','');
/*!40000 ALTER TABLE `enterprise` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `evaluation`
--

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
INSERT INTO `login_records` VALUES ('bgrg8ul5jj83bqe15570','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-10 16:38:18.205491','web','::1','PostmanRuntime/7.4.0'),('bgs10lt5jj8560mqq6rg','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-11 11:41:11.495297','web','::1','PostmanRuntime/7.4.0'),('bgu512l5jj8bqof636n0','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-14 17:04:10.480458','web','::1','PostmanRuntime/7.4.0'),('bgu572d5jj8bto51fb40','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-14 17:16:57.542707','web','::1','PostmanRuntime/7.4.0'),('bgulbud5jj8di574rqp0','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-15 11:39:37.623998','web','::1','PostmanRuntime/7.4.0'),('bgvk6td5jj83u5psg2jg','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-16 22:45:09.043456','web','::1','PostmanRuntime/7.6.0'),('bh01eal5jj859dlhvjs0','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-17 13:48:26.409852','web','::1','PostmanRuntime/7.6.0'),('bh1a93t5jj83rpeo4ie0','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-19 12:16:15.174783','web','::1','PostmanRuntime/7.6.0'),('bh1gn6t5jj84kc015ps0','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-19 19:35:55.826631','web','::1','PostmanRuntime/7.6.0'),('bh2jnbl5jj86lef4dblg','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-21 11:25:34.295456','web','::1','PostmanRuntime/7.6.0'),('bh2oe4d5jj88c0rtbo40','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-21 16:47:13.116086','web','::1','PostmanRuntime/7.6.0'),('bh38p5d5jj8a2vro5h9g','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-22 11:23:01.568421','web','::1','PostmanRuntime/7.6.0'),('bh5e2r55jj81h9nbv2tg','bgrg80l5jj83bqe154h0','bgrg80l5jj83bqe154fg','2019-01-25 18:14:04.540642','web','::1','PostmanRuntime/7.6.0');
/*!40000 ALTER TABLE `login_records` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `message`
--

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
-- Table structure for table `personal_config`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `personal_config` (
  `agent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `config_content` text COLLATE utf8mb4_unicode_ci,
  UNIQUE KEY `agent_id` (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `personal_config`
--

LOCK TABLES `personal_config` WRITE;
/*!40000 ALTER TABLE `personal_config` DISABLE KEYS */;
/*!40000 ALTER TABLE `personal_config` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `plan`
--

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
-- Table structure for table `promotion_msgs`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `promotion_msgs` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `enterprise_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `source` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci,
  `content_sdk` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `countdown` int(11) NOT NULL DEFAULT '0',
  `enabled` tinyint(1) NOT NULL DEFAULT '0',
  `summary` varchar(300) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `thumbnail` varchar(250) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `created_on` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_on` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`),
  KEY `idx_enterprise` (`enterprise_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `promotion_msgs`
--

LOCK TABLES `promotion_msgs` WRITE;
/*!40000 ALTER TABLE `promotion_msgs` DISABLE KEYS */;
/*!40000 ALTER TABLE `promotion_msgs` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `queue_config`
--

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

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `quickreply_group` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `ent_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '标题',
  `rank` int(11) NOT NULL DEFAULT '0',
  `created_by` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `creator_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `quickreply_group`
--

LOCK TABLES `quickreply_group` WRITE;
/*!40000 ALTER TABLE `quickreply_group` DISABLE KEYS */;
INSERT INTO `quickreply_group` VALUES ('bh1vldt5jj85k57df81g','bgrg80l5jj83bqe154fg','this is a test group',0,'bgrg80l5jj83bqe154h0','','2019-01-20 12:36:07.330720','2019-03-08 21:36:43.441080'),('bh1vm7t5jj85koc4s86g','bgrg80l5jj83bqe154fg','this is a test group',1000,'bgrg80l5jj83bqe154h0','','2019-01-20 12:37:51.949737','2019-03-08 21:36:43.441080'),('bh1vm9d5jj85koc4s870','bgrg80l5jj83bqe154fg','this is a test group111',1000,'bgrg80l5jj83bqe154h0','','2019-01-20 12:37:57.638688','2019-03-08 21:36:43.441080');
/*!40000 ALTER TABLE `quickreply_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `quickreply_item`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `quickreply_item` (
  `id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `quickreply_group_id` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `title` varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb4_unicode_ci,
  `content_type` varchar(50) COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '',
  `rich_content` text COLLATE utf8mb4_unicode_ci,
  `rank` int(11) NOT NULL DEFAULT '0',
  `created_by` char(20) COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  `updated_at` datetime(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `quickreply_item`
--

LOCK TABLES `quickreply_item` WRITE;
/*!40000 ALTER TABLE `quickreply_item` DISABLE KEYS */;
INSERT INTO `quickreply_item` VALUES ('bh1vob55jj85l0l434h0','bh1vm9d5jj85koc4s870','qk item','qk item 1','text',NULL,0,'bgrg80l5jj83bqe154h0','2019-01-20 12:42:20.737874','2019-01-20 12:42:20.737874'),('bh1vofd5jj85l0l434hg','bh1vm9d5jj85koc4s870','qk item','qk item 2','text',NULL,0,'bgrg80l5jj83bqe154h0','2019-01-20 12:42:37.507867','2019-01-20 12:42:37.507867'),('bh1vogd5jj85l0l434i0','bh1vm9d5jj85koc4s870','qk item','qk item 3','text',NULL,0,'bgrg80l5jj83bqe154h0','2019-01-20 12:42:41.608462','2019-01-20 12:42:41.608462'),('bh2ot7t5jj88gldi5if0','bgrg80l5jj83bqe154h0','qk item','qk item 3','text',NULL,0,'bgrg80l5jj83bqe154h0','2019-01-21 17:19:27.161973','2019-01-21 17:19:27.161973'),('bh3hnfd5jj8bjr1fe1t0','bgrg80l5jj83bqe154h0','qk item','qk item 3','text',NULL,0,'bgrg80l5jj83bqe154h0','2019-01-22 21:33:49.889209','2019-01-22 21:33:49.889209');
/*!40000 ALTER TABLE `quickreply_item` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role`
--

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
INSERT INTO `role` VALUES ('bgrg80l5jj83bqe154gg','bgrg80l5jj83bqe154fg','超管'),('bh02oit5jj859dlhvjsg','bgrg80l5jj83bqe154fg','test role'),('bh1a9vt5jj83rpeo4ieg','bgrg80l5jj83bqe154fg','test role'),('bh1aa3l5jj83rpeo4if0','bgrg80l5jj83bqe154fg','test role111');
/*!40000 ALTER TABLE `role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_perm`
--

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

-- Dump completed on 2019-03-15 17:40:51
