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
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid1)
    user_id = Column(UUID(as_uuid=True), nullable=False)
    title = Column(String, nullable=False)
    company = Column(String, nullable=False)
    description = Column(Text, nullable=True)
    url = Column(String, nullable=False, unique=True, index=True)
    location = Column(String, nullable=False)
    is_remote = Column(Boolean, default=False)
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid1)
    city = Column(String, nullable=True)
    source = Column(String, default="glints")
    skills = Column(ARRAY(String), default=[])
    scraped_at = Column(DateTime(timezone=True), server_default=func.now())
    
    #ini untuk indexing agar pencarian cepat
    __table_args__ = (
        Index("idx_jobs_is_remote_city", is_remote, city),
    )

    def __repr__ (self):
        return f"<Job(id={self.id}, title={self.title}, company={self.company})>"