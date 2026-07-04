from pydantic import BaseModel

class CreateAnalysisRequest(BaseModel):
    user_id : str
    job_title : str
    job_skills : list
    user_skills : list
    match_score : float

class GetAnalysisRequest(BaseModel):
    user_id : str

class EmbeddingRequest(BaseModel):
    hard_skills : list[str]
    title : str 