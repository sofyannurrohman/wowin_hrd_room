-- Migration 012: Seed modules and questions for testing
-- Run after all previous migrations.
-- Uses fixed UUIDs so this script is idempotent (safe to re-check).

-- ─── MODULES ──────────────────────────────────────────────────────────────
INSERT INTO modules (id, name, description, total_weight, created_by) VALUES
  (
    'a1000000-0000-0000-0000-000000000001',
    'Kemampuan Verbal',
    'Mengukur kemampuan memahami dan menggunakan bahasa secara efektif, termasuk sinonim, antonim, dan analogi kata.',
    100.00,
    '11111111-1111-1111-1111-111111111111'
  ),
  (
    'a2000000-0000-0000-0000-000000000002',
    'Kemampuan Numerik',
    'Mengukur kemampuan berhitung, penalaran angka, dan pemahaman pola matematika dasar.',
    100.00,
    '11111111-1111-1111-1111-111111111111'
  ),
  (
    'a3000000-0000-0000-0000-000000000003',
    'Kemampuan Logika & Analitis',
    'Mengukur kemampuan penalaran logis, deduksi, induksi, dan analisis situasi secara kritis.',
    100.00,
    '11111111-1111-1111-1111-111111111111'
  ),
  (
    'a4000000-0000-0000-0000-000000000004',
    'Tes Psikologi & Kepribadian',
    'Mengukur profil kepribadian, kecenderungan perilaku, dan kesesuaian dengan budaya kerja.',
    100.00,
    '11111111-1111-1111-1111-111111111111'
  )
ON CONFLICT (id) DO NOTHING;


-- ══════════════════════════════════════════════════════════════════════════
-- MODULE 1 · KEMAMPUAN VERBAL
-- ══════════════════════════════════════════════════════════════════════════

-- Q1 · Multiple Choice – Sinonim
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b1010000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001',
   'multiple_choice', 'Apa sinonim (persamaan kata) dari kata "EFISIEN"?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c1010001-0000-0000-0000-000000000001', 'b1010000-0000-0000-0000-000000000001', 'Boros',        FALSE),
  ('c1010002-0000-0000-0000-000000000001', 'b1010000-0000-0000-0000-000000000001', 'Hemat',        FALSE),
  ('c1010003-0000-0000-0000-000000000001', 'b1010000-0000-0000-0000-000000000001', 'Tepat Guna',   TRUE),  -- ← correct
  ('c1010004-0000-0000-0000-000000000001', 'b1010000-0000-0000-0000-000000000001', 'Lamban',       FALSE),
  ('c1010005-0000-0000-0000-000000000001', 'b1010000-0000-0000-0000-000000000001', 'Rumit',        FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q2 · Multiple Choice – Antonim
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b1020000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001',
   'multiple_choice', 'Apa antonim (lawan kata) dari "KONKRET"?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c1020001-0000-0000-0000-000000000001', 'b1020000-0000-0000-0000-000000000001', 'Nyata',      FALSE),
  ('c1020002-0000-0000-0000-000000000001', 'b1020000-0000-0000-0000-000000000001', 'Abstrak',    TRUE),  -- ← correct
  ('c1020003-0000-0000-0000-000000000001', 'b1020000-0000-0000-0000-000000000001', 'Pasti',      FALSE),
  ('c1020004-0000-0000-0000-000000000001', 'b1020000-0000-0000-0000-000000000001', 'Jelas',      FALSE),
  ('c1020005-0000-0000-0000-000000000001', 'b1020000-0000-0000-0000-000000000001', 'Terlihat',   FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q3 · Multiple Choice – Analogi
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b1030000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001',
   'multiple_choice', 'DOKTER : PASIEN = GURU : ?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c1030001-0000-0000-0000-000000000001', 'b1030000-0000-0000-0000-000000000001', 'Sekolah',    FALSE),
  ('c1030002-0000-0000-0000-000000000001', 'b1030000-0000-0000-0000-000000000001', 'Murid',      TRUE),  -- ← correct
  ('c1030003-0000-0000-0000-000000000001', 'b1030000-0000-0000-0000-000000000001', 'Buku',       FALSE),
  ('c1030004-0000-0000-0000-000000000001', 'b1030000-0000-0000-0000-000000000001', 'Kelas',      FALSE),
  ('c1030005-0000-0000-0000-000000000001', 'b1030000-0000-0000-0000-000000000001', 'Pensil',     FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q4 · True/False – EYD
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b1040000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001',
   'true_false', 'Penulisan "di mana" (dua kata) selalu benar dalam setiap konteks kalimat Bahasa Indonesia.', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c1040001-0000-0000-0000-000000000001', 'b1040000-0000-0000-0000-000000000001', 'Benar',  FALSE),
  ('c1040002-0000-0000-0000-000000000001', 'b1040000-0000-0000-0000-000000000001', 'Salah',  TRUE)   -- ← correct
ON CONFLICT (id) DO NOTHING;

-- Q5 · True/False – Sinonim
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b1050000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001',
   'true_false', '"AMBIGU" dan "MEMBINGUNGKAN" memiliki makna yang sama.', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c1050001-0000-0000-0000-000000000001', 'b1050000-0000-0000-0000-000000000001', 'Benar',  TRUE),  -- ← correct
  ('c1050002-0000-0000-0000-000000000001', 'b1050000-0000-0000-0000-000000000001', 'Salah',  FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q6 · Multiple Choice – Pemahaman bacaan
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b1060000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001',
   'multiple_choice',
   'Bacaan: "Perusahaan memutuskan untuk melakukan PHK massal karena kondisi pasar yang tidak menentu."\n\nKata "massal" pada kalimat di atas berarti...', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c1060001-0000-0000-0000-000000000001', 'b1060000-0000-0000-0000-000000000001', 'Sedikit',              FALSE),
  ('c1060002-0000-0000-0000-000000000001', 'b1060000-0000-0000-0000-000000000001', 'Besar-besaran',        TRUE),  -- ← correct
  ('c1060003-0000-0000-0000-000000000001', 'b1060000-0000-0000-0000-000000000001', 'Tiba-tiba',            FALSE),
  ('c1060004-0000-0000-0000-000000000001', 'b1060000-0000-0000-0000-000000000001', 'Berulang-ulang',       FALSE),
  ('c1060005-0000-0000-0000-000000000001', 'b1060000-0000-0000-0000-000000000001', 'Terencana',            FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q7 · Short Answer – Needs HR Review
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b1070000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001',
   'short_answer',
   'Buatlah sebuah kalimat efektif menggunakan kata "INOVASI" dan "BERKELANJUTAN" yang menggambarkan visi sebuah perusahaan modern!',
   2.0, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Q8 · Multiple Choice – Kalimat efektif
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b1080000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001',
   'multiple_choice',
   'Manakah kalimat yang PALING EFEKTIF di bawah ini?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c1080001-0000-0000-0000-000000000001', 'b1080000-0000-0000-0000-000000000001', 'Para karyawan-karyawan sedang rapat.',             FALSE),
  ('c1080002-0000-0000-0000-000000000001', 'b1080000-0000-0000-0000-000000000001', 'Karyawan sedang mengadakan rapat.',                TRUE),  -- ← correct
  ('c1080003-0000-0000-0000-000000000001', 'b1080000-0000-0000-0000-000000000001', 'Rapat diadakan oleh para karyawan-karyawan.',      FALSE),
  ('c1080004-0000-0000-0000-000000000001', 'b1080000-0000-0000-0000-000000000001', 'Para karyawan mengadakan rapat-rapat.',            FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q9 · True/False – Logika bahasa
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b1090000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001',
   'true_false', '"Saya sudah pernah makan nasi goreng" dan "Saya pernah makan nasi goreng" memiliki arti yang identik.', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c1090001-0000-0000-0000-000000000001', 'b1090000-0000-0000-0000-000000000001', 'Benar',  TRUE),  -- ← correct
  ('c1090002-0000-0000-0000-000000000001', 'b1090000-0000-0000-0000-000000000001', 'Salah',  FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q10 · Multiple Choice – Peribahasa
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b1100000-0000-0000-0000-000000000001', 'a1000000-0000-0000-0000-000000000001',
   'multiple_choice', 'Apa makna peribahasa "Bagai air di daun talas"?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c1100001-0000-0000-0000-000000000001', 'b1100000-0000-0000-0000-000000000001', 'Sangat setia dan konsisten',                   FALSE),
  ('c1100002-0000-0000-0000-000000000001', 'b1100000-0000-0000-0000-000000000001', 'Pendirian yang tidak tetap / mudah berubah',   TRUE),  -- ← correct
  ('c1100003-0000-0000-0000-000000000001', 'b1100000-0000-0000-0000-000000000001', 'Orang yang sangat dermawan',                  FALSE),
  ('c1100004-0000-0000-0000-000000000001', 'b1100000-0000-0000-0000-000000000001', 'Pekerjaan yang mudah dilakukan',              FALSE),
  ('c1100005-0000-0000-0000-000000000001', 'b1100000-0000-0000-0000-000000000001', 'Sesuatu yang terasa menyenangkan',            FALSE)
ON CONFLICT (id) DO NOTHING;


-- ══════════════════════════════════════════════════════════════════════════
-- MODULE 2 · KEMAMPUAN NUMERIK
-- ══════════════════════════════════════════════════════════════════════════

-- Q1 · MC – Aritmatika
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b2010000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002',
   'multiple_choice', 'Hasil dari 15% × 240 adalah...', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c2010001-0000-0000-0000-000000000002', 'b2010000-0000-0000-0000-000000000002', '24',  FALSE),
  ('c2010002-0000-0000-0000-000000000002', 'b2010000-0000-0000-0000-000000000002', '36',  TRUE),   -- ← correct
  ('c2010003-0000-0000-0000-000000000002', 'b2010000-0000-0000-0000-000000000002', '30',  FALSE),
  ('c2010004-0000-0000-0000-000000000002', 'b2010000-0000-0000-0000-000000000002', '40',  FALSE),
  ('c2010005-0000-0000-0000-000000000002', 'b2010000-0000-0000-0000-000000000002', '48',  FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q2 · MC – Deret angka
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b2020000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002',
   'multiple_choice', 'Lanjutan deret: 2, 6, 18, 54, ...', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c2020001-0000-0000-0000-000000000002', 'b2020000-0000-0000-0000-000000000002', '108', FALSE),
  ('c2020002-0000-0000-0000-000000000002', 'b2020000-0000-0000-0000-000000000002', '162', TRUE),   -- ← correct (×3)
  ('c2020003-0000-0000-0000-000000000002', 'b2020000-0000-0000-0000-000000000002', '144', FALSE),
  ('c2020004-0000-0000-0000-000000000002', 'b2020000-0000-0000-0000-000000000002', '180', FALSE),
  ('c2020005-0000-0000-0000-000000000002', 'b2020000-0000-0000-0000-000000000002', '216', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q3 · MC – Rasio
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b2030000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002',
   'multiple_choice', 'Perbandingan umur Andi dan Budi adalah 3:5. Jika umur Budi 25 tahun, berapa umur Andi?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c2030001-0000-0000-0000-000000000002', 'b2030000-0000-0000-0000-000000000002', '12 tahun', FALSE),
  ('c2030002-0000-0000-0000-000000000002', 'b2030000-0000-0000-0000-000000000002', '15 tahun', TRUE),   -- ← correct
  ('c2030003-0000-0000-0000-000000000002', 'b2030000-0000-0000-0000-000000000002', '18 tahun', FALSE),
  ('c2030004-0000-0000-0000-000000000002', 'b2030000-0000-0000-0000-000000000002', '20 tahun', FALSE),
  ('c2030005-0000-0000-0000-000000000002', 'b2030000-0000-0000-0000-000000000002', '10 tahun', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q4 · True/False – Fakta matematika
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b2040000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002',
   'true_false', 'Bilangan prima adalah bilangan yang hanya habis dibagi oleh 1 dan dirinya sendiri, dan angka 1 termasuk bilangan prima.', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c2040001-0000-0000-0000-000000000002', 'b2040000-0000-0000-0000-000000000002', 'Benar', FALSE),
  ('c2040002-0000-0000-0000-000000000002', 'b2040000-0000-0000-0000-000000000002', 'Salah', TRUE)   -- ← correct (1 bukan prima)
ON CONFLICT (id) DO NOTHING;

-- Q5 · MC – Persentase
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b2050000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002',
   'multiple_choice', 'Sebuah laptop dijual dengan harga Rp8.000.000 setelah diskon 20%. Berapa harga aslinya?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c2050001-0000-0000-0000-000000000002', 'b2050000-0000-0000-0000-000000000002', 'Rp9.000.000',  FALSE),
  ('c2050002-0000-0000-0000-000000000002', 'b2050000-0000-0000-0000-000000000002', 'Rp10.000.000', TRUE),   -- ← correct
  ('c2050003-0000-0000-0000-000000000002', 'b2050000-0000-0000-0000-000000000002', 'Rp9.600.000',  FALSE),
  ('c2050004-0000-0000-0000-000000000002', 'b2050000-0000-0000-0000-000000000002', 'Rp9.800.000',  FALSE),
  ('c2050005-0000-0000-0000-000000000002', 'b2050000-0000-0000-0000-000000000002', 'Rp8.800.000',  FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q6 · MC – Rata-rata
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b2060000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002',
   'multiple_choice', 'Nilai rata-rata 5 siswa adalah 78. Jika ditambahkan satu siswa baru dengan nilai 90, berapakah nilai rata-rata sekarang?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c2060001-0000-0000-0000-000000000002', 'b2060000-0000-0000-0000-000000000002', '80',  TRUE),   -- ← correct: (5×78+90)/6=480/6=80
  ('c2060002-0000-0000-0000-000000000002', 'b2060000-0000-0000-0000-000000000002', '82',  FALSE),
  ('c2060003-0000-0000-0000-000000000002', 'b2060000-0000-0000-0000-000000000002', '84',  FALSE),
  ('c2060004-0000-0000-0000-000000000002', 'b2060000-0000-0000-0000-000000000002', '79',  FALSE),
  ('c2060005-0000-0000-0000-000000000002', 'b2060000-0000-0000-0000-000000000002', '81',  FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q7 · True/False – Pecahan
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b2070000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002',
   'true_false', '3/4 lebih besar dari 7/10.', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c2070001-0000-0000-0000-000000000002', 'b2070000-0000-0000-0000-000000000002', 'Benar', TRUE),   -- ← correct: 0.75 > 0.70
  ('c2070002-0000-0000-0000-000000000002', 'b2070000-0000-0000-0000-000000000002', 'Salah', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q8 · MC – Kecepatan
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b2080000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002',
   'multiple_choice', 'Sebuah mobil menempuh 120 km dalam 1,5 jam. Berapa kecepatan rata-ratanya (km/jam)?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c2080001-0000-0000-0000-000000000002', 'b2080000-0000-0000-0000-000000000002', '60 km/jam',  FALSE),
  ('c2080002-0000-0000-0000-000000000002', 'b2080000-0000-0000-0000-000000000002', '80 km/jam',  TRUE),   -- ← correct
  ('c2080003-0000-0000-0000-000000000002', 'b2080000-0000-0000-0000-000000000002', '90 km/jam',  FALSE),
  ('c2080004-0000-0000-0000-000000000002', 'b2080000-0000-0000-0000-000000000002', '100 km/jam', FALSE),
  ('c2080005-0000-0000-0000-000000000002', 'b2080000-0000-0000-0000-000000000002', '75 km/jam',  FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q9 · Short Answer (HR Review) – Analisis data
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b2090000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002',
   'short_answer',
   'Penjualan sebuah toko selama 4 bulan terakhir berturut-turut adalah: Rp12 juta, Rp15 juta, Rp11 juta, dan Rp18 juta. Hitunglah rata-rata penjualan per bulan dan jelaskan tren penjualannya secara singkat!',
   2.0, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Q10 · MC – Pola gambar / angka
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b2100000-0000-0000-0000-000000000002', 'a2000000-0000-0000-0000-000000000002',
   'multiple_choice', 'Jika pola: 1, 1, 2, 3, 5, 8, 13, ... angka ke-8 adalah...', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c2100001-0000-0000-0000-000000000002', 'b2100000-0000-0000-0000-000000000002', '18', FALSE),
  ('c2100002-0000-0000-0000-000000000002', 'b2100000-0000-0000-0000-000000000002', '20', FALSE),
  ('c2100003-0000-0000-0000-000000000002', 'b2100000-0000-0000-0000-000000000002', '21', TRUE),   -- ← correct (Fibonacci)
  ('c2100004-0000-0000-0000-000000000002', 'b2100000-0000-0000-0000-000000000002', '24', FALSE),
  ('c2100005-0000-0000-0000-000000000002', 'b2100000-0000-0000-0000-000000000002', '26', FALSE)
ON CONFLICT (id) DO NOTHING;


-- ══════════════════════════════════════════════════════════════════════════
-- MODULE 3 · KEMAMPUAN LOGIKA & ANALITIS
-- ══════════════════════════════════════════════════════════════════════════

-- Q1 · MC – Silogisme
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b3010000-0000-0000-0000-000000000003', 'a3000000-0000-0000-0000-000000000003',
   'multiple_choice',
   'Premis 1: Semua karyawan berprestasi mendapat bonus.\nPremis 2: Budi adalah karyawan berprestasi.\nKesimpulan yang tepat adalah...', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c3010001-0000-0000-0000-000000000003', 'b3010000-0000-0000-0000-000000000003', 'Budi mungkin mendapat bonus.',          FALSE),
  ('c3010002-0000-0000-0000-000000000003', 'b3010000-0000-0000-0000-000000000003', 'Budi pasti mendapat bonus.',            TRUE),   -- ← correct
  ('c3010003-0000-0000-0000-000000000003', 'b3010000-0000-0000-0000-000000000003', 'Budi tidak mendapat bonus.',           FALSE),
  ('c3010004-0000-0000-0000-000000000003', 'b3010000-0000-0000-0000-000000000003', 'Tidak dapat disimpulkan.',             FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q2 · True/False – Logika deduktif
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b3020000-0000-0000-0000-000000000003', 'a3000000-0000-0000-0000-000000000003',
   'true_false',
   '"Jika semua A adalah B, dan semua B adalah C, maka semua A adalah C." Pernyataan ini BENAR secara logika.', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c3020001-0000-0000-0000-000000000003', 'b3020000-0000-0000-0000-000000000003', 'Benar', TRUE),   -- ← correct
  ('c3020002-0000-0000-0000-000000000003', 'b3020000-0000-0000-0000-000000000003', 'Salah', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q3 · MC – Urutan kejadian
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b3030000-0000-0000-0000-000000000003', 'a3000000-0000-0000-0000-000000000003',
   'multiple_choice',
   'Ana lebih tua dari Beni. Cindy lebih muda dari Ana tetapi lebih tua dari Beni. Siapa yang paling muda?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c3030001-0000-0000-0000-000000000003', 'b3030000-0000-0000-0000-000000000003', 'Ana',    FALSE),
  ('c3030002-0000-0000-0000-000000000003', 'b3030000-0000-0000-0000-000000000003', 'Beni',   TRUE),   -- ← correct
  ('c3030003-0000-0000-0000-000000000003', 'b3030000-0000-0000-0000-000000000003', 'Cindy',  FALSE),
  ('c3030004-0000-0000-0000-000000000003', 'b3030000-0000-0000-0000-000000000003', 'Tidak dapat ditentukan', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q4 · MC – Sebab-akibat
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b3040000-0000-0000-0000-000000000003', 'a3000000-0000-0000-0000-000000000003',
   'multiple_choice',
   'Penjualan turun 30% karena kualitas produk menurun. Langkah PERTAMA yang paling logis dilakukan manajemen adalah...', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c3040001-0000-0000-0000-000000000003', 'b3040000-0000-0000-0000-000000000003', 'Meningkatkan anggaran iklan',                      FALSE),
  ('c3040002-0000-0000-0000-000000000003', 'b3040000-0000-0000-0000-000000000003', 'Mengurangi harga jual produk',                     FALSE),
  ('c3040003-0000-0000-0000-000000000003', 'b3040000-0000-0000-0000-000000000003', 'Mengidentifikasi penyebab penurunan kualitas',      TRUE),   -- ← correct
  ('c3040004-0000-0000-0000-000000000003', 'b3040000-0000-0000-0000-000000000003', 'Mengganti seluruh tim pemasaran',                  FALSE),
  ('c3040005-0000-0000-0000-000000000003', 'b3040000-0000-0000-0000-000000000003', 'Menutup lini produksi',                            FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q5 · True/False – Logika negasi
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b3050000-0000-0000-0000-000000000003', 'a3000000-0000-0000-0000-000000000003',
   'true_false',
   'Negasi dari pernyataan "Semua karyawan hadir" adalah "Tidak ada karyawan yang hadir".', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c3050001-0000-0000-0000-000000000003', 'b3050000-0000-0000-0000-000000000003', 'Benar', FALSE),
  ('c3050002-0000-0000-0000-000000000003', 'b3050000-0000-0000-0000-000000000003', 'Salah', TRUE)   -- ← correct: negasi = "Ada karyawan yang tidak hadir"
ON CONFLICT (id) DO NOTHING;

-- Q6 · MC – Analogi logis
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b3060000-0000-0000-0000-000000000003', 'a3000000-0000-0000-0000-000000000003',
   'multiple_choice',
   'PANAS : API = DINGIN : ?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c3060001-0000-0000-0000-000000000003', 'b3060000-0000-0000-0000-000000000003', 'Angin', FALSE),
  ('c3060002-0000-0000-0000-000000000003', 'b3060000-0000-0000-0000-000000000003', 'Es',    TRUE),   -- ← correct
  ('c3060003-0000-0000-0000-000000000003', 'b3060000-0000-0000-0000-000000000003', 'Salju', FALSE),
  ('c3060004-0000-0000-0000-000000000003', 'b3060000-0000-0000-0000-000000000003', 'Air',   FALSE),
  ('c3060005-0000-0000-0000-000000000003', 'b3060000-0000-0000-0000-000000000003', 'Malam', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q7 · Short Answer – Analisa kasus (HR Review)
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b3070000-0000-0000-0000-000000000003', 'a3000000-0000-0000-0000-000000000003',
   'short_answer',
   'Suatu hari, Anda menemukan bahwa rekan kerja Anda secara rutin pulang lebih awal tanpa laporan kepada atasan. Hal ini berdampak pada beban kerja tim. Bagaimana cara Anda menghadapi situasi ini secara profesional?',
   2.0, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Q8 · MC – Pemecahan masalah
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b3080000-0000-0000-0000-000000000003', 'a3000000-0000-0000-0000-000000000003',
   'multiple_choice',
   'Ada 5 kotak: merah, biru, hijau, kuning, hitam. Kotak merah di sebelah kiri kotak biru. Kotak hijau tidak bersebelahan dengan merah. Kotak kuning di antara hitam dan biru. Manakah urutan yang mungkin (kiri ke kanan)?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c3080001-0000-0000-0000-000000000003', 'b3080000-0000-0000-0000-000000000003', 'Hijau – Merah – Biru – Kuning – Hitam', FALSE),
  ('c3080002-0000-0000-0000-000000000003', 'b3080000-0000-0000-0000-000000000003', 'Hijau – Hitam – Kuning – Biru – Merah', FALSE),
  ('c3080003-0000-0000-0000-000000000003', 'b3080000-0000-0000-0000-000000000003', 'Hijau – Merah – Biru – Hitam – Kuning', FALSE),
  ('c3080004-0000-0000-0000-000000000003', 'b3080000-0000-0000-0000-000000000003', 'Merah – Biru – Kuning – Hitam – Hijau', FALSE),
  ('c3080005-0000-0000-0000-000000000003', 'b3080000-0000-0000-0000-000000000003', 'Hitam – Kuning – Merah – Biru – Hijau', TRUE)   -- ← correct
ON CONFLICT (id) DO NOTHING;

-- Q9 · True/False – Deduksi
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b3090000-0000-0000-0000-000000000003', 'a3000000-0000-0000-0000-000000000003',
   'true_false',
   'Jika "P menyebabkan Q" benar dan "Q tidak terjadi", maka dapat disimpulkan "P tidak terjadi".', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c3090001-0000-0000-0000-000000000003', 'b3090000-0000-0000-0000-000000000003', 'Benar', TRUE),   -- ← correct (modus tollens)
  ('c3090002-0000-0000-0000-000000000003', 'b3090000-0000-0000-0000-000000000003', 'Salah', FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q10 · MC – Prioritas keputusan
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b3100000-0000-0000-0000-000000000003', 'a3000000-0000-0000-0000-000000000003',
   'multiple_choice',
   'Anda memiliki 4 tugas hari ini dengan tenggat waktu berbeda:\n- Laporan klien (deadline 2 jam)\n- Meeting internal (1 jam lagi)\n- Email follow-up (besok)\n- Proposal baru (3 hari lagi)\n\nUrutan prioritas yang paling tepat adalah...', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c3100001-0000-0000-0000-000000000003', 'b3100000-0000-0000-0000-000000000003', 'Proposal → Meeting → Laporan → Email',     FALSE),
  ('c3100002-0000-0000-0000-000000000003', 'b3100000-0000-0000-0000-000000000003', 'Meeting → Laporan → Email → Proposal',     TRUE),   -- ← correct (by urgency)
  ('c3100003-0000-0000-0000-000000000003', 'b3100000-0000-0000-0000-000000000003', 'Laporan → Proposal → Meeting → Email',     FALSE),
  ('c3100004-0000-0000-0000-000000000003', 'b3100000-0000-0000-0000-000000000003', 'Email → Meeting → Laporan → Proposal',     FALSE),
  ('c3100005-0000-0000-0000-000000000003', 'b3100000-0000-0000-0000-000000000003', 'Laporan → Meeting → Email → Proposal',     FALSE)
ON CONFLICT (id) DO NOTHING;


-- ══════════════════════════════════════════════════════════════════════════
-- MODULE 4 · TES PSIKOLOGI & KEPRIBADIAN
-- (All questions require_manual_review = TRUE as answers are subjective)
-- ══════════════════════════════════════════════════════════════════════════

-- Q1 · Psychological (HR Review)
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b4010000-0000-0000-0000-000000000004', 'a4000000-0000-0000-0000-000000000004',
   'psychological',
   'Dalam situasi kerja yang penuh tekanan dan tenggat waktu ketat, bagaimana Anda biasanya mengelola stres dan tetap produktif?',
   3.0, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Q2 · Multiple Choice – Preferensi kerja
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b4020000-0000-0000-0000-000000000004', 'a4000000-0000-0000-0000-000000000004',
   'multiple_choice',
   'Ketika Anda dihadapkan pada konflik dengan rekan kerja, pendekatan PERTAMA yang Anda lakukan adalah...', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c4020001-0000-0000-0000-000000000004', 'b4020000-0000-0000-0000-000000000004', 'Menghindari konflik dan tidak membicarakannya',               FALSE),
  ('c4020002-0000-0000-0000-000000000004', 'b4020000-0000-0000-0000-000000000004', 'Langsung melaporkan ke atasan',                               FALSE),
  ('c4020003-0000-0000-0000-000000000004', 'b4020000-0000-0000-0000-000000000004', 'Berbicara langsung dan terbuka dengan rekan kerja tersebut',   TRUE),   -- ← preferred
  ('c4020004-0000-0000-0000-000000000004', 'b4020000-0000-0000-0000-000000000004', 'Meminta rekan lain untuk menjadi perantara',                  FALSE),
  ('c4020005-0000-0000-0000-000000000004', 'b4020000-0000-0000-0000-000000000004', 'Diam dan mengalah agar situasi tidak memburuk',               FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q3 · Psychological (HR Review)
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b4030000-0000-0000-0000-000000000004', 'a4000000-0000-0000-0000-000000000004',
   'psychological',
   'Ceritakan sebuah situasi di mana Anda harus bekerja dalam tim yang memiliki perbedaan pendapat yang signifikan. Apa peran Anda dan bagaimana hasilnya?',
   3.0, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Q4 · True/False – Self-awareness
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b4040000-0000-0000-0000-000000000004', 'a4000000-0000-0000-0000-000000000004',
   'true_false', 'Seseorang dengan kecerdasan emosional tinggi selalu menghindari konfrontasi dalam situasi apapun.', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c4040001-0000-0000-0000-000000000004', 'b4040000-0000-0000-0000-000000000004', 'Benar', FALSE),
  ('c4040002-0000-0000-0000-000000000004', 'b4040000-0000-0000-0000-000000000004', 'Salah', TRUE)   -- ← correct: EI = manage, not avoid
ON CONFLICT (id) DO NOTHING;

-- Q5 · Multiple Choice – Motivasi
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b4050000-0000-0000-0000-000000000004', 'a4000000-0000-0000-0000-000000000004',
   'multiple_choice',
   'Apa faktor yang paling memotivasi Anda dalam bekerja?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c4050001-0000-0000-0000-000000000004', 'b4050000-0000-0000-0000-000000000004', 'Gaji dan benefit finansial semata',                                     FALSE),
  ('c4050002-0000-0000-0000-000000000004', 'b4050000-0000-0000-0000-000000000004', 'Pengakuan dan apresiasi dari rekan serta atasan',                        FALSE),
  ('c4050003-0000-0000-0000-000000000004', 'b4050000-0000-0000-0000-000000000004', 'Kesempatan untuk terus belajar dan berkembang secara profesional',        TRUE),   -- ← preferred
  ('c4050004-0000-0000-0000-000000000004', 'b4050000-0000-0000-0000-000000000004', 'Memiliki jabatan yang prestisius',                                       FALSE),
  ('c4050005-0000-0000-0000-000000000004', 'b4050000-0000-0000-0000-000000000004', 'Lingkungan kerja yang selalu menyenangkan',                              FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q6 · Psychological (HR Review)
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b4060000-0000-0000-0000-000000000004', 'a4000000-0000-0000-0000-000000000004',
   'psychological',
   'Jika Anda diminta memimpin sebuah proyek baru dengan tim yang belum pernah Anda kenal, apa yang akan Anda lakukan dalam 3 hari pertama?',
   3.0, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Q7 · True/False – Teamwork
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b4070000-0000-0000-0000-000000000004', 'a4000000-0000-0000-0000-000000000004',
   'true_false', 'Dalam tim yang efektif, setiap anggota harus memiliki tanggung jawab yang sama persis tanpa pembagian peran khusus.', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c4070001-0000-0000-0000-000000000004', 'b4070000-0000-0000-0000-000000000004', 'Benar', FALSE),
  ('c4070002-0000-0000-0000-000000000004', 'b4070000-0000-0000-0000-000000000004', 'Salah', TRUE)   -- ← correct: specialization improves team performance
ON CONFLICT (id) DO NOTHING;

-- Q8 · Multiple Choice – Adaptasi
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b4080000-0000-0000-0000-000000000004', 'a4000000-0000-0000-0000-000000000004',
   'multiple_choice',
   'Perusahaan tiba-tiba mengubah kebijakan kerja dari WFH menjadi WFO penuh. Reaksi Anda adalah...', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c4080001-0000-0000-0000-000000000004', 'b4080000-0000-0000-0000-000000000004', 'Langsung menolak kebijakan tersebut',                                                    FALSE),
  ('c4080002-0000-0000-0000-000000000004', 'b4080000-0000-0000-0000-000000000004', 'Menerima perubahan dan mengatur jadwal agar tetap produktif',                             TRUE),   -- ← preferred
  ('c4080003-0000-0000-0000-000000000004', 'b4080000-0000-0000-0000-000000000004', 'Resign karena tidak sesuai ekspektasi',                                                  FALSE),
  ('c4080004-0000-0000-0000-000000000004', 'b4080000-0000-0000-0000-000000000004', 'Mengabaikan kebijakan selama tidak ada sanksi',                                          FALSE),
  ('c4080005-0000-0000-0000-000000000004', 'b4080000-0000-0000-0000-000000000004', 'Menunggu rekan kerja lain bereaksi terlebih dahulu sebelum memutuskan sikap',            FALSE)
ON CONFLICT (id) DO NOTHING;

-- Q9 · Psychological (HR Review)
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b4090000-0000-0000-0000-000000000004', 'a4000000-0000-0000-0000-000000000004',
   'psychological',
   'Sebutkan 3 kelebihan dan 3 kekurangan diri Anda yang relevan dengan posisi yang Anda lamar, serta rencana Anda untuk mengatasi kekurangan tersebut.',
   3.0, TRUE)
ON CONFLICT (id) DO NOTHING;

-- Q10 · Multiple Choice – Integritas
INSERT INTO questions (id, module_id, type, content, weight, requires_manual_review) VALUES
  ('b4100000-0000-0000-0000-000000000004', 'a4000000-0000-0000-0000-000000000004',
   'multiple_choice',
   'Anda mengetahui bahwa rekan kerja Anda melakukan kesalahan kecil dalam laporan keuangan yang belum diketahui atasan. Apa yang Anda lakukan?', 1.0, FALSE)
ON CONFLICT (id) DO NOTHING;

INSERT INTO question_options (id, question_id, content, is_correct) VALUES
  ('c4100001-0000-0000-0000-000000000004', 'b4100000-0000-0000-0000-000000000004', 'Diam karena bukan urusan saya',                                                                   FALSE),
  ('c4100002-0000-0000-0000-000000000004', 'b4100000-0000-0000-0000-000000000004', 'Langsung melaporkan ke atasan tanpa memberi tahu rekan terlebih dahulu',                          FALSE),
  ('c4100003-0000-0000-0000-000000000004', 'b4100000-0000-0000-0000-000000000004', 'Memberi tahu rekan kerja agar segera diperbaiki sebelum dilaporkan ke atasan',                    TRUE),   -- ← preferred
  ('c4100004-0000-0000-0000-000000000004', 'b4100000-0000-0000-0000-000000000004', 'Mengancam rekan kerja agar memberikan imbalan atas kerahasiaannya',                               FALSE),
  ('c4100005-0000-0000-0000-000000000004', 'b4100000-0000-0000-0000-000000000004', 'Membetulkan laporan sendiri tanpa memberitahu siapapun',                                          FALSE)
ON CONFLICT (id) DO NOTHING;
