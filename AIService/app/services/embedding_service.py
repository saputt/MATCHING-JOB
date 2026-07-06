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

    def create_user_narrative(self, minat_posisi : str, user_skills : list) -> str:
        skills_str = ", ".join(user_skills)
        return f"Minat Posisi: {minat_posisi}. Memiliki keahlian teknis dan teknologi: {skills_str}."

    def generate_embedding(self, title : str, hard_skills : list, isUser : bool | None = None) -> list[float]:
        try:
            if (isUser):
                narrative_text = self.create_user_narrative(minat_posisi=title, user_skills=hard_skills)
            else:
                narrative_text = self.create_narrative(title=title, hard_skills=hard_skills)
            
            embedding = self.client.feature_extraction(
                narrative_text,
                model=self.model_name
            )

            if hasattr(embedding, "tolist"):
                return embedding.tolist()
                
            return list(embedding)
        except Exception as e:
            raise Exception(f"MiniLM Cloud embedding calculation failed: {str(e)}")

