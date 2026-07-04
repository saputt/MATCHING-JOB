import os
from google import genai
import json
from openai import OpenAI
from dotenv import load_dotenv

load_dotenv()

class LLMAdapter:
    def __init__(self):
        or_api_key = os.getenv("OR_API")
        self.or_client = OpenAI(
            base_url="https://openrouter.ai/api/v1",
            api_key=or_api_key
        )
        self.or_model_name = 'openrouter/owl-alpha'

        om_api_key = os.getenv("OM_API")
        self.om_client = OpenAI(
            base_url="https://api.openmodel.ai",
            api_key=om_api_key
        )
        self.om_model_name = 'deepseek-v4-flash'

    def extract_skills_from_desc(self, title : str, description : str, skills : list) -> dict:
        prompt = f"""
            Kamu adalah seorang Tech HRD, Data Architect, & Technical Recruiter Expert tingkat senior.
            Tugasmu adalah menganalisis secara mendalam data lowongan kerja berikut:
            Judul Pekerjaan: {title}
            Deskripsi Pekerjaan: {description}
            Skills_Mentah: {skills}

            Tugas Utamamu:
            1. KATEGORI IT (is_it): Tentukan apakah lowongan ini termasuk bidang Teknologi Informasi (IT) atau sejenisnya (true/false).
            
            2. LEVEL PEKERJAAN (job_level): Tentukan senioritas lowongan ini ke dalam salah satu dari 3 level berikut secara ketat berdasarkan judul dan deskripsi:
            - "Entry-Level / Internship": Jika ramah pemula, mencari anak magang, junior, fresh graduate, atau syarat pengalaman 0-1 tahun.
            - "Mid-Level": Jika mencari level menengah, syarat pengalaman 2-4 tahun, atau role mandiri tanpa memimpin tim besar.
            - "Senior-Level / Lead": Jika mencari level ahli, senior, lead engineer, software architect, atau syarat pengalaman 5 tahun ke atas.

            3. LOGIKA EVALUASI HYBRID (Deep Inspection):
            - JIKA `Skills_Mentah` KOSONG `[]`: Baca penuh `Deskripsi Pekerjaan` dan ekstrak semua hard_skills & soft_skills dari nol.
            3. LOGIKA EVALUASI HYBRID (Deep Inspection):
               - JIKA `Skills_Mentah` KOSONG `[]`: Baca penuh `Deskripsi Pekerjaan` dan ekstrak semua hard_skills & soft_skills dari nol.
               - JIKA `Skills_Mentah` ADA ISINYA (Tidak Kosong): 
                 a. Lakukan "Double-Check Filtering" (Penyaringan Ganda). `Skills_Mentah` yang dikirim dari database BUKANLAH DATA SUCI, melainkan data kotor yang wajib kamu evaluasi ulang.
                 b. Periksa setiap elemen di dalam `Skills_Mentah`, lalu WAJIB BUANG kata-kata yang melanggar aturan di Poin 4 (seperti menghapus "Web App", "Natural Language Processing", dll). Jangan ditelan mentah-mentah!
                 c. Tetap baca `Deskripsi Pekerjaan` untuk memburu/mencari apakah ada hard_skills atau soft_skills penting lainnya yang belum tercantum, lalu gabungkan ke dalam output setelah disaring ketat.
                d. Konsolidasi Akhir: Gabungkan hasil pembersihan data mentah dengan temuan baru dari deskripsi, lalu kembalikan array yang 100% bersih, ringkas, murni nama teknologi/konsep baku, dan bebas dari kata-kata sampah generik.
            4. STANDAR KAMAR SKILL (BATASAN KETAT):
               - hard_skills: Bahasa pemrograman, framework, database, konsep arsitektur IT, tools, atau metodologi teknis profesional bidang khusus. Standardisasi penulisan (Contoh: "golang" -> "Go", "reactjs" -> "React", "restapi" -> "REST API").
                 * ATURAN GARIS MIRING (/): Jika menemukan dua teknologi yang digabung menggunakan tanda garis miring atau kata "and/or" (Contoh: "Excel / Google Sheets", "MySQL / PostgreSQL"), WAJIB KAMU PECAH menjadi dua elemen terpisah di dalam array (Contoh menjadi: ["Excel", "Google Sheets"]). Jangan digabung dalam satu string!
                 * ATURAN SINONIM & AKRONIM: Jika menemukan istilah teknis panjang yang memiliki singkatan baku standar industri, wajib konversikan ke bentuk singkatannya saja (Contoh: "Object Oriented Programming" -> "OOP", "Search Engine Optimization" -> "SEO", "Cloud Computing" -> "Cloud").
                 * ATURAN HARAM: JANGAN PERNAH memasukkan jenis platform/output produk generik seperti: "Web App", "Web Application", "Mobile App", "Desktop App", "Software", "Application", "iOS", "Android" ke dalam list skill.
                 * ATURAN REDUNDANSI: Jika menemukan nama panjang dan singkatan, AMBIL VERSI SINGKATAN STANDAR INDUSTRI SAJA. (Contoh: "Natural Language Processing" dan "NLP" -> Cukup ambil "NLP". "Speech-to-Text" dan "Voice Recognition" -> Cukup ambil "Speech-to-Text").
               - soft_skills: Persyaratan karakter, kemampuan interpersonal, bahasa kalbu, atau keahlian fungsional kantoran yang terlalu umum (Contoh: "Administration", "Positive Attitude", "Good Communication", "Leadership").

            CONTOH KASUS 1 (Mode Ekstraksi Penuh - Lowongan IT Entry-Level, Skills Mentah Kosong):
            Input Title: "Junior Web Developer (Fresh Graduate Welcome)"
            Input Skills_Mentah: []
            Input Desc: "We are looking for an entry-level web developer. You will help us maintain websites using javascript and php. Must be a fast learner with a positive attitude. 0-1 years of experience or internship experience is acceptable."
            Output JSON Sah:
            {{
                "is_it": true,
                "job_level": "Entry-Level / Internship",
                "hard_skills": ["JavaScript", "PHP"],
                "soft_skills": ["Fast Learner", "Positive Attitude"]
            }}

            CONTOH KASUS 2 (Mode Ekstraksi Penuh - Lowongan IT Senior-Level, Skills Mentah Kosong):
            Input Title: "Lead Software Architect"
            Input Skills_Mentah: []
            Input Desc: "Seeking a Senior Lead Architect to design high-throughput microservices. Expert-level knowledge in Golang, Docker, and Kubernetes is mandatory. You must have at least 6 years of experience and proven leadership skills to manage a team of 10 engineers."
            Output JSON Sah:
            {{
                "is_it": true,
                "job_level": "Senior-Level / Lead",
                "hard_skills": ["Go", "Docker", "Kubernetes", "Microservices"],
                "soft_skills": ["Leadership"]
            }}

            CONTOH KASUS 3 (Mode Hybrid - Lowongan IT Mid-Level, Skills Mentah Sudah Ada):
            Input Title: "Backend Engineer (NodeJS)"
            Input Skills_Mentah: ["nodejs", "express", "postgresql"]
            Input Desc: "Requirements: 3 years of experience in backend development. Strong problem solving skills and good communication are required to sync with the product team."
            Output JSON Sah:
            {{
                "is_it": true,
                "job_level": "Mid-Level",
                "hard_skills": ["Node.js", "Express", "PostgreSQL"],
                "soft_skills": ["Problem Solving", "Good Communication"]
            }}

            CONTOH KASUS 4 (Mode Hybrid - Lowongan Murni Non-IT, Skills Mentah Sudah Ada):
            Input Title: "Senior Textile Designer"
            Input Skills_Mentah: ["textile design", "administration"]
            Input Desc: "Looking for a senior designer with 5+ years experience in garment factories. Leadership skills required."
            Output JSON Sah:
            {{
                "is_it": false,
                "job_level": "Senior-Level / Lead",
                "hard_skills": ["Textile Design"],
                "soft_skills": ["Administration", "Leadership"]
            }}

            CONTOH KASUS 5 (Pembersihan Teks Generik & Redundansi - Kasus Khusus Deep Inspection):
            Input Title: "AI Developer Intern"
            Input Skills_Mentah: ["python", "laravel"]
            Input Desc: "We need an intern to build Web App and Mobile App. Must understand Natural Language Processing / NLP and create Voice Assistant. Good communication is a must."
            Output JSON Sah:
            {{
                "is_it": true,
                "job_level": "Entry-Level / Internship",
                "hard_skills": ["Python", "Laravel", "NLP", "Voice Assistant"],
                "soft_skills": ["Good Communication"]
            }}
            *(Penjelasan: 'Web App' & 'Mobile App' HARAM dimasukkan karena produk generik. 'Natural Language Processing' dibuang karena sudah diwakili oleh singkatan bakunya yaitu 'NLP')*

            Aturan Output:
            Wajib kembalikan response HANYA dalam format JSON mentah murni tanpa tanda ```json maupun teks pembuka/penutup apapun:
            {{
                "is_it": true/false,
                "job_level": "Entry-Level / Internship" atau "Mid-Level" atau "Senior-Level / Lead",
                "hard_skills": ["NamaHardSkill1", "NamaHardSkill2"],
                "soft_skills": ["NamaSoftSkill1", "NamaSoftSkill2"]
            }}
        """
        return self._call_llm(client=self.or_client, model=self.or_model_name, prompt=prompt)
    
    def analyze_job_match(self, job_title : str, job_skills : list, user_skills : list, match_score : float) -> dict :
        prompt = f"""
            Kamu adalah seorang IT Career Mentor sekaligus Senior Technical Recruiter yang berpengalaman merekrut Software Engineer, Backend Engineer, Frontend Engineer, Fullstack Engineer, Data Engineer, dan AI Engineer.

            Tugasmu adalah memberikan analisis profesional terhadap kecocokan kandidat berdasarkan hasil semantic matching yang telah dihitung sebelumnya menggunakan embedding.

            =====================
            DATA INPUT
            =====================

            Target Posisi:
            {job_title}

            Job Skills:
            {job_skills}

            User Skills:
            {user_skills}

            Match Score (0-100):
            {match_score}

            KETERANGAN:
            - Match Score di atas SUDAH dihitung menggunakan semantic similarity embedding.
            - Angka tersebut merupakan tingkat kecocokan keseluruhan antara profil kandidat dan kebutuhan lowongan.
            - JANGAN menghitung ulang, mengubah, ataupun mengoreksi nilai Match Score tersebut.
            - Gunakan Match Score hanya sebagai referensi ketika membuat analisis.

            =====================
            TUGAS
            =====================

            1. MATCH SCORE
            - Gunakan nilai Match Score yang diberikan tanpa perubahan.
            - Jangan melakukan perhitungan ulang.

            2. MATCHED SKILLS
            - Identifikasi skill dari User Skills yang relevan atau memiliki makna yang sama dengan Job Skills.
            - Lakukan pencocokan secara case-insensitive.
            - Perbolehkan semantic matching sederhana.
            Contoh:
            - JS ≈ JavaScript
            - TS ≈ TypeScript
            - PostgreSQL ≈ Postgres
            - Express ≈ Express.js

            3. MISSING SKILLS
            - Identifikasi skill atau teknologi yang diminta pada Job Skills tetapi belum dimiliki kandidat.
            - Jangan memasukkan skill yang sudah dianggap cocok secara semantic.

            4. AI NARRATIVE
            - Buat evaluasi maksimal 3 kalimat.
            - Jelaskan arti Match Score tersebut.
            - Sebutkan kekuatan kandidat.
            - Sebutkan area yang masih perlu ditingkatkan.
            - Gunakan bahasa profesional tetapi tetap santai seperti mentor karier IT Indonesia.
            - Hindari kalimat yang terlalu berlebihan atau menjatuhkan kandidat.

            5. LEARNING ROADMAP
            - Susun roadmap belajar berdasarkan daftar missing_skills.
            - Minimal 3 langkah.
            - Urutkan dari fundamental menuju implementasi nyata.
            - Fokus pada skill yang memberikan dampak terbesar terhadap posisi tersebut.

            6. PROTIP
            - Berikan satu tips praktis agar kandidat tetap memiliki peluang besar dipanggil recruiter meskipun masih memiliki skill gap.
            - Contohnya:
                - project portfolio
                - GitHub
                - deployment
                - technical blog
                - open source
                - sertifikasi
                - optimasi CV

            =====================
            ATURAN OUTPUT
            =====================

            Kembalikan HANYA JSON mentah.
            Jangan gunakan markdown.
            Jangan gunakan ```json.
            Jangan menambahkan penjelasan apa pun.

            Format JSON wajib seperti berikut:

            {{
                "match_score": {match_score},
                "matched_skills": [
                    "React",
                    "Node.js"
                ],
                "missing_skills": [
                    "Docker",
                    "TypeScript"
                ],
                "ai_narrative": "Match score sebesar {match_score} menunjukkan kandidat memiliki tingkat kecocokan yang cukup baik terhadap posisi ini. Skill backend yang dimiliki sudah relevan, namun masih ada beberapa teknologi penting yang perlu dipelajari agar semakin kompetitif.",
                "learning_roadmap": [
                    "Pelajari dasar TypeScript lalu migrasikan salah satu project JavaScript ke TypeScript.",
                    "Pelajari Docker dan containerisasi aplikasi Node.js.",
                    "Deploy project yang sudah menggunakan Docker ke platform cloud seperti Render atau Railway."
                ],
                "protip": "Buat satu project end-to-end yang menggunakan seluruh tech stack utama pada lowongan dan tampilkan di GitHub beserta dokumentasi yang rapi."
            }}
            """
        return self._call_llm(client=self.om_client, model=self.om_model_name, prompt=prompt)

    def _call_llm(self, client, model, prompt) -> dict:
        try:
            response = client.chat.completions.create(
                model=model,
                messages=[
                    {"role": "system", "content": "You are a JSON-only API. Output raw JSON without markdown."},
                    {"role": "user", "content": prompt}
                ]
            )
            raw = response.choices[0].message.content.replace("```json", "").replace("```", "").strip()
            return json.loads(raw)
        except Exception as e:
            print(f"[LLM ERROR] Gagal di model {model}: {e}")
            return {}
