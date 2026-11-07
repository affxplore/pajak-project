-- ====================================================================
-- Skema Database PostgreSQL untuk Aplikasi Suara Rakyat: Kritik Pajak
-- ====================================================================

-- 1. Buat Tabel comments
CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50) NOT NULL DEFAULT 'Umum',
    comment TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Index untuk optimasi pengurutan data berdasarkan waktu terbaru
CREATE INDEX IF NOT EXISTS idx_comments_created_at ON comments(created_at DESC);

-- 3. Data Awal (Dummy / Seed Data)
INSERT INTO comments (name, category, comment, created_at) VALUES
('Rina S.', 'Pajak Pembalut', 'Pembalut itu barang primer, menstruasi bukan pilihan! Sangat tidak etis kalau kebutuhan biologis bulanan wanita dibebani PPN.', NOW() - INTERVAL '2 hours'),
('Budi K.', 'Sektor Informal (Cukur/Pijat)', 'Tukang cukur dan pijat keliling itu kerja keras sambung hidup. Seharusnya dibantu permodalannya, bukan malah diincar pajaknya.', NOW() - INTERVAL '1 hour'),
('Dewi A.', 'PPN 12%', 'Daya beli sedang turun dan banyak PHK, kenaikan PPN 12% semakin mendongkrak harga kebutuhan harian keluarga.', NOW() - INTERVAL '30 minutes')
ON CONFLICT DO NOTHING;

