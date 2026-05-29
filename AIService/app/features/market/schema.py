# file ini bertujuan untuk mendefinisikan bentuk data yang masuk dan keluar
from pydantic import BaseModel, ConfigDict
from typing import List, Optional
from app.models.job import Job

# data yang dikirim ke python untuk menganalisis market
class AnalyzeMarketRequest(BaseModel):
    skills: List[str] 

# ini respon dari top_skill pada market insight response
class SkillFrequencySchema(BaseModel):
    skill: str
    frequency: int

# ini respon dari atribut skillgap pada market insight response
class SkillGapSchema(BaseModel):
    skill_owned: List[str]
    skill_missing: List[str]
    match_score: int

# ini respon dari atribut top score pada market insight response
class JobWithScoreSchema(BaseModel):
    job_id: str
    title: str
    company: str
    location: str
    is_remote: bool
    city: Optional[str]
    match_score: int
    matched_skills: List[str]
    missing_skills: List[str]
    url: str

# ini respon yang diberikan ke FE, terkait hasil analisis market insight
class MarketInsightResponse(BaseModel):
    top_skills: List[SkillFrequencySchema]
    skills_gap: SkillGapSchema
    market_match_score: int
    market_narrative: str
    top_jobs: List[JobWithScoreSchema]

class CalculateScoreGapResult(BaseModel):
    match_score : int
    matched_skills : List[str]
    missed_skills : List[str] 

class TopJobResult(CalculateScoreGapResult):
    model_config=ConfigDict(arbitrary_types_allowed=True)
    job : Job

class MarketMatchSkill(CalculateScoreGapResult):
    pass
    