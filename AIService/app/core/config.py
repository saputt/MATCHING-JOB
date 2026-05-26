# file ini berfungsi sebagai settingan konfigurasi aplikasi
import os
from dotenv import load_dotenv

load_dotenv()

class Settings:
    DATABASE_URL: str = os.getenv("DATABASE_URL")
    PORT: int = int(os.getenv("PORT", "5000"))
