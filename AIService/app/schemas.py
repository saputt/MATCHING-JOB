# file ini bertujuan untuk mendefinisikan bentuk data yang masuk dan keluar

from pydantic import BaseModel
from typing import List, Optional

# data yang dikirim ke python untuk menganalisis market
class AnalyzeMarketRequest(BaseModel):
    user_id: str
    skills: List[str]

# data yang dikirim ke python ketika analisis job
class AnalyzeJobRequest(BaseModel):
    user_id: str
    job_id: str
    skills: List[str]    

# ini respon dari top_skill pada market insight response
class SkillFrequencySchema(BaseModel):
    skill: str
    frequncy: int

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
    skills_gap: SkillGap
    market_match_score: int
    market_narrative: str
    top_jobs: List[JobWithScoreSchema]

# ini respom yang diberikan ke fe, terkait hasil analisis job
class JobAnalysisResponse(BaseModel):
    job_id: str
    match_score: int
    matched_skills: List[str]
    missing_skills: List[str]
    ai_narrative: str
    learning_roadmap: List[str]
    pro_tip: str

