from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.features.market.router import router as market_router
import uvicorn
from app.core.config import Settings

app = FastAPI()
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000", "http://localhost:8080"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(market_router)

if __name__ == "__main__":
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=Settings.PORT | 5000,
        reload=True
    )