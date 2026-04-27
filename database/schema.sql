-- Drop tables if they exist to prevent errors on recreation
DROP TABLE IF EXISTS `verifikasi_penyewa`;
DROP TABLE IF EXISTS `pembayaran`;
DROP TABLE IF EXISTS `reservasi`;
DROP TABLE IF EXISTS `moderasi_kos`;
DROP TABLE IF EXISTS `kos`;
DROP TABLE IF EXISTS `penyewa`;
DROP TABLE IF EXISTS `pemilik_kos`;
DROP TABLE IF EXISTS `admin`;

-- 1. admin
CREATE TABLE `admin` (
  `id_admin` INT AUTO_INCREMENT PRIMARY KEY,
  `nama` VARCHAR(100) NOT NULL,
  `email` VARCHAR(100) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `nomor_telepon` VARCHAR(20)
);

-- 2. pemilik_kos
CREATE TABLE `pemilik_kos` (
  `id_pemilik` INT AUTO_INCREMENT PRIMARY KEY,
  `nama` VARCHAR(100) NOT NULL,
  `email` VARCHAR(100) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `nomor_telepon` VARCHAR(20),
  `alamat` TEXT
);

-- 3. penyewa
CREATE TABLE `penyewa` (
  `id_penyewa` INT AUTO_INCREMENT PRIMARY KEY,
  `nama` VARCHAR(100) NOT NULL,
  `email` VARCHAR(100) NOT NULL UNIQUE,
  `password` VARCHAR(255) NOT NULL,
  `no_telepon` VARCHAR(20),
  `alamat` TEXT
);

-- 4. kos
CREATE TABLE `kos` (
  `id_kos` INT AUTO_INCREMENT PRIMARY KEY,
  `id_pemilik` INT NOT NULL,
  `nama_kos` VARCHAR(150) NOT NULL,
  `alamat_kos` TEXT NOT NULL,
  `harga_sewa` DECIMAL(10, 2) NOT NULL,
  `deskripsi` TEXT,
  `fasilitas` TEXT,
  `status_kos` VARCHAR(50) DEFAULT 'Tersedia',
  FOREIGN KEY (`id_pemilik`) REFERENCES `pemilik_kos`(`id_pemilik`) ON DELETE CASCADE
);

-- 5. moderasi_kos
CREATE TABLE `moderasi_kos` (
  `id_moderasi` INT AUTO_INCREMENT PRIMARY KEY,
  `id_kos` INT NOT NULL,
  `status_moderasi` VARCHAR(50) NOT NULL,
  `pesan_admin` TEXT,
  `tanggal_moderasi` DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`id_kos`) REFERENCES `kos`(`id_kos`) ON DELETE CASCADE
);

-- 6. reservasi
CREATE TABLE `reservasi` (
  `id_reservasi` INT AUTO_INCREMENT PRIMARY KEY,
  `id_penyewa` INT NOT NULL,
  `id_kos` INT NOT NULL,
  `tanggal_reservasi` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `status_reservasi` VARCHAR(50) DEFAULT 'Pending',
  FOREIGN KEY (`id_penyewa`) REFERENCES `penyewa`(`id_penyewa`) ON DELETE CASCADE,
  FOREIGN KEY (`id_kos`) REFERENCES `kos`(`id_kos`) ON DELETE CASCADE
);

-- 7. pembayaran
CREATE TABLE `pembayaran` (
  `id_pembayaran` INT AUTO_INCREMENT PRIMARY KEY,
  `id_reservasi` INT NOT NULL,
  `metode_pembayaran` VARCHAR(50) NOT NULL,
  `jumlah_bayar` DECIMAL(10, 2) NOT NULL,
  `tanggal_bayar` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `status_pembayaran` VARCHAR(50) DEFAULT 'Pending',
  FOREIGN KEY (`id_reservasi`) REFERENCES `reservasi`(`id_reservasi`) ON DELETE CASCADE
);

-- 8. verifikasi_penyewa
CREATE TABLE `verifikasi_penyewa` (
  `id_verifikasi` INT AUTO_INCREMENT PRIMARY KEY,
  `id_reservasi` INT NOT NULL,
  `status_verifikasi` VARCHAR(50) DEFAULT 'Belum Terverifikasi',
  `tanggal_verifikasi` DATETIME DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (`id_reservasi`) REFERENCES `reservasi`(`id_reservasi`) ON DELETE CASCADE
);

-- ==========================================
-- SEEDER DATA
-- ==========================================

-- Insert Admin
INSERT INTO `admin` (`nama`, `email`, `password`, `nomor_telepon`) VALUES
('Admin Utama', 'admin@example.com', 'hashed_password_here', '081234567890');

-- Insert Pemilik Kos
INSERT INTO `pemilik_kos` (`nama`, `email`, `password`, `nomor_telepon`, `alamat`) VALUES
('Bapak Budi', 'budi@example.com', 'hashed_password_here', '081298765432', 'Jl. Sudirman No 1'),
('Ibu Siti', 'siti@example.com', 'hashed_password_here', '081311122233', 'Jl. Merdeka No 2');

-- Insert Penyewa
INSERT INTO `penyewa` (`nama`, `email`, `password`, `no_telepon`, `alamat`) VALUES
('Andi', 'andi@example.com', 'hashed_password_here', '085712341234', 'Jl. Ahmad Yani No 10'),
('Rina', 'rina@example.com', 'hashed_password_here', '085743214321', 'Jl. Pahlawan No 5');

-- Insert Kos
INSERT INTO `kos` (`id_pemilik`, `nama_kos`, `alamat_kos`, `harga_sewa`, `deskripsi`, `fasilitas`, `status_kos`) VALUES
(1, 'Kos Mawar (Putra)', 'Jl. Sudirman Gang 1', 1500000.00, 'Kos nyaman untuk putra', 'Kasur, Lemari, WiFi', 'Tersedia'),
(2, 'Kos Melati (Putri)', 'Jl. Merdeka Gang 2', 2000000.00, 'Kos eksklusif untuk putri', 'AC, Kamar Mandi Dalam, WiFi', 'Penuh');

-- Insert Moderasi Kos
INSERT INTO `moderasi_kos` (`id_kos`, `status_moderasi`, `pesan_admin`, `tanggal_moderasi`) VALUES
(1, 'Disetujui', 'Kos sesuai dengan ketentuan.', '2023-10-01 10:00:00'),
(2, 'Disetujui', 'Kos eksklusif terverifikasi.', '2023-10-02 11:30:00');

-- Insert Reservasi
INSERT INTO `reservasi` (`id_penyewa`, `id_kos`, `tanggal_reservasi`, `status_reservasi`) VALUES
(1, 1, '2023-10-15 14:00:00', 'Disetujui'),
(2, 2, '2023-10-16 09:00:00', 'Pending');

-- Insert Pembayaran
INSERT INTO `pembayaran` (`id_reservasi`, `metode_pembayaran`, `jumlah_bayar`, `tanggal_bayar`, `status_pembayaran`) VALUES
(1, 'Transfer Bank', 1500000.00, '2023-10-15 14:30:00', 'Lunas'),
(2, 'E-Wallet', 2000000.00, '2023-10-16 09:15:00', 'Pending');

-- Insert Verifikasi Penyewa
INSERT INTO `verifikasi_penyewa` (`id_reservasi`, `status_verifikasi`, `tanggal_verifikasi`) VALUES
(1, 'Terverifikasi', '2023-10-15 15:00:00'),
(2, 'Belum Terverifikasi', '2023-10-16 09:30:00');
