from sqlalchemy.orm import Session
from app.models.job import Job
from sqlalchemy import or_, func

# ini adalah class untuk melakukan interaksi ke database
class JobRepository:
    # melakukan inisial
    def __init__(self, db : Session):
        self.db = db

    def get_all_jobs(self):
        return self.db.query(Job).all()

    # method ini untuk mengambil data jobstreet, karna jobstreet kosong skillsnya
    def get_unprocessed_jobs(self):
        return self.db.query(Job).filter(
            Job.is_active == True,
            or_(
                Job.is_it == None,
                Job.skills == None,
                func.cardinality(Job.skills) == 0,
                Job.softskills == None,
                func.cardinality(Job.softskills) == 0,
                Job.embedding == None,
            )
        ).all()

    # method ini untuk melakukan update
    def update_is_it(self, job_id : str, is_it : bool):
        job = self.db.query(Job).filter(Job.id == job_id).first()

        if job:
            job.is_it = is_it
            self.db.commit()
            self.db.refresh(job)
        return job

    # method ini untuk melakukan update
    def update_job_skills_and_embedding(self, session, job_id : str, is_it : bool, skills : list, softskills : list, joblevel : str, is_ok : bool, embedding_vector : list):
        try:
            job = session.query(Job).filter(Job.id == job_id).first()

            if job:
                job.softskills = softskills
                job.is_it = is_it
                job.skills = skills
                job.joblevel = joblevel
                job.is_ok = is_ok
                job.embedding = embedding_vector
                
                session.commit()
                return True
                
            return False
        except Exception as e:
            session.rollback()
            print(f"[REPO ERROR] Gagal update Job ID {job_id}: {e}")
            raise e
    
    def delete_non_it(self):
        try:
            jobs = self.db.query(Job).filter(Job.is_it == False)

            count = jobs.count()

            if count > 0:
                jobs.delete(synchronize_session=False)
                self.db.commit()
            return count
        except Exception as e:
            self.db.rollback()
            raise e