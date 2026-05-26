from app.features.market.repository import MarketRepository
from app.features.market.service import MarketService
from fastapi import APIRouter, Depends, HTTPException
from sqlalchemy.orm import Session
from app.core.database import get_db
from app.features.market.schema import (
    MarketInsightResponse,
    AnalyzeMarketRequest
)

router = APIRouter(prefix="/analyze-market", tags=["Market Analysis"])

# endpoint untuk menganalisis pasar kerja berdasarkan skill user
# endpoint ini memiliki beberapa alur kerja
# 1. setup repository and service
# 2. panggil service untuk top skill
# 3. panggil service untuk top jobs
# 4. panggil service untuk skill gap
@router.post("", response_model=MarketInsightResponse)
async def analyze_market(
    req : AnalyzeMarketRequest,
    db : Session = Depends(get_db)
):
    repo = MarketRepository(db)
    service = MarketService(repo)
    return service.analyze_market(req.skills)
