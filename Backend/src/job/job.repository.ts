import { Injectable, InternalServerErrorException } from "@nestjs/common";
import { PrismaService } from "src/prisma/prisma.service";

@Injectable()
export class JobRepository {
    constructor(private readonly prisma : PrismaService) {}

    async findTopRecommendations(userVectorString: string): Promise<any[]> {
        return await this.prisma.$queryRawUnsafe<any[]>(`
            SELECT 
            id, 
            title, 
            company, 
            description, 
            url, 
            city, 
            skills, 
            salary,
            (1 - (embedding <=> CAST('${userVectorString}' AS vector))) * 100 AS "matchScore"
            FROM jobs
            ORDER BY embedding <=> CAST('${userVectorString}' AS vector) ASC
            LIMIT 10;
        `);
    }

    async calculateSingleStore(jobId : string, userVectorString : string) {
        return await this.prisma.$queryRaw<any[]>`
            SELECT 
                (1 - (embedding <=> ${userVectorString}::vector)) * 100 AS "matchScore"
            FROM jobs
            WHERE id = ${jobId}
            LIMIT 1;
        `;0
    }

    async findJobById(jobId : string) {
        return this.prisma.job.findUnique({
            where : {
                id : jobId
            }
        })
    }
}