import { Module } from "@nestjs/common";
import { AnalysisService } from "./analysis.service";
import { AnalysisRepository } from "./analysis.repository";
import { AnalysisController } from "./analysis.controller";
import { PrismaModule } from "src/prisma/prisma.module";
import { UserModule } from "src/user/user.module";
import { AiModule } from "src/ai/ai.module";
import { JobModule } from "src/job/job.module";

@Module({
    providers : [AnalysisService, AnalysisRepository],
    controllers : [AnalysisController],
    imports : [PrismaModule, UserModule, AiModule, JobModule]
})
export class AnalysisModule {}