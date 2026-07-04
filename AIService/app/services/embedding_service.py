import os
from huggingface_hub import InferenceClient
from dotenv import load_dotenv

load_dotenv()

class EmbeddingService:
    def __init__(self):
        self.client = InferenceClient(
            provider="hf-inference",
            api_key=os.environ.get("HF_TOKEN")
        )
        self.model_name="sentence-transformers/all-MiniLM-L6-v2"

    def create_narrative(self, title : str, hard_skills : list) -> str:
        skills_str = ", ".join(hard_skills)
        return f"Posisi: {title}. Membutuhkan keahlian teknis dan teknologi: {skills_str}."

    def generate_embedding(self, title : str, hard_skills : list) -> list[float]:
        try:
            narrative_text = self.create_narrative(title=title, hard_skills=hard_skills)
            embedding = self.client.feature_extraction(
                narrative_text,
                model=self.model_name
            )
            return embedding
        except Exception as e:
            raise Exception(f"❌ Gagal hitung embedding MiniLM Cloud: {str(e)}")

