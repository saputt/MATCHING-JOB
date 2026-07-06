from pydantic import BaseModel
from app.services.embedding_service import EmbeddingService
from fastapi import APIRouter
from app.schemas import EmbeddingRequest

router = APIRouter(
    prefix="/api",
    tags=["Text Embedding"]
)

@router.post("/embedding")
def generate_embedding(dto : EmbeddingRequest):
    service = EmbeddingService()
    text_embedding = service.generate_embedding(hard_skills=dto.hard_skills,title=dto.title, isUser=True)
    return {
        "status" : "",
        "message" : "",
        "data" : text_embedding
    }
