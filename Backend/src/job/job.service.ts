import { Injectable, InternalServerErrorException, NotFoundException } from "@nestjs/common";
import { JobRepository } from "./job.repository";

@Injectable()
export class JobService {
    constructor(private readonly repo : JobRepository) {}

    async isJobExist(jobId : string) {
        const job = await this.repo.findJobById(jobId)
        if (!job) throw new NotFoundException("job not found")
        return job
    }

    async findTopRecommendation(userVectorString : string) {
        try {
            const topJobs = await this.repo.findTopRecommendations(userVectorString)

            return topJobs.map(job => {
                const rawScore = job.matchScore ? Number(job.matchScore) : 0;
                
                return {
                    ...job,
                    matchScore: Math.round(rawScore),
                };
            });
        } catch (error) {
            throw new InternalServerErrorException('failed to get recommendation job');
        }
    }

    async calculateSingleScore(jobId : string, userVectorString : string) {
        try {
            const result = await this.repo.calculateSingleStore(jobId, userVectorString)
            
            if (!result || result.length === 0 || result[0].matchScore === null) {
                throw new NotFoundException(`Job with ID ${jobId} doesnt exist in database`);
            }
    
            const rawScore = result[0].matchScore
            
            if (typeof rawScore !== 'number') {
                return 0;
            }

            return Math.round(rawScore)
        } catch (error) {
            if (error instanceof NotFoundException) throw error;
            
            console.error("[SERVICE LLM ERROR]:", error);

            throw new InternalServerErrorException('Failed to calculate the vector similarity score.');
        }
    }
}