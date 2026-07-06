import { Injectable, InternalServerErrorException } from "@nestjs/common";
import { ConfigService } from "@nestjs/config";

export interface JobAnalysis {
    jobTitle : string;
    jobHardSkills : string[];
    jobSoftSkills : string[];
    userSkills : string[];
    matchScore : number;
    jobId : string
}

export interface GenerateEmbedding {
    hardSkills : string[];
    title : string;
}

@Injectable()
export class AiServiceClient {
    private readonly baseUrl : string

    constructor(private configService : ConfigService) {
        this.baseUrl = this.configService.get<string>('AI_URL') || 'http://127.0.0.1:8000/api'
    }

    async triggerJobAnalysis(userId : string, dto : JobAnalysis) {
        try {
            const res = await fetch(`${this.baseUrl}/analyze/${dto.jobId}`, {
                method : "POST",
                headers : {
                    'Content-Type': 'application/json',
                },
                body : JSON.stringify({
                    user_id : userId,
                    job_title : dto.jobTitle,
                    job_hard_skills : dto.jobHardSkills,
                    job_soft_skills : dto.jobSoftSkills,
                    user_skills : dto.userSkills,
                    match_score : dto.matchScore
                })
            })

            if (!res.ok) throw new Error(`Python AI Service returned status: ${res.status}`)

            return await res.json() 
        } catch (error) {
            throw new InternalServerErrorException("failed to communication with python service")
        }
    }

    async generateEmbedding(dto : GenerateEmbedding) {
        try {
            const res = await fetch(`${this.baseUrl}/embedding`, {
                method : "POST",
                headers : {
                    'Content-Type': 'application/json',
                },
                body : JSON.stringify({
                    hard_skills : dto.hardSkills,
                    title : dto.title
                })
            })

            if (!res.ok) throw new Error(`Python AI Service returned status: ${res.status}`)

            return await res.json() 
        } catch (error) {
            throw new InternalServerErrorException("failed to communication with python service")
        }
    }
}