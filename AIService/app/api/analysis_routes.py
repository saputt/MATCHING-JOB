from pydantic import BaseModel
from app.services.analyze_service import AnalyzeService
from fastapi import APIRouter, HTTPException, Depends
from app.repositories.analyze_repository import AnalyzeRepository
from app.adapters.llm_adapter import LLMAdapter
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.schemas import CreateAnalysisRequest, GetAnalysisRequest
from uuid import UUID
import traceback

router = APIRouter(
    prefix="/api",
    tags=["Text Embedding"]
)

@router.post("/analyze/{job_id}")
def analyze_job(job_id : str, req : CreateAnalysisRequest, db : Session = Depends(get_db)):
    try:
        repo = AnalyzeRepository(db=db)
        llm = LLMAdapter()
        service = AnalyzeService(repo=repo, llm=llm)
        res = service.analyze_job(req=req, job_id=job_id, user_id=str(req.user_id))
        return {
            "status" : "success",
            "message" : "create new analyze job success",
            "data" : res
        }
    except Exception as e:
        traceback.print_exc()
        raise HTTPException(status_code=500, detail=f"Gagal mengeksekusi pipeline: {str(e)}")

@router.get("/analyze/{job_id}")
def get_analyze_job(job_id : str, req : GetAnalysisRequest, db : Session = Depends(get_db)):
    try:
        repo = AnalyzeRepository(db=db)
        llm = LLMAdapter()
        service = AnalyzeService(repo=repo, llm=llm)
        res = service.get_analyze_by_job_id(job_id=job_id, user_id=req.user_id)
        return {
            "status" : "success",
            "message" : "get analyze job success",
            "data" : res
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Gagal mengeksekusi pipeline: {str(e)}")