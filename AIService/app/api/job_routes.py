from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.repositories.job_repository import JobRepository
from app.services.pipeline_service import PipelineService
from app.adapters.llm_adapter import LLMAdapter

router = APIRouter(
    prefix="/api",
    tags=["AI Data Processing"]
)

@router.post("/patcher-it", response_model=None)
def patcher_it(db : Session = Depends(get_db)):
    """
        
    """
    try:
        repo = JobRepository(db)
        llm = LLMAdapter()
        service = PipelineService(repo=repo, llm=llm)

        service.jobs_patcher(db=db)

        return {
            "status": "success",
            "message": f"Berhasil memproses maksimal antrean data di background"
        }

    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Gagal mengeksekusi pipeline: {str(e)}")

@router.delete("/non-it")
def delete_non_it_jobs(db : Session = Depends(get_db)):
    try:
        repo = JobRepository(db)
        llm = LLMAdapter()
        service = PipelineService(repo=repo, llm=llm)
        
        job_non_it = service.delete_non_it()

        return {
            "status" : "success",
            "message" : f"berhasil mengahapus job non it sebanyak {job_non_it}"
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Gagal mengeksekusi pipeline: {str(e)}")

