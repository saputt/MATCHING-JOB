import { Module } from "@nestjs/common";
import { JobRepository } from "./job.repository";
import { PrismaModule } from "src/prisma/prisma.module";
import { JobService } from "./job.service";
import { AuthModule } from "src/auth/auth.module";

@Module({
    providers : [JobRepository, JobService],
    imports : [PrismaModule, AuthModule],
    exports : [JobService]
})
export class JobModule {}