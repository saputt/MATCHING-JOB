import { Injectable } from "@nestjs/common";
import { PrismaService } from "src/prisma/prisma.service";

@Injectable()
export class AnalysisRepository {
    constructor(private prisma : PrismaService) {}

    async findAnalysis(userId : string, jobId : string) {
        return await this.prisma.analysis.findUnique({
            where : {
                userId_jobId : {
                    jobId,
                    userId
                }
            }
        })
    }
}