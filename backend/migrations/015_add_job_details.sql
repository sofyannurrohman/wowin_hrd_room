-- Migration 015: Add description and requirements to job_positions
ALTER TABLE job_positions ADD COLUMN description TEXT DEFAULT '';
ALTER TABLE job_positions ADD COLUMN requirements TEXT DEFAULT '';

-- Update default positions with samples
UPDATE job_positions SET 
    description = 'Bertanggung jawab dalam pengembangan dan pemeliharaan aplikasi berbasis web menggunakan teknologi terbaru.',
    requirements = '- Pengalaman minimal 1 tahun di bidang SE\n- Menguasai minimal satu bahasa pemrograman (Go/NodeJS)\n- Memahami SQL & NoSQL'
WHERE name = 'Software Engineer';

UPDATE job_positions SET 
    description = 'Merancang antarmuka pengguna yang menarik dan intuitif untuk meningkatkan pengalaman pengguna aplikasi kami.',
    requirements = '- Portofolio desain UI/UX yang kuat\n- Menguasai Figma atau Adobe XD\n- Memahami prinsip-prinsip desain modern'
WHERE name = 'UI/UX Designer';

UPDATE job_positions SET 
    description = 'Mengelola siklus pengembangan produk dari ide hingga peluncuran.',
    requirements = '- Pengalaman minimal 2 tahun sebagai PM\n- Kemampuan komunikasi & analisis yang baik'
WHERE name = 'Product Manager';
