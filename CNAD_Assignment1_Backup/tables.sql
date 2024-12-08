Create database user_management_db;

Create table users (
    id int auto_increment primary key,
    name varchar(255) not null unique,
    email varchar(255) not null unique,
    password varchar(255) not null,
    membership_tier ENUM('Basic','Premium','VIP') not null
);    
    
CREATE TABLE membership_benefits (
    tier ENUM('Basic', 'Premium', 'VIP') PRIMARY KEY,
    discount_rate DECIMAL(5,2) NOT NULL, -- discount in percentage (e.g., 10.00 for 10%)
    increased_limit INT NOT NULL,         -- increased booking limit for the tier
    priority_access BOOLEAN NOT NULL      -- true/false for priority vehicle access
);

INSERT INTO membership_benefits (tier, discount_rate, increased_limit, priority_access)
VALUES
    ('Basic', 0.05, 5, FALSE),
    ('Premium', 0.10, 10, TRUE),
    ('VIP', 0.15, 20, TRUE);


Create database vehicle_reservation_db;

Create table vehicles (
    id int auto_increment primary key,
    make varchar(255),
    model varchar(255),
    availability Boolean
)

CREATE TABLE reservations (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    vehicle_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    status ENUM('active', 'cancelled') DEFAULT 'active',
    FOREIGN KEY (vehicle_id) REFERENCES vehicles(id),
);

CREATE TABLE rental_history (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    vehicle_id INT NOT NULL,
    rental_start DATE NOT NULL,
    rental_end DATE NOT NULL,
    total_amount DECIMAL(10,2) NOT NULL
);


Create database billingpayment_db;

Create table promotions (
    id int auto_increment primary key,
    code varchar(255),
    discount_percentage int,
    expiration_date DATE not null
);

Create table billing (
    id int auto_increment primary key,
    reservation_id int not null,
    amount DECIMAL(10,2) not null,
    payment_status ENUM('Pending','Paid') not null
)

CREATE TABLE vehicle_pricing (
    id INT AUTO_INCREMENT PRIMARY KEY,
    vehicle_type VARCHAR(50) NOT NULL,
    base_rate_per_hour DECIMAL(10, 2) NOT NULL,
    discount_basic DECIMAL(5, 2) DEFAULT 0.00,  -- Percentage discount for Basic members
    discount_premium DECIMAL(5, 2) DEFAULT 10.00,  -- Percentage discount for Premium members
    discount_vip DECIMAL(5, 2) DEFAULT 20.00  -- Percentage discount for VIP members
);
