# file ini berfungsi sebagai jembatan ke market di database
from typing import List, Optional
from sqlalchemy.orm import Session
from app.models.job import Job

class MarketRepository:
    def __init__(self, db: Session):
        self.db = db
    
    # method ini untuk mengambil semua list jobs
    def get_all_jobs(self) -> List[Job]:
        return self.db.query(Job).all()
    
    # method ini berfungsi mengambil job berdasarkan id
    def get_job_by_id(self, job_id : str) -> Optional[Job]:
        return self.db.query(Job).filter(Job.id == job_id).first()
