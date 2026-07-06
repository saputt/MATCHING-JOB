from app.repositories.analyze_repository import AnalyzeRepository
from app.adapters.llm_adapter import LLMAdapter
from app.models.analysis import Analysis
from app.schemas import CreateAnalysisRequest
import uuid

class AnalyzeService:
    def __init__(self, repo : AnalyzeRepository, llm : LLMAdapter):
        self.repo = repo
        self.llm = llm

    def analyze_job(self, req : CreateAnalysisRequest, job_id : str, user_id : str):
        analyze_exist = self.repo.get_analyze_by_job_id(job_id=job_id, user_id=user_id)
        if not analyze_exist:
            ai_result = self.llm.analyze_job_match(
                job_title=req.job_title, 
                job_hard_skills=req.job_hard_skills,
                job_soft_skills=req.job_soft_skills, 
                user_skills=req.user_skills,
                match_score=req.match_score
            )

            if not ai_result:
                raise Exception("failed to get response from llm")
            
            new_analyze : Analysis = Analysis(
                job_id=job_id,
                user_id=user_id,
                match_score=req.match_score,
                matched_skills=ai_result.get("matched_skills", []),
                missing_skills=ai_result.get("missing_skills", []),
                ai_narrative=ai_result.get("ai_narrative", ""),
                learning_roadmap=ai_result.get("learning_roadmap", []),
                protip=ai_result.get("protip", "")
            )

            print(new_analyze)

            return self.repo.save(analysis_data=new_analyze)
        
        return analyze_exist

    def get_analyze_by_job_id(self, job_id : str, user_id : str):
        analyze_job = self.repo.get_analyze_by_job_id(job_id=job_id, user_id=user_id)

        return analyze_job
