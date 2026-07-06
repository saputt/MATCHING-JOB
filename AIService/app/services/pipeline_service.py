from app.repositories.job_repository import JobRepository
from app.adapters.llm_adapter import LLMAdapter
import time
from sqlalchemy.orm import Session, sessionmaker
from app.services.embedding_service import EmbeddingService
from concurrent.futures import ThreadPoolExecutor
import json

# ini adalah class mengatur pipeline 
class PipelineService:
    # ini method inisialisasi
    def __init__(self, repo : JobRepository, llm : LLMAdapter):
        self.repo = repo
        self.llm = llm
    
    def patch_all(self, db : Session):
        self.jobs_patcher(db=db)
        self.delete_non_it()

    def delete_non_it(self):
        """
        service ini untuk melakukan delete terhadap data non-it
        """
        jobs_deleted = self.repo.delete_non_it()
        return jobs_deleted

    # ini method untuk melakukan penambalan terhadap skills pengguna yang masih kosong dengan memanngil llm
    def jobs_patcher(self, db : Session):
        """
        service ini untuk melakukan penambalan terhadap atribut lowongan yang masing kosong seperti skills, is_it, joblevel
        """
        jobs_queue = self.repo.get_unprocessed_jobs()

        if not jobs_queue:
            print("[PIPELINE] Selesai! Tidak ditemukan lagi lowongan Jobstreet yang kosong jirr.")
            return
        
        print(f"[PIPELINE] Menemukan {len(jobs_queue)} lowongan yang memiliki atribut kosong. Mulai menambal...")
        
        MAX_WORKER = 3

        session_factory = sessionmaker(bind=db.get_bind())

        def worker_task(job):
            worker_session = session_factory()
            try:
                print(f" -> Sedang memproses Job ID: {job.id} | Title: {job.title} | Source: {job.source}")

                ai_result = self.llm.extract_skills_from_desc(
                    title=job.title,
                    description=job.description,    
                    skills=job.skills
                )

                embedding_text = EmbeddingService().generate_embedding(title=job.title, hard_skills=ai_result["hard_skills"])

                if isinstance(embedding_text, str):
                    embedding_text = json.loads(embedding_text)
                
                self.repo.update_job_skills_and_embedding(
                    job_id=job.id,
                    softskills=ai_result.get("soft_skills", []),
                    is_it=ai_result.get("is_it", False),
                    skills=ai_result["hard_skills"],
                    joblevel=ai_result["job_level"],
                    is_ok=bool(True),
                    embedding_vector=embedding_text,
                    session=worker_session
                )

                time.sleep(2)

            except Exception as e:
                print(f"   [GAGAL] Skip ID {job.id} karena error: {e}")
                worker_session.rollback()
            
            finally:
                worker_session.close()
        
        with ThreadPoolExecutor(max_workers=MAX_WORKER) as executor:
            executor.map(worker_task, jobs_queue)
        
        print("[PIPELINE SELESAI] Semua worker selesai tugasnya! Cek database")