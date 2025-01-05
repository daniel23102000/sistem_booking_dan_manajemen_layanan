-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Waktu pembuatan: 05 Jan 2025 pada 03.47
-- Versi server: 10.4.32-MariaDB
-- Versi PHP: 8.0.30

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `booking_system`
--

-- --------------------------------------------------------

--
-- Struktur dari tabel `bookings`
--

CREATE TABLE `bookings` (
  `booking_id` int(11) NOT NULL,
  `user_id` int(11) DEFAULT NULL,
  `room_id` int(11) DEFAULT NULL,
  `check_in_date` date DEFAULT NULL,
  `check_out_date` date DEFAULT NULL,
  `total_price` decimal(10,2) DEFAULT NULL,
  `status` enum('pending','confirmed','cancelled','completed') DEFAULT 'pending',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `booking_services`
--

CREATE TABLE `booking_services` (
  `booking_service_id` int(11) NOT NULL,
  `booking_id` int(11) DEFAULT NULL,
  `service_id` int(11) DEFAULT NULL,
  `quantity` int(11) DEFAULT 1,
  `total_price` decimal(10,2) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `payments`
--

CREATE TABLE `payments` (
  `payment_id` int(11) NOT NULL,
  `booking_id` int(11) DEFAULT NULL,
  `payment_method` enum('credit_card','debit_card','paypal','cash') DEFAULT NULL,
  `amount` decimal(10,2) DEFAULT NULL,
  `payment_date` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `properties`
--

CREATE TABLE `properties` (
  `property_id` int(11) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `address` text DEFAULT NULL,
  `description` text DEFAULT NULL,
  `contact_number` varchar(20) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `rooms`
--

CREATE TABLE `rooms` (
  `room_id` int(11) NOT NULL,
  `property_id` int(11) DEFAULT NULL,
  `room_name` varchar(50) DEFAULT NULL,
  `room_type` enum('single','double','suite','family') DEFAULT NULL,
  `price_per_night` decimal(10,2) DEFAULT NULL,
  `status` enum('available','booked','maintenance') DEFAULT 'available'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `services`
--

CREATE TABLE `services` (
  `service_id` int(11) NOT NULL,
  `property_id` int(11) DEFAULT NULL,
  `service_name` varchar(100) DEFAULT NULL,
  `price` decimal(10,2) DEFAULT NULL,
  `description` text DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Struktur dari tabel `users`
--

CREATE TABLE `users` (
  `user_id` int(11) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL,
  `password_hash` varchar(255) DEFAULT NULL,
  `phone_number` varchar(20) DEFAULT NULL,
  `role` enum('customer','staff','admin') DEFAULT 'customer',
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Indexes for dumped tables
--

--
-- Indeks untuk tabel `bookings`
--
ALTER TABLE `bookings`
  ADD PRIMARY KEY (`booking_id`),
  ADD KEY `user_id` (`user_id`),
  ADD KEY `room_id` (`room_id`);

--
-- Indeks untuk tabel `booking_services`
--
ALTER TABLE `booking_services`
  ADD PRIMARY KEY (`booking_service_id`),
  ADD KEY `booking_id` (`booking_id`),
  ADD KEY `service_id` (`service_id`);

--
-- Indeks untuk tabel `payments`
--
ALTER TABLE `payments`
  ADD PRIMARY KEY (`payment_id`),
  ADD KEY `booking_id` (`booking_id`);

--
-- Indeks untuk tabel `properties`
--
ALTER TABLE `properties`
  ADD PRIMARY KEY (`property_id`);

--
-- Indeks untuk tabel `rooms`
--
ALTER TABLE `rooms`
  ADD PRIMARY KEY (`room_id`),
  ADD KEY `property_id` (`property_id`);

--
-- Indeks untuk tabel `services`
--
ALTER TABLE `services`
  ADD PRIMARY KEY (`service_id`),
  ADD KEY `property_id` (`property_id`);

--
-- Indeks untuk tabel `users`
--
ALTER TABLE `users`
  ADD PRIMARY KEY (`user_id`),
  ADD UNIQUE KEY `email` (`email`);

--
-- AUTO_INCREMENT untuk tabel yang dibuang
--

--
-- AUTO_INCREMENT untuk tabel `bookings`
--
ALTER TABLE `bookings`
  MODIFY `booking_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `booking_services`
--
ALTER TABLE `booking_services`
  MODIFY `booking_service_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `payments`
--
ALTER TABLE `payments`
  MODIFY `payment_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `properties`
--
ALTER TABLE `properties`
  MODIFY `property_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `rooms`
--
ALTER TABLE `rooms`
  MODIFY `room_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `services`
--
ALTER TABLE `services`
  MODIFY `service_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT untuk tabel `users`
--
ALTER TABLE `users`
  MODIFY `user_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- Ketidakleluasaan untuk tabel pelimpahan (Dumped Tables)
--

--
-- Ketidakleluasaan untuk tabel `bookings`
--
ALTER TABLE `bookings`
  ADD CONSTRAINT `bookings_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
  ADD CONSTRAINT `bookings_ibfk_2` FOREIGN KEY (`room_id`) REFERENCES `rooms` (`room_id`);

--
-- Ketidakleluasaan untuk tabel `booking_services`
--
ALTER TABLE `booking_services`
  ADD CONSTRAINT `booking_services_ibfk_1` FOREIGN KEY (`booking_id`) REFERENCES `bookings` (`booking_id`),
  ADD CONSTRAINT `booking_services_ibfk_2` FOREIGN KEY (`service_id`) REFERENCES `services` (`service_id`);

--
-- Ketidakleluasaan untuk tabel `payments`
--
ALTER TABLE `payments`
  ADD CONSTRAINT `payments_ibfk_1` FOREIGN KEY (`booking_id`) REFERENCES `bookings` (`booking_id`);

--
-- Ketidakleluasaan untuk tabel `rooms`
--
ALTER TABLE `rooms`
  ADD CONSTRAINT `rooms_ibfk_1` FOREIGN KEY (`property_id`) REFERENCES `properties` (`property_id`) ON DELETE CASCADE;

--
-- Ketidakleluasaan untuk tabel `services`
--
ALTER TABLE `services`
  ADD CONSTRAINT `services_ibfk_1` FOREIGN KEY (`property_id`) REFERENCES `properties` (`property_id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
