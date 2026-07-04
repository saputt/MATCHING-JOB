from sqlalchemy.orm import Session
from app.models.analysis import Analysis

class AnalyzeRepository:
    def __init__(self, db : Session):
        self.db = db
    
    def get_analyze_by_job_id(self, job_id : str, user_id : str) -> Analysis:
        return self.db.query(Analysis).filter(
            Analysis.job_id == job_id,
            Analysis.user_id == user_id
        ).first()

    def save(self, analysis_data : Analysis) -> Analysis:
        try:
            self.db.add(analysis_data)
            self.db.commit()
            self.db.refresh(analysis_data)
        except Exception as e:
            self.db.rollback()
            raise e

    