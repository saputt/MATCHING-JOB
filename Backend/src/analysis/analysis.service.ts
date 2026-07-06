import { BadRequestException, Injectable } from "@nestjs/common";
import { JobService } from "src/job/job.service";
import { UserService } from "src/user/user.service";
import { AnalysisRepository } from "./analysis.repository";
import { AiServiceClient, JobAnalysis } from "src/ai/ai-service.client";

@Injectable()
export class AnalysisService {
    constructor(
        private readonly jobService : JobService,
        private readonly userService : UserService,
        private readonly analysisRepo : AnalysisRepository,
        private readonly aiService : AiServiceClient,
    ) {}

    async isAnalysisAlready(userId : string, jobId : string) {
        const analysis = await this.analysisRepo.findAnalysis(userId, jobId)
        if (analysis) return analysis
    }
 
    async getOrCreateAnalysis(userId : string, jobId : string) {
        await this.isAnalysisAlready(userId, jobId)
        const embedding = await this.userService.getUserEmbedding(userId)
        const user = await this.userService.isUserExist(userId)
        const job = await this.jobService.isJobExist(jobId)

        if (!embedding) {
            throw new BadRequestException("User doesnt have embedding. setup preference first");
        }

        const finalScore = await this.jobService.calculateSingleScore(jobId, embedding)

        const payloadAnalysis : JobAnalysis = {
            jobId,
            jobHardSkills : job.skills,
            jobSoftSkills : job.softskills,
            jobTitle : job.title!,
            matchScore : finalScore,
            userSkills : user.skills
        }

        const aiResult = await this.aiService.triggerJobAnalysis(userId, payloadAnalysis)

        return aiResult
    }
}