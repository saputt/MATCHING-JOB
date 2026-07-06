# file ini mendefiniskan model untuk tabel analysis
from sqlalchemy import Column, String, Boolean, Integer, DateTime, Text, Index
from sqlalchemy.dialects.postgresql import UUID, ARRAY
from sqlalchemy.sql import func
from app.core.database import Base
import uuid

class Analysis(Base):
    # ini merepresentasikan nama tabel
    __tablename__ = "analyses"

    # ini skema tabelny
    id = Column(String, primary_key=True, default=uuid.uuid1)
    match_score = Column(Integer, nullable=False)
    matched_skills = Column(ARRAY(String), default=[])
    missing_skills = Column(ARRAY(String), default=[])
    ai_narrative = Column(String, nullable=False)
    learning_roadmap = Column(ARRAY(String), default=[])
    protip = Column(Text, nullable=True)
    job_id = Column(String, nullable=False)
    user_id = Column(String, nullable=False)
    created_at = Column(DateTime(timezone=True), server_default=func.now())
    
    # untuk melakukan indexing
    __table_args__ = (
        Index("idx_analysis_user_job", user_id, job_id, unique=True),
    )
    def __repr__(self):
        return f"<Analysis(user_id={self.user_id}, job_id={self.job_id}, match_score={self.match_score})>"