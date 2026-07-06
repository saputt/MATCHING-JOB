import { Module } from '@nestjs/common';
import { PrismaModule } from './prisma/prisma.module';
import { UserModule } from './user/user.module';
import { AuthModule } from './auth/auth.module';
import { ConfigModule } from '@nestjs/config';
import { AiModule } from './ai/ai.module';
import { AnalysisModule } from './analysis/analysis.module';
import { JobModule } from './job/job.module';
import { MatchingModule } from './matching/matching.module';

@Module({
  imports: [PrismaModule, UserModule, AuthModule, AiModule, AnalysisModule, JobModule, MatchingModule, ConfigModule.forRoot({isGlobal : true})],
})
export class AppModule {}
