from typing import List, Dict
from collections import Counter
from app.features.market.repository import MarketRepository
from app.features.market.schema import CalculateScoreGapResult, MarketMatchSkill, TopJobResult, MarketInsightResponse, SkillGapSchema, JobWithScoreSchema

class MarketService:
    def __init__(self, repo: MarketRepository):
        self.repo = repo
    
    # method ini adalah method utama, yang akan dipanggil di router
    def analyze_market(self, user_skills: List[str], limit : int = 5) -> MarketInsightResponse:
        # memanggil internal method untuk menganalisis market fit sesuai dengan skill pengguna
        market_fit_analysis = self.analyze_skill_gap_with_market(user_skills)
        # memanggil internal method untuk mengambil top job berdasarkan skill pengguna
        top_jobs_recommendation = self.get_top_jobs(user_skills)
        # memanggil internal method untuk mengambil top skills yang ada di market
        top_skills_in_market = self.get_top_skills()

        # membuat inisialisasi untuk memasukkan job dengan schema yang sudah ditentukan
        job_with_schema = []

        # melakukan looping sebanyakan jumlah job recomendation, dan memasukkannya kedalam list yang sudah di inisialisasi
        for item in top_jobs_recommendation:
            # mempersingkat pemanggilan terhadap detail job
            job = item.job
            # memasukkan data kedalam list
            job_with_schema.append(JobWithScoreSchema(
                city=job.city,
                company=job.company,
                is_remote=job.is_remote,
                job_id=job.id,
                location=job.location,
                match_score=item.match_score,
                matched_skills=item.matched_skills,
                missing_skills=item.missed_skills,
                title=job.title,
                url=job.url
            ))

        return MarketInsightResponse(
            market_match_score=market_fit_analysis.match_score,
            skills_gap=SkillGapSchema(
                match_score=market_fit_analysis.match_score,
                skill_missing=market_fit_analysis.missed_skills,
                skill_owned=market_fit_analysis.matched_skills
            ),
            top_jobs=job_with_schema,
            market_narrative="",
            top_skills=top_skills_in_market
        )

    # method ini berfungsi untuk mendapatkan top skills yang sering muncul di market
    def get_top_skills(self, limit : int = 5) -> List[Dict[str, int]]:
        # melakukan query ke tabel jobs melalui repository
        jobs = self.repo.get_all_jobs()

        # jika tidak ada jobs sama sekali di database, kembalikan array kosong
        if not jobs:
            return []

        skill_counter = Counter()

        # melakukan looping, terhadap semua jobs
        for job in jobs:
            if job.skills is not None:
                # melakukan looping terhadap skills di setiap job yand didapatkan
                for skill in job.skills:
                    # jika job tidak kosong tambahkan counter 
                    if skill:
                        skill_counter[skill] += 1
        
        # menghitung total jobs yang tersedia di database
        total_jobs = len(jobs)

        # membuat list kosong, yang nantinya akan diisi dengan dictionary dari hasil looping analisis persentase skill
        result = []

        # mengambil skill yang paling banyak muncul, dan melakukan interasi
        for skill, count in skill_counter.most_common(limit):
            # menghitung persentase kemunculan dari total jobs yang ada
            percentage = round((count / total_jobs) * 100)

            # memasukkan hasil dari perhitugan kemunculan skill kedalam list yang sudah dibuat sebelumnya
            result.append({
                "skill" : skill,
                "frequency" : percentage
            })
        
        return result
    
    # fungsi ini untuk menghitung skor kecocokan dan gap dari skill user dan skill di job
    def calculate_score_gap(self, user_skills: List[str], job_skills: List[str]) -> CalculateScoreGapResult:
        # jika job skills kosong, kembalikan dict karna tidak ada yang cocok
        if not job_skills:
            return CalculateScoreGapResult(match_score=0, matched_skills=[], missed_skills=[])
        
        # membuat user set dan job set agar bisa dilakukan operasi matematika seperti intersection
        user_set = set([s.lower() for s in user_skills])
        job_set = set([j.lower() for j in job_skills])

        # menghitung job yang cocok antara skill user dan skill yang ada di job
        matched_skill = user_set.intersection(job_set)
        missing_skill = job_set.difference(user_set)
        
        # menghitung score job
        match_score = round(len(matched_skill) / len(job_set) * 100)

        return CalculateScoreGapResult(
            match_score=match_score,
            matched_skills=matched_skill,
            missed_skills=missing_skill
        )
    
    # fungsi ini untuk menganalisis seberapa cocok skill user dengan trend pasar
    def analyze_skill_gap_with_market(self, user_skills: List[str]) -> MarketMatchSkill:
        # mengambil skill yang paling sering muncul
        top_skills = self.get_top_skills(limit=5)

        # jika top skills kosong kembalikan nilai inis
        if not top_skills:
            return MarketMatchSkill(
                match_score=0,
                matched_skills=[],
                missed_skills=[]
            )
        
        # hanya mengambil daftar skillnya saja
        top_skills_name = [item["skill"] for item in top_skills]

        # memanggil fungsi internal untuk menghitung gap dan score
        result = self.calculate_score_gap(user_skills, top_skills_name)

        return MarketMatchSkill(
            match_score=result.match_score,
            matched_skills=result.matched_skills,
            missed_skills=result.missed_skills
        )
    
    # method ini berfungsi untuk mengambil top job teratas yang sesuai dengan skill user
    def get_top_jobs(self, user_skills: List[str], limit: int = 10) -> List[TopJobResult]:
        # mengambil semua jobs
        jobs = self.repo.get_all_jobs()

        result = []

        # melakukan perulangan ke semua jobs
        for job in jobs:
            # melakukan analysis untuk mendapatkan matchscore dan gap skill
            analysis = self.calculate_score_gap(user_skills, job.skills)
            result.append(TopJobResult(
                job=job,
                match_score=analysis.match_score,
                matched_skills=analysis.matched_skills,
                missed_skills=analysis.missed_skills
            ))
        
        # mengurutkan berdasarkan matchscore tertinggi
        result.sort(key=lambda x : x.match_score, reverse=True)

        # mengembalikan hasil sebanyak limit
        return result[:limit]   