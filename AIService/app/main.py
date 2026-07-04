from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api.job_routes import router as ai_router
from app.api.embedding_routes import router as embedding_router
from app.api.analysis_routes import router as analysis_router

import uvicorn
from app.core.config import Settings

app = FastAPI(
    title="Hyperlocal IT Matchmaker AI API",
    description="Backend AI Pipeline & Services Platform menggunakan Clean Architecture SoC",
    version="1.0.0"
)
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000", "http://localhost:8080"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(ai_router)
app.include_router(embedding_router)
app.include_router(analysis_router)

if __name__ == "__main__":
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=Settings.PORT or 5000,
        reload=True
    )