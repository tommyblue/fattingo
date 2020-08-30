CREATE TABLE `customers` (
    `id` int NOT NULL AUTO_INCREMENT,
    `user_id` int DEFAULT NULL,
    `title` varchar(255) DEFAULT NULL,
    `name` varchar(255) DEFAULT NULL,
    `surname` varchar(255) DEFAULT NULL,
    `address` varchar(255) DEFAULT NULL,
    `zip_code` varchar(255) DEFAULT NULL,
    `town` varchar(255) DEFAULT NULL,
    `province` varchar(255) DEFAULT NULL,
    `country` varchar(255) DEFAULT NULL,
    `tax_code` varchar(255) DEFAULT NULL,
    `vat` varchar(255) DEFAULT NULL,
    `created_at` datetime NOT NULL,
    `updated_at` datetime NOT NULL,
    `info` text,
    PRIMARY KEY (`id`),
    UNIQUE KEY `index_customers_on_title_and_user_id` (`title`, `user_id`),
    KEY `index_customers_on_title` (`title`),
    KEY `index_customers_on_name` (`name`),
    KEY `index_customers_on_surname` (`surname`),
    KEY `index_customers_on_user_id` (`user_id`)
) ENGINE = InnoDB AUTO_INCREMENT = 39 DEFAULT CHARSET = latin1;
