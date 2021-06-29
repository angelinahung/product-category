-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- 主機： 127.0.0.1
-- 產生時間： 2021-06-29 16:30:51
-- 伺服器版本： 10.4.19-MariaDB
-- PHP 版本： 8.0.7

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- 資料庫: `product_category`
--

-- --------------------------------------------------------

--
-- 資料表結構 `category`
--

CREATE TABLE `category` (
  `id` int(20) NOT NULL,
  `name` varchar(50) NOT NULL,
  `is_invisible` tinyint(1) NOT NULL DEFAULT 0,
  `parent_id` int(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='產品目錄';

--
-- 傾印資料表的資料 `category`
--

INSERT INTO `category` (`id`, `name`, `is_invisible`, `parent_id`) VALUES
(1, 'TEST', 1, 0),
(2, 'lin', 1, 0),
(3, 'ffff', 0, 1),
(4, 'ffffh', 1, 2),
(5, 'ffff', 0, 1),
(6, 'ffffh', 1, 2);

-- --------------------------------------------------------

--
-- 資料表結構 `product`
--

CREATE TABLE `product` (
  `id` int(20) NOT NULL,
  `name` varchar(50) NOT NULL,
  `budget` int(50) DEFAULT NULL,
  `price` int(50) DEFAULT NULL,
  `description` text NOT NULL,
  `is_sale` tinyint(1) NOT NULL DEFAULT 0,
  `start_sale_time` datetime NOT NULL,
  `end_sale_time` datetime NOT NULL,
  `category_id` int(20) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

--
-- 傾印資料表的資料 `product`
--

INSERT INTO `product` (`id`, `name`, `budget`, `price`, `description`, `is_sale`, `start_sale_time`, `end_sale_time`, `category_id`) VALUES
(1, 'aho', 100, 200, '3', 0, '2021-06-30 10:00:00', '2021-06-30 22:00:00', 3),
(2, 'angel', 100, 200, 'x', 1, '2021-06-30 10:00:00', '2021-06-30 22:00:00', 5),
(4, 'lin', 0, 0, 'x', 0, '2021-06-30 10:00:00', '2021-06-30 22:00:00', 4);

--
-- 已傾印資料表的索引
--

--
-- 資料表索引 `category`
--
ALTER TABLE `category`
  ADD PRIMARY KEY (`id`);

--
-- 資料表索引 `product`
--
ALTER TABLE `product`
  ADD PRIMARY KEY (`id`),
  ADD KEY `category_id` (`category_id`);

--
-- 已傾印資料表的限制式
--

--
-- 資料表的限制式 `product`
--
ALTER TABLE `product`
  ADD CONSTRAINT `product_ibfk_1` FOREIGN KEY (`category_id`) REFERENCES `category` (`id`);
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
