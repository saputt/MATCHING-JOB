# file ini mendefiniskan model untuk tabel job
from sqlalchemy import Column, String, Boolean, Integer, DateTime, Text, Index
from sqlalchemy.dialects.postgresql import UUID, ARRAY
from sqlalchemy.sql import func
from app.core.database import Base
import uuid

# mendefinisikan model pada database
class Job(Base):
    # ini merepresentasikan nama table
    # nama table harus sama seperti di database
    __tablename__ = "jobs"

    # ini skema tabelnya
    id = Column(String, primary_key=True, default=uuid.uuid1)
    title = Column(String, nullable=False)
    company = Column(String, nullable=False)
    description = Column(Text, nullable=True)
    url = Column(String, nullable=False, index=True)
    location = Column(String, nullable=False)
    is_it = Column(Boolean, nullable=True)
    is_active = Column(Boolean, default=True)
    city = Column(String, nullable=True)
    source = Column(String, default="glints")
    salary = Column(String, nullable=False)
    skills = Column(ARRAY(String), default=[])
    softskills = Column(ARRAY(String), default=[])
    scraped_at = Column(DateTime(timezone=True), server_default=func.now())
    joblevel = Column(String, nullable=True)
    is_ok = Column(Boolean, default=False)
    embedding = Column(String, nullable=True)

    #ini untuk indexing agar pencarian cepat
    __table_args__ = (
        Index("idx_job_title_remote", title, company, unique=True),
    )

    def __repr__ (self):
        return f"<Job(id={self.id}, title={self.title}, company={self.company})>"