import { Module } from "@nestjs/common";
import { MatchingController } from "./matching.controller";
import { MatchingService } from "./matching.service";
import { AiModule } from "src/ai/ai.module";
import { UserModule } from "src/user/user.module";
import { JobModule } from "src/job/job.module";

@Module({
    controllers : [MatchingController],
    providers : [MatchingService],
    imports : [AiModule, UserModule, JobModule]
})
export class MatchingModule {}