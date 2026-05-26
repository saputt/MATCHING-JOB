# file ini bertanggung jawab untuk:
# 1. Membuat koneksi ke database
# 2. Menyediakan session factory untuk operasi ke database
# 3. Menyediakan base class untuk SQLAlchemy models
import os
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker, DeclarativeBase
from dotenv import load_dotenv
from app.core.config import Settings

#membuat database engine atau koneksi pool
#echo untuk menampilkan query sql di console untuk melakukan debugging
engine = create_engine(Settings.DATABASE_URL, echo=True, future=True)

# session adalah pabrik untuk membuat session database
# session digunakan untuk melakukan CRUD
SessionLocal = sessionmaker(autoflush=False, bind=engine)

# Base adalah class yang akan di inherit oleh models yang akan dibuat
class Base(DeclarativeBase):
    pass

# fungsi ini akan dipanggil setiap ada request ke fastpi
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()